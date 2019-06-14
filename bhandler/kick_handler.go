package bhandler

type KickHandler struct {
}

func (*KickHandler) GetType() byte {
	return 0xA9
}
func (h *KickHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)
	response.OpData = []byte{0x01}
	return &response
}
