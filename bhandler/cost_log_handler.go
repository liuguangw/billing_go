package bhandler

import "github.com/liuguangw/billing_go/common"

type CostLogHandler struct {
}

func (*CostLogHandler) GetType() byte {
	return 0xC5
}
func (h *CostLogHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	opData := request.OpData[0:21]
	//tools.LogMessage(fmt.Sprintf("CostLog - mSerial key=%s", string(opData)))
	opData = append(opData, 0x01)
	response.OpData = opData
	return response
}
