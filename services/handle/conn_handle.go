package handle

import (
	"github.com/liuguangw/billing_go/common"
	"go.uber.org/zap"
)

// ConnHandle tcp连接处理工具
type ConnHandle struct {
	server          ServerInterface
	logger          *zap.Logger
	allowIps        []string
	handlers        map[byte]common.PacketHandler
	isCommandClient bool //是否为命令连接(由stop、show_users发起的连接)
}

// NewConnHandle 初始化tcp连接处理工具
func NewConnHandle(server ServerInterface, logger *zap.Logger,
	allowIps []string, handlers map[byte]common.PacketHandler) *ConnHandle {
	return &ConnHandle{
		server:   server,
		logger:   logger,
		allowIps: allowIps,
		handlers: handlers,
	}
}
