package models

import (
	"testing"

	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/stretchr/testify/assert"
)

func TestLockoutInfo(t *testing.T) {
	db := SetupTestDB()
	secret := utils.NewRandomHash()
	secretHash := utils.ShaSecret(secret[:])
	scToken := utils.NewRandomAddress()
	data := &LockoutInfo{
		Secret:         secret,
		SecretHash:     secretHash,
		SCTokenAddress: scToken,
		SCUserAddress:  utils.NewRandomAddress(),
		Amount:         big.NewInt(5),
		MCLockStatus:   LockStatusLock,
	}
	var err error
	list, err := db.GetAllLockoutInfo()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	err = db.UpdateLockoutInfo(data)
	fmt.Println("1----", err)
	assert.NotNil(t, err)

	err = db.NewLockoutInfo(data)
	assert.Nil(t, err)

	err = db.NewLockoutInfo(data)
	fmt.Println("3----", err)
	assert.NotNil(t, err)

	d1, err := db.GetLockoutInfo(secretHash)
	assert.Nil(t, err)
	assert.EqualValues(t, data, d1)

	d1.Amount = big.NewInt(10)
	err = db.UpdateLockoutInfo(d1)
	assert.Nil(t, err)

	d2, err := db.GetLockoutInfo(secretHash)
	assert.Nil(t, err)
	assert.EqualValues(t, d1, d2)
	assert.EqualValues(t, big.NewInt(10), d2.Amount)

	list, err = db.GetAllLockoutInfo()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))

	list2, err := db.GetAllLockoutInfoBySCToken(scToken)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list2))

	list3, err := db.GetAllLockoutInfoBySCToken(utils.EmptyAddress)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list3))

}
