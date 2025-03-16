package billing

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

// Run 运行billing
func (s *Server) Run(logFilePath string) error {
	//初始化日志系统
	if err := s.initLogger(logFilePath); err != nil {
		return errors.New("init logger failed: " + err.Error())
	}
	//退出前,执行清理任务
	defer s.clean()
	//log
	s.logger.Info("(###) billing service starting ...")
	s.logger.Info("log file: " + s.logFile.Name())
	//初始化tcp Listener
	if err := s.initListener(); err != nil {
		s.logger.Fatal("init listener failed: " + err.Error())
	}
	//初始化数据库连接
	if err := s.initDatabase(); err != nil {
		s.logger.Fatal("init database connection failed: " + err.Error())
	}
	//标记为已启动
	s.running = true
	s.logger.Info("billing server run at " + s.listener.Addr().String())
	ctx, cancel := context.WithCancel(context.Background())
	//载入handlers
	s.loadHandlers(cancel)
	//循环accept tcp
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
	s.logger.Info("billing server stopped")
	return nil
}
