package params

import "github.com/ethereum/go-ethereum/common"

// Config :
type Config struct {
	DataBasePath       string
	Address            common.Address
	KeystorePath       string
	Password           string
	NotaryConfFilePath string
	SmcRPCEndPoint     string
	EthRPCEndPoint     string
	BtcRPCEndPoint     string
	BtcRPCUser         string
	BtcRPCPass         string
	BtcRPCCertFilePath string
	UserAPIListen      string
	NotaryAPIListen    string
	NonceServerHost    string
}
