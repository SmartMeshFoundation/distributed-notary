package userapi

import (
	"testing"
)

func TestUserAPI_Start(t *testing.T) {
	ua := NewUserAPI("127.0.0.1:8888", nil)
	ua.Start()
}
