package bhandler

import "github.com/liuguangw/billing_go/common"

type ConnectHandler struct {
}

func (*ConnectHandler) GetType() byte {
	return 0xA0
}
func (h *ConnectHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	response.OpData = []byte{0x20, 0x00}
	return response
}
