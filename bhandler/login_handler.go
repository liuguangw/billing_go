package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/tools"
)

type LoginHandler struct {
	Db      *sql.DB
	AutoReg bool
}

func (*LoginHandler) GetType() byte {
	return 0xA2
}
func (h *LoginHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)
	var opData []byte
	//用户名
	offset := 0
	usernameLength := request.OpData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := request.OpData[offset : offset+tmpLength]
	//密码
	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	password := string(request.OpData[offset : offset+tmpLength])
	//登录IP
	offset += tmpLength
	tmpLength = int(request.OpData[offset])
	offset++
	loginIP := string(request.OpData[offset : offset+tmpLength])
	loginResult, loginErr := models.GetLoginResult(h.Db, string(username), password)
	loginResultTxt := "success"
	if loginErr != nil {
		loginResultTxt = loginErr.Error()
	}
	// 如果未开启自动注册,当用户不存在时会返回密码错误
	if (!h.AutoReg) && (loginResult == 9) {
		loginResult = 3
		loginResultTxt = "user " + string(username) + " password error"
	}
	tools.LogMessage(fmt.Sprintf("user [%v] try to login from %v : %v", string(username), loginIP, loginResultTxt))
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, loginResult)
	response.OpData = opData
	return &response
}
