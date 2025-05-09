package handle

import (
	"github.com/liuguangw/billing_go/common"
	"io"
	"net"
)

// readPacketFromClient 读取billing包,并把它放到packetChan中
func (h *ConnHandle) readPacketFromClient(tcpConn *net.TCPConn, packetChan chan<- *common.BillingPacket) {
	defer close(packetChan)
	var (
		clientData = make([]byte, 0, 1024) //所有数据
		buff       = make([]byte, 1024)    //每次读取的缓冲区
	)
	for {
		n, err := tcpConn.Read(buff)
		if err != nil {
			//读取错误
			if h.server.Running() && !h.isCommandClient {
				clientAddrStr := tcpConn.RemoteAddr().String()
				if err == io.EOF {
					h.logger.Info("client " + clientAddrStr + " disconnected")
				} else {
					//记录读取错误
					h.logger.Error("read from client " + clientAddrStr + " failed: " + err.Error())
				}
			}
			return
		}
		//读取到数据了
		if n > 0 {
			//append到尾部
			clientData = append(clientData, buff[:n]...)
			packTotalSize, readErr := h.readPacket(clientData, packetChan)
			if readErr != nil {
				h.logger.Error(readErr.Error())
				return
			}
			//删除已经读取过的数据
			if packTotalSize > 0 {
				clientDataSize := len(clientData)
				if packTotalSize < clientDataSize {
					//把剩余数据copy到前面
					copy(clientData, clientData[packTotalSize:])
				}
				//调整length
				clientData = clientData[:clientDataSize-packTotalSize]
			}
		}
	}
}

// readPacket 循环读取数据包,并把数据包放入channel中,直到没有完整的数据可以读取,或者数据包格式错误
//
// 返回本次读取到的所有数据包的长度之和
func (h *ConnHandle) readPacket(clientData []byte, packetChan chan<- *common.BillingPacket) (int, error) {
	packTotalSize := 0
	for {
		//解析数据包
		packet, readErr := common.ReadBillingPacket(clientData[packTotalSize:])
		if readErr == common.ErrorPacketNotFull {
			break
		}
		//其它错误
		if readErr != nil {
			return 0, readErr
		}
		// +数据包长度
		packTotalSize += packet.FullLength()
		//标记
		if packet.OpType == 0 {
			h.isCommandClient = true
		}
		packetChan <- packet
	}
	return packTotalSize, nil
}
