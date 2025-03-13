package bhandler

import (
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
)

// KeepHandler keep
type KeepHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*KeepHandler) GetType() byte {
	return packetTypeKeep
}

// GetResponse 根据请求获得响应
func (h *KeepHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//读取请求信息
	packetReader := services.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//等级
	playerLevel := packetReader.ReadUint16()
	//标记在线
	clientInfo := &common.ClientInfo{}
	markOnline(h.Resource.LoginUsers, h.Resource.OnlineUsers, h.Resource.MacCounters, string(username), clientInfo)
	h.Resource.Logger.Info(fmt.Sprintf("keep: user [%s] level %d", username, playerLevel))
	//Packets::BLRetBillingKeep
	opData := make([]byte, 0, usernameLength+2+12)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x01)
	//额外数据
	//mLeftTime: 4U
	//mStorePoint: 4U
	//mUserPoint: 4U
	extraData := make([]byte, 12)
	//fake mLeftTime: 100
	extraData[3] = 100
	opData = append(opData, extraData...)
	response.OpData = opData
	return response
}
