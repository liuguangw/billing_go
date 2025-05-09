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

// loadHandlers 载入handlers
func (s *Server) loadHandlers(cancel context.CancelFunc) {
	resource := &common.HandlerResource{
		Db:          s.database,
		Logger:      s.logger,
		LoginUsers:  make(map[string]*common.ClientInfo),
		OnlineUsers: make(map[string]*common.ClientInfo),
		MacCounters: make(map[string]int),
	}
	s.handlers = make(map[byte]common.PacketHandler)
	s.addHandler(
		&bhandler.CommandHandler{
			Resource: resource,
			Cancel:   cancel,
		},
		&bhandler.ConnectHandler{},
		&bhandler.PingHandler{
			Resource: resource,
		},
		&bhandler.LoginHandler{
			Resource:         resource,
			AutoReg:          s.config.AutoReg,
			BillType:         s.config.BillType,
			MaxClientCount:   s.config.MaxClientCount,
			PcMaxClientCount: s.config.PcMaxClientCount,
		},
		&bhandler.EnterGameHandler{
			Resource: resource,
		},
		&bhandler.LogoutHandler{
			Resource: resource,
		},
		&bhandler.KeepHandler{
			Resource: resource,
		},
		&bhandler.KickHandler{},
		&bhandler.CostLogHandler{
			Resource: resource,
		},
		&bhandler.ConvertPointHandler{
			Resource: resource,
		},
		&bhandler.QueryPointHandler{
			Resource: resource,
			PointFix: s.config.PointFix,
		},
		&bhandler.RegisterHandler{
			Resource: resource,
		},
	)
}
