package pbft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGobMessage(t *testing.T) {
	ast := assert.New(t)
	sm := newStartMessage("aa")
	// Encode (send) the value.
	data := EncodeMsg(sm)
	sm2 := DecodeMsg(data)
	ast.EqualValues(sm2, sm)
}
