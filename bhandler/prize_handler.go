package bhandler

import (
	"strconv"

	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
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
	var prizeType byte
	if h.BillType == common.BillTypeHuaiJiu {
		prizeType = packetReader.ReadByteValue()
	}
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
	if h.BillType == common.BillTypeHuaiJiu {
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
	} else {
		h.processCheck1(username, response)
	}
	return response
}

// processCheck1 处理查询奖励(经典端)
func (h *PrizeHandler) processCheck1(username []byte, response *common.BillingPacket) {
	var checkResult byte
	prizeList, err := models.FindAccountPrizeList(h.Resource.Db, string(username), 0, 0, 30)
	if err != nil {
		checkResult = 3
		h.Resource.Logger.Error("FindAccountPrizeList failed: " + err.Error())
	}
	//记录条数
	prizeCount := len(prizeList)
	if prizeCount == 0 {
		//没有奖励可以领取
		checkResult = 5
	}
	usernameLength := len(username)
	opDataLength := usernameLength + 2
	if checkResult == 0 {
		opDataLength += (21*prizeCount + 1)
	}
	//Packets::LBLNewCheckPrize
	opData := make([]byte, 0, opDataLength)
	opData = append(opData, byte(usernameLength))
	opData = append(opData, username...)
	opData = append(opData, checkResult)
	if checkResult == 0 {
		opData = append(opData, byte(prizeCount))
		for _, prizeItem := range prizeList {
			//itemid
			itemIdData := make([]byte, 20)
			itemIdStrData := []byte("item " + strconv.Itoa(prizeItem.ItemID))
			copy(itemIdData, itemIdStrData)
			opData = append(opData, itemIdData...)
			//itemNum
			itemNum := prizeItem.ItemNum
			opData = append(opData, byte(itemNum&0xFF))
		}
	}
	response.OpData = opData
}

// processCheck 处理怀旧端查询奖励
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

// processFetch 处理怀旧端领取奖励
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
	opData := make([]byte, 0, opDataLength)
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
			opData = services.AppendDataLeUint64(opData, uint64(prizeItem.ID))
			//world
			opData = services.AppendDataLeUint16(opData, uint16(inputWorldId))
			opData = append(opData, padValue, 0xB2)
			//charguid
			opData = services.AppendDataLeUint32(opData, uint32(inputCharguid))
			//itemid
			opData = services.AppendDataLeUint32(opData, uint32(prizeItem.ItemID))
			//itemNum
			opData = services.AppendDataLeUint16(opData, uint16(prizeItem.ItemNum))
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
