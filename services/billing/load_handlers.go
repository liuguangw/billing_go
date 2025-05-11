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
			BillType: s.config.BillType,
		},
		&bhandler.LogoutHandler{
			Resource: resource,
		},
		&bhandler.KeepHandler{
			Resource: resource,
		},
		&bhandler.KickHandler{},
		&bhandler.PrizeFetchHandler{
			Resource: resource,
		},
		&bhandler.CostLogHandler{
			Resource: resource,
		},
		&bhandler.PrizeCardHandler{
			Resource: resource,
			BillType: s.config.BillType,
		},
		&bhandler.ConvertPointHandler{
			Resource: resource,
			BillType: s.config.BillType,
		},
		&bhandler.QueryPointHandler{
			Resource: resource,
			PointFix: s.config.PointFix,
			BillType: s.config.BillType,
		},
		&bhandler.PrizeHandler{
			Resource: resource,
			BillType: s.config.BillType,
		},
		&bhandler.RegisterHandler{
			Resource: resource,
		},
	)
}
