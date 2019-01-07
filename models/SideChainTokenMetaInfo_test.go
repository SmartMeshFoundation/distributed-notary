package models

import (
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/stretchr/testify/assert"
)

func TestDB_NewSCTokenMetaInfo(t *testing.T) {
	db := SetupTestDB()
	defer db.Close()

	list, err := db.GetSCTokenMetaInfoList()
	assert.Empty(t, err)
	assert.EqualValues(t, 0, len(list))

	sc := &SideChainTokenMetaInfo{
		SCToken:                  utils.NewRandomAddress(),
		SCTokenName:              "sc",
		SCTokenOwnerKey:          utils.NewRandomHash(),
		MCLockedContractAddress:  utils.NewRandomAddress(),
		MCName:                   "mc",
		MCLockedContractOwnerKey: utils.NewRandomHash(),
	}
	err = db.NewSCTokenMetaInfo(sc)
	assert.Empty(t, err)

	list, err = db.GetSCTokenMetaInfoList()
	assert.Empty(t, err)
	assert.EqualValues(t, 1, len(list))

	sc2 := list[0]
	assert.EqualValues(t, utils.ToJSONString(sc), utils.ToJSONString(sc2))
}
