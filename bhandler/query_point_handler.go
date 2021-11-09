package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// QueryPointHandler 查询点数
type QueryPointHandler struct {
	Db     *sql.DB
	Logger *zap.Logger
}

// GetType 可以处理的消息类型
func (*QueryPointHandler) GetType() byte {
	return 0xE2
}

// GetResponse 根据请求获得响应
func (h *QueryPointHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := services.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//登录IP
	tmpLength = int(packetReader.ReadByteValue())
	loginIP := string(packetReader.ReadBytes(tmpLength))
	//角色名
	tmpLength = int(packetReader.ReadByteValue())
	charNameGbkData := packetReader.ReadBytes(tmpLength)
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	charName, err := gbkDecoder.Bytes(charNameGbkData)
	if err != nil {
		h.Logger.Error("decode char name failed: " + err.Error())
		charName = []byte("?")
	}
	account, err := models.GetAccountByUsername(h.Db, string(username))
	if err != nil {
		h.Logger.Error("get account:" + string(username) + " info failed: " + err.Error())
	}
	var accountPoint = 0
	if account != nil {
		accountPoint = (account.Point + 1) * 1000
	}
	h.Logger.Info(fmt.Sprintf("user [%s] %s query point (%d) at %s", username, charName, account.Point, loginIP))
	var opData []byte
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
