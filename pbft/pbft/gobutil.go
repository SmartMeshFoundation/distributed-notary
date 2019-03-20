package pbft

import (
	"bytes"
	"encoding/gob"
)

type hs struct {
	Msg interface{}
}

//EncodeMsg 将pbft消息进行gob编码
func EncodeMsg(msg interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&hs{msg})
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

//DecodeMsg 借助hs进行gob解码
func DecodeMsg(data []byte) interface{} {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var q hs
	err := dec.Decode(&q)
	if err != nil {
		panic(err)
	}
	return q.Msg
}
