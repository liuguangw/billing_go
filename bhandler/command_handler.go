package bhandler

import (
	"bytes"
	"github.com/liuguangw/billing_go/common"
	"strconv"
)

//CommandHandler 处理发送过来的命令
type CommandHandler struct {
	Resource *common.HandlerResource
}

// GetType 可以处理的消息类型
func (*CommandHandler) GetType() byte {
	return packetTypeCommand
}

// GetResponse 根据请求获得响应
func (h *CommandHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	response.OpData = []byte{0, 0}
	//./billing show_users
	//获取billing中用户列表状态
	if bytes.Compare(request.OpData, []byte("show_users")) == 0 {
		h.showUsers(response)
	} else {
		//./billing stop
		//关闭billing服务
		//执行cancel后, Server.Run()中的ctx.Done()会达成,主协程会退出
		h.Resource.Cancel()
	}
	return response
}

//showUsers 将用户列表状态写入response
func (h *CommandHandler) showUsers(response *common.BillingPacket) {
	content := "login users:"
	if len(h.Resource.LoginUsers) == 0 {
		content += " empty"
	} else {
		for username, clientInfo := range h.Resource.LoginUsers {
			content += "\n\t" + username + ": " + clientInfo.String()
		}
	}
	//
	content += "\n\nonline users:"
	if len(h.Resource.OnlineUsers) == 0 {
		content += " empty"
	} else {
		for username, clientInfo := range h.Resource.OnlineUsers {
			content += "\n\t" + username + ": " + clientInfo.String()
		}
	}
	//
	content += "\n\nmac counters:"
	if len(h.Resource.MacCounters) == 0 {
		content += " empty"
	} else {
		for macMd5, counterValue := range h.Resource.MacCounters {
			content += "\n\t" + macMd5 + ": " + strconv.Itoa(counterValue)
		}
	}
	response.OpData = []byte(content)
}
