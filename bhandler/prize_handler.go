package bhandler

import (
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	prizeTypeCheck    byte = 2 //查询奖励
	prizeTypeCheckRet byte = 3
	prizeTypeFetch    byte = 4 //领取奖励
	prizeTypeFetchRet byte = 5
)

// PrizeHandler 活动奖励查询/领取, 脚本ID: 808062
type PrizeHandler struct {
	Resource *common.HandlerResource
	BillType int //billing类型
}

// GetType 可以处理的消息类型
func (*PrizeHandler) GetType() byte {
	return packetTypePrize
}

// GetResponse 根据请求获得响应
func (h *PrizeHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := services.NewPacketDataReader(request.OpData)
	prizeType := packetReader.ReadByteValue()
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
	//标记在线
	clientInfo := &common.ClientInfo{
		IP:       loginIP,
		CharName: string(charName),
	}
	markOnline(h.Resource.LoginUsers, h.Resource.OnlineUsers, h.Resource.MacCounters, string(username), clientInfo)
	//世界id
	worldId := packetReader.ReadLeUint16()
	//角色id
	charguid := packetReader.ReadLeInt()
	if prizeType == prizeTypeCheck {
		h.processCheck(username, worldId, charguid, response)
	} else if prizeType == prizeTypeFetch {
		//领取类型
		fetchType := packetReader.ReadByteValue()
		h.processFetch(fetchType, charName, username, worldId, charguid, response)
	}
	//debug
	//h.Resource.Logger.Info("dump response\n" + response.String())
	return response
}

// processCheck 处理查询奖励
func (h *PrizeHandler) processCheck(username []byte, worldId uint16, charguid int, response *common.BillingPacket) {
	usernameLength := len(username)
	opData := make([]byte, 0, usernameLength+15)
	opData = append(opData, prizeTypeCheckRet, byte(usernameLength))
	opData = append(opData, username...)
	opData = append(opData, 0)
	//查询具体奖励状态
	usernameStr := string(username)
	var state byte
	for i := range 3 {
		inputWorldId := 0
		inputCharguid := 0
		//限定区服
		if i > 0 {
			inputWorldId = int(worldId)
		}
		//限定角色
		if i > 1 {
			inputCharguid = charguid
		}
		//state初始化
		state = 0
		if tmp, err := models.CheckAccountPrizeState(h.Resource.Db, usernameStr, inputWorldId, inputCharguid); err != nil {
			h.Resource.Logger.Error("CheckAccountPrizeState failed: " + err.Error())
		} else {
			state = tmp
		}
		opData = append(opData, state, 0, 0, 0)
	}
	response.OpData = opData
}

// processFetch 处理领取奖励
func (h *PrizeHandler) processFetch(fetchType byte, charName, username []byte, worldId uint16, charguid int, response *common.BillingPacket) {
	inputWorldId := 0
	inputCharguid := 0
	var padValue byte = 0x3F
	//限定区服
	if fetchType > 0 {
		inputWorldId = int(worldId)
		padValue = 0x2F
	}
	//限定角色
	if fetchType > 1 {
		inputCharguid = charguid
		padValue = 0x1F
	}
	prizeList, err := models.FindAccountPrizeList(h.Resource.Db, string(username), inputWorldId, inputCharguid, 30)
	if err != nil {
		h.Resource.Logger.Error("FindAccountPrizeList failed: " + err.Error())
	}
	//记录条数
	prizeCount := len(prizeList)
	usernameLength := len(username)
	opDataLength := usernameLength + 3
	if prizeCount > 0 {
		opDataLength += (22*prizeCount + 1)
	}
	opData := make([]byte, 0, usernameLength+3)
	opData = append(opData, prizeTypeFetchRet, byte(usernameLength))
	opData = append(opData, username...)
	var fetchResult byte
	if prizeCount == 0 {
		fetchResult = 4
	}
	opData = append(opData, fetchResult)
	if prizeCount > 0 {
		opData = append(opData, byte(prizeCount))
		var itemIdList []int64
		for _, prizeItem := range prizeList {
			//itemid
			itemIdList = append(itemIdList, prizeItem.ID)
			//id
			for i := range 8 {
				itemValue := prizeItem.ID
				movePos := i * 8
				if movePos > 0 {
					itemValue >>= int64(movePos)
				}
				itemValue = itemValue & 0xFF
				opData = append(opData, byte(itemValue))
			}
			//world
			opData = append(opData, byte(inputWorldId&0xFF), byte((inputWorldId>>8)&0xFF))
			opData = append(opData, padValue, 0xB2)
			//charguid
			for i := range 4 {
				itemValue := inputCharguid
				movePos := i * 8
				if movePos > 0 {
					itemValue >>= movePos
				}
				itemValue = itemValue & 0xFF
				opData = append(opData, byte(itemValue))
			}
			//itemid
			for i := range 4 {
				itemValue := prizeItem.ItemID
				movePos := i * 8
				if movePos > 0 {
					itemValue >>= movePos
				}
				itemValue = itemValue & 0xFF
				opData = append(opData, byte(itemValue))
			}
			//itemNum
			itemNum := prizeItem.ItemNum
			opData = append(opData, byte(itemNum&0xFF), byte((itemNum>>8)&0xFF))
			h.Resource.Logger.Info("add prize item for "+string(charName)+": ",
				zap.Int64("id", prizeItem.ID),
				zap.String("username", string(username)),
				zap.Int("world", int(worldId)),
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
}
