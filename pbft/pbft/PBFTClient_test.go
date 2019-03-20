package pbft

import (
	"math/rand"
	"testing"
	"time"
)

type mockServer struct {
	id       int //公证人编号
	replicas []int
	f        int
	m        map[string]int
	msgChan  chan interface{}
	sender   MessageSender
}

func NewMockServer(rid, f, initSeq int, msgChan chan interface{}, sender MessageSender, nodes []int) *mockServer {
	ms := &mockServer{
		id:       rid,
		sender:   sender,
		msgChan:  msgChan,
		replicas: nodes,
		f:        f,
		m:        make(map[string]int),
	}
	go ms.loop()
	return ms
}
func (s *mockServer) loop() {
	seq := 0
	for {
		r2 := <-s.msgChan
		if s.id != 0 {
			continue
		}
		r, ok := r2.(*RequestMessage)
		if !ok {
			continue
		}
		_, ok = s.m[r.Arg.Op]
		if ok {
			continue //稍后自会回复
		}
		seq++
		s.m[r.Arg.Op] = seq

		go func(seq2 int, arg RequestArgs) {
			delay := rand.Int63n(5 * int64(retransmitTimeout))
			time.Sleep(time.Duration(delay))
			for i := 0; i <= s.f; i++ {
				s.sender.SendMessage(newResponseMessage(&ResponseArgs{
					View: 0,
					Seq:  seq2,
					Cid:  arg.ID,
					Rid:  i,
					Res:  arg.Op,
				}), arg.ID)
			}

		}(seq, *r.Arg)

	}
}
func TestRequestOrder(t *testing.T) {
	/*
		保证客户端是在队列中有请求处理完毕才填充进去新的请求
	*/
	testBasic(t, 1, outstanding*10, 0, 2*time.Millisecond, true)
}
