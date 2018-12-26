package notaryapi

import (
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nkbai/log"
)

// APIName :
type APIName string

/* #no gosec */
const (
	APINamePhase1PubKeyProof = "Phase1PubKeyProof"
)

// APIName2URLMap :
var APIName2URLMap map[string]string

func init() {
	APIName2URLMap = make(map[string]string)
	APIName2URLMap[APINamePhase1PubKeyProof] = "/api/1/private-key/phase1"
}

/*
NotaryAPI :
提供给其他公证人节点的API
*/
type NotaryAPI struct {
	api.BaseAPI
}

// NewNotaryAPI :
func NewNotaryAPI(host string, db *models.DB) *NotaryAPI {
	var userAPI NotaryAPI
	router, err := rest.MakeRouter(
	// TODO
	/*
		api about private key
	*/
	//rest.Post(APINamePhase1PubKeyProof, userAPI.CreatePrivateKey),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	userAPI.BaseAPI = api.NewBaseAPI(host, router)
	return &userAPI
}
