package bhandler

import (
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// EnterGameHandler 进入游戏
type EnterGameHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*EnterGameHandler) GetType() byte {
	return packetTypeEnterGame
}

// GetResponse 根据请求获得响应
func (h *EnterGameHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//读取请求信息
	packetReader := services.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//角色名
	tmpLength = int(packetReader.ReadByteValue())
	charNameGbkData := packetReader.ReadBytes(tmpLength)
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	charName, err := gbkDecoder.Bytes(charNameGbkData)
	if err != nil {
		h.Resource.Logger.Error("decode char name failed: " + err.Error())
		charName = []byte("?")
	}
	//标记在线
	clientInfo := &common.ClientInfo{
		CharName: string(charName),
	}
	markOnline(h.Resource.LoginUsers, h.Resource.OnlineUsers, h.Resource.MacCounters, string(username), clientInfo)
	//
	h.Resource.Logger.Info("user [" + string(username) + "] " + string(charName) + " entered game")
	//Packets::BLRetBillingStart
	opData := make([]byte, 0, usernameLength+2+13)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	//额外数据
	//mFeeType: 1u
	//mLeftTime: 4u
	//mStorePoint: 4u
	//mUserPoint: 4u
	extraData := make([]byte, 13)
	opData = append(opData, extraData...)
	response.OpData = opData
	return response
}
