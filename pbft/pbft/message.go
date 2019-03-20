package pbft

import "encoding/gob"

type MessageSender interface {
	SendMessage(req interface{}, target int)
}

type ClientMessager interface {
	IsClientMessage() bool
}
type ServerMessager interface {
	IsServerMessage() bool
}

type ClientMessage struct {
}

func (cm *ClientMessage) IsClientMessage() bool {
	return true
}

type ServerMessage struct {
}

func (sm *ServerMessage) IsServerMessage() bool {
	return true
}

//client message
type StartMessage struct {
	ClientMessage
	Arg string
}

var _ ClientMessager = newStartMessage("")
var _ ServerMessager = newRequestMessage(nil)

func newStartMessage(arg string) *StartMessage {
	return &StartMessage{
		Arg: arg,
	}
}

type ResponseMessage struct {
	ClientMessage
	Arg *ResponseArgs
}

func newResponseMessage(arg *ResponseArgs) *ResponseMessage {
	return &ResponseMessage{
		Arg: arg,
	}
}

type RequestMessage struct {
	ServerMessage
	Arg *RequestArgs
}

func newRequestMessage(arg *RequestArgs) *RequestMessage {
	return &RequestMessage{
		Arg: arg,
	}
}

type PrePrepareMessage struct {
	ServerMessage
	Arg *PrePrepareArgs
}

func newPrePrepareMessage(arg *PrePrepareArgs) *PrePrepareMessage {
	return &PrePrepareMessage{
		Arg: arg,
	}
}

type PrepareMessage struct {
	ServerMessage
	Arg *PrepareArgs
}

func newPrepareMessage(arg *PrepareArgs) *PrepareMessage {
	return &PrepareMessage{
		Arg: arg,
	}
}

type CommitMessage struct {
	ServerMessage
	Arg *CommitArgs
}

func newCommitMessage(arg *CommitArgs) *CommitMessage {
	return &CommitMessage{
		Arg: arg,
	}
}

type CheckPointMessage struct {
	ServerMessage
	Arg *CheckPointArgs
}

func newCheckPointMessage(arg *CheckPointArgs) *CheckPointMessage {
	return &CheckPointMessage{
		Arg: arg,
	}
}

type ViewChangeMessage struct {
	ServerMessage
	Arg *ViewChangeArgs
}

func newViewChangeMessage(arg *ViewChangeArgs) *ViewChangeMessage {
	return &ViewChangeMessage{
		Arg: arg,
	}
}

type NewViewMessage struct {
	ServerMessage
	Arg *NewViewArgs
}

func newNewViewMessage(arg *NewViewArgs) *NewViewMessage {
	return &NewViewMessage{
		Arg: arg,
	}
}

type internalTimerMessage struct {
	ServerMessage
}

func init() {
	gob.Register(&ClientMessage{})
	gob.Register(&StartMessage{})
	gob.Register(&ResponseMessage{})
	gob.Register(&ServerMessage{})
	gob.Register(&RequestMessage{})
	gob.Register(&PrePrepareMessage{})
	gob.Register(&PrepareMessage{})
	gob.Register(&CommitMessage{})
	gob.Register(&CheckPointMessage{})
	gob.Register(&ViewChangeMessage{})
	gob.Register(&NewViewMessage{})
}
