package bhandler

import (
	"encoding/hex"
	"fmt"

	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/models"
	"github.com/liuguangw/billing_go/services"
)

// 登录结果定义
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
	Resource         *common.HandlerResource
	AutoReg          bool //自动注册
	BillType         int  //billing类型
	MaxClientCount   int  //最多允许进入的用户数量(0表示无限制)
	PcMaxClientCount int  //每台电脑最多允许进入的用户数量(0表示无限制)
}

// GetType 可以处理的消息类型
func (*LoginHandler) GetType() byte {
	return packetTypeLogin
}

// GetResponse 根据请求获得响应
func (h *LoginHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	packetReader := services.NewPacketDataReader(request.OpData)
	var (
		usernameLength byte   //用户名长度
		username       []byte //用户名
		password       string //密码
		loginIP        string //登录IP
		macMd5         string //mac哈希
	)
	if h.BillType == common.BillTypeCommon {
		//用户名
		usernameLength = packetReader.ReadByteValue()
		tmpLength := int(usernameLength)
		username = packetReader.ReadBytes(tmpLength)
		//密码
		tmpLength = int(packetReader.ReadByteValue())
		password = string(packetReader.ReadBytes(tmpLength))
		//登录IP
		tmpLength = int(packetReader.ReadByteValue())
		loginIP = string(packetReader.ReadBytes(tmpLength))
		//跳过level,密码卡数据
		packetReader.Skip(2 + 6 + 6)
		macMd5 = string(packetReader.ReadBytes(32))
	} else {
		//怀旧版
		packetReader.Skip(4)
		//用户名
		usernameLength = packetReader.ReadByteValue()
		tmpLength := int(usernameLength)
		username = packetReader.ReadBytes(tmpLength)
		//登录IP
		tmpLength = int(packetReader.ReadByteValue())
		loginIP = string(packetReader.ReadBytes(tmpLength))
		packetReader.Skip(46)
		//mac
		tmpLength = int(packetReader.ReadByteValue())
		macMd5 = string(packetReader.ReadBytes(tmpLength))
		//password
		passwordData := packetReader.ReadBytes(16)
		password = hex.EncodeToString(passwordData)
	}
	//初始化
	var (
		loginResult    = loginCodeSuccess
		loginResultTxt = "success"
	)
	if err := models.CheckLogin(h.Resource.Db, h.Resource.OnlineUsers, string(username), password); err != nil {
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
		currentCount := len(h.Resource.OnlineUsers)
		if currentCount >= h.MaxClientCount {
			loginResult = loginCodeOtherError
			loginResultTxt = "reach max_client_count limit"
		}
	}
	//判断此电脑的连接数是否达到限制
	if loginResult == loginCodeSuccess && h.PcMaxClientCount > 0 {
		macCounter := 0
		if value, valueExists := h.Resource.MacCounters[macMd5]; valueExists {
			macCounter = value
		}
		if macCounter >= h.PcMaxClientCount {
			loginResult = loginCodeOtherError
			loginResultTxt = "reach pc_max_client_count limit"
		}
	}
	//将用户信息写入LoginUsers
	if loginResult == loginCodeSuccess || loginResult == loginCodeShowRegister {
		h.Resource.LoginUsers[string(username)] = &common.ClientInfo{
			IP:     loginIP,
			MacMd5: macMd5,
		}
	}
	h.Resource.Logger.Info(fmt.Sprintf("user [%s] try to login from %s(Mac_md5=%s) : %s", username, loginIP, macMd5, loginResultTxt))
	//Packets::BLRetAccCheck
	opData := make([]byte, 0, usernameLength+2)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, loginResult)
	//额外数据
	if loginResult == loginCodeSuccess {
		//mCardPoint: 4U
		//mCardDay: 2U
		//mIsFatigue: 1U
		//mAccTotalOnlineSecond: 4U
		//mIsPhoneBind, 1u
		//mIsIPBind, 1u
		//mIsMiBaoBind, 1u
		//mIsMacBind, 1u
		//mIsRealNameBind, 1u
		//mIsInputNameBind, 1u
		//mIsPhoneMiBaoBind, 1u
		//mIsPilferedAccount, 2u
		extraData := make([]byte, 4+2+1+4+7+2)
		extraData[6] = 'N'
		//fake mAccTotalOnlineSecond: 100s
		//extraData[10] = 100
		for i := 11; i < 18; i++ {
			extraData[i] = 'N'
		}
		opData = append(opData, extraData...)
	}
	response.OpData = opData
	return response
}
