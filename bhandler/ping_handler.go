package bhandler

import "github.com/liuguangw/billing_go/common"

// PingHandler ping
type PingHandler struct {
}

// GetType 可以处理的消息类型
func (*PingHandler) GetType() byte {
	return 0xA1
}

// GetResponse 根据请求获得响应
func (h *PingHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	response.OpData = []byte{0x01, 0x00}
	return response
}
