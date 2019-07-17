package server

import (
	"database/sql"
	"github.com/liuguangw/billing_go/bhandler"
	"github.com/liuguangw/billing_go/config"
	"github.com/liuguangw/billing_go/tools"
	"io"
	"net"
)

func createHandle(sConfig *config.ServerConfig, db *sql.DB,
	conn *net.TCPConn, listener *net.TCPListener) *BillingDataHandle {
	handle := BillingDataHandle{
		Handlers: map[byte]bhandler.BillingHandler{},
		Conn:     conn,
		Config:   sConfig,
	}
	handle.AddHandler(
		&bhandler.CloseHandler{
			Listener: listener,
		},
		&bhandler.ConnectHandler{},
		&bhandler.PingHandler{},
		&bhandler.KeepHandler{},
		&bhandler.LoginHandler{
			Db:      db,
			AutoReg: sConfig.AutoReg},
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
			ConvertNumber: sConfig.ConvertNumber},
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
		_ = handle.Conn.Close()
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
		//tools.LogMessage("start read")
		readBytes, err := handle.Conn.Read(buff)
		//tools.LogMessage("after read")
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
		// 当读取到数据时,将数据append到clientData后面
		if readBytes > 0 {
			clientData = append(clientData, buff[:readBytes]...)
		}
		//tools.LogMessage("read bytes " + strconv.Itoa(readBytes))
		//binary data
		//fmt.Println(clientData)
		// 循环读取(防止粘包导致数据包沉积)
		for {
			// 尝试解析数据包
			request, resultMask, packLength := bhandler.ReadBillingData(clientData)
			if resultMask == bhandler.BillingDataError {
				//包结构不正确
				tools.ShowErrorInfoStr("billing data struct error")
				return
			} else if resultMask == bhandler.BillingReadOk {
				//成功读取到一个完整包
				// 从缓冲的clientData中移除此包
				clientData = clientData[packLength:]
				//处理包
				err = handle.ProcessRequest(request)
				if err != nil {
					tools.ShowErrorInfo("response failed", err)
					return
				}
			} else {
				//数据包不完整，跳出解析数据包循环
				break
			}
		}
		/*tools.LogMessage("clientData left size:" + strconv.Itoa(len(clientData)))
		if len(clientData) > 0 {
			fmt.Println(clientData)
		}*/
	} //end for
}
