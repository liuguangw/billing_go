package billing

import (
	"context"
)

// runAcceptLoop 运行accept tcp循环
func (s *Server) runAcceptLoop(cancel context.CancelFunc) {
	for {
		//接受TCP connect
		tcpConn, err := s.Listener.AcceptTCP()
		if err != nil {
			if s.running {
				s.Logger.Error("accept tcp client failed: " + err.Error())
				//丢弃异常连接,等待下个连接的进入
				continue
			}
			//已收到stop命令
			//s.Logger.Info("accept loop stoped")
			return
		}
		connHandle := newTcpConnection(cancel, s, tcpConn)
		go connHandle.HandleConnection()
	}
}
