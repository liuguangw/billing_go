package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"go.uber.org/zap"
)

type RegisterHandler struct {
	Db     *sql.DB
	Logger *zap.Logger
}

func (*RegisterHandler) GetType() byte {
	return 0xF1
}

func (h *RegisterHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	//读取请求信息
	packetReader := common.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByte()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//超级密码
	tmpLength = int(packetReader.ReadByte())
	superPassword := string(packetReader.ReadBytes(tmpLength))
	//密码
	tmpLength = int(packetReader.ReadByte())
	password := string(packetReader.ReadBytes(tmpLength))
	//注册IP
	tmpLength = int(packetReader.ReadByte())
	registerIP := string(packetReader.ReadBytes(tmpLength))
	//email
	tmpLength = int(packetReader.ReadByte())
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
	if err := models.RegisterAccount(h.Db, account); err != nil {
		regResult = 4
		regResultTxt = err.Error()
	}
	h.Logger.Info(fmt.Sprintf("user [%s](%s) try to register from %s : %s", username, email, registerIP, regResultTxt))
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, regResult)
	response.OpData = opData
	return response
}
