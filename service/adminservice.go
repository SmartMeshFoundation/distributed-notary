package service

import (
	"errors"
	"time"

	"fmt"

	"bytes"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	ethevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/service/messagetosign"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nkbai/log"
)

// AdminService :
type AdminService struct {
	db              *models.DB
	notaryService   *NotaryService
	dispatchService dispatchServiceBackend
}

// NewAdminService :
func NewAdminService(db *models.DB, notaryService *NotaryService, dispatchService dispatchServiceBackend) (ns *AdminService, err error) {
	ns = &AdminService{
		db:              db,
		notaryService:   notaryService,
		dispatchService: dispatchService,
	}
	return
}

// OnEvent 链上事件逻辑处理
func (as *AdminService) OnEvent(e chain.Event) {
	// TODO 处理新块事件,保存各链最新块号
}

// OnRequest restful请求处理
func (as *AdminService) OnRequest(req api.Request) {
	switch r := req.(type) {
	/*
		user api
	*/
	case *userapi.GetNotaryListRequest:
		as.onGetNotaryListRequest(r)
	case *userapi.GetSCTokenListRequest:
		as.onGetSCTokenListRequest(r)
	/*
		admin api
	*/
	case *userapi.GetPrivateKeyListRequest:
		as.onGetPrivateKeyListRequest(r)
	case *userapi.CreatePrivateKeyRequest:
		as.onCreatePrivateKeyRequest(r)
	case *userapi.RegisterSCTokenRequest:
		as.onRegisterSCTokenRequest(r)
	case *notaryapi.NewSCTokenRequest:
		as.onNewSCTokenRequest(r)
	/*
		debug api
	*/
	case *userapi.DebugTransferToAccountRequest:
		as.onDebugTransferToAccountRequest(r)
	default:
		req.WriteErrorResponse(api.ErrorCodeParamsWrong)
		return
	}
	return
}

// 公证人列表查询
func (as *AdminService) onGetNotaryListRequest(req *userapi.GetNotaryListRequest) {
	if req.GetRequestName() != userapi.APIUserNameGetNotaryList {
		req.WriteErrorResponse(api.ErrorCodeParamsWrong)
		return
	}
	notaries, err := as.db.GetNotaryInfo()
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
	}
	req.WriteSuccessResponse(notaries)
}

type scTokenInfoToResponse struct {
	SCToken                  common.Address `json:"sc_token"`                               // 侧链token地址
	SCTokenName              string         `json:"sc_token_name"`                          // 侧链Token名
	SCTokenOwnerKey          common.Hash    `json:"sc_token_owner_key"`                     // 侧链token合约owner的key
	MCLockedContractAddress  common.Address `json:"mc_locked_contract_address"`             // 对应主链锁定合约地址
	MCName                   string         `json:"mc_name"`                                // 对应主链名
	MCLockedContractOwnerKey common.Hash    `json:"mc_locked_contract_owner_key,omitempty"` // 对应主链锁定合约owner的key
	CreateTime               string         `json:"create_time"`                            // 创建时间
	OrganiserID              int            `json:"organiser_id"`                           // 发起人ID
}

func newSCTokenInfoToResponse(s *models.SideChainTokenMetaInfo) (r *scTokenInfoToResponse) {
	r = &scTokenInfoToResponse{
		SCToken:                  s.SCToken,
		SCTokenName:              s.SCTokenName,
		SCTokenOwnerKey:          s.SCTokenOwnerKey,
		MCLockedContractAddress:  s.MCLockedContractAddress,
		MCName:                   s.MCName,
		MCLockedContractOwnerKey: s.MCLockedContractOwnerKey,
		OrganiserID:              s.OrganiserID,
		CreateTime:               time.Unix(s.CreateTime, 0).String(),
	}
	return r
}

// SCToken列表查询
func (as *AdminService) onGetSCTokenListRequest(req *userapi.GetSCTokenListRequest) {
	if req.GetRequestName() != userapi.APIUserNameGetSCTokenList {
		req.WriteErrorResponse(api.ErrorCodeParamsWrong)
		return
	}
	tokens, err := as.db.GetSCTokenMetaInfoList()
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
	}
	var resp []*scTokenInfoToResponse
	for _, token := range tokens {
		resp = append(resp, newSCTokenInfoToResponse(token))
	}
	req.WriteSuccessResponse(resp)
}

type privateKeyInfoToResponse struct {
	ID         string `json:"id"`
	Address    string `json:"address,omitempty"`
	Status     int    `json:"status"`
	StatusMsg  string `json:"status_msg"`
	CreateTime string `json:"create_time"`
}

