package billing

import (
	"github.com/liuguangw/billing_go/services/handle"
)

// runAcceptLoop 运行accept tcp循环
func (s *Server) runAcceptLoop() {
	for {
		//接受TCP connect
		tcpConn, err := s.listener.AcceptTCP()
		if err != nil {
			if s.running {
				s.logger.Error("accept tcp client failed: " + err.Error())
				//丢弃异常连接,等待下个连接的进入
				continue
			}
			//已收到stop命令
			//s.logger.Info("accept loop stoped")
			return
		}
		connHandle := handle.NewConnHandle(s, s.logger, s.config.AllowIps, s.handlers)
		go connHandle.HandleConnection(tcpConn)
	}
}
