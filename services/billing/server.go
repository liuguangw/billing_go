package billing

import (
	"database/sql"
	"github.com/liuguangw/billing_go/common"
	"github.com/liuguangw/billing_go/services"
	"go.uber.org/zap"
	"net"
	"os"
)

// Server billing 服务
type Server struct {
	config      *common.ServerConfig //配置
	running     bool                 //是否正在运行
	database    *sql.DB              //数据库连接
	listener    *net.TCPListener     //tcp listener
	logFile     *os.File             //已打开的日志文件
	logger      *zap.Logger          //日志对象
	handlers    map[byte]common.PacketHandler
	loginUsers  map[string]*common.ClientInfo //已登录,还未进入游戏的用户
	onlineUsers map[string]*common.ClientInfo //已进入游戏的用户
	macCounters map[string]int                //已进入游戏的用户的mac地址计数器
}

// NewServer 创建一个Server对象
func NewServer() (*Server, error) {
	//加载配置
	serverConfig, err := services.LoadServerConfig()
	if err != nil {
		return nil, err
	}
	return &Server{config: serverConfig}, nil
}

// Running 判断是否为运行中,当调用stop停止服务时,此值为false
func (s *Server) Running() bool {
	return s.running
}
