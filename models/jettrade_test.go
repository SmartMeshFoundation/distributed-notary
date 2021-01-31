package models

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestJettradeEventInfo(t *testing.T) {
	db := SetupTestDB()
	defer db.Close()
	ast := assert.New(t)
	j := &JettradeEventInfo{
		ChainName:        "eth",
		FromAddress:      utils.NewRandomAddress(),
		BlockNumber:      20,
		EventName:        "test",
		From:             utils.NewRandomAddress(),
		To:               common.Address{},
		TokenID:          big.NewInt(22),
		NotaryIDInCharge: 0,
		TxHash:           common.Hash{},
	}
	err := db.NewJettradeEventInfo(j)
	ast.Nil(err)
	ls, err := db.GetAllJettradeEventInfo()
	ast.Nil(err)
	ast.Len(ls, 1)
	j2, err := db.GetJettradeEventInfo(j.ChainName, j.EventName, j.FromAddress, j.TokenID)
	ast.Nil(err)
	ast.Equal(j2.TokenID, j.TokenID)
	j2.NotaryIDInCharge = 2
	err = db.UpdateJettradeEventInfo(j2)
	ast.Nil(err)
	ls, err = db.GetAllJettradeEventInfo()
	ast.Nil(err)
	ast.Len(ls, 1)
}
