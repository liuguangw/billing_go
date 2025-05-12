package bhandler

import (
	"strconv"

	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
)

// PrizeFetchHandler 活动奖励领取(仅经典), 脚本ID: 808062
type PrizeFetchHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*PrizeFetchHandler) GetType() byte {
	return packetTypePrizeFetch
}

// GetResponse 根据请求获得响应
func (h *PrizeFetchHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
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
	charName := packetReader.ReadGbkString(tmpLength)
	//标记在线
	clientInfo := &common.ClientInfo{
		IP:       loginIP,
		CharName: string(charName),
	}
	markOnline(h.Resource.LoginUsers, h.Resource.OnlineUsers, h.Resource.MacCounters, string(username), clientInfo)

	//角色id
	charguid := packetReader.ReadInt()
	//等级
	charLv := packetReader.ReadUint16()
	//订单号
	prizeSerial := packetReader.ReadBytes(21)
	//初始化result
	var fetchResult byte
	prizeList, err := models.FindAccountPrizeList(h.Resource.Db, string(username), 0, 0, 30)
	if err != nil {
		fetchResult = 3
		h.Resource.Logger.Error("FindAccountPrizeList failed: " + err.Error())
	}
	//记录条数
	prizeCount := len(prizeList)
	if prizeCount == 0 {
		//没有奖励可以领取
		fetchResult = 5
	}
	//Packets::LBLNewPrize
	opDataLength := int(usernameLength) + 23
	if fetchResult == 0 {
		opDataLength += (21*prizeCount + 1)
	}
	opData := make([]byte, 0, opDataLength)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, prizeSerial...)
	opData = append(opData, fetchResult)
	if fetchResult == 0 {
		opData = append(opData, byte(prizeCount))
		var itemIdList []int64
		for _, prizeItem := range prizeList {
			//id
			itemIdList = append(itemIdList, prizeItem.ID)
			//itemid
			itemIdData := make([]byte, 20)
			itemIdStrData := []byte("item " + strconv.Itoa(prizeItem.ItemID))
			copy(itemIdData, itemIdStrData)
			opData = append(opData, itemIdData...)
			//itemNum
			itemNum := prizeItem.ItemNum
			opData = append(opData, byte(itemNum&0xFF))
			h.Resource.Logger.Info("add prize item for "+string(charName)+": ",
				zap.Int64("id", prizeItem.ID),
				zap.Uint16("charLv", charLv),
				zap.String("username", string(username)),
				zap.Int("charguid", charguid),
				zap.Int("itemID", prizeItem.ItemID),
				zap.Int("itemNum", prizeItem.ItemNum),
			)
		}
		//标记为已使用
		if err := models.MarkGetAccountPrize(h.Resource.Db, itemIdList); err != nil {
			h.Resource.Logger.Error("MarkGetAccountPrize failed: " + err.Error())
		}
	}
	response.OpData = opData
	return response
}
