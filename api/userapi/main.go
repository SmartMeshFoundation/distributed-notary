package userapi

import (
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nkbai/log"
)

/*
UserAPI :
提供给用户的API
暂时把SystemRequest和NotaryRequest都放在UserAPI
*/
type UserAPI struct {
	api.BaseAPI
}

// NewUserAPI :
func NewUserAPI(host string) *UserAPI {
	var userAPI UserAPI
	router, err := rest.MakeRouter(
		/*
			api about private key
		*/
		rest.Put("/api/1/private-key", userAPI.CreatePrivateKey),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	userAPI.BaseAPI = api.NewBaseAPI(host, router)
	return &userAPI
}
