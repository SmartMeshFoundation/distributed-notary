package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

// KeyGenerationPhase1MessageRequest :
type KeyGenerationPhase1MessageRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	Msg *models.KeyGenBroadcastMessage1 `json:"msg"`
}

// NewKeyGenerationPhase1MessageRequest :
func NewKeyGenerationPhase1MessageRequest(sessionID common.Hash, self *models.NotaryInfo, msg *models.KeyGenBroadcastMessage1) *KeyGenerationPhase1MessageRequest {
	return &KeyGenerationPhase1MessageRequest{
		BaseRequest:       api.NewBaseRequest(APINamePKNPhase1PubKeyProof),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		Msg:               msg,
	}
}

func (na *NotaryAPI) keyGenerationPhase1Message(w rest.ResponseWriter, r *rest.Request) {
	req := &KeyGenerationPhase1MessageRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req))
}

// KeyGenerationPhase2MessageRequest :
type KeyGenerationPhase2MessageRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	Msg *models.KeyGenBroadcastMessage2 `json:"msg"`
}

// NewKeyGenerationPhase2MessageRequest :
func NewKeyGenerationPhase2MessageRequest(sessionID common.Hash, self *models.NotaryInfo, msg *models.KeyGenBroadcastMessage2) *KeyGenerationPhase2MessageRequest {
	return &KeyGenerationPhase2MessageRequest{
		BaseRequest:       api.NewBaseRequest(APINamePKNPhase2PaillierKeyProof),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		Msg:               msg,
	}
}

func (na *NotaryAPI) keyGenerationPhase2Message(w rest.ResponseWriter, r *rest.Request) {
	req := &KeyGenerationPhase2MessageRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req))
}

// KeyGenerationPhase3MessageRequest :
type KeyGenerationPhase3MessageRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	Msg *models.KeyGenBroadcastMessage3 `json:"msg"`
}

// NewKeyGenerationPhase3MessageRequest :
func NewKeyGenerationPhase3MessageRequest(sessionID common.Hash, self *models.NotaryInfo, msg *models.KeyGenBroadcastMessage3) *KeyGenerationPhase3MessageRequest {
	return &KeyGenerationPhase3MessageRequest{
		BaseRequest:       api.NewBaseRequest(APINamePKNPhase3SecretShare),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		Msg:               msg,
	}
}

func (na *NotaryAPI) keyGenerationPhase3Message(w rest.ResponseWriter, r *rest.Request) {
	req := &KeyGenerationPhase3MessageRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req))
}

// KeyGenerationPhase4MessageRequest :
type KeyGenerationPhase4MessageRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	Msg *models.KeyGenBroadcastMessage4 `json:"msg"`
}

// NewKeyGenerationPhase4MessageRequest :
func NewKeyGenerationPhase4MessageRequest(sessionID common.Hash, self *models.NotaryInfo, msg *models.KeyGenBroadcastMessage4) *KeyGenerationPhase4MessageRequest {
	return &KeyGenerationPhase4MessageRequest{
		BaseRequest:       api.NewBaseRequest(APINamePKNPhase4PubKeyProof),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		Msg:               msg,
	}
}

func (na *NotaryAPI) keyGenerationPhase4Message(w rest.ResponseWriter, r *rest.Request) {
	req := &KeyGenerationPhase4MessageRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req))
}
