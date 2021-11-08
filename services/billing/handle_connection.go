package billing

import (
	"github.com/liuguangw/billing_go/common"
)

//HandleConnection 处理TCP连接
func (h *TcpConnection) HandleConnection() {
	clientAddrStr := h.tcpConn.RemoteAddr().String()
	h.server.Logger.Info("client " + clientAddrStr + " connected")
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
			h.server.Logger.Info("====packet====\n" + packet.String())
		}
		//fmt.Printf("%+v\n", packet)
		if handler, handlerExists := h.handlers[packet.OpType]; handlerExists {
			response := handler.GetResponse(packet)
			//把response放到输出channel中
			outputPacketChan <- response
		} else {
			//无法处理的消息类型
			h.server.Logger.Error("unknown packet: \n" + packet.String())
		}
	}
	//来到这一步时,说明inputPacketChan已经关闭(读取出现错误,或者手动关闭了服务)
	//关闭输出通道
	close(outputPacketChan)
}
