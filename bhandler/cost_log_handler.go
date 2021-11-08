package bhandler

import "github.com/liuguangw/billing_go/common"

type CostLogHandler struct {
}

func (*CostLogHandler) GetType() byte {
	return 0xC5
}
func (h *CostLogHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := common.NewPacketDataReader(request.OpData)
	mSerialKeyLength := 21
	mSerialKey := packetReader.ReadBytes(mSerialKeyLength)
	//tools.LogMessage(fmt.Sprintf("CostLog - mSerial key=%s", string(opData)))
	opData := make([]byte, 0, mSerialKeyLength+1)
	opData = append(opData, mSerialKey...)
	opData = append(opData, 0x01)
	response.OpData = opData
	return response
}
