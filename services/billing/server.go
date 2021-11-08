package billing

import (
	"database/sql"
	"github.com/liuguangw/billing_go/common"
	"go.uber.org/zap"
	"net"
	"os"
)

// Server billing 服务
type Server struct {
	Config   *common.ServerConfig //配置
	running  bool                 //是否正在运行
	Database *sql.DB              //数据库连接
	Listener *net.TCPListener     //tcp Listener
	logFiles []*os.File           //已打开的日志文件
	Logger   *zap.Logger          //日志对象
}

// NewServer 创建一个Server对象
func NewServer(serverConfig *common.ServerConfig) *Server {
	return &Server{Config: serverConfig}
}

// Running 判断是否为运行中,当调用stop停止服务时,此值为false
func (s *Server) Running() bool {
	return s.running
}
