package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
)

// KeyGenerationPhase1MessageRequest :
type KeyGenerationPhase1MessageRequest struct {
	api.BaseReq
	api.BaseReqWithSessionID
	api.BaseReqWithSignature
	api.BaseReqWithResponse
	Msg *models.KeyGenBroadcastMessage1 `json:"msg"`
}

// NewKeyGenerationPhase1MessageRequest :
func NewKeyGenerationPhase1MessageRequest(sessionID common.Hash, self *models.NotaryInfo, msg *models.KeyGenBroadcastMessage1) *KeyGenerationPhase1MessageRequest {
	return &KeyGenerationPhase1MessageRequest{
		BaseReq:              api.NewBaseReq(APINamePKNPhase1PubKeyProof),
		BaseReqWithSessionID: api.NewBaseReqWithSessionID(sessionID, self.ID),
		BaseReqWithSignature: api.NewBaseReqWithSignature(self.GetAddress()),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		Msg:                  msg,
	}
}

// KeyGenerationPhase2MessageRequest :
type KeyGenerationPhase2MessageRequest struct {
	api.BaseReq
	api.BaseReqWithSessionID
	api.BaseReqWithSignature
	api.BaseReqWithResponse
	Msg *models.KeyGenBroadcastMessage2 `json:"msg"`
}

// NewKeyGenerationPhase2MessageRequest :
func NewKeyGenerationPhase2MessageRequest(sessionID common.Hash, self *models.NotaryInfo, msg *models.KeyGenBroadcastMessage2) *KeyGenerationPhase2MessageRequest {
	return &KeyGenerationPhase2MessageRequest{
		BaseReq:              api.NewBaseReq(APINamePKNPhase2PaillierKeyProof),
		BaseReqWithSessionID: api.NewBaseReqWithSessionID(sessionID, self.ID),
		BaseReqWithSignature: api.NewBaseReqWithSignature(self.GetAddress()),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		Msg:                  msg,
	}
}

// KeyGenerationPhase3MessageRequest :
type KeyGenerationPhase3MessageRequest struct {
	api.BaseReq
	api.BaseReqWithSessionID
	api.BaseReqWithSignature
	api.BaseReqWithResponse
	Msg *models.KeyGenBroadcastMessage3 `json:"msg"`
}

// NewKeyGenerationPhase3MessageRequest :
func NewKeyGenerationPhase3MessageRequest(sessionID common.Hash, self *models.NotaryInfo, msg *models.KeyGenBroadcastMessage3) *KeyGenerationPhase3MessageRequest {
	return &KeyGenerationPhase3MessageRequest{
		BaseReq:              api.NewBaseReq(APINamePKNPhase3SecretShare),
		BaseReqWithSessionID: api.NewBaseReqWithSessionID(sessionID, self.ID),
		BaseReqWithSignature: api.NewBaseReqWithSignature(self.GetAddress()),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		Msg:                  msg,
	}
}

// KeyGenerationPhase4MessageRequest :
type KeyGenerationPhase4MessageRequest struct {
	api.BaseReq
	api.BaseReqWithSessionID
	api.BaseReqWithSignature
	api.BaseReqWithResponse
	Msg *models.KeyGenBroadcastMessage4 `json:"msg"`
}

// NewKeyGenerationPhase4MessageRequest :
func NewKeyGenerationPhase4MessageRequest(sessionID common.Hash, self *models.NotaryInfo, msg *models.KeyGenBroadcastMessage4) *KeyGenerationPhase4MessageRequest {
	return &KeyGenerationPhase4MessageRequest{
		BaseReq:              api.NewBaseReq(APINamePKNPhase4PubKeyProof),
		BaseReqWithSessionID: api.NewBaseReqWithSessionID(sessionID, self.ID),
		BaseReqWithSignature: api.NewBaseReqWithSignature(self.GetAddress()),
		BaseReqWithResponse:  api.NewBaseReqWithResponse(),
		Msg:                  msg,
	}
}
