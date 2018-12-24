package transport

import "github.com/SmartMeshFoundation/distributed-notary/models"

/*
Transport :
负责一组公证人之间的消息通讯
*/
type Transport struct {
	self     *models.NotaryInfo // 自己的信息
	notaries []*models.NotaryInfo
}
