package mecdsa

import "github.com/SmartMeshFoundation/distributed-notary/params"

type NotaryService struct {
	NotaryShareArg *params.NotaryShareArg
	Notaries       map[string]*params.NotatoryInfo
}
