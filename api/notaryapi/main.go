package notaryapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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

	APINameDSMAsk             = "DSM-Ask"
	APINameDSMNotifySelection = "DSM-NotifySelection"
	APINameDSMPhase1Broadcast = "DSM-Phase1Broadcast"
	APINameDSMPhase2MessageA  = "DSM-Phase2MessageA" // response中带Phase2MessageB
	APINameDSMPhase3DeltaI    = "DSM-Phase3DeltaI"
	APINameDSMPhase5A5BProof  = "DSM-Phase5A5BProof"
	APINameDSMPhase5CProof    = "DSM-Phase5CProof"
	APINameDSMPhase6ReceiveSI = "DSM-Phase6ReceiveSI"

	APIAdminNameNewSCToken = "NotaryNewSCToken" // 该接口在公证人参与签名部署合约时,由合约部署操作发起人在合约部署成功后,将合约信息广播给所有公证人
)

// APIName2URLMap :
var APIName2URLMap map[string]string

func init() {
	APIName2URLMap = make(map[string]string)
	/*
		api about private key generation
	*/
	APIName2URLMap[APINamePKNPhase1PubKeyProof] = "/api/1/private-key/phase1"
	APIName2URLMap[APINamePKNPhase2PaillierKeyProof] = "/api/1/private-key/phase2"
	APIName2URLMap[APINamePKNPhase3SecretShare] = "/api/1/private-key/phase3"
	APIName2URLMap[APINamePKNPhase4PubKeyProof] = "/api/1/private-key/phase4"

	/*
		api about distributed sign message
	*/
	APIName2URLMap[APINameDSMAsk] = "/api/1/sign/ask"
	APIName2URLMap[APINameDSMNotifySelection] = "/api/1/sign/notify-selection"
	APIName2URLMap[APINameDSMPhase1Broadcast] = "/api/1/sign/phase1"
	APIName2URLMap[APINameDSMPhase2MessageA] = "/api/1/sign/phase2"
	APIName2URLMap[APINameDSMPhase3DeltaI] = "/api/1/sign/phase3"
	APIName2URLMap[APINameDSMPhase5A5BProof] = "/api/1/sign/phase5A5B"
	APIName2URLMap[APINameDSMPhase5CProof] = "/api/1/sign/phase5C"
	APIName2URLMap[APINameDSMPhase6ReceiveSI] = "/api/1/sign/phase6"

	/*
		admin api
	*/
	APIName2URLMap[APIAdminNameNewSCToken] = "/api/1/admin/sctoken"
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
		/*
			api about distributed sign message
		*/
		rest.Post(APIName2URLMap[APINameDSMAsk], notaryAPI.postRequestWithBody),
		rest.Post(APIName2URLMap[APINameDSMNotifySelection], notaryAPI.postRequestWithBody),
		rest.Post(APIName2URLMap[APINameDSMPhase1Broadcast], notaryAPI.postRequestWithBody),
		rest.Post(APIName2URLMap[APINameDSMPhase2MessageA], notaryAPI.postRequestWithBody),
		rest.Post(APIName2URLMap[APINameDSMPhase3DeltaI], notaryAPI.postRequestWithBody),
		rest.Post(APIName2URLMap[APINameDSMPhase5A5BProof], notaryAPI.postRequestWithBody),
		rest.Post(APIName2URLMap[APINameDSMPhase5CProof], notaryAPI.postRequestWithBody),
		rest.Post(APIName2URLMap[APINameDSMPhase6ReceiveSI], notaryAPI.postRequestWithBody),
		/*
			others
		*/
		rest.Post(APIName2URLMap[APIAdminNameNewSCToken], notaryAPI.postRequestWithBody),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	notaryAPI.BaseAPI = api.NewBaseAPI("NotaryAPI-Server", host, router)
	return &notaryAPI
}

func (na *NotaryAPI) postRequestWithBody(w rest.ResponseWriter, r *rest.Request) {
	// 1. 读取body content
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.Return(w, api.NewFailResponse("", api.ErrorCodeParamsWrong))
		return
	}
	err = r.Body.Close()
	if err != nil {
		api.Return(w, api.NewFailResponse("", api.ErrorCodeParamsWrong))
		return
	}
	if len(content) == 0 {
		api.Return(w, api.NewFailResponse("", api.ErrorCodeParamsWrong))
		return
	}
	// 2. 基础解析,获取api name
	req := &api.BaseRequest{}
	err = json.Unmarshal(content, &req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	// 3. 根据api name 分别解析
	var req2 api.Request
	switch req.Name {
	/*
		api about private key generation
	*/
	// TODO
	/*
		api about distributed sign message
	*/
	case APINameDSMAsk:
		req2 = &DSMAskRequest{}
		err = json.Unmarshal(content, &req2)
	case APINameDSMNotifySelection:
		req2 = &DSMNotifySelectionRequest{}
		err = json.Unmarshal(content, &req2)
	case APINameDSMPhase1Broadcast:
		req2 = &DSMPhase1BroadcastRequest{}
		err = json.Unmarshal(content, &req2)
	case APINameDSMPhase2MessageA:
		req2 = &DSMPhase2MessageARequest{}
		err = json.Unmarshal(content, &req2)
	case APINameDSMPhase3DeltaI:
		req2 = &DSMPhase3DeltaIRequest{}
		err = json.Unmarshal(content, &req2)
	case APINameDSMPhase5A5BProof:
		req2 = &DSMPhase5A5BProofRequest{}
		err = json.Unmarshal(content, &req2)
	case APINameDSMPhase5CProof:
		req2 = &DSMPhase5CProofRequest{}
		err = json.Unmarshal(content, &req2)
	case APINameDSMPhase6ReceiveSI:
		req2 = &DSMPhase6ReceiveSIRequest{}
		err = json.Unmarshal(content, &req2)
	case APIAdminNameNewSCToken:
		req2 = &NewSCTokenRequest{}
		err = json.Unmarshal(content, &req2)
	}
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong, err.Error()))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req2))
}
