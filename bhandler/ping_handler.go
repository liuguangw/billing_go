package bhandler

type PingHandler struct {
}

func (*PingHandler) GetType() byte {
	return 0xA1
}
func (h *PingHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)
	response.OpData = []byte{0x01, 0x00}
	return &response
}
