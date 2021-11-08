package bhandler

import (
	"database/sql"
	"github.com/liuguangw/billing_go/common"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type EnterGameHandler struct {
	Db     *sql.DB
	Logger *zap.Logger
}

func (*EnterGameHandler) GetType() byte {
	return 0xA3
}
func (h *EnterGameHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//读取请求信息
	packetReader := common.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByte()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//角色名
	tmpLength = int(packetReader.ReadByte())
	charNameGbkData := packetReader.ReadBytes(tmpLength)
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	charName, err := gbkDecoder.Bytes(charNameGbkData)
	if err != nil {
		h.Logger.Error("decode char name failed: " + err.Error())
		charName = []byte("?")
	}
	//todo 更新在线状态
	//
	h.Logger.Info("user [" + string(username) + "] " + string(charName) + " entered game")
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, 0x1)
	response.OpData = opData
	return response
}
