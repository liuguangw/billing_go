package bhandler

import (
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// CostLogHandler 元宝消息记录
type CostLogHandler struct {
	Logger      *zap.Logger
	LoginUsers  map[string]*common.ClientInfo //已登录,还未进入游戏的用户
	OnlineUsers map[string]*common.ClientInfo //已进入游戏的用户
	MacCounters map[string]int                //已进入游戏的用户的mac地址计数器
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
	charNameGbkData := packetReader.ReadBytes(tmpLength)
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	charName, err := gbkDecoder.Bytes(charNameGbkData)
	if err != nil {
		h.Logger.Error("decode char name failed: " + err.Error())
		charName = []byte("?")
	}
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
	markOnline(h.LoginUsers, h.OnlineUsers, h.MacCounters, string(username), clientInfo)
	//
	opData := make([]byte, 0, mSerialKeyLength+1)
	opData = append(opData, mSerialKey...)
	opData = append(opData, 0x01)
	response.OpData = opData
	return response
}
