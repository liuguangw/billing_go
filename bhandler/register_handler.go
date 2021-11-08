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
	var opData []byte
	//用户名
	offset := 0
	usernameLength := request.OpData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := request.OpData[offset : offset+tmpLength]
	//超级密码
	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	superPassword := string(request.OpData[offset : offset+tmpLength])
	//密码
	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	password := string(request.OpData[offset : offset+tmpLength])
	//注册IP
	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	registerIP := string(request.OpData[offset : offset+tmpLength])
	//email
	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	email := string(request.OpData[offset : offset+tmpLength])
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
	h.Logger.Info(fmt.Sprintf("user [%v](%v) try to register from %v : %v", string(username), email, registerIP, regResultTxt))
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, regResult)
	response.OpData = opData
	return response
}
