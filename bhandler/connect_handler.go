package bhandler

type ConnectHandler struct {
}

func (*ConnectHandler) GetType() byte {
	return 0xA0
}
func (h *ConnectHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)
	response.OpData = []byte{0x20, 0x00}
	return &response
}
