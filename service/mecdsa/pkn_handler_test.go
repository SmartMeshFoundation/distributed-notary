package mecdsa

import (
	"testing"

	"sync"

	"time"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/cfg"
	"github.com/SmartMeshFoundation/distributed-notary/curv/share"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/stretchr/testify/assert"
)

func init() {
	cfg.Notaries = &cfg.NotaryCfg{
		ShareCount:     5,
		ThresholdCount: 2,
	}
}

type notaryClientForLocalPKNTest struct {
	lock     sync.Mutex
	handlers map[int]*PKNHandler
}

func (c *notaryClientForLocalPKNTest) registerPKNHandlers(phs ...*PKNHandler) {
	for _, pk := range phs {
		c.handlers[pk.selfNotaryID] = pk
	}
	fmt.Printf("register %d notaries\n", len(phs))
}

func (c *notaryClientForLocalPKNTest) WSBroadcast(req api.Req, targetNotaryIDs ...int) {
	for _, notaryID := range targetNotaryIDs {
		c.SendWSReqToNotary(req, notaryID)
	}
}
func (c *notaryClientForLocalPKNTest) SendWSReqToNotary(req api.Req, targetNotaryID int) {
	c.lock.Lock()
	ph := c.handlers[targetNotaryID]
	c.lock.Unlock()
	go ph.OnRequest(req)
}

func (c *notaryClientForLocalPKNTest) WaitWSResponse(requestID string, timeout ...time.Duration) (resp *api.BaseResponse, err error) {
	// TODO
	return
}

func newFivePKNHandler() (l0, l1, l2, l3, l4 *PKNHandler) {
	sessionID := utils.NewRandomHash()
	c := &notaryClientForLocalPKNTest{handlers: make(map[int]*PKNHandler)}
	l0 = NewPKNHandler(nil, &models.NotaryInfo{ID: 0}, []int{1, 2, 3, 4}, sessionID, c)
	l1 = NewPKNHandler(nil, &models.NotaryInfo{ID: 1}, []int{0, 2, 3, 4}, sessionID, c)
	l2 = NewPKNHandler(nil, &models.NotaryInfo{ID: 2}, []int{1, 0, 3, 4}, sessionID, c)
	l3 = NewPKNHandler(nil, &models.NotaryInfo{ID: 3}, []int{1, 2, 0, 4}, sessionID, c)
	l4 = NewPKNHandler(nil, &models.NotaryInfo{ID: 4}, []int{1, 2, 3, 0}, sessionID, c)
	c.registerPKNHandlers(l0, l1, l2, l3, l4)
	return
}

func TestPKNBenchmark(t *testing.T) {
	num := 10
	wg := sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			TestPKN(t)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestPKN(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(5)
	//步骤1 -----------------------------------------------
	l0, l1, l2, l3, l4 := newFivePKNHandler()
	go startPKN(t, l0, wg)
	go startPKN(t, l1, wg)
	go startPKN(t, l2, wg)
	go startPKN(t, l3, wg)
	go startPKN(t, l4, wg)

	/*================================= final ======================*/
	wg.Wait()
	p0 := l0.privateKeyInfo
	p1 := l1.privateKeyInfo
	p2 := l2.privateKeyInfo
	p3 := l3.privateKeyInfo
	p4 := l4.privateKeyInfo
	assert.EqualValues(t, models.PrivateKeyNegotiateStatusFinished, p0.Status)
	assert.EqualValues(t, models.PrivateKeyNegotiateStatusFinished, p1.Status)
	assert.EqualValues(t, models.PrivateKeyNegotiateStatusFinished, p2.Status)
	assert.EqualValues(t, models.PrivateKeyNegotiateStatusFinished, p3.Status)
	assert.EqualValues(t, models.PrivateKeyNegotiateStatusFinished, p4.Status)

	//校验私钥是否正确分配
	var xi []share.SPrivKey
	xi = append(xi, p0.XI)
	xi = append(xi, p1.XI)
	xi = append(xi, p2.XI)
	xi = append(xi, p3.XI)
	xi = append(xi, p4.XI)

	//totalPrivKey := msg30[1].vss.Reconstruct([]int{0, 1, 2}, xi[0:3])
	sum := share.PrivKeyZero.Clone()
	share.ModAdd(sum, p0.UI)
	share.ModAdd(sum, p1.UI)
	share.ModAdd(sum, p2.UI)
	share.ModAdd(sum, p3.UI)
	share.ModAdd(sum, p4.UI)

	//if sum.D.Cmp(totalPrivKey.D) != 0 {
	//	t.Error("not equal")
	//}
	pubx, puby := share.S.ScalarBaseMult(sum.Bytes())
	if pubx.Cmp(p0.PublicKeyX) != 0 || puby.Cmp(p0.PublicKeyY) != 0 {
		t.Error("pub key error")
	}
}

func startPKN(t *testing.T, ph *PKNHandler, wg *sync.WaitGroup) {
	_, err := ph.StartPKNAndWaitFinish(nil)
	assert.Empty(t, err)
	wg.Done()
}
