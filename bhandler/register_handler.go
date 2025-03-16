package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
)

// RegisterHandler 用户注册
type RegisterHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*RegisterHandler) GetType() byte {
	return packetTypeRegister
}

// GetResponse 根据请求获得响应
func (h *RegisterHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//读取请求信息
	packetReader := services.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//超级密码
	tmpLength = int(packetReader.ReadByteValue())
	superPassword := string(packetReader.ReadBytes(tmpLength))
	//密码
	tmpLength = int(packetReader.ReadByteValue())
	password := string(packetReader.ReadBytes(tmpLength))
	//注册IP
	tmpLength = int(packetReader.ReadByteValue())
	registerIP := string(packetReader.ReadBytes(tmpLength))
	//email
	tmpLength = int(packetReader.ReadByteValue())
	email := string(packetReader.ReadBytes(tmpLength))
	//
	account := &models.Account{
		Name:     string(username),
		Password: password,
		Question: sql.NullString{
			String: superPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
	}
	var (
		regResult    byte = 1
		regResultTxt      = "success"
	)
	if err := models.RegisterAccount(h.Resource.Db, account); err != nil {
		regResult = 6
		regResultTxt = err.Error()
	}
	h.Resource.Logger.Info(fmt.Sprintf("user [%s](%s) try to register from %s : %s", username, email, registerIP, regResultTxt))
	//Packets::BLRetRegPassPort
	opData := make([]byte, 0, usernameLength+2)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, regResult)
	response.OpData = opData
	return response
}
