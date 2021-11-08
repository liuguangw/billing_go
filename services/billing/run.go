package billing

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Run 运行billing
func (s *Server) Run() {
	//初始化日志系统
	if err := s.initLogger(); err != nil {
		log.Fatalln("init Logger failed: " + err.Error())
	}
	defer func() {
		//同步日志
		s.Logger.Sync()
		//关闭打开的日志文件
		for _, logFile := range s.logFiles {
			logFile.Close()
		}
	}()
	//初始化tcp连接
	if err := s.initListener(); err != nil {
		s.Logger.Fatal("init Listener failed: " + err.Error())
	}
	//初始化数据库连接
	if err := s.initDatabase(); err != nil {
		s.Logger.Fatal("init Database connection failed: " + err.Error())
	}
	//标记为已启动
	s.running = true
	s.Logger.Info("billing server run at " + s.Listener.Addr().String())
	ctx, cancel := context.WithCancel(context.Background())
	go s.runAcceptLoop(cancel)
	//关注signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//signal和stop命令都可以触发关闭
	select {
	case <-c:
		s.Logger.Info("using signal to stop server")
	case <-ctx.Done():
		s.Logger.Info("using command to stop server")
	}
	if err := s.processStop(); err != nil {
		s.Logger.Fatal("stop service failed: " + err.Error())
	}
	s.Logger.Info("billing server stoped")
}
