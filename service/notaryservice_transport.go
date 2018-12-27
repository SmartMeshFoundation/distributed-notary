package service

import (
	"io"
	"net/http"
	"strings"

	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api/notaryapi"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
BroadcastMsg :
群发消息到各个notary的指定api TODO 目前全是同步
*/
func (ns *NotaryService) BroadcastMsg(sessionID common.Hash, apiName string, msg interface{}, isSync bool) (err error) {
	for _, notary := range ns.notaries {
		err = ns.SendMsg(sessionID, apiName, notary.ID, msg)
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
func (ns *NotaryService) SendMsg(sessionID common.Hash, apiName string, notaryID int, msg interface{}) (err error) {
	url := ns.getNotaryHostByID(notaryID) + notaryapi.APIName2URLMap[apiName]
	var payload string
	switch m := msg.(type) {
	case *models.KeyGenBroadcastMessage1:
		req := notaryapi.NewKeyGenerationPhase1MessageRequest(sessionID, crypto.PubkeyToAddress(ns.self.PublicKey), m)
		req.Sign(ns.privateKey)
		payload = req.ToJSONString()
	case *models.KeyGenBroadcastMessage2:
		req := notaryapi.NewKeyGenerationPhase2MessageRequest(sessionID, crypto.PubkeyToAddress(ns.self.PublicKey), m)
		req.Sign(ns.privateKey)
		payload = req.ToJSONString()
	case *models.KeyGenBroadcastMessage3:
		req := notaryapi.NewKeyGenerationPhase3MessageRequest(sessionID, crypto.PubkeyToAddress(ns.self.PublicKey), m)
		req.Sign(ns.privateKey)
		payload = req.ToJSONString()
	case *models.KeyGenBroadcastMessage4:
		req := notaryapi.NewKeyGenerationPhase4MessageRequest(sessionID, crypto.PubkeyToAddress(ns.self.PublicKey), m)
		req.Sign(ns.privateKey)
		payload = req.ToJSONString()
	}
	return doPost(url, payload)
}

func (ns *NotaryService) getNotaryHostByID(notaryID int) string {
	for _, v := range ns.notaries {
		if v.ID == notaryID {
			return v.Host
		}
	}
	return ""
}

func doPost(url string, payload string) (err error) {
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
	return
}