func newPrivateKeyInfoToResponse(p *models.PrivateKeyInfo) (r *privateKeyInfoToResponse) {
	r = &privateKeyInfoToResponse{
		ID:         p.Key.String(),
		Status:     p.Status,
		StatusMsg:  models.PrivateKeyInfoStatusMsgMap[p.Status],
		CreateTime: time.Unix(p.CreateTime, 0).String(),
	}
	if p.Status == models.PrivateKeyNegotiateStatusFinished {
		r.Address = p.ToAddress().String()
	}
	return
}

// 私钥列表查询
func (as *AdminService) onGetPrivateKeyListRequest(req *userapi.GetPrivateKeyListRequest) {
	if req.GetRequestName() != userapi.APIAdminNameGetPrivateKeyList {
		req.WriteErrorResponse(api.ErrorCodeParamsWrong)
		return
	}
	privateKeyInfoList, err := as.db.GetPrivateKeyList()
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
	}
	var respList []*privateKeyInfoToResponse
	for _, privateKeyInfo := range privateKeyInfoList {
		respList = append(respList, newPrivateKeyInfoToResponse(privateKeyInfo))
	}
	req.WriteSuccessResponse(respList)
}

/*
发起一次公钥-私钥片协商过程,并等待协商结果
*/
func (as *AdminService) onCreatePrivateKeyRequest(req *userapi.CreatePrivateKeyRequest) {
	if req.GetRequestName() != userapi.APIAdminNameCreatePrivateKey {
		req.WriteErrorResponse(api.ErrorCodeParamsWrong)
		return
	}
	// 1. 调用自己的notaryService,生成KeyGenerator,并开始协商过程
	privateKeyID, err := as.notaryService.startNewPrivateKeyNegotiation()
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	// 2. 使用PrivateKeyID轮询数据库,直到该key协商并生成完成
	times := 0
	for {
		time.Sleep(time.Second) // TODO 这里轮询周期设置为多少合适,是否需要设置超时
		privateKey, err := as.db.LoadPrivateKeyInfo(privateKeyID)
		if err != nil {
			log.Error(err.Error())
			req.WriteErrorResponse(api.ErrorCodeException, err.Error())
			return
		}
		if privateKey.Status != models.PrivateKeyNegotiateStatusFinished {
			if times%10 == 0 {
				log.Trace(SessionLogMsg(privateKeyID, "waiting for PrivateKeyNegotiate..."))
			}
			times++
			continue
		}
		req.WriteSuccessResponse(newPrivateKeyInfoToResponse(privateKey))
		return
	}
}

