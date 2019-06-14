package server

import (
	"billing/bhandler"
	"billing/config"
	"billing/tools"
	"database/sql"
	"io"
	"net"
	"strings"
	"time"
)

func IsClientIpAllowed(serverConfig *config.ServerConfig, clientIP string) bool {
	ipAllowed := false
	for _, allowIP := range serverConfig.AllowIps {
		if allowIP == clientIP {
			ipAllowed = true
			break
		}
	}
	return ipAllowed
}

//处理连接
func handleConnection(serverConfig *config.ServerConfig, db *sql.DB, conn *net.TCPConn, listener *net.TCPListener,
	handlers []bhandler.BillingHandler) {
	//获取连接者的IP
	remoteAddr := conn.RemoteAddr().String()
	clientIP := remoteAddr[:strings.LastIndex(remoteAddr, ":")]
	// 当数组不为空时,只允许指定的ip连接
	if len(serverConfig.AllowIps) > 0 {
		if !IsClientIpAllowed(serverConfig, clientIP) {
			//断开连接
			_ = conn.Close()
			tools.ShowErrorInfoStr("client ip " + clientIP + " is not allowed !")
			return
		}
	}
	//保持socket长连接
	_ = conn.SetKeepAlive(true)
	keepAlivePeriod, err := time.ParseDuration("30s")
	if err == nil {
		_ = conn.SetKeepAlivePeriod(keepAlivePeriod)
	}
	tools.LogMessage("client ip " + clientIP + " connected")
	//读取传入的数据包
	var clientData []byte
	//定义每次读取的缓冲区
	var buff = make([]byte, 500)
	for {
		readBytes, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				// 连接意外断开
				tools.ShowErrorInfo("client ip "+clientIP+" disconnected", err)
			} else if !tools.ServerStoped {
				// 读取错误
				tools.ShowErrorInfo("read error", err)
			}
			return
		}
		// 当读取到数据时,将数据append到clientData后面
		if readBytes > 0 {
			clientData = append(clientData, buff[:readBytes]...)
		}
		//binary data
		//fmt.Println(clientData)
		// 尝试读取数据包
		billingData, resultMask, packLength := bhandler.ReadBillingData(clientData)
		if resultMask == bhandler.BillingDataError {
			//包结构不正确
			tools.ShowErrorInfoStr("billing data struct error")
			return
		}
		//成功读取到一个完整包
		if resultMask == bhandler.BillingReadOk {
			// 从缓冲的clientData中移除此包
			clientData = clientData[packLength:]
			//处理包
			err = processBillingData(billingData, db, conn, handlers)
			if err != nil {
				tools.ShowErrorInfo("response failed", err)
				return
			}
		}
	}
}
