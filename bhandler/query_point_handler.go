package bhandler

import (
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// QueryPointHandler 查询点数
type QueryPointHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*QueryPointHandler) GetType() byte {
	return packetTypeQueryPoint
}

// GetResponse 根据请求获得响应
func (h *QueryPointHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
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
	account, err := models.GetAccountByUsername(h.Resource.Db, string(username))
	if err != nil {
		h.Resource.Logger.Error("get account:" + string(username) + " info failed: " + err.Error())
	}
	//标记在线
	clientInfo := &common.ClientInfo{
		IP:       loginIP,
		CharName: string(charName),
	}
	markOnline(h.Resource.LoginUsers, h.Resource.OnlineUsers, h.Resource.MacCounters, string(username), clientInfo)
	//
	var accountPoint = 0
	if account != nil {
		accountPoint = account.Point
	}
	h.Resource.Logger.Info(fmt.Sprintf("user [%s] %s query point (%d) at %s", username, charName, accountPoint, loginIP))
	//Packets::BLRetAskPoint
	opData := make([]byte, 0, usernameLength+5)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	accountPoint *= 1000
	for i := 0; i < 4; i++ {
		tmpValue := accountPoint
		movePos := (3 - i) * 8
		if movePos > 0 {
			tmpValue >>= movePos
		}
		opData = append(opData, byte(tmpValue&0xff))
	}
	response.OpData = opData
	return response
}
