package service

import (
	"encoding/hex"
	"runtime"
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/api"

	japi "github.com/SmartMeshFoundation/distributed-notary/api/userapi/jettradeapi"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	jchain "github.com/SmartMeshFoundation/distributed-notary/chainjettrade"
	"github.com/nkbai/log"
)

//为了实现方便,2号公证人,收到相应的事件后,自动发起请求,
type jhelper struct {
	js *JettradeService
}

func NewJHelper(js *JettradeService) *jhelper {
	return &jhelper{js}
}
func (j *jhelper) onEvent(e chain.Event) {
	//不是2号公证人,什么都不做
	if j.js.ds.notaryService.self.ID != 2 {
		return
	}
	switch e2 := e.(type) {
	case jchain.IssueDocumentPOEvent:
		if e2.ChainName == cfg.ETH.Name {
			j.createAndIssuePOOnSpectrum(&e2)
		}
	case jchain.SignDocumentPOEvent:
		if e2.ChainName == cfg.SMC.Name {
			j.signPOOnEth(&e2)
		}
	case jchain.SignDocumentDOFFEvent:
		if e2.ChainName == cfg.SMC.Name {
			j.createAndSignDOFFOnEthereum(&e2)
		}
	case jchain.IssueDocumentINVEvent:
		if e2.ChainName == cfg.SMC.Name {
			j.issueINV(&e2)
		}
	case jchain.SignDocumentDOBuyerEvent:
		if e2.ChainName == cfg.ETH.Name {
			j.signDOBuyer(&e2)
		}
	}
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
func (j *jhelper) doreq(req api.ReqWithResponse) {
	req.NewResponseChan()
	time.Sleep(time.Second * 5) //给其他公证人处理的时间.
	log.Info("send req=%s", log.StringInterface(req, 3))
	j.js.ds.userAPI.SendToService(req)
	resp := j.js.ds.userAPI.WaitServiceResponse(req)
	log.Info("jhepler req=%s got response=%s", req.GetRequestID(), log.StringInterface(resp, 3))
}
func (j *jhelper) createAndIssuePOOnSpectrum(e *jchain.IssueDocumentPOEvent) {
	log.Info("jhelper will %s because of event=%s", runFuncName(), log.StringInterface(e, 3))
	/*
		在以太坊上查询,获取信息,然后在spectrum上重复
	*/
	docInfo, err := j.js.ds.ethJettradeService.Proxy.DocInfo(nil, e.TokenID)
	if err != nil {
		log.Error("get doc info err %s", err)
		return
	}

	ponum, err := j.js.ds.ethJettradeService.Proxy.PONum(nil, e.TokenID)
	if err != nil {
		log.Error("get PONum info err %s", err)
		return
	}
	buyer, err := j.js.ds.ethJettradeService.Proxy.Buyer(nil, e.TokenID)
	if err != nil {
		log.Error("get buyer info err %s", err)
		return
	}
	famer, err := j.js.ds.ethJettradeService.Proxy.Farmer(nil, e.TokenID)
	if err != nil {
		log.Error("get famer info err %s", err)
		return
	}
	req := &japi.IssuePOOnSpectrumRequest{
		TokenId:      e.TokenID,
		DocumentInfo: docInfo,
		PONUm:        hex.EncodeToString(ponum[:]),
		Buyer:        buyer,
		Farmer:       famer,
	}
	j.doreq(req)
}
func (j *jhelper) signPOOnEth(e *jchain.SignDocumentPOEvent) {
	log.Info("jhelper will %s because of event=%s", runFuncName(), log.StringInterface(e, 3))
	req := &japi.SignPOONEthereumRequest{
		TokenId: e.TokenID,
	}
	j.doreq(req)
}
func (j *jhelper) createAndSignDOFFOnEthereum(e *jchain.SignDocumentDOFFEvent) {
	log.Info("jhelper will %s because of event=%s", runFuncName(), log.StringInterface(e, 3))
	/*
		在smt上查询,获取信息,然后在spectrum上重复
	*/
	docInfo, err := j.js.ds.spectrumJettradeService.Proxy.DocInfo(nil, e.TokenID)
	if err != nil {
		log.Error("get doc info err %s", err)
		return
	}

	ponum, err := j.js.ds.spectrumJettradeService.Proxy.PONum(nil, e.TokenID)
	if err != nil {
		log.Error("get PONum info err %s", err)
		return
	}
	donum, err := j.js.ds.spectrumJettradeService.Proxy.DONum(nil, e.TokenID, ponum)
	if err != nil {
		log.Error("get donum info err %s", err)
		return
	}
	buyer, err := j.js.ds.spectrumJettradeService.Proxy.Buyer(nil, e.TokenID)
	if err != nil {
		log.Error("get buyer info err %s", err)
		return
	}
	famer, err := j.js.ds.spectrumJettradeService.Proxy.Farmer(nil, e.TokenID)
	if err != nil {
		log.Error("get famer info err %s", err)
		return
	}
	ff, err := j.js.ds.spectrumJettradeService.Proxy.FreightForward(nil, e.TokenID)
	if err != nil {
		log.Error("get FreightForward info err %s", err)
		return
	}
	req := &japi.SignDOFFOnEthereumRequest{
		TokenId:        e.TokenID,
		DocumentInfo:   docInfo,
		PONUm:          hex.EncodeToString(ponum[:]),
		DONum:          hex.EncodeToString(donum[:]),
		Buyer:          buyer,
		Farmer:         famer,
		FreightForward: ff,
	}
	j.doreq(req)
}
func (j *jhelper) signDOBuyer(e *jchain.SignDocumentDOBuyerEvent) {
	log.Info("jhelper will %s because of event=%s", runFuncName(), log.StringInterface(e, 3))
	req := &japi.SignDOBuyerOnSpectrumRequest{
		TokenId: e.TokenID,
	}
	j.doreq(req)
}
func (j *jhelper) issueINV(e *jchain.IssueDocumentINVEvent) {
	log.Info("jhelper will %s because of event=%s", runFuncName(), log.StringInterface(e, 3))
	docInfo, err := j.js.ds.spectrumJettradeService.Proxy.DocInfo(nil, e.TokenID)
	if err != nil {
		log.Error("get doc info err %s", err)
		return
	}

	ponum, err := j.js.ds.spectrumJettradeService.Proxy.PONum(nil, e.TokenID)
	if err != nil {
		log.Error("get PONum info err %s", err)
		return
	}
	donum, err := j.js.ds.spectrumJettradeService.Proxy.DONum(nil, e.TokenID, ponum)
	if err != nil {
		log.Error("get donum info err %s", err)
		return
	}
	invnum, err := j.js.ds.spectrumJettradeService.Proxy.INVNum(nil, e.TokenID, ponum, donum)
	if err != nil {
		log.Error("get INVNum info err %s", err)
		return
	}
	buyer, err := j.js.ds.spectrumJettradeService.Proxy.Buyer(nil, e.TokenID)
	if err != nil {
		log.Error("get buyer info err %s", err)
		return
	}
	famer, err := j.js.ds.spectrumJettradeService.Proxy.Farmer(nil, e.TokenID)
	if err != nil {
		log.Error("get famer info err %s", err)
		return
	}
	req := &japi.IssueINVOnEthereumRequest{
		TokenId:      e.TokenID,
		DocumentInfo: docInfo,
		PONUm:        hex.EncodeToString(ponum[:]),
		DONum:        hex.EncodeToString(donum[:]),
		INVNUm:       hex.EncodeToString(invnum[:]),
		Buyer:        buyer,
		Farmer:       famer,
	}
	j.doreq(req)
}
