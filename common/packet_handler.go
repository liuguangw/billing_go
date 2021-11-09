package common

// PacketHandler billing包handler定义
type PacketHandler interface {
	// GetType 可以处理的消息类型
	GetType() byte
	// GetResponse 根据请求获得响应
	GetResponse(request *BillingPacket) *BillingPacket
}
