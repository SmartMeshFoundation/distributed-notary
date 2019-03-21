package notaryapi

import (
	"github.com/SmartMeshFoundation/distributed-notary/api"
)

//PBFTMessage 用于协商nonce,Key用于区分是协商哪条链上的哪个账户
type PBFTMessage struct {
	api.BaseReq
	Key string
	Msg []byte
}
