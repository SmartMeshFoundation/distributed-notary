package service

/*
BroadcastMsg :
群发消息到各个notary的指定api
*/
func (ns *NotaryService) BroadcastMsg(apiName string, msg interface{}, isSync bool) (err error) {
	// TODO
	return
}

/*
SendMsg :
同步请求
*/
func (ns *NotaryService) SendMsg(apiName string, notaryID int, msg interface{}) (err error) {
	// TODO
	return
}
