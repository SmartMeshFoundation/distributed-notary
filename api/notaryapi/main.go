package notaryapi

import (
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nkbai/log"
)

// APIName :
type APIName string

/* #nogosec */
const (
	APINamePKNPrefix                 = "PKN-"
	APINamePKNPhase1PubKeyProof      = APINamePKNPrefix + "Phase1PubKeyProof"
	APINamePKNPhase2PaillierKeyProof = APINamePKNPrefix + "Phase2PaillierKeyProof"
	APINamePKNPhase3SecretShare      = APINamePKNPrefix + "Phase3SecretShare"
	APINamePKNPhase4PubKeyProof      = APINamePKNPrefix + "Phase4PubKeyProof"
)

// APIName2URLMap :
var APIName2URLMap map[string]string

func init() {
	APIName2URLMap = make(map[string]string)
	APIName2URLMap[APINamePKNPhase1PubKeyProof] = "/api/1/private-key/phase1"
	APIName2URLMap[APINamePKNPhase2PaillierKeyProof] = "/api/1/private-key/phase2"
	APIName2URLMap[APINamePKNPhase3SecretShare] = "/api/1/private-key/phase3"
	APIName2URLMap[APINamePKNPhase4PubKeyProof] = "/api/1/private-key/phase4"
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
		rest.Post(APIName2URLMap[APINamePKNPhase1PubKeyProof], notaryAPI.keyGenerationPhase1Message),
		rest.Post(APIName2URLMap[APINamePKNPhase2PaillierKeyProof], notaryAPI.keyGenerationPhase2Message),
		rest.Post(APIName2URLMap[APINamePKNPhase3SecretShare], notaryAPI.keyGenerationPhase3Message),
		rest.Post(APIName2URLMap[APINamePKNPhase4PubKeyProof], notaryAPI.keyGenerationPhase4Message),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	notaryAPI.BaseAPI = api.NewBaseAPI("NotaryAPI-Server", host, router)
	return &notaryAPI
}
