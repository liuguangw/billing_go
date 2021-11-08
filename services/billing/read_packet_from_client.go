package billing

import (
	"github.com/liuguangw/billing_go/common"
	"io"
)

//readPacketFromClient 读取billing包,并把它放到packetChan中
func (h *TcpConnection) readPacketFromClient(packetChan chan<- *common.BillingPacket) {
	defer close(packetChan)
	var (
		clientData = make([]byte, 0, 1024) //所有数据
		buff       = make([]byte, 1024)    //每次读取的缓冲区
	)
	for {
		n, err := h.tcpConn.Read(buff)
		if err != nil {
			//读取错误
			if h.server.Running() {
				clientAddrStr := h.tcpConn.RemoteAddr().String()
				if err == io.EOF {
					h.server.Logger.Info("client " + clientAddrStr + " disconnected")
				} else {
					//记录读取错误
					h.server.Logger.Error("read from client " + clientAddrStr + " failed: " + err.Error())
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
				h.server.Logger.Error(readErr.Error())
				return
			}
			//删除已经读取过的数据
			if packTotalSize > 0 {
				clientData = clientData[packTotalSize:]
			}
		}
	}
}

func (h *TcpConnection) readPacket(clientData []byte, packetChan chan<- *common.BillingPacket) (int, error) {
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
		packTotalSize += packet.FullLength()
		packetChan <- packet
	}
	return packTotalSize, nil
}
