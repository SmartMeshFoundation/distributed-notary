package pbft

import (
	"fmt"

	"github.com/nkbai/log"
)

//for test only
type MessageHub struct {
	ClientAddresses []chan interface{}
	ServerAddresses []chan interface{}
	ni              *Network
}

//不再同步了,完全是异步
func (c *MessageHub) SendMessage(req interface{}, receiver int) {
	// Remote server is offline
	if c.ni != nil && !c.ni.connected[receiver] {
		return
	}
	log.Trace(fmt.Sprintf("clientEnd SendMessage to %v", req))
	switch req.(type) {
	case ClientMessager:
		c.ClientAddresses[receiver] <- req
	case ServerMessager:
		c.ServerAddresses[receiver] <- req
	default:
		panic("unkown req ")

	}
}

func NewMessageHub(serverAddrs []chan interface{}, clientAddres []chan interface{}, ni *Network) *MessageHub {
	return &MessageHub{
		ServerAddresses: serverAddrs,
		ClientAddresses: clientAddres,
		ni:              ni,
	}
}

// Network is a struct use to hold all network information, every rpc client have a pointer
// to this function
type Network struct {
	connected []bool
	drop      int // Drop rate is drop/100
	latency   int // Latency in millsecond
	n         int
}

// Enable is use to set if a server is enabled
func (ni *Network) Enable(index int, enable bool) {
	ni.connected[index] = enable
}

// SetLatency is use to set max latency in this network
func (ni *Network) SetLatency(lat int) {
	ni.latency = lat
}

// SetDrop is use to set drop rate between every pair of cs
func (ni *Network) SetDrop(drop int) {
	ni.drop = drop
}

// NewNetwork is use to create a network struct, which hold all connect,drop,latency information
func NewNetwork(n int) *Network {
	ni := Network{
		connected: make([]bool, n),
		drop:      0,
		latency:   0,
		n:         n,
	}
	for i := 0; i < n; i++ {
		ni.connected[i] = true
	}
	return &ni
}
