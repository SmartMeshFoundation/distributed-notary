package notaryapi

import (
	"fmt"
	"time"

	"strings"

	"errors"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/nkbai/log"
	"golang.org/x/net/websocket"
)

/*
NotaryClient :
提供给service层的消息发送接口
*/
type NotaryClient interface {
	WSBroadcast(req api.Req, targetNotaryIDs ...int)
	SendWSReqToNotary(req api.Req, targetNotaryID int)
	WaitWSResponse(requestID string, timeout ...time.Duration) (resp *api.BaseResponse, err error)
}

/*
WSBroadcast :
广播消息
*/
func (na *NotaryAPI) WSBroadcast(req api.Req, targetNotaryIDs ...int) {
	if targetNotaryIDs == nil || len(targetNotaryIDs) == 0 {
		return
	}
	// 1. 如果是需要签名的请求,签名
	if reqWithSignature, ok := req.(api.ReqWithSignature); ok && reqWithSignature.GetSignature() == nil {
		reqWithSignature.Sign(reqWithSignature, na.selfPrivateKey)
	}
	for _, notaryID := range targetNotaryIDs {
		na.SendWSReqToNotary(req, notaryID)
	}

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
		na.waitingResponseMap.LoadOrStore(reqWithResponse.GetRequestID(), reqWithResponse)
	}
	sendingQueueInterface, ok := na.sendingChanMap.Load(targetNotaryID)
	if ok {
		// 存在现成的发送队列,直接投递
		sendingQueue := sendingQueueInterface.(chan api.Req)
		sendingQueue <- req
		return
	}
	panic("never happen")
}

/*
WaitWSResponse :
	阻塞等待某个请求的返回
*/
func (na *NotaryAPI) WaitWSResponse(requestID string, timeout ...time.Duration) (resp *api.BaseResponse, err error) {
	reqInterface, ok := na.waitingResponseMap.Load(requestID)
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
			notary = n
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

/*
常驻发送线程,对应一个公证人
*/
func (na *NotaryAPI) notaryMsgSendLoop(targetNotaryID int, sendingQueueChan chan api.Req) {
	log.Trace("notaryMsgSendLoop with notary[ID=%d] start...", targetNotaryID)
	for {
		select {
		case req := <-sendingQueueChan:
			// 1. 如果是需要签名的请求,签名
			if reqWithSignature, ok := req.(api.ReqWithSignature); ok && reqWithSignature.GetSignature() == nil {
				reqWithSignature.Sign(reqWithSignature, na.selfPrivateKey)
			}
			for {
				// 2. 获取连接
				wsInterface, ok := na.wsConnMap.Load(targetNotaryID)
				if ok {
					//已有连接,直接使用
					api.WSWriteJSON(wsInterface.(*websocket.Conn), req)
					break
				}
				// 3. 没有连接,获取url origin,失败丢弃
				wsURL, origin, err := na.getNotaryWSURLAndOrigin(targetNotaryID)
				if err != nil {
					log.Error("send request to notary[ID=%d] err : %s, \nRequest:\n%s ignore", targetNotaryID, err.Error(), utils.ToJSONStringFormat(req))
					break
				}
				// 4. 建立连接,失败重试
				ws, err := websocket.Dial(wsURL, "", origin)
				if err != nil {
					log.Error("websocket connect to notary[ID=%d] err : %s, retry ...", targetNotaryID, err.Error())
					time.Sleep(time.Second)
					continue
				}
				// 5.保存连接,这里存在其他线程创建成功的情况,概率较低,如果出现,使用保存成功的连接,并关闭新创建的
				wsSavedInterface, loaded := na.wsConnMap.LoadOrStore(targetNotaryID, ws)
				wsSaved := wsSavedInterface.(*websocket.Conn)
				if loaded {
					log.Warn("websocket connection with notary[ID=%d] repeat create,use old one and close new", targetNotaryID)
					err = ws.Close()
					if err != nil {
						log.Error("close err : %s", err.Error())
					}
				} else {
					// 新建成功,启动消息处理线程
					go na.notaryMsgReceiveLoop(wsSaved, targetNotaryID)
				}
				// 6.使用保存成功的连接处理消息
				api.WSWriteJSON(wsSaved, req)
				break
			}
		}
	}
}
