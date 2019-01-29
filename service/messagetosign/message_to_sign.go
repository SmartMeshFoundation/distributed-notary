package messagetosign

// MessageToSign 待签名的消息体
type MessageToSign interface {
	GetName() string
	GetTransportBytes() []byte
	GetSignBytes() []byte
	Parse(buf []byte) error
}
