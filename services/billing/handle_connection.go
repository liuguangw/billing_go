package billing

import (
	"github.com/liuguangw/billing_go/common"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(1)
	go h.writePacketToClient(&wg, outputPacketChan)
	//处理数据包
	for packet := range inputPacketChan {
		//记录packet
		if packet.OpType != 0xA1 {
			h.server.Logger.Info("====packet====\n" + packet.String())
		}
		//fmt.Printf("%+v\n", packet)
		if handler, handlerExists := h.handlers[packet.OpType]; handlerExists {
			response := handler.GetResponse(packet)
			outputPacketChan <- response
		} else {
			//无法处理的消息类型
			h.server.Logger.Error("unknown packet: \n" + packet.String())
		}
	}
	close(outputPacketChan)
	wg.Wait()
	h.tcpConn.Close()
}
