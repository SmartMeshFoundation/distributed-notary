package notaryapi

import (
	"fmt"
	"time"

	"strings"

	"errors"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/nkbai/log"
	"golang.org/x/net/websocket"
)

/*
NotaryClient :
提供给service层的消息发送接口
*/
type NotaryClient interface {
	SendWSReqToNotary(req api.Req, targetNotaryID int)
	WaitWSResponse(requestID string, timeout ...time.Duration) (resp *api.BaseResponse, err error)
}

/*
SendWSReqToNotary :
	发送请求到目标公证人,该方法理论上不阻塞,且不关心返回值
	但需要保证发给相同公证人的多条消息是有序的,不同公证人之间无所谓
	注意: sendingQueueChan缓冲区满了还是会阻塞的,为了保证有序性,缓冲区必须足够大,一旦出现多个线程阻塞在写sendingQueueChan上,就有出现消息乱序的可能
*/
func (na *NotaryAPI) SendWSReqToNotary(req api.Req, targetNotaryID int) {
	// 1. 如果是需要返回的消息,存储到内存
	if reqWithResponse, ok2 := req.(api.ReqWithResponse); ok2 {
		na.dealingWSReqMap.LoadOrStore(reqWithResponse.GetRequestID(), reqWithResponse)
	}
	sendingQueueInterface, ok := na.wsReqSendingQueueMap.Load(targetNotaryID)
	if ok {
		// 存在现成的发送队列,直接投递
		sendingQueue := sendingQueueInterface.(chan api.Req)
		sendingQueue <- req
		return
	}
	// 1. 创建发送队列,缓冲区大小取10倍公证人数量,应该够用了
	sendingQueue := make(chan api.Req, 10*params.ShareCount)
	// 2. 保存队列
	na.wsReqSendingQueueMap.Store(targetNotaryID, sendingQueue)
	// 3. 启动该队列对应的常驻发送线程
	go na.notaryMsgSendLoop(targetNotaryID, sendingQueue)
	// 4. 投递消息
	sendingQueue <- req
}

/*
常驻发送线程,对应一个公证人
*/
func (na *NotaryAPI) notaryMsgSendLoop(targetNotaryID int, sendingQueueChan chan api.Req) {
	log.Trace("notaryMsgSendLoop with notary[ID=%d] start...", targetNotaryID)
	for {
		select {
		case req, ok := <-sendingQueueChan:
			if !ok {
				return
			}
			// 2. 查询ws连接,如果有直接使用
			na.notaryWSConnMapLock.Lock()
			wsInterface, ok := na.notaryWSConnMap.Load(targetNotaryID)
			if ok {
				na.notaryWSConnMapLock.Unlock()
				ws := wsInterface.(*websocket.Conn)
				//fmt.Printf("=========> send %s to %d\n", req.GetRequestName(), targetNotaryID)
				api.WSWriteJSON(ws, req)
				continue
			}
			// 2. 如果没有连接,创建并保存
			wsURL, origin, err := na.getNotaryWSURLAndOrigin(targetNotaryID)
			if err != nil {
				log.Error("notaryMsgSendLoop with notary[ID=%d] end because getNotaryHostAndOrigin err : %s", targetNotaryID, err.Error())
				na.notaryWSConnMapLock.Unlock()
				return
			}
			ws, err := websocket.Dial(wsURL, "", origin)
			if err != nil {
				log.Error("notaryMsgSendLoop with notary[ID=%d] end because websocket connect to %s with origin %s err = %s", targetNotaryID, wsURL, origin, err.Error())
				na.notaryWSConnMapLock.Unlock()
				return
			}
			na.notaryWSConnMap.Store(targetNotaryID, ws)
			na.notaryWSConnMapLock.Unlock()
			// 3. 为这个连接启动消息接收线程,需要在发送之前启动,否则第一次消息的回执可能收不到
			go func(ws *websocket.Conn, targetNotaryID int) {
				na.notaryMsgReceiveLoop(ws, targetNotaryID)
				// 如果这里返回了,说明连接被对方关闭了,线程也可以直接退出,但是需要close连接
				err := ws.Close()
				if err != nil {
					log.Error("websocket connection with notary[ID=%d] close err : %s", targetNotaryID, err.Error())
				}
			}(ws, targetNotaryID)
			//time.Sleep(500 * time.Millisecond)
			// 4. 发送
			//fmt.Printf("=========> first send %s to %d wsURL=%s origin=%s\n", req.GetRequestName(), targetNotaryID, wsURL, origin)
			api.WSWriteJSON(ws, req)
		}
	}
}

/*
WaitWSResponse :
	阻塞等待某个请求的返回
*/
func (na *NotaryAPI) WaitWSResponse(requestID string, timeout ...time.Duration) (resp *api.BaseResponse, err error) {
	reqInterface, ok := na.dealingWSReqMap.Load(requestID)
	if !ok {
		err = fmt.Errorf("can not find req[requestID=%s] in dealingWSReqMap", requestID)
		return
	}
	req := reqInterface.(api.ReqWithResponse)

	requestTimeout := na.Timeout
	if len(timeout) > 0 && timeout[0] > 0 {
		requestTimeout = timeout[0]
	}
	if requestTimeout > 0 {
		select {
		case resp = <-req.GetResponseChan():
		case <-time.After(requestTimeout):
			resp = api.NewFailResponse(requestID, api.ErrorCodeTimeout)
		}
	} else {
		resp = <-req.GetResponseChan()
	}
	if resp.ErrorCode != api.ErrorCodeSuccess {
		err = errors.New(resp.ErrorMsg)
	}
	return
}

func (na *NotaryAPI) getNotaryWSURLAndOrigin(notaryID int) (wsURL string, origin string, err error) {
	var notary *models.NotaryInfo
	for _, n := range na.notaries {
		if n.ID == notaryID {
			notary = &n
			break
		}
	}
	if notary == nil {
		err = fmt.Errorf("can not find notary info with id : %d", notaryID)
		return
	}
	host := notary.Host
	if strings.HasPrefix(host, "ws") {
		wsURL = host
		origin = strings.Replace(host, "ws", "http", 1)
	} else if strings.HasPrefix(host, "http") {
		origin = host
		wsURL = strings.Replace(host, "http", "ws", 1)
	} else {
		wsURL = "ws://" + host
		origin = "http://" + host
	}
	return
}
