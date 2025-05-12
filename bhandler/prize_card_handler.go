package bhandler

import (
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
)

// PrizeCardHandler 物品卡奖励, 脚本ID: 808078
type PrizeCardHandler struct {
	Resource *common.HandlerResource
	BillType int //billing类型
}

// GetType 可以处理的消息类型
func (*PrizeCardHandler) GetType() byte {
	return packetTypePrizeCard
}

// GetResponse 根据请求获得响应
func (h *PrizeCardHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := services.NewPacketDataReader(request.OpData)
	//卡号
	cardLenth := packetReader.ReadByteValue()
	cardData := packetReader.ReadBytes(int(cardLenth))
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//登录IP
	tmpLength = int(packetReader.ReadByteValue())
	loginIP := string(packetReader.ReadBytes(tmpLength))
	//角色名
	tmpLength = int(packetReader.ReadByteValue())
	charName := packetReader.ReadGbkString(tmpLength)
	//标记在线
	clientInfo := &common.ClientInfo{
		IP:       loginIP,
		CharName: string(charName),
	}
	markOnline(h.Resource.LoginUsers, h.Resource.OnlineUsers, h.Resource.MacCounters, string(username), clientInfo)
	h.Resource.Logger.Info(string(charName) + "(" + string(username) + ") use card: " + string(cardData))
	opData := make([]byte, 0, usernameLength+2)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	//默认拒绝
	opData = append(opData, 1)
	response.OpData = opData
	return response
}
