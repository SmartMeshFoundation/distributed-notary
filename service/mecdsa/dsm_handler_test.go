package mecdsa

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"errors"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/stretchr/testify/assert"
)

type notaryClientForLocalDSMTest struct {
	lock       sync.Mutex
	handlers   map[int]*DSMHandler
	waitingMap *sync.Map
}

func (c *notaryClientForLocalDSMTest) registerDSMHandlers(dhs ...*DSMHandler) {
	for _, dh := range dhs {
		c.handlers[dh.selfNotaryID] = dh
	}
	fmt.Printf("register %d notaries\n", len(dhs))
}

func (c *notaryClientForLocalDSMTest) WSBroadcast(req api.Req, targetNotaryIDs ...int) {
	for _, notaryID := range targetNotaryIDs {
		c.SendWSReqToNotary(req, notaryID)
	}
}
func (c *notaryClientForLocalDSMTest) SendWSReqToNotary(req api.Req, targetNotaryID int) {
	c.lock.Lock()
	dh := c.handlers[targetNotaryID]
	c.lock.Unlock()
	if reqWithResponse, ok := req.(api.ReqWithResponse); ok {
		c.waitingMap.Store(reqWithResponse.GetRequestID(), reqWithResponse)
	}
	go dh.OnRequest(req)
}

func (c *notaryClientForLocalDSMTest) WaitWSResponse(requestID string, timeout ...time.Duration) (resp *api.BaseResponse, err error) {
	reqInterface, ok := c.waitingMap.Load(requestID)
	if !ok {
		err = errors.New("not found")
		return
	}
	reqWithResponse := reqInterface.(api.ReqWithResponse)
	resp = <-reqWithResponse.GetResponseChan()
	return
}

type testMessage struct {
	data []byte
}

func (tm *testMessage) GetSignBytes() []byte {
	return tm.data
}

func (tm *testMessage) GetTransportBytes() []byte {
	return tm.data
}

func (tm *testMessage) GetName() string {
	return "testMessage"
}

func (tm *testMessage) Parse(buf []byte) error {
	tm.data = buf
	return nil
}

func newFiveDSMHandler(t *testing.T) (d0, d1, d2, d3, d4 *DSMHandler) {
	//生成一个key
	wg := &sync.WaitGroup{}
	wg.Add(5)
	l0, l1, l2, l3, l4 := newFivePKNHandler()
	go startPKN(t, l0, wg)
	go startPKN(t, l1, wg)
	go startPKN(t, l2, wg)
	go startPKN(t, l3, wg)
	go startPKN(t, l4, wg)
	wg.Wait()

	sessionID := utils.NewRandomHash()
	message := &testMessage{
		data: []byte{1, 2, 3},
	}
	c := &notaryClientForLocalDSMTest{handlers: make(map[int]*DSMHandler), waitingMap: new(sync.Map)}
	d0 = NewDSMHandler(nil, &models.NotaryInfo{ID: 0}, []int{1, 2, 3, 4}, message, sessionID, l0.privateKeyInfo, c)
	d1 = NewDSMHandler(nil, &models.NotaryInfo{ID: 1}, []int{0, 2, 3, 4}, message, sessionID, l1.privateKeyInfo, c)
	d2 = NewDSMHandler(nil, &models.NotaryInfo{ID: 2}, []int{1, 0, 3, 4}, message, sessionID, l2.privateKeyInfo, c)
	d3 = NewDSMHandler(nil, &models.NotaryInfo{ID: 3}, []int{1, 2, 0, 4}, message, sessionID, l3.privateKeyInfo, c)
	d4 = NewDSMHandler(nil, &models.NotaryInfo{ID: 4}, []int{1, 2, 3, 0}, message, sessionID, l4.privateKeyInfo, c)
	c.registerDSMHandlers(d0, d1, d2, d3, d4)
	return
}

func TestDSM(t *testing.T) {

	wg := &sync.WaitGroup{}
	wg.Add(5)
	//步骤1 -----------------------------------------------
	d0, d1, d2, d3, d4 := newFiveDSMHandler(t)
	go startDSM(t, d0, wg)
	go startDSM(t, d1, wg)
	go startDSM(t, d2, wg)
	go startDSM(t, d3, wg)
	go startDSM(t, d4, wg)
	wg.Wait()
	//验证
}

func startDSM(t *testing.T, dh *DSMHandler, wg *sync.WaitGroup) {
	_, err := dh.StartDSMAndWaitFinish()
	assert.Empty(t, err)
	wg.Done()
}
