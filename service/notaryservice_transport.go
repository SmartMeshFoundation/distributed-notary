package service

import (
	"io"
	"net/http"
	"strings"

	"fmt"

	"encoding/json"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kataras/iris/core/errors"
	"github.com/nkbai/log"
)

/*
BroadcastMsg :
群发消息到各个notary的指定api TODO 目前全是同步
*/
func (ns *NotaryService) BroadcastMsg(sessionID common.Hash, apiName string, msg interface{}, isSync bool, notaryIDs ...int) (err error) {
	if len(notaryIDs) > 0 {
		for _, notaryID := range notaryIDs {
			if notaryID == ns.self.ID {
				continue
			}
			_, err = ns.SendMsg(sessionID, apiName, notaryID, msg)
			if err != nil {
				return
			}
		}
		return
	}
	for _, notary := range ns.notaries {
		if notary.ID == ns.self.ID {
			continue
		}
		_, err = ns.SendMsg(sessionID, apiName, notary.ID, msg)
		if err != nil {
			return
		}
	}
	return
}

/*
SendMsg :
同步请求
*/
func (ns *NotaryService) SendMsg(sessionID common.Hash, apiName string, notaryID int, msg interface{}) (response api.BaseResponse, err error) {
	url := ns.getNotaryHostByID(notaryID) + notaryapi.APIName2URLMap[apiName]
	var payload string
	switch m := msg.(type) {
	/*
		pkn
	*/
	case *models.KeyGenBroadcastMessage1:
		req := notaryapi.NewKeyGenerationPhase1MessageRequest(sessionID, ns.self.GetAddress(), m)
		api.NotarySign(req, ns.privateKey)
		payload = utils.ToJSONString(req)
	case *models.KeyGenBroadcastMessage2:
		req := notaryapi.NewKeyGenerationPhase2MessageRequest(sessionID, ns.self.GetAddress(), m)
		api.NotarySign(req, ns.privateKey)
		payload = utils.ToJSONString(req)
	case *models.KeyGenBroadcastMessage3:
		req := notaryapi.NewKeyGenerationPhase3MessageRequest(sessionID, ns.self.GetAddress(), m)
		api.NotarySign(req, ns.privateKey)
		payload = utils.ToJSONString(req)
	case *models.KeyGenBroadcastMessage4:
		req := notaryapi.NewKeyGenerationPhase4MessageRequest(sessionID, ns.self.GetAddress(), m)
		api.NotarySign(req, ns.privateKey)
		payload = utils.ToJSONString(req)
	/*
		dsm
	*/
	case *notaryapi.DSMAskRequest:
		api.NotarySign(m, ns.privateKey)
		payload = utils.ToJSONString(m)
	case *notaryapi.DSMNotifySelectionRequest:
		api.NotarySign(m, ns.privateKey)
		payload = utils.ToJSONString(m)
	case *notaryapi.DSMPhase1BroadcastRequest:
		api.NotarySign(m, ns.privateKey)
		payload = utils.ToJSONString(m)
	case *notaryapi.DSMPhase2MessageARequest:
		api.NotarySign(m, ns.privateKey)
		payload = utils.ToJSONString(m)
	case *notaryapi.DSMPhase3DeltaIRequest:
		api.NotarySign(m, ns.privateKey)
		payload = utils.ToJSONString(m)
	case *notaryapi.DSMPhase5A5BProofRequest:
		api.NotarySign(m, ns.privateKey)
		payload = utils.ToJSONString(m)
	case *notaryapi.DSMPhase5CProofRequest:
		api.NotarySign(m, ns.privateKey)
		payload = utils.ToJSONString(m)
	case *notaryapi.DSMPhase6ReceiveSIRequest:
		api.NotarySign(m, ns.privateKey)
		payload = utils.ToJSONString(m)
	default:
		err = errors.New("api call not expect")
		return
	}
	return doPost(sessionID, url, payload)
}

func (ns *NotaryService) getNotaryHostByID(notaryID int) string {
	for _, v := range ns.notaries {
		if v.ID == notaryID {
			return v.Host
		}
	}
	return ""
}

func doPost(sessionID common.Hash, url string, payload string) (response api.BaseResponse, err error) {
	//log.Trace(SessionLogMsg(sessionID, "post to %s, payload : %s", url, payload))
	log.Trace(SessionLogMsg(sessionID, "post to %s", url))
	var reqBody io.Reader
	if payload == "" {
		reqBody = nil
	} else {
		reqBody = strings.NewReader(payload)
	}
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if req.Body != nil {
			err = req.Body.Close()
		}
		if resp != nil && resp.Body != nil {
			err = resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http request err : status code = %d", resp.StatusCode)
	}
	var buf [4096 * 1024]byte
	n := 0
	n, err = resp.Body.Read(buf[:])
	if err != nil && err.Error() == "EOF" {
		err = nil
	}
	respBody := buf[:n]
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}
	log.Trace(SessionLogMsg(sessionID, "get response %s", utils.ToJSONString(response)))
	if response.ErrorCode != api.ErrorCodeSuccess {
		log.Error(SessionLogMsg(sessionID, response.ErrorMsg))
		err = errors.New(response.ErrorMsg)
	}
	return
}
