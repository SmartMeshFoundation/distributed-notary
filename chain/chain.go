package chain

/*
Chain :

*/
type Chain interface {
	GetChainName() string
	GetEventChan() <-chan Event
	StartEventListener() error
	StopEventListener()
}
