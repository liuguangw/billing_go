package billing

// processStop 处理停止服务
func (s *Server) processStop() error {
	//标记为停止
	s.running = false
	//关闭listener
	if err := s.Listener.Close(); err != nil {
		return err
	}
	//关闭数据库连接
	if err := s.Database.Close(); err != nil {
		return err
	}
	return nil
}
