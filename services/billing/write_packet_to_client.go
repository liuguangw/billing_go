package billing

import (
	"github.com/liuguangw/billing_go/common"
	"sync"
)

//writePacketToClient 从packetChan取出需要发送的billing包,发送给client
func (h *TcpConnection) writePacketToClient(wg *sync.WaitGroup, packetChan <-chan *common.BillingPacket) {
	defer wg.Done()
	for response := range packetChan {
		responseData := response.PackData()
		if _, err := h.tcpConn.Write(responseData); err != nil {
			if h.server.Running() {
				clientAddrStr := h.tcpConn.RemoteAddr().String()
				h.server.Logger.Error("write packet to client " + clientAddrStr + " failed: " + err.Error())
			}
			return
		}
	}
}
