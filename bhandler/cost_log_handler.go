package bhandler

import (
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
)

// CostLogHandler 元宝消息记录
type CostLogHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*CostLogHandler) GetType() byte {
	return packetTypeCostLog
}

// GetResponse 根据请求获得响应
func (h *CostLogHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := services.NewPacketDataReader(request.OpData)
	mSerialKeyLength := 21
	mSerialKey := packetReader.ReadBytes(mSerialKeyLength)
	//skip zoneId(u2)
	//     +mWorldId(u4)+mServerId(u4)+mSceneId(u4)
	//     +mUserGUID(u4)+mCostTime(u4)+mYuanBao(u4)
	packetReader.Skip(26)
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//角色名
	tmpLength = int(packetReader.ReadByteValue())
	charName := packetReader.ReadGbkString(tmpLength)
	//skip level(u2)
	packetReader.Skip(2)
	//登录IP
	tmpLength = int(packetReader.ReadByteValue())
	loginIP := string(packetReader.ReadBytes(tmpLength))
	//标记在线
	clientInfo := &common.ClientInfo{
		IP:       loginIP,
		CharName: string(charName),
	}
	markOnline(h.Resource.LoginUsers, h.Resource.OnlineUsers, h.Resource.MacCounters, string(username), clientInfo)
	//Packets::LBLCostLog
	opData := make([]byte, 0, mSerialKeyLength+1)
	opData = append(opData, mSerialKey...)
	opData = append(opData, 0x01)
	response.OpData = opData
	return response
}
