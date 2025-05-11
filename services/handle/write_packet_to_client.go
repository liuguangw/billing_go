package handle

import (
	"net"

	"github.com/liuguangw/billing_go/common"
)

// writePacketToClient 从packetChan取出需要发送的billing包,发送给client
func (h *ConnHandle) writePacketToClient(tcpConn *net.TCPConn, packetChan <-chan *common.BillingPacket) {
	for response := range packetChan {
		responseData := response.PackData()
		//debug 测试粘包时是否有问题
		/*if response.OpType == 0xC6 {
			responseData = append(responseData, responseData...)
		}*/
		if _, err := tcpConn.Write(responseData); err != nil {
			//写入错误
			if h.server.Running() && !h.isCommandClient {
				clientAddrStr := tcpConn.RemoteAddr().String()
				h.logger.Error("write packet to client " + clientAddrStr + " failed: " + err.Error())
			}
			return
		}
	}
}
