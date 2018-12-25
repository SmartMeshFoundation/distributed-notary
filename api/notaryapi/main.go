package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
)

/*
NotaryAPI :
提供给其他公证人节点的API
*/
type NotaryAPI struct {
	host string // notary api监听的ip及端口
	db   *models.DB
}

// NewNotaryAPI :
func NewNotaryAPI(host string, db *models.DB) *NotaryAPI {
	// TODO
	return &NotaryAPI{
		host: host,
		db:   db,
	}
}

// Start 启动监听
func (na *NotaryAPI) Start() {

}

// GetRequestChan :
func (na *NotaryAPI) GetRequestChan() <-chan api.Request {
	//  TODO
	return nil
}
