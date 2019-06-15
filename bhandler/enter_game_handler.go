package bhandler

import (
	"database/sql"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/tools"
)

type EnterGameHandler struct {
	Db *sql.DB
}

func (*EnterGameHandler) GetType() byte {
	return 0xA3
}
func (h *EnterGameHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)
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
	// 更新在线状态
	err := models.UpdateOnlineStatus(h.Db, string(username), true)
	if err != nil {
		tools.ShowErrorInfo("update username:"+string(username)+" to online failed", err)
	}
	tools.LogMessage("user [" + string(username) + "] " + charName + " entered game")
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	response.OpData = opData
	return &response
}
