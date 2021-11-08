package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"go.uber.org/zap"
)

type QueryPointHandler struct {
	Db     *sql.DB
	Logger *zap.Logger
}

func (*QueryPointHandler) GetType() byte {
	return 0xE2
}
func (h *QueryPointHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
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
	// todo 更新在线状态
	account, err := models.GetAccountByUsername(h.Db, string(username))
	if err != nil {
		h.Logger.Error("get account:" + string(username) + " info failed: " + err.Error())
	}
	var accountPoint = 0
	if account != nil {
		accountPoint = (account.Point + 1) * 1000
	}
	h.Logger.Info(fmt.Sprintf("user [%v] %v query point (%v) at %v", string(username), charName, account.Point, loginIP))
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
	return response
}
