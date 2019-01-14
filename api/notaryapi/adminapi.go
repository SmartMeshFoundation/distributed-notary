package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
)

/*
NewSCTokenRequest :
该接口在公证人参与签名部署合约时,由合约部署操作发起人在合约部署成功后,将合约信息广播给所有公证人, 但不仅限于发起人调用
该请求由AdminService处理
*/
type NewSCTokenRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	api.BaseCrossChainRequest
	SCTokenMetaInfo *models.SideChainTokenMetaInfo `json:"sc_token_meta_info"`
}

// NewNewSCTokenRequest :
func NewNewSCTokenRequest(self *models.NotaryInfo, scTokenMetaInfo *models.SideChainTokenMetaInfo) *NewSCTokenRequest {
	sessionID := utils.NewRandomHash()
	return &NewSCTokenRequest{
		BaseRequest:           api.NewBaseRequest(APINameNewSCToken),
		BaseNotaryRequest:     api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		BaseCrossChainRequest: api.NewBaseCrossChainRequest(scTokenMetaInfo.SCToken),
		SCTokenMetaInfo:       scTokenMetaInfo,
	}
}
