package bhandler

import "github.com/liuguangw/billing_go/common"

// KickHandler Kick处理
type KickHandler struct {
}

// GetType 可以处理的消息类型
func (*KickHandler) GetType() byte {
	return packetTypeKick
}

// GetResponse 根据请求获得响应
func (h *KickHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//Packets::BLRetKickALL
	response.OpData = []byte{0x01}
	return response
}
