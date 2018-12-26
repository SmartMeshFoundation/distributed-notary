package userapi

import (
	"testing"
	"time"
)

func TestUserAPI_Start(t *testing.T) {
	ua := NewUserAPI("127.0.0.1:8888")
	ua.SetTimeout(3 * time.Second)
	ua.Start()
}

func TestUserAPI_CreatePrivateKey(t *testing.T) {
	ua := NewUserAPI("127.0.0.1:8888")
	ua.SetTimeout(3 * time.Second)
	ua.CreatePrivateKey(nil, nil)
}
