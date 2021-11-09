package bhandler

import (
	"database/sql"
	"fmt"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
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
	return 0xA2
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
	//判断用户是否在线
	if loginResult == 1 {
		if _, userOnline := h.OnlineUsers[string(username)]; userOnline {
			loginResult = 4
			loginResultTxt = models.ErrorLoginAccountOnline.Error()
		}
	}
	//判断连接的用户数是否达到限制
	if loginResult == 1 && h.MaxClientCount > 0 {
		currentCount := len(h.OnlineUsers)
		if currentCount >= h.MaxClientCount {
			loginResult = 6
			loginResultTxt = "reach max_client_count limit"
		}
	}
	//判断此电脑的连接数是否达到限制
	if loginResult == 1 && h.PcMaxClientCount > 0 {
		macCounter := 0
		if value, valueExists := h.MacCounters[macMd5]; valueExists {
			macCounter = value
		}
		if macCounter >= h.PcMaxClientCount {
			loginResult = 6
			loginResultTxt = "reach pc_max_client_count limit"
		}
	}
	//将用户信息写入LoginUsers
	if loginResult == 1 || loginResult == 9 {
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