/*
使用某个私钥片创建一组新的合约
*/
func (as *AdminService) onRegisterSCTokenRequest(req *userapi.RegisterSCTokenRequest) {
	if req.GetRequestName() != userapi.APIAdminNameRegisterNewSCToken {
		req.WriteErrorResponse(api.ErrorCodeParamsWrong)
		return
	}
	// 1. 校验私钥ID可用性
	privateKeyInfo, err := as.db.LoadPrivateKeyInfo(common.HexToHash(req.PrivateKeyID))
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeDataNotFound, err.Error())
		return
	}
	// 3. 部署主链合约
	mcContractAddress, mcDeploySessionID, err := as.distributedDeployMCContact(req.MainChainName, privateKeyInfo)
	if err != nil {
		err = fmt.Errorf("err when distributedDeployMCContact : %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	log.Info("deploy MCContract success, contract address = %s", mcContractAddress.String())
	// TODO 这里是否需要中途存储
	// 4. 部署侧链合约
	scTokenAddress, scDeploySessionID, err := as.distributedDeploySCToken(privateKeyInfo)
	if err != nil {
		err = fmt.Errorf("err when distributedDeploySCToken : %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	log.Info("deploy SCTokenContract success, contract address = %s", scTokenAddress.String())
	// 5. 构造SideChainTokenMetaInfo并存储
	var scTokenMetaInfo models.SideChainTokenMetaInfo
	scTokenMetaInfo.MCName = req.MainChainName
	scTokenMetaInfo.SCTokenName = req.MainChainName + "-Token"
	scTokenMetaInfo.SCToken = scTokenAddress
	scTokenMetaInfo.SCTokenOwnerKey = privateKeyInfo.Key
	scTokenMetaInfo.SCTokenDeploySessionID = scDeploySessionID
	scTokenMetaInfo.MCLockedContractAddress = mcContractAddress
	scTokenMetaInfo.MCLockedContractOwnerKey = privateKeyInfo.Key
	scTokenMetaInfo.MCLockedContractDeploySessionID = mcDeploySessionID
	scTokenMetaInfo.CreateTime = time.Now().Unix()
	scTokenMetaInfo.OrganiserID = as.notaryService.self.ID
	err = as.db.NewSCTokenMetaInfo(&scTokenMetaInfo)
	if err != nil {
		log.Error("err when NewSCTokenMetaInfo : %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	// 5. 向DispatchService注册新SCToken
	err = as.dispatchService.registerNewSCToken(&scTokenMetaInfo)
	if err != nil {
		log.Error("err when registerNewSCToken on %s: %s", scTokenMetaInfo.MCName, err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	// 6. 通知其余公证人
	newSCTokenReq := notaryapi.NewNewSCTokenRequest(&scTokenMetaInfo)
	err = as.notaryService.BroadcastMsg(utils.EmptyHash, notaryapi.APIAdminNameNewSCToken, newSCTokenReq, true)
	if err != nil {
		log.Error("err when broadcast NewSCTokenRequest to other notaries err=", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	// 7. 返回
	req.WriteSuccessResponse(scTokenMetaInfo)
}

func (as *AdminService) distributedDeployMCContact(chainName string, privateKeyInfo *models.PrivateKeyInfo) (contractAddress common.Address, sessionID common.Hash, err error) {
	if chainName != ethevents.ChainName {
		err = errors.New("only support ethereum as main chain now")
		return
	}
	var c chain.Chain
	c, err = as.dispatchService.getChainByName(chainName)
	// 暂时主链只有ethereum,复用spcetrum的signer
	return as.distributedDeployOnSpectrum(c, privateKeyInfo)
}

func (as *AdminService) distributedDeploySCToken(privateKeyInfo *models.PrivateKeyInfo) (contractAddress common.Address, sessionID common.Hash, err error) {
	var c chain.Chain
	c, err = as.dispatchService.getChainByName(smcevents.ChainName)
	if err != nil {
		return
	}
	tokenName := c.GetChainName() + params.SCTokenNameSuffix
	return as.distributedDeployOnSpectrum(c, privateKeyInfo, tokenName)
}

func (as *AdminService) distributedDeployOnSpectrum(c chain.Chain, privateKeyInfo *models.PrivateKeyInfo, params ...string) (contractAddress common.Address, sessionID common.Hash, err error) {
	// 1. 获取待签名的数据
	var msgToSign mecdsa.MessageToSign
	msgToSign = messagetosign.NewSpectrumContractDeployTX(c, privateKeyInfo.ToAddress(), params...)
	// 2. 签名
	var signature []byte
	signature, sessionID, err = as.notaryService.startDistributedSignAndWait(msgToSign, privateKeyInfo)
	if err != nil {
		return
	}
	log.Info("deploy contract on %s with account=%s, signature=%s", c.GetChainName(), privateKeyInfo.ToAddress().String(), common.Bytes2Hex(signature))
	// 4. 部署合约
	transactor := &bind.TransactOpts{
		From: privateKeyInfo.ToAddress(),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != privateKeyInfo.ToAddress() {
				return nil, errors.New("not authorized to sign this account")
			}
			msgToSign2 := signer.Hash(tx).Bytes()
			if bytes.Compare(msgToSign.GetSignBytes(), msgToSign2) != 0 {
				err = fmt.Errorf("txbytes when deploy contract step1 and step2 does't match")
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
	contractAddress, err = c.DeployContract(transactor, params...)
	return
}

func (as *AdminService) onNewSCTokenRequest(req *notaryapi.NewSCTokenRequest) {
	var err error
	scTokenMetaInfo := req.SCTokenMetaInfo
	// 1. 校验信息 TODO 需要验证两个sessionID在本地是否存在且状态为签名完成
	/*
		TODO 需要验证的信息 :
		a. 主链合约状态,对应签名在本地是否有记录,跟请求是否匹配
		b. 侧链合约状态,对应签名在本地是否有记录,跟请求是否匹配
	*/
	// 2. 保存
	err = as.db.NewSCTokenMetaInfo(scTokenMetaInfo)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException)
		return
	}
	// 3. 注册到DispatchService并开始提供服务
	err = as.dispatchService.registerNewSCToken(scTokenMetaInfo)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException)
		return
	}
	req.WriteSuccessResponse(nil)
}
