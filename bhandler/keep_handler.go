package bhandler

import (
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
)

// KeepHandler keep
type KeepHandler struct {
	Logger *zap.Logger
}

// GetType 可以处理的消息类型
func (*KeepHandler) GetType() byte {
	return 0xA6
}

// GetResponse 根据请求获得响应
func (h *KeepHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//读取请求信息
	packetReader := services.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//等级
	playerLevel := packetReader.ReadUint16()
	h.Logger.Info(fmt.Sprintf("keep: user [%s] level %d", username, playerLevel))
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x01)
	response.OpData = opData
	return response
}
