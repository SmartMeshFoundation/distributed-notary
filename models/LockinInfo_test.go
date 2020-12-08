package models

import (
	"testing"

	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/stretchr/testify/assert"
)

func TestLockinInfo(t *testing.T) {
	db := SetupTestDB()
	secret := utils.NewRandomHash()
	secretHash := utils.ShaSecret(secret[:])
	scToken := utils.NewRandomAddress()
	data := &LockinInfo{
		Secret:         secret,
		SecretHash:     secretHash,
		SCTokenAddress: scToken,
		SCUserAddress:  utils.NewRandomAddress(),
		Amount:         big.NewInt(5),
		MCLockStatus:   LockStatusLock,
		CrossFee:       big.NewInt(1),
		//SCExpiration:   1<<64 - 1, uint64最高位不能为1
	}
	var err error
	list, err := db.GetAllLockinInfo()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	err = db.UpdateLockinInfo(data)
	fmt.Println("1----", err)
	assert.NotNil(t, err)

	err = db.NewLockinInfo(data)
	assert.Nil(t, err)

	err = db.NewLockinInfo(data)
	fmt.Println("3----", err)
	assert.NotNil(t, err)

	d1, err := db.GetLockinInfo(secretHash)
	assert.Nil(t, err)
	assert.EqualValues(t, data, d1)

	d1.Amount = big.NewInt(10)
	err = db.UpdateLockinInfo(d1)
	assert.Nil(t, err)

	d2, err := db.GetLockinInfo(secretHash)
	assert.Nil(t, err)
	assert.EqualValues(t, d1, d2)
	assert.EqualValues(t, big.NewInt(10), d2.Amount)

	list, err = db.GetAllLockinInfo()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))

	list2, err := db.GetAllLockinInfoBySCToken(scToken)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list2))

	list3, err := db.GetAllLockinInfoBySCToken(utils.EmptyAddress)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list3))

}
