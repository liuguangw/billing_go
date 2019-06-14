package server

import (
	"billing/bhandler"
	"billing/config"
	"billing/database"
	"billing/tools"
	"fmt"
	"net"
)

//启动函数
func RunBilling(c *config.ServerConfig) {
	//初始化数据库连接
	db, dbVersion, err := database.GetConnection(c)
	if err != nil {
		tools.ShowErrorInfo("Database Error", err)
		return
	}
	tools.LogMessage("mysql version: " + dbVersion)
	//监听端口
	listenAddress := fmt.Sprintf("%s:%d", c.Ip, c.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		tools.ShowErrorInfo("resolve TCPAddr failed", err)
		return
	}
	ln, err := net.ListenTCP("tcp", serverEndpoint)
	if err != nil {
		// handle error
		tools.ShowErrorInfo("failed to listen at "+listenAddress, err)
		return
	}
	// 监听端口成功
	tools.LogMessage("billing server run at " + listenAddress)
	tools.ServerStoped = false
	//加载handler
	var handlers = []bhandler.BillingHandler{
		&bhandler.CloseHandler{
			Listener: ln,
		},
		&bhandler.ConnectHandler{},
		&bhandler.PingHandler{},
		&bhandler.KeepHandler{},
		&bhandler.LoginHandler{
			Db:      db,
			AutoReg: c.AutoReg},
		&bhandler.RegisterHandler{
			Db: db},
		&bhandler.EnterGameHandler{
			Db: db},
		&bhandler.LogoutHandler{
			Db: db},
		&bhandler.KickHandler{},
		&bhandler.QueryPointHandler{
			Db: db},
	}
	for {
		//接受connect
		conn, err := ln.AcceptTCP()
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
		go handleConnection(c, db, conn, ln, handlers)
	}
}
