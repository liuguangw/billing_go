package bhandler

import (
	"database/sql"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// EnterGameHandler 进入游戏
type EnterGameHandler struct {
	Db          *sql.DB
	Logger      *zap.Logger
	LoginUsers  map[string]*common.ClientInfo //已登录,还未进入游戏的用户
	OnlineUsers map[string]*common.ClientInfo //已进入游戏的用户
	MacCounters map[string]int                //已进入游戏的用户的mac地址计数器
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
		h.Logger.Error("decode char name failed: " + err.Error())
		charName = []byte("?")
	}
	//标记在线
	clientInfo := &common.ClientInfo{
		CharName: string(charName),
	}
	markOnline(h.LoginUsers, h.OnlineUsers, h.MacCounters, string(username), clientInfo)
	//
	h.Logger.Info("user [" + string(username) + "] " + string(charName) + " entered game")
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	response.OpData = opData
	return response
}
