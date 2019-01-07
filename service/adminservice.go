package service

import (
	"errors"
	"time"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nkbai/log"
)

// AdminService :
type AdminService struct {
	db            *models.DB
	notaryService *NotaryService
	chainMap      map[string]chain.Chain
}

// NewAdminService :
func NewAdminService(db *models.DB, notaryService *NotaryService, chainMap map[string]chain.Chain) (ns *AdminService, err error) {
	ns = &AdminService{
		db:            db,
		notaryService: notaryService,
		chainMap:      chainMap,
	}
	return
}

// OnEvent 链上事件逻辑处理
func (as *AdminService) OnEvent(e chain.Event) {
	// TODO 处理新块事件,保存各链最新块号
}

// OnRequest restful请求处理
func (as *AdminService) OnRequest(req api.Request) {
	//TODO
	switch r := req.(type) {
	/*
		user api
	*/
	case *userapi.GetNotaryListRequest:
		as.onGetNotaryListRequest(r)
	/*
		admin api
	*/
	case *userapi.GetPrivateKeyListRequest:
		as.onGetPrivateKeyListRequest(r)
	case *userapi.CreatePrivateKeyRequest:
		as.onCreatePrivateKeyRequest(r)
	case *userapi.RegisterSCTokenRequest:
		as.onRegisterSCTokenRequest(r)
	}
	return
}

// 公证人列表查询
func (as *AdminService) onGetNotaryListRequest(req *userapi.GetNotaryListRequest) {
	notaries, err := as.db.GetNotaryInfo()
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
	}
	req.WriteSuccessResponse(notaries)
}

type privateKeyInfoToResponse struct {
	ID        string `json:"id"`
	Address   string `json:"address,omitempty"`
	Status    int    `json:"status"`
	StatusMsg string `json:"status_msg"`
}

func newPrivateKeyInfoToResponse(p *models.PrivateKeyInfo) (r *privateKeyInfoToResponse) {
	r = &privateKeyInfoToResponse{
		ID:        p.Key.String(),
		Status:    p.Status,
		StatusMsg: models.PrivateKeyInfoStatusMsgMap[p.Status],
	}
	if p.Status == models.PrivateKeyNegotiateStatusFinished {
		r.Address = p.ToAddress().String()
	}
	return
}

// 私钥列表查询
func (as *AdminService) onGetPrivateKeyListRequest(req *userapi.GetPrivateKeyListRequest) {
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
	// 1. 校验私钥ID可用性
	privateKeyInfo, err := as.db.LoadPrivateKeyInfo(common.HexToHash(req.PrivateKeyID))
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeDataNotFound, err.Error())
		return
	}
	// 3. 部署主链合约
	mcContractAddress, err := as.distributedDeployMCContact(req.MainChainName, privateKeyInfo)
	if err != nil {
		err = fmt.Errorf("err when distributedDeployMCContact : %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	// TODO 这里是否需要中途存储
	// 4. 部署侧链合约
	scTokenAddress, err := as.distributedDeploySCToken(privateKeyInfo)
	if err != nil {
		err = fmt.Errorf("err when distributedDeploySCToken : %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	// 5. 构造SideChainTokenMetaInfo并存储
	var scTokenMetaInfo models.SideChainTokenMetaInfo
	scTokenMetaInfo.MCName = req.MainChainName
	scTokenMetaInfo.SCTokenName = req.MainChainName + "-Token"
	scTokenMetaInfo.SCToken = scTokenAddress
	scTokenMetaInfo.SCTokenOwnerKey = privateKeyInfo.Key
	scTokenMetaInfo.MCLockedContractAddress = mcContractAddress
	scTokenMetaInfo.MCLockedContractOwnerKey = privateKeyInfo.Key
	err = as.db.NewSCTokenMetaInfo(&scTokenMetaInfo)
	if err != nil {
		log.Error("err when NewSCTokenMetaInfo : %s", err.Error())
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	// 6. 构造CrossChainService开始提供服务
	// TODO
	// 7. 返回
	req.WriteSuccessResponse(scTokenMetaInfo)
}

func (as *AdminService) distributedDeployMCContact(chainName string, privateKeyInfo *models.PrivateKeyInfo) (contractAddress common.Address, err error) {
	var c chain.Chain
	c, ok := as.chainMap[chainName]
	if !ok && chainName != smcevents.ChainName {
		err = errors.New("only support ethereum as main chain now")
		return
	}
	// 暂时主链只有ethereum,复用spcetrum的signer
	return as.distributedDeployOnSpectrum(c, privateKeyInfo)
}

func (as *AdminService) distributedDeploySCToken(privateKeyInfo *models.PrivateKeyInfo) (contractAddress common.Address, err error) {
	c := as.chainMap[smcevents.ChainName]
	return as.distributedDeployOnSpectrum(c, privateKeyInfo)
}

func (as *AdminService) distributedDeployOnSpectrum(c chain.Chain, privateKeyInfo *models.PrivateKeyInfo) (contractAddress common.Address, err error) {
	// 1. 获取待签名的数据
	var msgToSign mecdsa.MessageToSign
	msgToSign = NewSpectrumContractDeployTX(c, privateKeyInfo.ToAddress())
	// 2. 签名
	var signature []byte
	signature, err = as.notaryService.startDistributedSignAndWait(msgToSign, privateKeyInfo)
	if err != nil {
		return
	}
	// 3. 部署合约
	transactor := &bind.TransactOpts{
		From: privateKeyInfo.ToAddress(),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != privateKeyInfo.ToAddress() {
				return nil, errors.New("not authorized to sign this account")
			}
			return tx.WithSignature(signer, signature)
		},
	}
	return c.DeployContract(transactor)
}
