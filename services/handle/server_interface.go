package handle

// ServerInterface server状态接口
type ServerInterface interface {
	// Running 判断server是否正在运行中
	Running() bool
}
