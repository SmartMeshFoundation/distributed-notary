package service

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/service/messagetosign"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	sevents "github.com/SmartMeshFoundation/distributed-notary/chainjettrade/spectrum/events"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	japi "github.com/SmartMeshFoundation/distributed-notary/api/userapi/jettradeapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
)

/*
1. 处理请求
2. 记录连上事件
*/
type JettradeService struct {
	ds                   *DispatchService
	notaryAddress        common.Address
	notaryPrivateKeyInfo *models.PrivateKeyInfo
	jh                   *jhelper
}

func newJettradeService(ds *DispatchService) *JettradeService {
	js := &JettradeService{ds: ds}
	js.jh = NewJHelper(js)
	return js
}
func (js *JettradeService) init() (err error) {
	notoary1, err := js.ds.spectrumJettradeService.Proxy.Notary(nil)
	if err != nil {
		return
	}
	notary2, err := js.ds.ethJettradeService.Proxy.Notary(nil)
	if err != nil {
		return
	}
	if bytes.Compare(notoary1.Bytes(), notary2.Bytes()) != 0 {
		panic("contract not match")
	}
	js.notaryAddress = notoary1
	lists, err := js.ds.db.GetPrivateKeyList()
	if err != nil {
		return
	}
	found := false
	for _, l := range lists {
		if bytes.Compare(l.Address.Bytes(), js.notaryAddress.Bytes()) == 0 {
			found = true
			js.notaryPrivateKeyInfo = l
		}
	}
	if !found {
		err = errors.New("cannot found notary private key")
		return
	}
	return
}
func (js *JettradeService) onEvent(e chain.Event) {
	je, isJe := e.(chainjettrade.IsJettradeEvent)
	if !isJe {
		log.Error("dispatch to JettradeService, e=%s", log.StringInterface(e, 3))
		return
	}
	s := je.GetShareData()
	var txHash common.Hash
	if io, ok := e.(*chainjettrade.IssueDocumentPOEvent); ok {
		txHash = io.TxHash
	}
	err := js.ds.db.NewJettradeEventInfo(models.NewJettradeEventInfo(e.GetChainName(), string(e.GetEventName()), e.GetFromAddress(), s.From, s.To, e.GetBlockNumber(), s.TokenID, txHash))
	if err != nil {
		log.Error("NewJettradeEventInfo err %s", err)
	}
	go js.jh.onEvent(e)
}

