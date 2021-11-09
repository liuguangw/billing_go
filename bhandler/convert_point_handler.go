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

// ConvertPointHandler 处理点数兑换
type ConvertPointHandler struct {
	Db            *sql.DB
	Logger        *zap.Logger
	ConvertNumber int
}

// GetType 可以处理的消息类型
func (*ConvertPointHandler) GetType() byte {
	return 0xE1
}

// GetResponse 根据请求获得响应
func (h *ConvertPointHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
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
	//orderId 21u
	orderIDBytes := packetReader.ReadBytes(21)
	// extraData 6u
	extraDataBytes := packetReader.ReadBytes(6)
	//跳过兑换点数的后2字节
	packetReader.Skip(2)
	//获取需要兑换的点数:4u
	needPoint := packetReader.ReadInt()
	needPoint /= h.ConvertNumber
	if needPoint < 0 {
		needPoint = 0
	}
	//每次兑换点数上限 u2
	var maxPoint = 0xffff //65535
	if needPoint > maxPoint {
		needPoint = maxPoint
	}
	userPoint := 0
	//获取用户当前点数总额
	account, err := models.GetAccountByUsername(h.Db, string(username))
	if err != nil {
		h.Logger.Error("get account:" + string(username) + " info failed: " + err.Error())
	}
	if account != nil {
		userPoint = account.Point
		if userPoint < 0 {
			userPoint = 0
		}
	}
	//最终可兑换的点数
	var realPoint int
	if needPoint > userPoint {
		realPoint = userPoint
	} else {
		realPoint = needPoint
	}
	// 执行兑换
	err = models.ConvertUserPoint(h.Db, string(username), realPoint)
	if err != nil {
		h.Logger.Error("convert point failed: " + err.Error())
		realPoint = 0
	} else {
		h.Logger.Info(fmt.Sprintf("user [%s] %s(ip: %s) point total [%d], need point [%d]: %d-%d=%d",
			username, charName, loginIP, userPoint, needPoint,
			userPoint, realPoint, userPoint-realPoint))
	}
	// 数据包组合
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, orderIDBytes...)
	tmpBytes := []byte{0x00, 0x00, 0x00, 0x03, 0xE8}
	opData = append(opData, tmpBytes...)
	opData = append(opData, extraDataBytes...)
	// 点数 2u
	opData = append(opData, byte((realPoint&0xff00)>>8), byte(realPoint&0xff))
	response.OpData = opData
	return response
}
