package billing

import (
	"context"
	"github.com/liuguangw/billing_go/bhandler"
	"github.com/liuguangw/billing_go/common"
	"net"
)

//TcpConnection tcp连接处理工具
type TcpConnection struct {
	server   *Server
	tcpConn  *net.TCPConn //tcp连接
	handlers map[byte]common.PacketHandler
}

//newTcpConnection 初始化连接对象
func newTcpConnection(cancel context.CancelFunc, server *Server, tcpConn *net.TCPConn) *TcpConnection {
	h := &TcpConnection{
		server:   server,
		tcpConn:  tcpConn,
		handlers: make(map[byte]common.PacketHandler),
	}
	h.addHandler(
		&bhandler.CloseHandler{
			Cancel: cancel,
		},
		&bhandler.ConnectHandler{},
		&bhandler.PingHandler{},
		&bhandler.KeepHandler{
			Logger: server.Logger,
		},
		&bhandler.LoginHandler{
			Db:      server.Database,
			Logger:  server.Logger,
			AutoReg: server.Config.AutoReg},
		&bhandler.RegisterHandler{
			Db:     server.Database,
			Logger: server.Logger},
		&bhandler.EnterGameHandler{
			Db:     server.Database,
			Logger: server.Logger},
		&bhandler.LogoutHandler{
			Db:     server.Database,
			Logger: server.Logger},
		&bhandler.KickHandler{},
		&bhandler.QueryPointHandler{
			Db:     server.Database,
			Logger: server.Logger},
		&bhandler.ConvertPointHandler{
			Db:            server.Database,
			Logger:        server.Logger,
			ConvertNumber: server.Config.TransferNumber},
		&bhandler.CostLogHandler{},
	)
	return h
}

// AddHandler 添加handler
func (h *TcpConnection) addHandler(handlers ...common.PacketHandler) {
	for _, handler := range handlers {
		h.handlers[handler.GetType()] = handler
	}
}
