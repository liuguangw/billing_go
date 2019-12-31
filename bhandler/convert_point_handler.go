package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/tools"
)

type ConvertPointHandler struct {
	Db            *sql.DB
	ConvertNumber int
}

func (*ConvertPointHandler) GetType() byte {
	return 0xE1
}
func (h *ConvertPointHandler) GetResponse(request *BillingData) *BillingData {
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
	loginIP := string(request.OpData[offset : offset+tmpLength])

	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	charName := string(request.OpData[offset : offset+tmpLength])

	//orderId 21u
	offset += tmpLength
	orderIDBytes := request.OpData[offset : offset+21]
	// extraData 6u
	offset += 21
	extraDataBytes := request.OpData[offset : offset+6]
	//跳过本身6字节+兑换点数的后2字节
	offset += 8
	//获取需要兑换的点数:4u
	needPoint := 0
	for i := 0; i < 4; i++ {
		tmpInt := int(request.OpData[offset])
		offset++
		if i < 3 {
			tmpInt = tmpInt << uint((3-i)*8)
		}
		needPoint += tmpInt
	}
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
	if err!=nil{
		tools.ShowErrorInfo("get account:"+string(username)+" info failed", err)
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
		tools.ShowErrorInfo("convert point failed", err)
		realPoint = 0
	} else {
		tools.LogMessage(fmt.Sprintf("user [%s] %v(ip: %v) point total [%v], need point [%v]: %v-%v=%v",
			username, charName, loginIP, userPoint, needPoint,
			userPoint, realPoint, userPoint-realPoint))
	}
	// 数据包组合
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, orderIDBytes...)
	tmpBytes := []byte{0x00, 0x00, 0x00, 0x03, 0xE8}
	opData = append(opData, tmpBytes...)
	opData = append(opData, extraDataBytes...)
	// 点数 2u
	opData = append(opData, byte((realPoint&0xff00)>>8), byte(realPoint&0xff))
	response.OpData = opData
	return &response
}
