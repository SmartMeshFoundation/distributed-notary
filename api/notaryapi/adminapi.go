package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
)

/*
NewSCTokenRequest :
该接口在公证人参与签名部署合约时,由合约部署操作发起人在合约部署成功后,将合约信息广播给所有公证人, 但不仅限于发起人调用
该请求由AdminService处理
*/
type NewSCTokenRequest struct {
	api.BaseRequest
	SCTokenMetaInfo *models.SideChainTokenMetaInfo
}

// NewNewSCTokenRequest :
func NewNewSCTokenRequest(scTokenMetaInfo *models.SideChainTokenMetaInfo) *NewSCTokenRequest {
	return &NewSCTokenRequest{
		BaseRequest:     api.NewBaseRequest(APIAdminNameNewSCToken),
		SCTokenMetaInfo: scTokenMetaInfo,
	}
}
