package server

import (
	"github.com/liuguangw/billing_go/config"
	"github.com/liuguangw/billing_go/database"
	"github.com/liuguangw/billing_go/tools"
	"net"
	"strconv"
)

//启动函数
func RunBilling(serverConfig *config.ServerConfig) {
	//初始化数据库连接
	db, dbVersion, err := database.GetConnection(serverConfig)
	if err != nil {
		tools.ShowErrorInfo("Database Error", err)
		return
	}
	tools.LogMessage("mysql version: " + dbVersion)
	//监听端口
	listenAddress := serverConfig.Ip + ":" + strconv.Itoa(serverConfig.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		tools.ShowErrorInfo("resolve TCPAddr failed", err)
		return
	}
	listener, err := net.ListenTCP("tcp", serverEndpoint)
	if err != nil {
		// handle error
		tools.ShowErrorInfo("failed to listen at "+listenAddress, err)
		return
	}
	// 监听端口成功
	tools.LogMessage("billing server run at " + listenAddress)
	tools.ServerStoped = false
	for {
		//接受connect
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			// handle error
			if !tools.ServerStoped {
				//异常
				tools.ShowErrorInfo("accept client failed", err)
				//丢弃异常连接,等待下个连接的进入
				continue
			} else {
				// 服务端停止
				tools.LogMessage("billing server stoped ok")
				return
			}
		}
		handle := createHandle(serverConfig, db, tcpConn, listener)
		go handleConnection(handle)
	}
}
