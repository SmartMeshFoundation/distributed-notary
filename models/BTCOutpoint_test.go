package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBTCOutPoint(t *testing.T) {
	db := SetupTestDB()
	defer db.Close()
	o := &BTCOutpoint{
		TxHashStr: "11111",
		Status:    BTCOutpointStatusUsable,
	}
	err := db.NewBTCOutpoint(o)
	if err != nil {
		t.Error(err)
		return
	}
	l := db.GetBTCOutpointList(-1)
	assert.EqualValues(t, 1, len(l))

	l = db.GetBTCOutpointList(BTCOutpointStatusUsable)
	assert.EqualValues(t, 1, len(l))

	l = db.GetBTCOutpointList(BTCOutpointStatusUsed)
	assert.EqualValues(t, 0, len(l))
}
