package bhandler

import (
	"billing/models"
	"billing/tools"
	"database/sql"
	"fmt"
)

type QueryPointHandler struct {
	Db *sql.DB
}

func (*QueryPointHandler) GetType() byte {
	return 0xE2
}
func (h *QueryPointHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)
	var opData []byte
	//用户名
	offset := 0
	usernameLength := request.OpData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := request.OpData[offset : offset+tmpLength]
	//登录IP
	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	loginIP := string(request.OpData[offset : offset+tmpLength])
	//角色名
	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	charName := string(request.OpData[offset : offset+tmpLength])
	// 更新在线状态
	err := models.UpdateOnlineStatus(h.Db, string(username), true)
	if err != nil {
		tools.ShowErrorInfo("update username:"+string(username)+" to online failed", err)
	}
	account, queryOp := models.GetAccountByUsername(h.Db, string(username))
	var accountPoint int32
	if queryOp == models.UserFound {
		accountPoint = (account.Point + 1) * 1000
	}
	tools.LogMessage(fmt.Sprintf("user [%v] %v query point (%v) at %v", string(username), charName, account.Point, loginIP))
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	var tmpByte byte
	tmpByte = byte(accountPoint >> 24)
	opData = append(opData, tmpByte)
	tmpByte = byte((accountPoint >> 16) & 0xff)
	opData = append(opData, tmpByte)
	tmpByte = byte((accountPoint >> 8) & 0xff)
	opData = append(opData, tmpByte)
	tmpByte = byte(accountPoint & 0xff)
	opData = append(opData, tmpByte)
	response.OpData = opData
	return &response
}
