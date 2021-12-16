package bhandler

import (
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
	"golang.org/x/text/encoding/simplifiedchinese"
)

//兑换结果
const (
	convertSuccess byte = 0 //兑换成功
	convertFailed  byte = 1 //兑换失败
)

// ConvertPointHandler 处理点数兑换
type ConvertPointHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*ConvertPointHandler) GetType() byte {
	return packetTypeConvertPoint
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
		h.Resource.Logger.Error("decode char name failed: " + err.Error())
		charName = []byte("?")
	}
	//orderId 21u
	orderIDBytes := packetReader.ReadBytes(21)
	mGoodsTypeNum := packetReader.ReadUint16() //始终为1
	//物品类型: 999表示买元宝
	mGoodsType := packetReader.ReadInt()
	//fmt.Println(mGoodsType)
	//物品数量(输入的点数)
	mGoodsNumber := packetReader.ReadUint16()
	//需要兑换的点数
	needPoint := int(mGoodsNumber)
	//每次兑换点数上限 u2
	var maxPoint = 0xffff //65535
	if needPoint > maxPoint {
		needPoint = maxPoint
	}
	//初始化兑换的结果
	convertResult := convertSuccess
	convertResultText := "success"
	//查询数据库获取用户当前点数余额
	userPoint := 0
	account, err := models.GetAccountByUsername(h.Resource.Db, string(username))
	if err != nil {
		convertResult = convertFailed
		convertResultText = "get account info error"
		h.Resource.Logger.Error("get account:" + string(username) + " info failed: " + err.Error())
	}
	if account != nil {
		userPoint = account.Point
	}
	//初始化剩余点数
	leftPoint := userPoint
	if leftPoint < 0 {
		leftPoint = 0
	}
	//初始化兑换的点数
	realPoint := needPoint
	//兑换的点数不能超过用户拥有的点数,如果超过则截断
	if realPoint > leftPoint {
		realPoint = leftPoint
	}
	//兑换的点数必须是正整数
	if realPoint <= 0 {
		convertResult = convertFailed
		convertResultText = "invalid need point value"
	}
	// 执行兑换
	if convertResult == convertSuccess {
		if err := models.ConvertUserPoint(h.Resource.Db, string(username), realPoint); err != nil {
			h.Resource.Logger.Error("convert point failed: " + err.Error())
			convertResult = convertFailed
			convertResultText = "convert point failed"
		} else {
			leftPoint -= realPoint
		}
	}
	//日志记录
	h.Resource.Logger.Info(fmt.Sprintf("user [%s] %s(ip: %s) "+
		"point total [%d], need point [%d],"+
		" (%d - %d = %d): %s",
		username, charName, loginIP,
		userPoint, needPoint,
		userPoint, realPoint, leftPoint,
		convertResultText))
	// 数据包组合
	opData := make([]byte, 0, 1+usernameLength+21+1+4+2+4+2)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, orderIDBytes...)
	opData = append(opData, convertResult)
	//写入剩余点数:u4(此值不会被用到,服务端以购买的数量来进行计算)
	leftPointU4 := leftPoint * 1000
	for i := 0; i < 4; i++ {
		tmpValue := leftPointU4
		movePos := (3 - i) * 8
		if movePos > 0 {
			tmpValue >>= movePos
		}
		opData = append(opData, byte(tmpValue&0xff))
	}
	//mGoodsTypeNum:u2
	opData = append(opData, byte((mGoodsTypeNum&0xff00)>>8), byte(mGoodsTypeNum&0xff))
	// 写入mGoodsType:u4
	for i := 0; i < 4; i++ {
		tmpValue := mGoodsType
		movePos := (3 - i) * 8
		if movePos > 0 {
			tmpValue >>= movePos
		}
		opData = append(opData, byte(tmpValue&0xff))
	}
	//写入mGoodsNumber(实际购买的数量):u2
	if convertResult != convertSuccess {
		realPoint = 0
	}
	opData = append(opData, byte((realPoint&0xff00)>>8), byte(realPoint&0xff))
	response.OpData = opData
	return response
}
