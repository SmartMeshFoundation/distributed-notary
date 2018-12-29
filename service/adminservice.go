package service

import (
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/nkbai/log"
)

// AdminService :
type AdminService struct {
	db            *models.DB
	notaryService *NotaryService
}

// NewAdminService :
func NewAdminService(db *models.DB, notaryService *NotaryService) (ns *AdminService, err error) {
	ns = &AdminService{
		db:            db,
		notaryService: notaryService,
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
	case *userapi.CreatePrivateKeyRequest:
		as.onCreatePrivateKeyRequest(r)
	case *userapi.GetNotaryListRequest:
		as.onGetNotaryListRequest(r)
	case *userapi.GetPrivateKeyListRequest:
		as.onGetPrivateKeyListRequest(r)
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

// 私钥列表查询
func (as *AdminService) onGetPrivateKeyListRequest(req *userapi.GetPrivateKeyListRequest) {
	privateKeyInfoList, err := as.db.GetPrivateKeyList()
	if err != nil {
		req.WriteErrorResponse(api.ErrorCodeException, err.Error())
	}
	req.WriteSuccessResponse(privateKeyInfoList)
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
		req.WriteSuccessResponse(privateKey)
		return
	}
}
