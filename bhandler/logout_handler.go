package bhandler

import (
	"database/sql"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
)

// LogoutHandler 退出游戏
type LogoutHandler struct {
	Db          *sql.DB
	Logger      *zap.Logger
	LoginUsers  map[string]*common.ClientInfo //已登录,还未进入游戏的用户
	OnlineUsers map[string]*common.ClientInfo //已进入游戏的用户
	MacCounters map[string]int                //已进入游戏的用户的mac地址计数器
}

// GetType 可以处理的消息类型
func (*LogoutHandler) GetType() byte {
	return 0xA4
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
	if clientInfo, userOnline := h.OnlineUsers[usernameStr]; userOnline {
		delete(h.OnlineUsers, usernameStr)
		macMd5 := clientInfo.MacMd5
		macCounter := 0
		if value, valueExists := h.MacCounters[macMd5]; valueExists {
			macCounter = value
		}
		macCounter--
		if macCounter < 0 {
			macCounter = 0
		}
		h.MacCounters[macMd5] = macCounter
	}
	//
	h.Logger.Info("user [" + string(username) + "] logout game")
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	response.OpData = opData
	return response
}
