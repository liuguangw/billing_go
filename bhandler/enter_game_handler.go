package bhandler

import (
	"database/sql"
	"github.com/liuguangw/billing_go/common"
	"go.uber.org/zap"
)

type EnterGameHandler struct {
	Db     *sql.DB
	Logger *zap.Logger
}

func (*EnterGameHandler) GetType() byte {
	return 0xA3
}
func (h *EnterGameHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	var opData []byte
	offset := 0
	usernameLength := request.OpData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := request.OpData[offset : offset+tmpLength]

	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	charName := string(request.OpData[offset : offset+tmpLength])
	//todo 更新在线状态
	//
	h.Logger.Info("user [" + string(username) + "] " + charName + " entered game")
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	response.OpData = opData
	return response
}
