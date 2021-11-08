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
	packetReader := common.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByte()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//密码
	tmpLength = int(packetReader.ReadByte())
	password := string(packetReader.ReadBytes(tmpLength))
	//登录IP
	tmpLength = int(packetReader.ReadByte())
	loginIP := string(packetReader.ReadBytes(tmpLength))
	//跳过level,密码卡数据
	packetReader.Skip(2 + 6 + 6)
	macMd5 := string(packetReader.ReadBytes(32))
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
	h.Logger.Info(fmt.Sprintf("user [%s] try to login from %s(Mac_md5=%s) : %s", username, loginIP, macMd5, loginResultTxt))
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, loginResult)
	response.OpData = opData
	return response
}
