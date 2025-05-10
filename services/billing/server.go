package billing

import (
	"database/sql"
	"net"
	"os"

	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
)

// Server billing 服务
type Server struct {
	config   *common.ServerConfig          //配置
	running  bool                          //是否正在运行
	database *sql.DB                       //数据库连接
	listener *net.TCPListener              //tcp listener
	logFile  *os.File                      //已打开的日志文件
	logger   *zap.Logger                   //日志对象
	handlers map[byte]common.PacketHandler //数据包handlers
}

// NewServer 创建一个Server对象
func NewServer() (*Server, error) {
	//加载配置
	serverConfig, err := services.LoadServerConfig()
	if err != nil {
		return nil, err
	}
	//初始化头部标识
	common.InitBillingPacketHead(serverConfig.BillType)
	return &Server{config: serverConfig}, nil
}

// Running 判断是否为运行中,当调用stop停止服务时,此值为false
func (s *Server) Running() bool {
	return s.running
}
