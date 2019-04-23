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

func TestBTCOutPointPBFT(t *testing.T) {
	db := SetupTestDB()
	defer db.Close()
	ast := assert.New(t)
	os := []*BTCOutpoint{
		{
			TxHashStr: "1",
			Amount:    3000,
		},
		{
			TxHashStr: "2",
			Amount:    100,
		},
		{
			TxHashStr: "5",
			Amount:    800,
		},
	}
	for _, o := range os {
		err := db.NewBTCOutpoint(o)
		ast.Nil(err)
	}
	os2, err := db.GetAvailableUTXO(0, 100)
	ast.Nil(err)
	ast.Len(os2, 1)
	//应该选100而不是选3000或者800,多了会造成不必要的浪费
	ast.EqualValues(os2[0], os[1])

	err = db.PreSelectUTXO("2", 0, "test")
	ast.Nil(err)
	err = db.PreSelectUTXO("2", 0, "test")
	ast.EqualValues(err, errUTXODuplciateUase)
	err = db.PreSelectUTXO("2", 2, "test")
	ast.Nil(err)
	err = db.CommitUTXO("2", 2, "test")
	ast.Nil(err)

	os2, err = db.GetAvailableUTXO(0, 100)
	ast.Nil(err)
	ast.Len(os2, 1)
	//100已经被消费了
	ast.EqualValues(os2[0], os[2])
	os2, err = db.GetAvailableUTXO(0, 4000)
	ast.EqualValues(err, errUTXONoAvailabe)
	os2, err = db.GetAvailableUTXO(0, 3800)
	ast.Nil(err)
	ast.Len(os2, 2)
}
