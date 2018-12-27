package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ant0ine/go-json-rest/rest"
)

// KeyGenerationPhase1MessageRequest :
type KeyGenerationPhase1MessageRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	Msg *models.KeyGenBroadcastMessage1
}

func (na *NotaryAPI) keyGenerationPhase1Message(w rest.ResponseWriter, r *rest.Request) {
	req := &KeyGenerationPhase1MessageRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if !req.VerifySignature() {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodePermissionDenied))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req))
}

// KeyGenerationPhase2MessageRequest :
type KeyGenerationPhase2MessageRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	Msg *models.KeyGenBroadcastMessage2
}

func (na *NotaryAPI) keyGenerationPhase2Message(w rest.ResponseWriter, r *rest.Request) {
	req := &KeyGenerationPhase2MessageRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if !req.VerifySignature() {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodePermissionDenied))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req))
}

// KeyGenerationPhase3MessageRequest :
type KeyGenerationPhase3MessageRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	Msg *models.KeyGenBroadcastMessage3
}

func (na *NotaryAPI) keyGenerationPhase3Message(w rest.ResponseWriter, r *rest.Request) {
	req := &KeyGenerationPhase3MessageRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if !req.VerifySignature() {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodePermissionDenied))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req))
}

// KeyGenerationPhase4MessageRequest :
type KeyGenerationPhase4MessageRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	Msg *models.KeyGenBroadcastMessage4
}

func (na *NotaryAPI) keyGenerationPhase4Message(w rest.ResponseWriter, r *rest.Request) {
	req := &KeyGenerationPhase4MessageRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodeParamsWrong))
		return
	}
	if !req.VerifySignature() {
		api.Return(w, api.NewFailResponse(req.RequestID, api.ErrorCodePermissionDenied))
		return
	}
	api.Return(w, na.SendToServiceAndWaitResponse(req))
}
