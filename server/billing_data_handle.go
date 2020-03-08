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
	TcpConn  *net.TCPConn
	Config   *config.ServerConfig
}

// 添加handler
func (handle *BillingDataHandle) AddHandler(handlers ...bhandler.BillingHandler) {
	for _, handler := range handlers {
		handle.Handlers[handler.GetType()] = handler
	}
}

// 处理request
func (handle *BillingDataHandle) ProcessRequest(request *bhandler.BillingData) error {
	var response *bhandler.BillingData = nil
	for opType, handler := range handle.Handlers {
		if request.OpType == opType {
			response = handler.GetResponse(request)
			break
		}
	}
	if response != nil {
		//响应
		responseData := response.PackData()
		_, err := handle.TcpConn.Write(responseData)
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
func (handle *BillingDataHandle) GetClientIp() string {
	remoteAddr := handle.TcpConn.RemoteAddr().String()
	return remoteAddr[:strings.LastIndex(remoteAddr, ":")]
}

// 判断ip是否允许连接
func (handle *BillingDataHandle) IsClientIpAllowed(clientIP string) bool {
	// 当配置的白名单为空时,表示允许所有ip连接
	if len(handle.Config.AllowIps) == 0 {
		return true
	}
	// 当数组不为空时,只允许指定的ip连接
	ipAllowed := false
	for _, allowIP := range handle.Config.AllowIps {
		if allowIP == clientIP {
			ipAllowed = true
			break
		}
	}
	return ipAllowed
}

// TCP keepalive
func (handle *BillingDataHandle) SetKeepAlive() error {
	err := handle.TcpConn.SetKeepAlive(true)
	if err != nil {
		return err
	}
	keepAlivePeriod, err := time.ParseDuration("30s")
	if err == nil {
		return err
	}
	return handle.TcpConn.SetKeepAlivePeriod(keepAlivePeriod)
}
