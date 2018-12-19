package models

import (
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/stretchr/testify/assert"
)

func TestDB_NewLockedout(t *testing.T) {
	l := &SignMessage{
		Key:            utils.NewRandomHash(),
		UsedPrivateKey: utils.NewRandomHash(),
		S:              []int{1, 2, 3},
		Sigma:          share.RandomPrivateKey(),
		Phase3Delta: map[int]*DeltaPhase3{
			1: {
				share.RandomPrivateKey(),
			},
		},
	}
	db := SetupTestDB()
	err := db.NewSignMessage(l)
	if err != nil {
		t.Error(err)
		return
	}
	l2, err := db.LoadSignMessage(l.Key)
	if err != nil {
		t.Error(err)
		return
	}
	//t.Logf("l=%s", utils.StringInterface(l, 5))
	//t.Logf("l2=%s", utils.StringInterface(l2, 5))
	assert.EqualValues(t, l, l2)

	l.AlphaGamma = make(map[int]share.SPrivKey)
	l.AlphaGamma[3] = share.RandomPrivateKey()
	err = db.UpdateSignMessage(l)
	if err != nil {
		t.Error(err)
		return
	}
	l2, err = db.LoadSignMessage(l.Key)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, l, l2)

}
