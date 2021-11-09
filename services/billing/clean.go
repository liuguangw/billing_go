package billing

// clean 释放资源
func (s *Server) clean() {
	//关闭listener
	if s.listener != nil {
		s.listener.Close()
	}
	//关闭数据库连接
	if s.database != nil {
		s.database.Close()
	}
	//同步日志
	if s.logger != nil {
		s.logger.Sync()
	}
	//关闭打开的日志文件
	if s.logFile != nil {
		s.logFile.Close()
	}
}
