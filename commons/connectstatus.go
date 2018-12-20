package commons

import "time"

// ConnectStatus :
type ConnectStatus int

const (
	//Disconnected init status
	Disconnected = ConnectStatus(iota)
	//Connected connection status
	Connected
	//Closed user closed
	Closed
	//Reconnecting connection error
	Reconnecting
)

// ConnectStatusChange :
type ConnectStatusChange struct {
	OldStatus  ConnectStatus
	NewStatus  ConnectStatus
	ChangeTime time.Time
}
