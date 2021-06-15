package chainservice

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"

	"github.com/SmartMeshFoundation/distributed-notary/chain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/spectrum/contracts"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nkbai/log"
)

type SMTProxy struct {
	*contracts.Doc721
	c *chainjettrade.SafeEthClient
}

func (p *SMTProxy) CreateAndIssuePO2(opts *bind.TransactOpts, tokenId *big.Int, documentInfo string, PONum string, buyer common.Address, farmer common.Address) (tx *types.Transaction, err error) {
	ponum, err := hex.DecodeString(PONum)
	if err != nil {
		return
	}
	var _ponum [32]byte
	copy(_ponum[:], ponum)
	tx, err = p.CreateAndIssuePO(opts, tokenId, documentInfo, _ponum, buyer, farmer)
	if err != nil {
		log.Error("IssuePOWait %s ", err)
		return
	}
	return
}
func (p *SMTProxy) IssuePOWait(opts *bind.TransactOpts, tokenId *big.Int, documentInfo string, PONum string, buyer common.Address, farmer common.Address) (err error) {
	tx, err := p.CreateAndIssuePO2(opts, tokenId, documentInfo, PONum, buyer, farmer)
	log.Info("spectrum IssuePO tx=%s", tx.Hash().String())
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.c, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract IssuePO success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}
func (p *SMTProxy) SignDOBuyerWait(opts *bind.TransactOpts, tokenID *big.Int) (err error) {
	tx, err := p.SignDOBuyer(opts, tokenID)
	if err != nil {
		log.Error("SignDOBuyerWait %s", err)
		return
	}
	log.Info("spectrum SignDOBuyerWait tx=%s", tx.Hash().String())
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.c, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract SignDOBuyerWait success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}

type ChainSmtService struct {
	*ChainService
	Proxy SMTProxy
}

func NewChainSmtService(host string, contractAddress common.Address) (ces *ChainSmtService, err error) {
	ss, err := NewChainService(host, "jettrade-smtservice", contractAddress)
	if err != nil {
		return
	}
	proxy, err := contracts.NewDoc721(contractAddress, ss.c)
	if err != nil {
		return
	}
	ces = &ChainSmtService{
		ChainService: ss,
		Proxy:        SMTProxy{proxy, ss.c},
	}
	ss.co = ces
	return
}
func (ss *ChainSmtService) GetChainName() string {
	return cfg.SMC.Name
}
func (ss *ChainSmtService) parserLogsToEventsAndSort(logs []types.Log) (es []chain.Event, err error) {
	if len(logs) == 0 {
		return
	}
	for _, l := range logs {
		eventName := events.TopicToEventName[l.Topics[0]]
		// 根据已处理流水去重
		if doneBlockNumber, ok := ss.eventsDone[l.TxHash.String()+eventName]; ok {
			if doneBlockNumber == l.BlockNumber {
				//log.Trace(fmt.Sprintf("get event txhash=%s repeated,ignore...", l.TxHash.String()))
				continue
			}
			log.Warn(fmt.Sprintf("SmcService.EventListener event tx=%s happened at %d, but now happend at %d ", l.TxHash.String(), doneBlockNumber, l.BlockNumber))
		}
		switch eventName {
		case events.EventNameIssueDocumentPO:
			{
				var e chainjettrade.IssueDocumentPOEvent
				e, err = events.CreateIssueDocumentPOEvent(l)
				if err != nil {
					return
				}
				es = append(es, e)
			}
		case events.EventNameSignDocumentPO:
			{
				var e chainjettrade.SignDocumentPOEvent
				e, err = events.CreateSignDocumentPOEvent(l)
				if err != nil {
					return
				}
				es = append(es, e)
			}
		case events.EventNameIssueDocumentDO:
			{
				var e chainjettrade.IssueDocumentDOEvent
				e, err = events.CreateIssueDocumentDOEvent(l)
				if err != nil {
					return
				}
				es = append(es, e)
			}
		case events.EventNameSignDocumentDOFF:
			{
				var e chainjettrade.SignDocumentDOFFEvent
				e, err = events.CreateSignDocumentDOFFEvent(l)
				if err != nil {
					return
				}
				es = append(es, e)
			}
		case events.EventNameSignDocumentDOBuyer:
			{
				var e chainjettrade.SignDocumentDOBuyerEvent
				e, err = events.CreateSignDocumentDOBuyerEvent(l)
				if err != nil {
					return
				}
				es = append(es, e)
			}
		case events.EventNameIssueDocumentINV:
			{
				var e chainjettrade.IssueDocumentINVEvent
				e, err = events.CreateIssueDocumentINVEvent(l)
				if err != nil {
					return
				}
				es = append(es, e)
			}
		default:
			//erc 721有很多event是我们不需要关心的
			//log.Trace(fmt.Sprintf("SmcService.EventListener receive unkonwn type event from chain : \n%s\n", utils.ToJSONStringFormat(l)))
		}
		ss.eventsDone[l.TxHash.String()+eventName] = l.BlockNumber
	}
	return
}
