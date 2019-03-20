package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonce(t *testing.T) {
	ast := assert.New(t)
	db := SetupTestDB()
	n, err := db.GetNonce("aaa")
	ast.EqualValues(0, n)
	ast.EqualValues(err, nil)
	err = db.UpdateNonce("aaa", 3)
	ast.EqualValues(err, nil)
	n, err = db.GetNonce("aaa")
	ast.EqualValues(err, nil)
	ast.EqualValues(n, 3)
}
