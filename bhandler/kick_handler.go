package bhandler

import "github.com/liuguangw/billing_go/common"

type KickHandler struct {
}

func (*KickHandler) GetType() byte {
	return 0xA9
}
func (h *KickHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	response.OpData = []byte{0x01}
	return response
}
