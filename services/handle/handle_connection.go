package handle

import (
	"github.com/liuguangw/billing_go/common"
)

//HandleConnection 处理TCP连接
func (h *ConnHandle) HandleConnection() {
	clientAddr := h.tcpConn.RemoteAddr()
	//判断是否允许此IP连接
	if !h.allowAddr(clientAddr.String()) {
		h.logger.Warn("client " + clientAddr.String() + " is not allowed to connect")
		return
	}
	h.logger.Info("client " + clientAddr.String() + " connected")
	//keepalive
	if err := h.tcpConn.SetKeepAlive(true); err != nil {
		h.logger.Error("SetKeepAlive failed: " + err.Error())
	}
	var (
		inputPacketChan  = make(chan *common.BillingPacket, 50)
		outputPacketChan = make(chan *common.BillingPacket, 50)
	)
	go h.readPacketFromClient(inputPacketChan)
	go h.writePacketToClient(outputPacketChan)
	//处理inputPacketChan中的数据包
	for packet := range inputPacketChan {
		//记录packet
		if packet.OpType != 0xA1 {
			h.logger.Info("====packet====\n" + packet.String())
		}
		//fmt.Printf("%+v\n", packet)
		if handler, handlerExists := h.handlers[packet.OpType]; handlerExists {
			response := handler.GetResponse(packet)
			//把response放到输出channel中
			outputPacketChan <- response
		} else {
			//无法处理的消息类型
			h.logger.Error("unknown packet: \n" + packet.String())
		}
	}
	//来到这一步时,说明inputPacketChan已经关闭(读取出现错误,或者手动关闭了服务)
	//关闭输出通道
	close(outputPacketChan)
}
