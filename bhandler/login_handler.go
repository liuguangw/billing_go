package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"go.uber.org/zap"
)

type LoginHandler struct {
	Db      *sql.DB
	Logger  *zap.Logger
	AutoReg bool
}

func (*LoginHandler) GetType() byte {
	return 0xA2
}
func (h *LoginHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
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
	var (
		loginResult    byte = 1
		loginResultTxt      = "success"
	)
	if err := models.CheckLogin(h.Db, string(username), password); err != nil {
		loginResultTxt = err.Error()
		if err == models.ErrorLoginUserNotFound {
			//用户不存在,展示注册
			loginResult = 9
		} else if err == models.ErrorLoginInvalidPassword {
			//密码错误
			loginResult = 3
		} else if err == models.ErrorLoginAccountLocked {
			//停权
			loginResult = 7
		} else if err == models.ErrorLoginAccountOnline {
			//有角色在线
			loginResult = 4
		} else {
			//数据库异常
			loginResult = 6
		}
	}
	// 如果未开启自动注册,当用户不存在时会返回密码错误
	if (!h.AutoReg) && (loginResult == 9) {
		loginResult = 3
		loginResultTxt = models.ErrorLoginInvalidPassword.Error()
	}
	h.Logger.Info(fmt.Sprintf("user [%v] try to login from %v : %v", string(username), loginIP, loginResultTxt))
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, loginResult)
	response.OpData = opData
	return response
}
