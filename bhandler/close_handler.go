package bhandler

import (
	"context"
	"github.com/liuguangw/billing_go/common"
)

type CloseHandler struct {
	Cancel context.CancelFunc
}

func (*CloseHandler) GetType() byte {
	return 0
}

func (h *CloseHandler) GetResponse(request *common.BillingPacket) *common.BillingPacket {
	response := request.PrepareResponse()
	response.OpData = []byte{0, 0}
	//执行cancel后, Server.Run()中的ctx.Done()会达成,主协程会清理资源(processStop)并退出
	h.Cancel()
	return response
}
