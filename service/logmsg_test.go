package service

import (
	"fmt"
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
)

func TestSessionLogMsg(t *testing.T) {
	fmt.Println(SessionLogMsg(utils.NewRandomHash(), "123... %s %s", "aaa", "bbbb"))
}
