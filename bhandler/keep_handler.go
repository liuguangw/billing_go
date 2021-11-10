package bhandler

import (
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
)

// KeepHandler keep
type KeepHandler struct {
	Logger      *zap.Logger
	LoginUsers  map[string]*common.ClientInfo //已登录,还未进入游戏的用户
	OnlineUsers map[string]*common.ClientInfo //已进入游戏的用户
	MacCounters map[string]int                //已进入游戏的用户的mac地址计数器
}

// GetType 可以处理的消息类型
func (*KeepHandler) GetType() byte {
	return 0xA6
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
	markOnline(h.LoginUsers, h.OnlineUsers, h.MacCounters, string(username), clientInfo)
	h.Logger.Info(fmt.Sprintf("keep: user [%s] level %d", username, playerLevel))
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x01)
	response.OpData = opData
	return response
}
