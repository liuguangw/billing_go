package billing

import (
	"context"
	"github.com/liuguangw/billing_go/bhandler"
	"github.com/liuguangw/billing_go/common"
)

// addHandler 添加handler
func (s *Server) addHandler(handlers ...common.PacketHandler) {
	for _, handler := range handlers {
		s.handlers[handler.GetType()] = handler
	}
}

//loadHandlers 载入handlers
func (s *Server) loadHandlers(cancel context.CancelFunc) {
	s.handlers = make(map[byte]common.PacketHandler)
	s.addHandler(
		&bhandler.CommandHandler{
			Cancel:      cancel,
			Logger:      s.logger,
			LoginUsers:  s.loginUsers,
			OnlineUsers: s.onlineUsers,
			MacCounters: s.macCounters,
		},
		&bhandler.ConnectHandler{},
		&bhandler.PingHandler{
			Logger: s.logger,
		},
		&bhandler.KeepHandler{
			Logger: s.logger,
		},
		&bhandler.LoginHandler{
			Db:               s.database,
			Logger:           s.logger,
			AutoReg:          s.config.AutoReg,
			MaxClientCount:   s.config.MaxClientCount,
			PcMaxClientCount: s.config.PcMaxClientCount,
			LoginUsers:       s.loginUsers,
			OnlineUsers:      s.onlineUsers,
			MacCounters:      s.macCounters,
		},
		&bhandler.RegisterHandler{
			Db:     s.database,
			Logger: s.logger},
		&bhandler.EnterGameHandler{
			Db:          s.database,
			Logger:      s.logger,
			LoginUsers:  s.loginUsers,
			OnlineUsers: s.onlineUsers,
			MacCounters: s.macCounters,
		},
		&bhandler.LogoutHandler{
			Db:          s.database,
			Logger:      s.logger,
			LoginUsers:  s.loginUsers,
			OnlineUsers: s.onlineUsers,
			MacCounters: s.macCounters,
		},
		&bhandler.KickHandler{},
		&bhandler.QueryPointHandler{
			Db:     s.database,
			Logger: s.logger},
		&bhandler.ConvertPointHandler{
			Db:            s.database,
			Logger:        s.logger,
			ConvertNumber: s.config.TransferNumber},
		&bhandler.CostLogHandler{},
	)
}
