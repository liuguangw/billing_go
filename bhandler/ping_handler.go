package bhandler

import (
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
)

// PingHandler ping
type PingHandler struct {
	Resource           *common.HandlerResource
	currentPlayerCount uint16
}

// GetType 可以处理的消息类型
func (*PingHandler) GetType() byte {
	return packetTypePing
}

// GetResponse 根据请求获得响应
func (h *PingHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//读取请求信息
	packetReader := services.NewPacketDataReader(request.OpData)
	zoneID := packetReader.ReadUint16()
	worldID := packetReader.ReadUint16()
	playerCount := packetReader.ReadUint16()
	//当玩家数发生变化时,记录信息
	if h.currentPlayerCount != playerCount {
		h.currentPlayerCount = playerCount
		h.Resource.Logger.Info("server status: ",
			zap.Uint16("zoneID", zoneID),
			zap.Uint16("worldID", worldID),
			zap.Uint16("playerCount", playerCount))
	}
	//Packets::BLRetKeepLive
	response.OpData = []byte{0, 0}
	return response
}
