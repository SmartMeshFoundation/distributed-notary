package chainjettrade

import (
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/ethereum/go-ethereum/common"
)

type IsJettradeEvent interface {
	IsJettradeEvent() bool
	GetShareData() *ShareData
}
type ShareData struct {
	From    common.Address
	To      common.Address
	TokenID *big.Int
}

func (s ShareData) IsJettradeEvent() bool {
	return true
}
func (s ShareData) GetShareData() *ShareData {
	return &s
}

type IssueDocumentPOEvent struct {
	*chain.BaseEvent
	ShareData
}
type SignDocumentPOEvent struct {
	*chain.BaseEvent
	ShareData
}
type IssueDocumentDOEvent struct {
	*chain.BaseEvent
	ShareData
}

type SignDocumentDOFFEvent struct {
	*chain.BaseEvent
	ShareData
}

type SignDocumentDOBuyerEvent struct {
	*chain.BaseEvent
	ShareData
}
type IssueDocumentINVEvent struct {
	*chain.BaseEvent
	ShareData
}
