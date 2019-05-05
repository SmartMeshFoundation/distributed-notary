package service

import (
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
)

// SessionLogMsg :
func SessionLogMsg(sessionID common.Hash, formatter string, a ...interface{}) string {
	formatter = fmt.Sprintf("[SessionID=%s] %s", utils.HPex(sessionID), formatter)
	if len(a) == 0 {
		return formatter
	}
	return fmt.Sprintf(formatter, a...)
}

// SCTokenLogMsg :
func SCTokenLogMsg(scToken *models.SideChainTokenMetaInfo, formatter string, a ...interface{}) string {
	formatter = fmt.Sprintf("[%s-Token] %s", scToken.MCName, formatter)
	if len(a) == 0 {
		return formatter
	}
	return fmt.Sprintf(formatter, a...)
}
