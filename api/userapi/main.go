package userapi

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nkbai/log"
)

/*
UserAPI :
提供给用户的API
*/
type UserAPI struct {
	host string // user api监听的ip及端口
	db   *models.DB
	api  *rest.Api
}

// NewUserAPI :
func NewUserAPI(host string, db *models.DB) *UserAPI {
	// TODO
	return &UserAPI{
		host: host,
		db:   db,
	}
}

// Start 启动监听线程
func (ua *UserAPI) Start() {
	ua.api = rest.NewApi()
	ua.api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(

		/*
			api about private key
		*/
		rest.Put("/api/1/private-key", ua.CreatePrivateKey),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	ua.api.SetApp(router)
	log.Crit(fmt.Sprintf("http listen and serve :%s", http.ListenAndServe(ua.host, ua.api.MakeHandler())))
}

// GetRequestChan :
func (ua *UserAPI) GetRequestChan() <-chan api.Request {
	//  TODO
	return nil
}