func (js *JettradeService) onRequest(req api.Req) {
	var ei *models.JettradeEventInfo
	var err error
	switch req2 := req.(type) {
	case *japi.IssuePOOnSpectrumRequest:
		ei, err = js.callIssuePOOnSpectrum(req2)
	case *japi.SignDOBuyerOnSpectrumRequest:
		ei, err = js.callSignDOBuyerOnSpectrum(req2)
	case *japi.SignPOONEthereumRequest:
		ei, err = js.callSignPOOnEth(req2)
	case *japi.SignDOFFOnEthereumRequest:
		ei, err = js.callSignDOFFOnEth(req2)
	case *japi.IssueINVOnEthereumRequest:
		ei, err = js.callIssueINVOnEth(req2)
	default:
		panic(fmt.Sprintf("unkonw req=%s", log.StringInterface(req, 3)))
	}
	//一定是ReqWithResponse
	reqWithResponse := req.(api.ReqWithResponse)
	if err != nil {
		reqWithResponse.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	ei.NotaryIDInCharge = js.ds.getSelfNotaryInfo().ID
	err = js.ds.db.UpdateJettradeEventInfo(ei)
	if err != nil {
		reqWithResponse.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	reqWithResponse.WriteSuccessResponse(ei)
}
func (js *JettradeService) onIssuePORequest(req *japi.IssuePOOnSpectrumRequest) {

	//2. 发起合约调用
	ei, err := js.callIssuePOOnSpectrum(req)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	ei.NotaryIDInCharge = js.ds.getSelfNotaryInfo().ID
	err = js.ds.db.UpdateJettradeEventInfo(ei)
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
		return
	}
	req.WriteSuccessResponse(ei)
}
func (js *JettradeService) callIssuePOOnSpectrum(req *japi.IssuePOOnSpectrumRequest) (ei *models.JettradeEventInfo, err error) {
	//0. 校验相应的事件存在, 发生在eth上
	ei, err = js.ds.db.GetJettradeEventInfo(cfg.ETH.Name, sevents.EventNameIssueDocumentPO, js.ds.ethJettradeService.ContractAddress, req.TokenId)
	if err != nil {
		return
	}
	//1. 获取nonce
	nonce, err := js.ds.applyNonceFromNonceServer(cfg.SMC.Name, js.notaryPrivateKeyInfo.Key, fmt.Sprintf("%s-%s", sevents.EventNameIssueDocumentPO), nil)
	if err != nil {
		return
	}
	//2. 构造MessageToSign
	msgToSign, err := messagetosign.NewJettradeTxIssuePOData(req, nonce, js.notaryPrivateKeyInfo.Address, &js.ds.spectrumJettradeService.Proxy)
	if err != nil {
		return
	}
	//3. 发起分布式签名
	var signature []byte
	signature, _, err = js.ds.getNotaryService().startDistributedSignAndWait(msgToSign, js.notaryPrivateKeyInfo)
	if err != nil {
		return
	}
	//4. 调用合约
	transactor := &bind.TransactOpts{
		From:  js.notaryPrivateKeyInfo.ToAddress(),
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != js.notaryPrivateKeyInfo.Address {
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
	err = js.ds.spectrumJettradeService.Proxy.IssuePOWait(transactor, req.TokenId, req.DocumentInfo, req.PONUm, req.Buyer, req.Farmer)
	return

}

func (js *JettradeService) callSignDOBuyerOnSpectrum(req *japi.SignDOBuyerOnSpectrumRequest) (ei *models.JettradeEventInfo, err error) {
	//0. 校验相应的事件存在, 发生在eth上
	ei, err = js.ds.db.GetJettradeEventInfo(cfg.ETH.Name, sevents.EventNameSignDocumentDOBuyer, js.ds.ethJettradeService.ContractAddress, req.TokenId)
	if err != nil {
		return
	}
	//1. 获取nonce
	nonce, err := js.ds.applyNonceFromNonceServer(cfg.SMC.Name, js.notaryPrivateKeyInfo.Key, fmt.Sprintf("%s-%s", sevents.EventNameIssueDocumentPO), nil)
	if err != nil {
		return
	}
	//2. 构造MessageToSign
	msgToSign, err := messagetosign.NewJettradeTxSignDOBuyerData(req, nonce, js.notaryPrivateKeyInfo.Address, &js.ds.spectrumJettradeService.Proxy)
	if err != nil {
		return
	}
	//3. 发起分布式签名
	var signature []byte
	signature, _, err = js.ds.getNotaryService().startDistributedSignAndWait(msgToSign, js.notaryPrivateKeyInfo)
	if err != nil {
		return
	}
	//4. 调用合约
	transactor := &bind.TransactOpts{
		From:  js.notaryPrivateKeyInfo.ToAddress(),
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != js.notaryPrivateKeyInfo.Address {
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
	err = js.ds.spectrumJettradeService.Proxy.SignDOBuyerWait(transactor, req.TokenId)
	return

}

func (js *JettradeService) callSignPOOnEth(req *japi.SignPOONEthereumRequest) (ei *models.JettradeEventInfo, err error) {
	//0. 校验相应的事件存在, spectrum
	ei, err = js.ds.db.GetJettradeEventInfo(cfg.SMC.Name, sevents.EventNameSignDocumentPO, js.ds.spectrumJettradeService.ContractAddress, req.TokenId)
	if err != nil {
		return
	}
	//1. 获取nonce
	nonce, err := js.ds.applyNonceFromNonceServer(cfg.ETH.Name, js.notaryPrivateKeyInfo.Key, fmt.Sprintf("%s-%s", sevents.EventNameIssueDocumentPO), nil)
	if err != nil {
		return
	}
	//2. 构造MessageToSign
	msgToSign, err := messagetosign.NewJettradeTxSignPOData(req, nonce, js.notaryPrivateKeyInfo.Address, js.ds.ethJettradeService.Proxy)
	if err != nil {
		return
	}
	//3. 发起分布式签名
	var signature []byte
	signature, _, err = js.ds.getNotaryService().startDistributedSignAndWait(msgToSign, js.notaryPrivateKeyInfo)
	if err != nil {
		return
	}
	//4. 调用合约
	transactor := &bind.TransactOpts{
		From:  js.notaryPrivateKeyInfo.ToAddress(),
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != js.notaryPrivateKeyInfo.Address {
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
	err = js.ds.ethJettradeService.Proxy.SignPOWait(transactor, req.TokenId)
	return

}

func (js *JettradeService) callSignDOFFOnEth(req *japi.SignDOFFOnEthereumRequest) (ei *models.JettradeEventInfo, err error) {
	//0. 校验相应的事件存在, spectrum
	ei, err = js.ds.db.GetJettradeEventInfo(cfg.SMC.Name, sevents.EventNameSignDocumentDOFF, js.ds.spectrumJettradeService.ContractAddress, req.TokenId)
	if err != nil {
		return
	}
	//1. 获取nonce
	nonce, err := js.ds.applyNonceFromNonceServer(cfg.ETH.Name, js.notaryPrivateKeyInfo.Key, fmt.Sprintf("%s-%s", sevents.EventNameIssueDocumentPO), nil)
	if err != nil {
		return
	}
	//2. 构造MessageToSign
	msgToSign, err := messagetosign.NewJettradeTxSignDOFFData(req, nonce, js.notaryPrivateKeyInfo.Address, js.ds.ethJettradeService.Proxy)
	if err != nil {
		return
	}
	//3. 发起分布式签名
	var signature []byte
	signature, _, err = js.ds.getNotaryService().startDistributedSignAndWait(msgToSign, js.notaryPrivateKeyInfo)
	if err != nil {
		return
	}
	//4. 调用合约
	transactor := &bind.TransactOpts{
		From:  js.notaryPrivateKeyInfo.ToAddress(),
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != js.notaryPrivateKeyInfo.Address {
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
	err = js.ds.ethJettradeService.Proxy.SignDOFFWait(transactor, req.TokenId, req.DocumentInfo, req.PONUm, req.DONum, req.Farmer, req.Buyer, req.FreightForward)
	return
}

func (js *JettradeService) callIssueINVOnEth(req *japi.IssueINVOnEthereumRequest) (ei *models.JettradeEventInfo, err error) {
	//0. 校验相应的事件存在, spectrum
	ei, err = js.ds.db.GetJettradeEventInfo(cfg.SMC.Name, sevents.EventNameIssueDocumentINV, js.ds.spectrumJettradeService.ContractAddress, req.TokenId)
	if err != nil {
		return
	}
	//1. 获取nonce
	nonce, err := js.ds.applyNonceFromNonceServer(cfg.ETH.Name, js.notaryPrivateKeyInfo.Key, fmt.Sprintf("%s-%s", sevents.EventNameIssueDocumentPO), nil)
	if err != nil {
		return
	}
	//2. 构造MessageToSign
	msgToSign, err := messagetosign.NewJettradeTxIssueINVData(req, nonce, js.notaryPrivateKeyInfo.Address, js.ds.ethJettradeService.Proxy)
	if err != nil {
		return
	}
	//3. 发起分布式签名
	var signature []byte
	signature, _, err = js.ds.getNotaryService().startDistributedSignAndWait(msgToSign, js.notaryPrivateKeyInfo)
	if err != nil {
		return
	}
	//4. 调用合约
	transactor := &bind.TransactOpts{
		From:  js.notaryPrivateKeyInfo.ToAddress(),
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != js.notaryPrivateKeyInfo.Address {
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
	err = js.ds.ethJettradeService.Proxy.IssueINVWait(transactor, req.TokenId, req.DocumentInfo, req.PONUm, req.DONum, req.INVNUm, req.Farmer, req.Buyer)
	return
}
