package mecdsa

import (
	"testing"

	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/stretchr/testify/assert"
)

func TestLockout(t *testing.T) {
	var finish bool
	var err error
	li0, _, _, li3, li4 := newTestLockin(t)
	message := []byte{1, 2, 3}
	key := utils.NewRandomHash()
	s := []int{0, 3, 4}
	l0, err := NewLockout(li0.db, li0.srv, message, key, li0.Key, s)
	if err != nil {
		t.Error(err)
		return
	}
	l3, err := NewLockout(li3.db, li3.srv, message, key, li3.Key, s)
	assert.EqualValues(t, err, nil)
	l4, err := NewLockout(li4.db, li4.srv, message, key, li4.Key, s)
	//第一步: 生成--------------------------
	assert.EqualValues(t, err, nil)
	msg01, err := l0.GeneratePhase1Broadcast()
	assert.EqualValues(t, err, nil)
	msg31, err := l3.GeneratePhase1Broadcast()
	assert.EqualValues(t, err, nil)
	msg41, err := l4.GeneratePhase1Broadcast()
	assert.EqualValues(t, err, nil)

	//0 收到来自3,4的证明
	finish, err = l0.ReceivePhase1Broadcast(msg31, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase1Broadcast(msg41, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 收到来自0,4的证明
	finish, err = l3.ReceivePhase1Broadcast(msg01, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase1Broadcast(msg41, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 收到来自3,0的证明
	finish, err = l4.ReceivePhase1Broadcast(msg31, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase1Broadcast(msg01, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//第二步 messageA,B
	msg20, err := l0.GeneratePhase2MessageA()
	assert.EqualValues(t, err, nil)
	msg23, err := l3.GeneratePhase2MessageA()
	assert.EqualValues(t, err, nil)
	msg24, err := l4.GeneratePhase2MessageA()
	assert.EqualValues(t, err, nil)

	var mb *models.MessageBPhase2

	// 0收到3的messageA,立即给3回复,然后3处理相应的消息
	mb, err = l0.ReceivePhase2MessageA(msg23, 3)
	assert.EqualValues(t, err, nil)
	finish, err = l3.ReceivePhase2MessageB(mb, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	assert.EqualValues(t, mb.MessageBWi.BProof.PK, l0.L.SignedKey.Gwi)

	//  0收到4的messageA,立即给4回复,然后4处理相应的消息
	mb, err = l0.ReceivePhase2MessageA(msg24, 4)
	assert.EqualValues(t, err, nil)
	finish, err = l4.ReceivePhase2MessageB(mb, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	assert.EqualValues(t, mb.MessageBWi.BProof.PK, l0.L.SignedKey.Gwi)

	//  3收到0的messageA,立即给0回复,然后0处理相应的消息
	mb, err = l3.ReceivePhase2MessageA(msg20, 0)
	assert.EqualValues(t, err, nil)
	finish, err = l0.ReceivePhase2MessageB(mb, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	assert.EqualValues(t, mb.MessageBWi.BProof.PK, l3.L.SignedKey.Gwi)

	//  3收到4的messageA,立即给4回复,然后4处理相应的消息
	mb, err = l3.ReceivePhase2MessageA(msg24, 4)
	assert.EqualValues(t, err, nil)
	finish, err = l4.ReceivePhase2MessageB(mb, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)
	assert.EqualValues(t, mb.MessageBWi.BProof.PK, l3.L.SignedKey.Gwi)

	//  4收到0的messageA,立即给0回复,然后0处理相应的消息
	mb, err = l4.ReceivePhase2MessageA(msg20, 0)
	assert.EqualValues(t, err, nil)
	finish, err = l0.ReceivePhase2MessageB(mb, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)
	assert.EqualValues(t, mb.MessageBWi.BProof.PK, l4.L.SignedKey.Gwi)

	//  4收到3的messageA,立即给3回复,然后3处理相应的消息
	mb, err = l4.ReceivePhase2MessageA(msg23, 3)
	assert.EqualValues(t, err, nil)
	finish, err = l3.ReceivePhase2MessageB(mb, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)
	assert.EqualValues(t, mb.MessageBWi.BProof.PK, l4.L.SignedKey.Gwi)

	//第三步 生成deltaI
	msg30, err := l0.GeneratePhase3DeltaI()
	assert.EqualValues(t, err, nil)
	msg33, err := l3.GeneratePhase3DeltaI()
	assert.EqualValues(t, err, nil)
	msg34, err := l4.GeneratePhase3DeltaI()
	assert.EqualValues(t, err, nil)

	//0 收到来自3,4的deltaI
	finish, err = l0.ReceivePhase3DeltaI(msg33, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase3DeltaI(msg34, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3 收到来自0,4的deltaI
	finish, err = l3.ReceivePhase3DeltaI(msg30, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase3DeltaI(msg34, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4 收到来自3,0的deltaI
	finish, err = l4.ReceivePhase3DeltaI(msg33, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase3DeltaI(msg30, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	assert.EqualValues(t, l0.L.Delta, l3.L.Delta)
	assert.EqualValues(t, l3.L.Delta, l4.L.Delta)
	assert.EqualValues(t, l0.L.Phase1BroadCast, l3.L.Phase1BroadCast)
	assert.EqualValues(t, l3.L.Phase1BroadCast, l4.L.Phase1BroadCast)
	//return
	//第四步  各自生成R
	r0, err := l0.GeneratePhase4R()
	assert.EqualValues(t, err, nil)
	r3, err := l3.GeneratePhase4R()
	assert.EqualValues(t, err, nil)
	r4, err := l4.GeneratePhase4R()
	assert.EqualValues(t, err, nil)

	assert.EqualValues(t, r0, r3)
	assert.EqualValues(t, r3, r4)

	//第五步 开始各自签名

	msg50, err := l0.GeneratePhase5a5bZkProof()
	assert.EqualValues(t, err, nil)
	msg53, err := l3.GeneratePhase5a5bZkProof()
	assert.EqualValues(t, err, nil)
	msg54, err := l4.GeneratePhase5a5bZkProof()
	assert.EqualValues(t, err, nil)

	//0,接受来自3,4的Proof
	finish, err = l0.ReceivePhase5A5BProof(msg53, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase5A5BProof(msg54, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3,接受来自0,4的Proof
	finish, err = l3.ReceivePhase5A5BProof(msg50, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase5A5BProof(msg54, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4,接受来自3,-的Proof
	finish, err = l4.ReceivePhase5A5BProof(msg53, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase5A5BProof(msg50, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//5.2 ------
	msg60, err := l0.GeneratePhase5CProof()
	assert.EqualValues(t, err, nil)
	msg63, err := l3.GeneratePhase5CProof()
	assert.EqualValues(t, err, nil)
	msg64, err := l4.GeneratePhase5CProof()
	assert.EqualValues(t, err, nil)

	finish, err = l0.ReceivePhase5cProof(msg63, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l0.ReceivePhase5cProof(msg64, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	finish, err = l3.ReceivePhase5cProof(msg60, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l3.ReceivePhase5cProof(msg64, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	finish, err = l4.ReceivePhase5cProof(msg63, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	finish, err = l4.ReceivePhase5cProof(msg60, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//最后一步验证并分发si
	si0, err := l0.Generate5dProof()
	assert.EqualValues(t, err, nil)
	si3, err := l3.Generate5dProof()
	assert.EqualValues(t, err, nil)
	si4, err := l4.Generate5dProof()
	assert.EqualValues(t, err, nil)

	//0接受3,4的签名片
	_, finish, err = l0.RecevieSI(si3, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	_, finish, err = l0.RecevieSI(si4, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//3接受0,4的签名片
	_, finish, err = l3.RecevieSI(si0, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	_, finish, err = l3.RecevieSI(si4, 4)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)

	//4接受3,0的签名片
	_, finish, err = l4.RecevieSI(si3, 3)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, false)
	_, finish, err = l4.RecevieSI(si0, 0)
	assert.EqualValues(t, err, nil)
	assert.EqualValues(t, finish, true)
}
