package bhandler

import (
	"database/sql"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/tools"
)

type LogoutHandler struct {
	Db *sql.DB
}

func (*LogoutHandler) GetType() byte {
	return 0xA4
}
func (h *LogoutHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)
	var opData []byte
	offset := 0
	usernameLength := request.OpData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := request.OpData[offset : offset+tmpLength]
	// 更新在线状态
	err := models.UpdateOnlineStatus(h.Db, string(username), false)
	if err != nil {
		tools.ShowErrorInfo("update username:"+string(username)+" to offline failed", err)
	}
	tools.LogMessage("user [" + string(username) + "] logout game")
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	response.OpData = opData
	return &response
}
