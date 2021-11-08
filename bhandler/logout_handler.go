package bhandler

import (
	"database/sql"
	"github.com/liuguangw/billing_go/common"
	"go.uber.org/zap"
)

type LogoutHandler struct {
	Db     *sql.DB
	Logger *zap.Logger
}

func (*LogoutHandler) GetType() byte {
	return 0xA4
}
func (h *LogoutHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := common.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByte()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	// todo 更新在线状态
	h.Logger.Info("user [" + string(username) + "] logout game")
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	response.OpData = opData
	return response
}
