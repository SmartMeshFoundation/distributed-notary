package userapi

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/ant0ine/go-json-rest/rest"
)

// CreatePrivateKey :
func (ua *UserAPI) CreatePrivateKey(w rest.ResponseWriter, r *rest.Request) {
	err := w.WriteJson("ok")
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
