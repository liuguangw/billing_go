package bhandler

import (
	"strconv"

	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
)

// EnterGameHandler 进入游戏
type EnterGameHandler struct {
	Resource *common.HandlerResource
	BillType int //billing类型
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
	charName := packetReader.ReadGbkString(tmpLength)
	//标记在线
	clientInfo := &common.ClientInfo{
		CharName: string(charName),
	}
	markOnline(h.Resource.LoginUsers, h.Resource.OnlineUsers, h.Resource.MacCounters, string(username), clientInfo)
	//角色id
	charguid := packetReader.ReadInt()
	//
	h.Resource.Logger.Info("user [" + string(username) + "] " + string(charName) + " entered game")
	//Packets::BLRetBillingStart
	opData := make([]byte, 0, usernameLength+2+14)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	//额外数据
	//mFeeType: 1u
	//mLeftTime: 4u
	//mStorePoint: 4u
	//mUserPoint: 4u
	//mWhyFlag: 1u
	padLen := 14
	if h.BillType == common.BillTypeHuaiJiu {
		padLen = 16
	}
	extraData := make([]byte, padLen)
	//检查角色是否为gm(仅限怀旧)
	if h.BillType == common.BillTypeHuaiJiu {
		if isGm, err := models.CheckIsGm(h.Resource.Db, charguid); err != nil {
			h.Resource.Logger.Error("check is gm failed: " + err.Error())
		} else {
			// 登录的角色为gm
			if isGm {
				extraData[padLen-1] = 1
				h.Resource.Logger.Info(string(charName) + "(guid: " + strconv.Itoa(charguid) + ") is GM")
			}
		}
		//检测活动奖励状态
		var state byte
		//检测角色活动奖励
		if accountPrize, err := models.FindFirstAccountPrize(h.Resource.Db, string(username)); err != nil {
			h.Resource.Logger.Error("FindFirstAccountPrize failed: " + err.Error())
		} else {
			//有奖励记录
			if accountPrize != nil {
				state = 1
				if accountPrize.Charguid != 0 {
					state = 2
				}
			}
		}
		if state > 0 {
			extraData[padLen-2] = state
		}
	}
	opData = append(opData, extraData...)
	response.OpData = opData
	return response
}
