package bhandler

import "github.com/liuguangw/billing_go/common"

type PingHandler struct {
}

func (*PingHandler) GetType() byte {
	return 0xA1
}
func (h *PingHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	response.OpData = []byte{0x01, 0x00}
	return response
}
