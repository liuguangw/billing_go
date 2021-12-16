package common

import (
	"database/sql"
	"go.uber.org/zap"
)

// HandlerResource handler所需的资源
type HandlerResource struct {
	Db          *sql.DB                //数据库连接
	Logger      *zap.Logger            //日志对象
	LoginUsers  map[string]*ClientInfo //已登录,还未进入游戏的用户
	OnlineUsers map[string]*ClientInfo //已进入游戏的用户
	MacCounters map[string]int         //已进入游戏的用户的mac地址计数器
}
