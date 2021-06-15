package chainservice

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/cfg"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/ethereum/contracts"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nkbai/log"
)

type EthProxy struct {
	*contracts.Doc721
	c *chainjettrade.SafeEthClient
}

func (p *EthProxy) SignPOWait(opts *bind.TransactOpts, tokenId *big.Int) (err error) {

	tx, err := p.SignPO(opts, tokenId)
	if err != nil {
		log.Error("SignPOWait %s ", err)
		return
	}
	log.Info("spectrum SignPOWait tx=%s", tx.Hash().String())
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.c, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract SignPOWait success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}
func (p *EthProxy) CreateDoAndSignFf2(opts *bind.TransactOpts, tokenId *big.Int, documentInfo, PoNum, DoNum string, farmer, buyer, freightForward common.Address) (tx *types.Transaction, err error) {
	poNum2, err := hex.DecodeString(PoNum)
	if err != nil {
		return
	}
	doNum2, err := hex.DecodeString(DoNum)
	if err != nil {
		return
	}
	var ponum3 [32]byte
	var donum3 [32]byte
	copy(ponum3[:], poNum2)
	copy(donum3[:], doNum2)
	tx, err = p.CreateDoAndSignFf(opts, tokenId, documentInfo, ponum3, donum3, farmer, buyer, freightForward)
	return
}
func (p *EthProxy) SignDOFFWait(opts *bind.TransactOpts, tokenId *big.Int, documentInfo, PoNum, DoNum string, farmer, buyer, freightForward common.Address) (err error) {
	tx, err := p.CreateDoAndSignFf2(opts, tokenId, documentInfo, PoNum, DoNum, farmer, buyer, freightForward)
	if err != nil {
		log.Error("SignDOFFWait %s", err)
		return
	}
	log.Info("eth SignDOFFWait tx=%s", tx.Hash().String())
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.c, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract SignDOFFWait success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}
func (p *EthProxy) CreateAndIssueINV2(opts *bind.TransactOpts, tokenId *big.Int, documentInfo, PoNum, DoNum, invNum string, farmer, buyer common.Address) (tx *types.Transaction, err error) {
	poNum2, err := hex.DecodeString(PoNum)
	if err != nil {
		return
	}
	doNum2, err := hex.DecodeString(DoNum)
	if err != nil {
		return
	}
	invNum2, err := hex.DecodeString(invNum)
	if err != nil {
		return
	}
	var ponum3 [32]byte
	var donum3 [32]byte
	var invnum3 [32]byte
	copy(ponum3[:], poNum2)
	copy(donum3[:], doNum2)
	copy(invnum3[:], invNum2)
	tx, err = p.CreateAndIssueINV(opts, tokenId, documentInfo, ponum3, donum3, invnum3, buyer, farmer)
	return
}
func (p *EthProxy) IssueINVWait(opts *bind.TransactOpts, tokenId *big.Int, documentInfo, PoNum, DoNum, invNum string, farmer, buyer common.Address) (err error) {
	tx, err := p.CreateAndIssueINV2(opts, tokenId, documentInfo, PoNum, DoNum, invNum, farmer, buyer)
	if err != nil {
		log.Error("CreateAndIssueINV %s", err)
		return
	}
	log.Info("spectrum IssueINVWait tx=%s", tx.Hash().String())
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, p.c, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("call contract IssueINVWait success but tx %s failed", r.TxHash.String())
		log.Error("failed tx :\n%s", utils.ToJSONStringFormat(tx))
		log.Error("failed receipt :\n%s", utils.ToJSONStringFormat(r))
	}
	return
}

type ChainEthService struct {
	*ChainService
	Proxy *EthProxy
}

func NewChainEthService(host string, contractAddress common.Address) (ces *ChainEthService, err error) {
	ss, err := NewChainService(host, "jettrade-ethservice", contractAddress)
	if err != nil {
		return
	}
	proxy, err := contracts.NewDoc721(contractAddress, ss.c)
	if err != nil {
		return
	}
	ces = &ChainEthService{
		ChainService: ss,
		Proxy: &EthProxy{
			Doc721: proxy,
			c:      ss.c,
		},
	}
	ss.co = ces
	return
}
func (ss *ChainEthService) GetChainName() string {
	return cfg.ETH.Name
}
func (ss *ChainEthService) parserLogsToEventsAndSort(logs []types.Log) (es []chain.Event, err error) {
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
				log.Info("jettrade eth servcie new event txhash=%s,blocknumber=%d,e=%s,log=%s", l.TxHash.String(), l.BlockNumber, log.StringInterface(e, 3), log.StringInterface(l, 4))
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
			//log.Trace(fmt.Sprintf("SmcService.EventListener receive unkonwn type event from chain : \n%s\n", utils.ToJSONStringFormat(l)))
		}
		ss.eventsDone[l.TxHash.String()+eventName] = l.BlockNumber
	}
	return
}
