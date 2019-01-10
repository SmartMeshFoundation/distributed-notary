package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/mecdsa"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
)

// DSMAskRequest :
type DSMAskRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
}

// NewDSMAskRequest :
func NewDSMAskRequest(sessionID common.Hash, self *models.NotaryInfo) *DSMAskRequest {
	return &DSMAskRequest{
		BaseRequest:       api.NewBaseRequest(APINameDSMAsk),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
	}
}

// DSMNotifySelectionRequest :
type DSMNotifySelectionRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	NotaryIDs               []int       `json:"notary_ids"`
	PrivateKeyID            common.Hash `json:"private_key_id"`
	MsgName                 string      `json:"msg_name"`
	MsgToSignTransportBytes []byte      `json:"msg_to_sign_transport_bytes"`
}

// NewDSMNotifySelectionRequest :
func NewDSMNotifySelectionRequest(sessionID common.Hash, self *models.NotaryInfo, notaryIDs []int, privateKeyID common.Hash, msgToSign mecdsa.MessageToSign) *DSMNotifySelectionRequest {
	return &DSMNotifySelectionRequest{
		BaseRequest:             api.NewBaseRequest(APINameDSMNotifySelection),
		BaseNotaryRequest:       api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		NotaryIDs:               notaryIDs,
		PrivateKeyID:            privateKeyID,
		MsgToSignTransportBytes: msgToSign.GetTransportBytes(),
		MsgName:                 msgToSign.GetName(),
	}
}

// DSMPhase1BroadcastRequest :
type DSMPhase1BroadcastRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	PrivateKeyID common.Hash                 `json:"private_key_id"`
	Msg          *models.SignBroadcastPhase1 `json:"msg"`
}

// NewDSMPhase1BroadcastRequest :
func NewDSMPhase1BroadcastRequest(sessionID common.Hash, self *models.NotaryInfo, privateKeyID common.Hash, msg *models.SignBroadcastPhase1) *DSMPhase1BroadcastRequest {
	return &DSMPhase1BroadcastRequest{
		BaseRequest:       api.NewBaseRequest(APINameDSMPhase1Broadcast),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		PrivateKeyID:      privateKeyID,
		Msg:               msg,
	}
}

// DSMPhase2MessageARequest :
type DSMPhase2MessageARequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	PrivateKeyID common.Hash      `json:"private_key_id"`
	Msg          *models.MessageA `json:"msg"`
}

// NewDSMPhase2MessageARequest :
func NewDSMPhase2MessageARequest(sessionID common.Hash, self *models.NotaryInfo, privateKeyID common.Hash, msg *models.MessageA) *DSMPhase2MessageARequest {
	return &DSMPhase2MessageARequest{
		BaseRequest:       api.NewBaseRequest(APINameDSMPhase2MessageA),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		PrivateKeyID:      privateKeyID,
		Msg:               msg,
	}
}

// DSMPhase3DeltaIRequest :
type DSMPhase3DeltaIRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	PrivateKeyID common.Hash         `json:"private_key_id"`
	Msg          *models.DeltaPhase3 `json:"msg"`
}

// NewDSMPhase3DeltaIRequest :
func NewDSMPhase3DeltaIRequest(sessionID common.Hash, self *models.NotaryInfo, privateKeyID common.Hash, msg *models.DeltaPhase3) *DSMPhase3DeltaIRequest {
	return &DSMPhase3DeltaIRequest{
		BaseRequest:       api.NewBaseRequest(APINameDSMPhase3DeltaI),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		PrivateKeyID:      privateKeyID,
		Msg:               msg,
	}
}

// DSMPhase5A5BProofRequest :
type DSMPhase5A5BProofRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	PrivateKeyID common.Hash     `json:"private_key_id"`
	Msg          *models.Phase5A `json:"msg"`
}

// NewDSMPhase5A5BProofRequest :
func NewDSMPhase5A5BProofRequest(sessionID common.Hash, self *models.NotaryInfo, privateKeyID common.Hash, msg *models.Phase5A) *DSMPhase5A5BProofRequest {
	return &DSMPhase5A5BProofRequest{
		BaseRequest:       api.NewBaseRequest(APINameDSMPhase5A5BProof),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		PrivateKeyID:      privateKeyID,
		Msg:               msg,
	}
}

// DSMPhase5CProofRequest :
type DSMPhase5CProofRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	PrivateKeyID common.Hash     `json:"private_key_id"`
	Msg          *models.Phase5C `json:"msg"`
}

// NewDSMPhase5CProofRequest :
func NewDSMPhase5CProofRequest(sessionID common.Hash, self *models.NotaryInfo, privateKeyID common.Hash, msg *models.Phase5C) *DSMPhase5CProofRequest {
	return &DSMPhase5CProofRequest{
		BaseRequest:       api.NewBaseRequest(APINameDSMPhase5CProof),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		PrivateKeyID:      privateKeyID,
		Msg:               msg,
	}
}

// DSMPhase6ReceiveSIRequest :
type DSMPhase6ReceiveSIRequest struct {
	api.BaseRequest
	api.BaseNotaryRequest
	PrivateKeyID common.Hash    `json:"private_key_id"`
	Msg          share.SPrivKey `json:"msg"`
}

// NewDSMPhase6ReceiveSIRequest :
func NewDSMPhase6ReceiveSIRequest(sessionID common.Hash, self *models.NotaryInfo, privateKeyID common.Hash, msg share.SPrivKey) *DSMPhase6ReceiveSIRequest {
	return &DSMPhase6ReceiveSIRequest{
		BaseRequest:       api.NewBaseRequest(APINameDSMPhase6ReceiveSI),
		BaseNotaryRequest: api.NewBaseNotaryRequest(sessionID, self.GetAddress(), self.ID),
		PrivateKeyID:      privateKeyID,
		Msg:               msg,
	}
}
