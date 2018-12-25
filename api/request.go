package api

import (
	"github.com/ethereum/go-ethereum/common"
)

// RequestName :
type RequestName string

// Request 普通请求,不带唯一key不带SCTokenAddress
type Request interface {
	GetRequestName() RequestName
}

// NotaryRequest 公证人之前的请求,带唯一key
type NotaryRequest interface {
	GetKey() common.Hash
}

// CrossChainRequest 跨链交易相关请求,带SCTokenAddress
type CrossChainRequest interface {
	GetSCTokenAddress() common.Address
}

// BaseRequest :
type BaseRequest struct {
	Name RequestName `json:"name"`
}

// GetRequestName :
func (br *BaseRequest) GetRequestName() RequestName {
	return br.Name
}

// BaseNotaryRequest :
type BaseNotaryRequest struct {
	Key common.Hash `json:"key"`
}

// GetKey :
func (bnr *BaseNotaryRequest) GetKey() common.Hash {
	return bnr.Key
}

// BaseCrossChainRequest :
type BaseCrossChainRequest struct {
	SCTokenAddress common.Address `json:"sc_token_address"`
}

// GetSCTokenAddress :
func (bcr *BaseCrossChainRequest) GetSCTokenAddress() common.Address {
	return bcr.SCTokenAddress
}
