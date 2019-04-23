package pbft

import "encoding/gob"

//MessageSender pbft向外发送消息
type MessageSender interface {
	SendMessage(req interface{}, target int)
}

//ClientMessager 客户端消息
type ClientMessager interface {
	IsClientMessage() bool
	//CheckMessageSender 用于消息的发送方是否和他声称的一致
	CheckMessageSender(id int) bool
}

//ServerMessager pbftServer消息
type ServerMessager interface {
	IsServerMessage() bool
	//CheckMessageSender 用于消息的发送方是否和他声称的一致
	CheckMessageSender(id int) bool
}

//ClientMessage 客户端消息,只是为了符合接口
type ClientMessage struct {
}

func (cm *ClientMessage) IsClientMessage() bool {
	return true
}

//ServerMessage 服务短消息,只是为了符合接口
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

//CheckMessageSender 不用交易交易发起消息,只会在进程内部流通
func (s *StartMessage) CheckMessageSender(_ int) bool {
	return true
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

//CheckMessageSender  校验发送方
func (s *ResponseMessage) CheckMessageSender(id int) bool {
	return s.Arg.Rid == id
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

//CheckMessageSender  校验发送方
func (s *RequestMessage) CheckMessageSender(id int) bool {
	return s.Arg.ID == id
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

//CheckMessageSender  校验发送方
func (s *PrePrepareMessage) CheckMessageSender(id int) bool {
	return s.Arg.View == id
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

//CheckMessageSender  校验发送方
func (s *PrepareMessage) CheckMessageSender(id int) bool {
	return s.Arg.Rid == id
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

//CheckMessageSender  校验发送方
func (s *CommitMessage) CheckMessageSender(id int) bool {
	return s.Arg.Rid == id
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

//CheckMessageSender  校验发送方
func (s *CheckPointMessage) CheckMessageSender(id int) bool {
	return s.Arg.Rid == id
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

//CheckMessageSender  校验发送方
func (s *ViewChangeMessage) CheckMessageSender(id int) bool {
	return s.Arg.Rid == id
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

//CheckMessageSender  校验发送方
func (s *NewViewMessage) CheckMessageSender(id int) bool {
	return s.Arg.View == id
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
