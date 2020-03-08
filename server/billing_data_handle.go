package server

import (
	"errors"
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

func (handle *BillingDataHandle) processData(clientData []byte) (int, error) {
	//已处理的数据大小
	processSize := 0
	for {
		// 尝试解析数据包
		request, resultMask, packLength := bhandler.ReadBillingData(clientData)
		if resultMask == bhandler.BillingDataError {
			//包结构不正确
			return processSize, errors.New("billing data struct error")
		} else if resultMask == bhandler.BillingReadOk {
			//成功读取到一个完整包
			// 从缓冲的clientData中移除此包
			clientData = clientData[packLength:]
			//调用handler处理包
			err := handle.ProcessRequest(request)
			if err != nil {
				return processSize, errors.New("response failed: " + err.Error())
			}
			processSize += packLength
			//->继续尝试解析下一个请求包(line 91)
		} else {
			//数据包不完整，跳出解析数据包循环
			break
		}
	}
	return processSize, nil
}
