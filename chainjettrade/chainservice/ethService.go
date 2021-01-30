package chainservice

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/ethereum/contracts"
	"github.com/SmartMeshFoundation/distributed-notary/chainjettrade/ethereum/events"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nkbai/log"
)

type ChainEthService struct {
	*ChainService
	proxy           *contracts.Doc721
	contractAddress common.Address
}

func NewChainEthService(host string, contractAddress common.Address) (ces *ChainEthService, err error) {
	ss, err := NewChainService(host)
	if err != nil {
		return
	}
	proxy, err := contracts.NewDoc721(contractAddress, ss.c)
	if err != nil {
		return
	}
	ces = &ChainEthService{
		ChainService:    ss,
		proxy:           proxy,
		contractAddress: contractAddress,
	}
	ss.co = ces
	return
}
func (ss *ChainEthService) parserLogsToEventsAndSort(logs []types.Log) (es []chainjettrade.Event, err error) {
	if len(logs) == 0 {
		return
	}
	for _, l := range logs {
		eventName := events.TopicToEventName[l.Topics[0]]
		// 根据已处理流水去重
		if doneBlockNumber, ok := ss.eventsDone[l.TxHash]; ok {
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
			log.Warn(fmt.Sprintf("SmcService.EventListener receive unkonwn type event from chain : \n%s\n", utils.ToJSONStringFormat(l)))
		}
		// 记录处理流水
		ss.eventsDone[l.TxHash] = l.BlockNumber
	}
	return
}
func (ss *ChainEthService) CreateNewBlockEvent(blockNumber uint64) chainjettrade.NewBlockEvent {
	return events.CreateNewBlockEvent(blockNumber)
}

// DeployContract : impl chaintmp.Chain 这里暂时只有EthereumToken一个合约,后续优化该接口为支持多主链
func (ss *ChainEthService) DeployContract(opts *bind.TransactOpts, notaryAddress common.Address) (contractAddress common.Address, err error) {
	contractAddress, tx, _, err := contracts.DeployDoc721(opts, ss.c, notaryAddress)
	if err != nil {
		return
	}
	ctx := context.Background()
	return bind.WaitDeployed(ctx, ss.c, tx)
}
