package handle

import (
	"github.com/liuguangw/billing_go/common"
)

//writePacketToClient 从packetChan取出需要发送的billing包,发送给client
func (h *ConnHandle) writePacketToClient(packetChan <-chan *common.BillingPacket) {
	for response := range packetChan {
		responseData := response.PackData()
		if _, err := h.tcpConn.Write(responseData); err != nil {
			//写入错误
			if h.server.Running() && !h.isCommandClient {
				clientAddrStr := h.tcpConn.RemoteAddr().String()
				h.logger.Error("write packet to client " + clientAddrStr + " failed: " + err.Error())
			}
			return
		}
	}
}
