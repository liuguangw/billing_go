package server

import (
	"database/sql"
	"github.com/liuguangw/billing_go/bhandler"
	"github.com/liuguangw/billing_go/config"
	"github.com/liuguangw/billing_go/tools"
	"io"
	"net"
)

func createHandle(serverConfig *config.ServerConfig, db *sql.DB,
	tcpConn *net.TCPConn, listener *net.TCPListener) *BillingDataHandle {
	handle := BillingDataHandle{
		Handlers: map[byte]bhandler.BillingHandler{},
		TcpConn:  tcpConn,
		Config:   serverConfig,
	}
	//添加handler
	handle.AddHandler(
		&bhandler.CloseHandler{
			Listener: listener,
		},
		&bhandler.ConnectHandler{},
		&bhandler.PingHandler{},
		&bhandler.KeepHandler{},
		&bhandler.LoginHandler{
			Db:      db,
			AutoReg: serverConfig.AutoReg},
		&bhandler.RegisterHandler{
			Db: db},
		&bhandler.EnterGameHandler{
			Db: db},
		&bhandler.LogoutHandler{
			Db: db},
		&bhandler.KickHandler{},
		&bhandler.QueryPointHandler{
			Db: db},
		&bhandler.ConvertPointHandler{
			Db:            db,
			ConvertNumber: serverConfig.TransferNumber},
		&bhandler.CostLogHandler{},
	)
	return &handle
}

//处理连接
func handleConnection(handle *BillingDataHandle) {
	//获取连接者的IP
	clientIP := handle.GetClientIp()
	// 判断是否允许此ip进行连接
	if !handle.IsClientIpAllowed(clientIP) {
		_ = handle.TcpConn.Close()
		tools.ShowErrorInfoStr("client ip " + clientIP + " is not allowed !")
		return
	}
	tools.LogMessage("client ip " + clientIP + " connected")
	//保持socket长连接
	err := handle.SetKeepAlive()
	if err != nil {
		tools.ShowErrorInfo("SetKeepAlive Failed", err)
	}
	//读取传入的数据包
	var clientData []byte
	//定义每次读取的缓冲区
	var buff = make([]byte, 1024)
	for {
		//读取到的字节数
		readBytes, err := handle.TcpConn.Read(buff)
		if err != nil {
			if err == io.EOF {
				// 连接意外断开
				tools.ShowErrorInfo("client ip "+clientIP+" disconnected", err)
			} else if !tools.ServerStoped {
				// 读取错误
				tools.ShowErrorInfo("read data error", err)
			}
			//读取异常，结束循环读取操作
			return
		}
		// 当读取到数据数>0时
		if readBytes > 0 {
			//缓存的数据+新读取到的数据
			clientData = append(clientData, buff[:readBytes]...)
			//处理数据
			processSize, err := handle.processData(clientData)
			if err != nil {
				tools.ShowErrorInfoStr(err.Error())
				return
			}
			//移出已经处理过的数据
			clientData = clientData[processSize:]
		}
	} //end for
}
