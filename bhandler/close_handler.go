package bhandler

import (
	"github.com/liuguangw/billing_go/tools"
	"net"
)

type CloseHandler struct {
	Listener *net.TCPListener
}

func (*CloseHandler) GetType() byte {
	return 0
}
func (h *CloseHandler) GetResponse(request *BillingData) *BillingData {
	var response BillingData
	response.PrepareResponse(request)
	response.OpData = []byte{0, 0}
	// 标记为停止
	tools.ServerStoped = true
	// 关闭服务端监听
	_ = h.Listener.Close()
	return &response
}
