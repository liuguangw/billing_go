package bhandler

import (
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
)

// LogoutHandler 退出游戏
type LogoutHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*LogoutHandler) GetType() byte {
	return packetTypeLogout
}

// GetResponse 根据请求获得响应
func (h *LogoutHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := services.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//更新在线状态
	usernameStr := string(username)
	if clientInfo, userOnline := h.Resource.OnlineUsers[usernameStr]; userOnline {
		delete(h.Resource.OnlineUsers, usernameStr)
		macMd5 := clientInfo.MacMd5
		if macMd5 != "" {
			macCounter := 0
			if value, valueExists := h.Resource.MacCounters[macMd5]; valueExists {
				macCounter = value
			}
			macCounter--
			if macCounter < 0 {
				macCounter = 0
			}
			h.Resource.MacCounters[macMd5] = macCounter
		}
	}
	//
	h.Resource.Logger.Info("user [" + string(username) + "] logout game")
	//Packets::BLRetBillingEnd
	opData := make([]byte, 0, usernameLength+2)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	response.OpData = opData
	return response
}
