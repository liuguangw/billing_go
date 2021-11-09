package bhandler

import (
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
)

// CostLogHandler 元宝消息记录
type CostLogHandler struct {
}

// GetType 可以处理的消息类型
func (*CostLogHandler) GetType() byte {
	return 0xC5
}

// GetResponse 根据请求获得响应
func (h *CostLogHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := services.NewPacketDataReader(request.OpData)
	mSerialKeyLength := 21
	mSerialKey := packetReader.ReadBytes(mSerialKeyLength)
	//tools.LogMessage(fmt.Sprintf("CostLog - mSerial key=%s", string(opData)))
	opData := make([]byte, 0, mSerialKeyLength+1)
	opData = append(opData, mSerialKey...)
	opData = append(opData, 0x01)
	response.OpData = opData
	return response
}
