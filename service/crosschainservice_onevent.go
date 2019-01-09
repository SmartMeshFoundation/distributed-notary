package service

import (
	"errors"
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	ethevents "github.com/SmartMeshFoundation/distributed-notary/chain/ethereum/events"
	smcevents "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/events"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/nkbai/log"
)

// OnEvent 链上事件逻辑处理 TODO
func (ns *CrossChainService) OnEvent(e chain.Event) {
	logPrefix := fmt.Sprintf("CrossChainService[SCToken=%s] ", utils.APex(ns.meta.SCToken))
	var err error
	switch event := e.(type) {
	/*
		events about block number
	*/
	case ethevents.NewBlockEvent:
		err = onEthereumNewBlockEvent(ns, event)
	case smcevents.NewBlockEvent:
	/*
		events about lockin
	*/
	case ethevents.PrepareLockinEvent: // MCPLI
	case smcevents.PrepareLockinEvent: // SCPLI
	case smcevents.LockinSecretEvent: //  SCLIS
	case ethevents.LockinEvent: // MCLI
	case ethevents.CancelLockinEvent: // MCCancelLI
	case smcevents.CancelLockinEvent: // SCCancelLI
	/*
		events about lockout
	*/
	case smcevents.PrepareLockoutEvent: // SCPLO
	case ethevents.PrepareLockoutEvent: // MCPLO
	case ethevents.LockoutSecretEvent: // MCLOS
	case smcevents.LockoutEvent: // SCLO
	case ethevents.CancelLockoutEvent: // MCCancelLO
	case smcevents.CancelLockoutEvent: // SCCancelLO
	default:
		err = errors.New("unknow event")
	}
	if err != nil {
		log.Error(fmt.Sprintf("%s deal event err=%s,event:\n%s\n", logPrefix, err.Error(), utils.ToJSONStringFormat(e)))
	}
	return
}

func onEthereumNewBlockEvent(ns *CrossChainService, event ethevents.NewBlockEvent) (err error) {
	return
}
