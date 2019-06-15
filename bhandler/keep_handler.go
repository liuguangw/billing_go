package bhandler

import (
	"fmt"
	"github.com/liuguangw/billing_go/tools"
)

type KeepHandler struct {
}

func (*KeepHandler) GetType() byte {
	return 0xA6
}
func (h *KeepHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)

	usernameLength := request.OpData[0]
	username := request.OpData[1 : 1+usernameLength]
	offset := 1 + usernameLength
	playerLevel := uint16(request.OpData[offset])
	offset++
	playerLevel += uint16(request.OpData[offset])
	tools.LogMessage(fmt.Sprintf("keep: user [%v] level %v", string(username), playerLevel))
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x01)
	response.OpData = opData
	return &response
}
