package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/cmd/nonce_server/nonceapi"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nkbai/log"
)

type chainName string

const (
	chainNameEthereum = "ethereum"
	chainNameSpectrum = "spectrum"
)

/*
由ethclient.Client实现
*/
type chainClient interface {
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type nonceService struct {
	nsAPI           *nonceapi.NonceServerAPI
	db              *models.DB
	quitChan        chan struct{}
	nonceManagerMap map[string]*nonceManager
	chainMap        map[string]chainClient
}

func newNonceService(db *models.DB, nsAPI *nonceapi.NonceServerAPI, smcRPCEndPoint, ethRPCEndPoint string) *nonceService {
	ns := &nonceService{
		db:              db,
		nsAPI:           nsAPI,
		quitChan:        make(chan struct{}),
		nonceManagerMap: make(map[string]*nonceManager),
		chainMap:        make(map[string]chainClient),
	}
	smcClient, err := ethclient.Dial(smcRPCEndPoint)
	if err != nil {
		panic(err)
	}
	ns.chainMap[chainNameSpectrum] = smcClient
	ethClient, err := ethclient.Dial(ethRPCEndPoint)
	if err != nil {
		panic(err)
	}
	ns.chainMap[chainNameEthereum] = ethClient
	return ns
}

/*
单线程处理
*/
func (ns *nonceService) start() {
	logPrefix := "RestfulRequest : "
	log.Info(fmt.Sprintf("%s dispather start ...", logPrefix))
	requestChan := ns.nsAPI.GetRequestChan()
	for {
		select {
		case req, ok := <-requestChan:
			if !ok {
				err := fmt.Errorf("%s read request chan err ", logPrefix)
				panic(err)
			}
			ns.dispatchRestfulRequest(req)
		case <-ns.quitChan:
			log.Info(fmt.Sprintf("%s nonceService stop success", logPrefix))
			return
		}
	}
}

func (ns *nonceService) dispatchRestfulRequest(req api.Req) {
	switch r := req.(type) {
	case *nonceapi.ApplyNonceReq:
		ns.applyNonce(r)
	default:
		r2, ok := req.(api.ReqWithResponse)
		if ok {
			r2.WriteErrorResponse(api.ErrorCodeParamsWrong)
			return
		}
	}
}

func (ns *nonceService) applyNonce(req *nonceapi.ApplyNonceReq) {
	c, ok := ns.chainMap[req.ChainName]
	if !ok {
		req.WriteErrorResponse(api.ErrorCodeDataNotFound, "unknown chain")
		return
	}
	account := common.HexToAddress(req.Account)
	key := account.String() + req.ChainName
	nm, ok := ns.nonceManagerMap[key]
	if !ok {
		nm = newNonceManager(account, req.ChainName, c)
		ns.nonceManagerMap[key] = nm
	}
	req.WriteSuccessResponse(&nonceapi.ApplyNonceResponse{
		Nonce: nm.applyNonce(req.CancelURL),
	})
}
