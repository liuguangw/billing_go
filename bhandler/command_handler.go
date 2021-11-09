package bhandler

import (
	"bytes"
	"context"
	"github.com/liuguangw/billing_go/common"
	"go.uber.org/zap"
	"strconv"
)

//CommandHandler 处理发送过来的命令
type CommandHandler struct {
	Cancel      context.CancelFunc
	Logger      *zap.Logger
	LoginUsers  map[string]*common.ClientInfo //已登录,还未进入游戏的用户
	OnlineUsers map[string]*common.ClientInfo //已进入游戏的用户
	MacCounters map[string]int                //已进入游戏的用户的mac地址计数器
}

// GetType 可以处理的消息类型
func (*CommandHandler) GetType() byte {
	return 0
}

// GetResponse 根据请求获得响应
func (h *CommandHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	response.OpData = []byte{0, 0}
	if bytes.Compare(request.OpData, []byte("show_users")) == 0 {
		h.showUsers()
	} else {
		//执行cancel后, Server.Run()中的ctx.Done()会达成,主协程会退出
		h.Cancel()
	}
	return response
}

//showUsers 日志打印用户列表状态
func (h *CommandHandler) showUsers() {
	if len(h.LoginUsers) == 0 {
		h.Logger.Info("login users: empty")
	} else {
		content := "login users:"
		for username, clientInfo := range h.LoginUsers {
			content += "\n" + username + ": " + clientInfo.String()
		}
		h.Logger.Info(content)
	}
	//
	if len(h.OnlineUsers) == 0 {
		h.Logger.Info("online users: empty")
	} else {
		content := "online users:"
		for username, clientInfo := range h.OnlineUsers {
			content += "\n" + username + ": " + clientInfo.String()
		}
		h.Logger.Info(content)
	}
	//
	if len(h.MacCounters) == 0 {
		h.Logger.Info("mac counters: empty")
	} else {
		content := "mac counters:"
		for macMd5, counterValue := range h.MacCounters {
			content += "\n" + macMd5 + ": " + strconv.Itoa(counterValue)
		}
		h.Logger.Info(content)
	}
}
