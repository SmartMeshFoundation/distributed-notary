package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
)

/*
NotifySCTokenDeployedRequest :
该接口在公证人参与签名部署合约时,由合约部署操作发起人在合约部署成功后,将合约信息广播给所有公证人, 但不仅限于发起人调用
该请求由AdminService处理
*/
type NotifySCTokenDeployedRequest struct {
	api.BaseReq
	api.BaseReqWithSessionID
	api.BaseReqWithSignature
	api.BaseReqWithSCToken
	api.BaseReqWithResponse
	SCTokenMetaInfo *models.SideChainTokenMetaInfo `json:"sc_token_meta_info"`
}

// NewNotifySCTokenDeployedRequest :
func NewNotifySCTokenDeployedRequest(self *models.NotaryInfo, scTokenMetaInfo *models.SideChainTokenMetaInfo) *NotifySCTokenDeployedRequest {
	sessionID := utils.NewRandomHash()
	req := &NotifySCTokenDeployedRequest{
		BaseReq:              api.NewBaseReq(APINameNotifySCTokenDeployed),
		BaseReqWithSessionID: api.NewBaseReqWithSessionID(sessionID, self.ID),
		BaseReqWithSignature: api.NewBaseReqWithSignature(self.GetAddress()),
		BaseReqWithSCToken:   api.NewBaseReqWithSCToken(scTokenMetaInfo.SCToken),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		SCTokenMetaInfo:      scTokenMetaInfo,
	}
	return req
}
