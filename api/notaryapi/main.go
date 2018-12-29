package notaryapi

import (
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nkbai/log"
)

// APIName :
type APIName string

/* #no gosec */
const (
	APINamePhase1PubKeyProof      = "Phase1PubKeyProof"
	APINAMEPhase2PaillierKeyProof = "Phase2PaillierKeyProof"
	APINAMEPhase3SecretShare      = "Phase3SecretShare"
	APINamePhase4PubKeyProof      = "Phase4PubKeyProof"
)

// APIName2URLMap :
var APIName2URLMap map[string]string

func init() {
	APIName2URLMap = make(map[string]string)
	APIName2URLMap[APINamePhase1PubKeyProof] = "/api/1/private-key/phase1"
	APIName2URLMap[APINAMEPhase2PaillierKeyProof] = "/api/1/private-key/phase2"
	APIName2URLMap[APINAMEPhase3SecretShare] = "/api/1/private-key/phase3"
	APIName2URLMap[APINamePhase4PubKeyProof] = "/api/1/private-key/phase4"
}

/*
NotaryAPI :
提供给其他公证人节点的API
*/
type NotaryAPI struct {
	api.BaseAPI
}

// NewNotaryAPI :
func NewNotaryAPI(host string) *NotaryAPI {
	var notaryAPI NotaryAPI
	router, err := rest.MakeRouter(
		/*
			api about private key generation
		*/
		rest.Post(APIName2URLMap[APINamePhase1PubKeyProof], notaryAPI.keyGenerationPhase1Message),
		rest.Post(APIName2URLMap[APINAMEPhase2PaillierKeyProof], notaryAPI.keyGenerationPhase2Message),
		rest.Post(APIName2URLMap[APINAMEPhase3SecretShare], notaryAPI.keyGenerationPhase3Message),
		rest.Post(APIName2URLMap[APINamePhase4PubKeyProof], notaryAPI.keyGenerationPhase4Message),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	notaryAPI.BaseAPI = api.NewBaseAPI("NotaryAPI-Server", host, router)
	return &notaryAPI
}
