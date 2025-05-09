package bhandler

import "github.com/liuguangw/billing_go/common"

// ConnectHandler 处理Connect
type ConnectHandler struct {
}

// GetType 可以处理的消息类型
func (*ConnectHandler) GetType() byte {
	return packetTypeConnect
}

// GetResponse 根据请求获得响应
func (h *ConnectHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//Packets::BLRetConnect
	response.OpData = []byte{0, 0}
	return response
}
