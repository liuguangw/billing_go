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

	usernameLength := request.OpData[0]
	username := request.OpData[1 : 1+usernameLength]
	offset := 1 + usernameLength
	playerLevel := uint16(request.OpData[offset])
	offset++
	playerLevel += uint16(request.OpData[offset])
	h.Logger.Info(fmt.Sprintf("keep: user [%v] level %v", string(username), playerLevel))
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x01)
	response.OpData = opData
	return response
}
