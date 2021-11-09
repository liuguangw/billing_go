package billing

import (
	"context"
	"github.com/liuguangw/billing_go/common"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// Run 运行billing
func (s *Server) Run() {
	//初始化日志系统
	if err := s.initLogger(); err != nil {
		log.Fatalln("init logger failed: " + err.Error())
	}
	defer s.clean()
	//输出build信息
	s.logger.Info("powered by liuguang @github https://github.com/liuguangw")
	s.logger.Info("build by " + runtime.Version())
	//初始化tcp连接
	if err := s.initListener(); err != nil {
		s.logger.Fatal("init listener failed: " + err.Error())
	}
	//初始化数据库连接
	if err := s.initDatabase(); err != nil {
		s.logger.Fatal("init database connection failed: " + err.Error())
	}
	//初始化map
	s.loginUsers = make(map[string]*common.ClientInfo)
	s.onlineUsers = make(map[string]*common.ClientInfo)
	s.macCounters = make(map[string]int)
	//标记为已启动
	s.running = true
	s.logger.Info("billing server run at " + s.listener.Addr().String())
	ctx, cancel := context.WithCancel(context.Background())
	s.loadHandlers(cancel)
	go s.runAcceptLoop()
	//关注signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//signal和stop命令都可以触发关闭
	select {
	case <-c:
		s.logger.Info("stop server using signal...")
	case <-ctx.Done():
		s.logger.Info("stop server using command...")
	}
	//标记为停止
	s.running = false
	s.logger.Info("billing server stoped")
}
