package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
)

//登录结果定义
const (
	// loginCodeSuccess 登录成功
	loginCodeSuccess byte = 0x01
	// loginCodeNoAccount 账号不存在
	loginCodeNoAccount byte = 0x02
	// loginCodeWrongPassword 密码错误
	loginCodeWrongPassword byte = 0x03
	// loginCodeUserOnline 用户在线
	loginCodeUserOnline byte = 0x04
	// loginCodeOtherError 其它错误
	loginCodeOtherError byte = 0x06
	// loginCodeForbit 禁止登录
	loginCodeForbit byte = 0x07
	// loginCodeShowRegister 显示注册窗口
	loginCodeShowRegister byte = 0x09
)

// LoginHandler 登录
type LoginHandler struct {
	Db               *sql.DB
	Logger           *zap.Logger
	AutoReg          bool
	MaxClientCount   int                           //最多允许进入的用户数量(0表示无限制)
	PcMaxClientCount int                           //每台电脑最多允许进入的用户数量(0表示无限制)
	LoginUsers       map[string]*common.ClientInfo //已登录,还未进入游戏的用户
	OnlineUsers      map[string]*common.ClientInfo //已进入游戏的用户
	MacCounters      map[string]int                //已进入游戏的用户的mac地址计数器
}

// GetType 可以处理的消息类型
func (*LoginHandler) GetType() byte {
	return packetTypeLogin
}

// GetResponse 根据请求获得响应
func (h *LoginHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := services.NewPacketDataReader(request.OpData)
	//用户名
	usernameLength := packetReader.ReadByteValue()
	tmpLength := int(usernameLength)
	username := packetReader.ReadBytes(tmpLength)
	//密码
	tmpLength = int(packetReader.ReadByteValue())
	password := string(packetReader.ReadBytes(tmpLength))
	//登录IP
	tmpLength = int(packetReader.ReadByteValue())
	loginIP := string(packetReader.ReadBytes(tmpLength))
	//跳过level,密码卡数据
	packetReader.Skip(2 + 6 + 6)
	macMd5 := string(packetReader.ReadBytes(32))
	//初始化
	var (
		loginResult    = loginCodeSuccess
		loginResultTxt = "success"
	)
	if err := models.CheckLogin(h.Db, h.OnlineUsers, string(username), password); err != nil {
		loginResultTxt = err.Error()
		if err == models.ErrorLoginUserNotFound {
			//用户不存在
			loginResult = loginCodeNoAccount
		} else if err == models.ErrorLoginInvalidPassword {
			//密码错误
			loginResult = loginCodeWrongPassword
		} else if err == models.ErrorLoginAccountLocked {
			//停权
			loginResult = loginCodeForbit
		} else if err == models.ErrorLoginAccountOnline {
			//用户已在线
			loginResult = loginCodeUserOnline
		} else {
			//数据库异常
			loginResult = loginCodeOtherError
		}
	}
	// 如果开启了自动注册
	if loginResult == loginCodeNoAccount && h.AutoReg {
		loginResult = loginCodeShowRegister
		loginResultTxt = "show register dialog"
	}
	//判断连接的用户数是否达到限制
	if loginResult == loginCodeSuccess && h.MaxClientCount > 0 {
		currentCount := len(h.OnlineUsers)
		if currentCount >= h.MaxClientCount {
			loginResult = loginCodeOtherError
			loginResultTxt = "reach max_client_count limit"
		}
	}
	//判断此电脑的连接数是否达到限制
	if loginResult == loginCodeSuccess && h.PcMaxClientCount > 0 {
		macCounter := 0
		if value, valueExists := h.MacCounters[macMd5]; valueExists {
			macCounter = value
		}
		if macCounter >= h.PcMaxClientCount {
			loginResult = loginCodeOtherError
			loginResultTxt = "reach pc_max_client_count limit"
		}
	}
	//将用户信息写入LoginUsers
	if loginResult == loginCodeSuccess || loginResult == loginCodeShowRegister {
		h.LoginUsers[string(username)] = &common.ClientInfo{
			IP:     loginIP,
			MacMd5: macMd5,
		}
	}
	h.Logger.Info(fmt.Sprintf("user [%s] try to login from %s(Mac_md5=%s) : %s", username, loginIP, macMd5, loginResultTxt))
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, loginResult)
	response.OpData = opData
	return response
}
