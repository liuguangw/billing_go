package bhandler

import (
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"go.uber.org/zap"
)

type KeepHandler struct {
	Logger *zap.Logger
}

func (*KeepHandler) GetType() byte {
	return 0xA6
}
func (h *KeepHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//读取请求信息
	packetReader := common.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByte()
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
