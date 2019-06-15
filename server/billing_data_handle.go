package server

import (
	"fmt"
	"github.com/liuguangw/billing_go/bhandler"
	"github.com/liuguangw/billing_go/config"
	"github.com/liuguangw/billing_go/tools"
	"net"
	"strings"
	"time"
)

type BillingDataHandle struct {
	Handlers map[byte]bhandler.BillingHandler
	Conn     *net.TCPConn
	Config   *config.ServerConfig
}

// 添加handler
func (h *BillingDataHandle) AddHandler(handlers ...bhandler.BillingHandler) {
	for _, handler := range handlers {
		h.Handlers[handler.GetType()] = handler
	}
}

// 处理request
func (h *BillingDataHandle) ProcessRequest(request *bhandler.BillingData) error {
	var response *bhandler.BillingData = nil
	for opType, handler := range h.Handlers {
		if request.OpType == opType {
			response = handler.GetResponse(request)
			break
		}
	}
	if response != nil {
		//响应
		responseData := response.PackData()
		_, err := h.Conn.Write(responseData)
		if err != nil {
			return err
		}
	} else {
		//无法处理当前请求类型
		tools.ShowErrorInfoStr(fmt.Sprintf("unknown BillingData \n\tOpType: 0x%X \n\tOpData: %v",
			int(request.OpType), request.OpData))
	}
	return nil
}

// 获取连接者的IP
func (h *BillingDataHandle) GetClientIp() string {
	remoteAddr := h.Conn.RemoteAddr().String()
	return remoteAddr[:strings.LastIndex(remoteAddr, ":")]
}

// 判断ip是否允许连接
func (h *BillingDataHandle) IsClientIpAllowed(clientIP string) bool {
	// 当配置的白名单为空时,表示允许所有ip连接
	if len(h.Config.AllowIps) == 0 {
		return true
	}
	// 当数组不为空时,只允许指定的ip连接
	ipAllowed := false
	for _, allowIP := range h.Config.AllowIps {
		if allowIP == clientIP {
			ipAllowed = true
			break
		}
	}
	return ipAllowed
}

// TCP keepalive
func (h *BillingDataHandle) SetKeepAlive() error {
	err := h.Conn.SetKeepAlive(true)
	if err != nil {
		return err
	}
	keepAlivePeriod, err := time.ParseDuration("30s")
	if err == nil {
		return err
	}
	return h.Conn.SetKeepAlivePeriod(keepAlivePeriod)
}
