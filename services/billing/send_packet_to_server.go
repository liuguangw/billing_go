package billing

import (
	"errors"
	"github.com/liuguangw/billing_go/common"
	"net"
	"strconv"
)

// sendPacketToServer 发送Packet到server
func (s *Server) sendPacketToServer(packet *common.BillingPacket) (*common.BillingPacket, error) {
	//如果监听的是0.0.0.0, 发送命令时应该发送到127.0.0.1
	serverIP := s.config.IP
	if serverIP == "0.0.0.0" {
		serverIP = "127.0.0.1"
	}
	listenAddress := serverIP + ":" + strconv.Itoa(s.config.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		return nil, errors.New("resolve TCPAddr failed: " + err.Error())
	}
	tcpConn, err := net.DialTCP("tcp", nil, serverEndpoint)
	if err != nil {
		return nil, errors.New("connect to billing failed: " + err.Error())
	}
	defer tcpConn.Close()
	if _, err = tcpConn.Write(packet.PackData()); err != nil {
		return nil, errors.New("send packet failed: " + err.Error())
	}
	//读取response
	return readPacket(tcpConn)
}

// 读取一个响应包
func readPacket(tcpConn *net.TCPConn) (*common.BillingPacket, error) {
	var (
		clientData = make([]byte, 0, 1024) //所有数据
		buff       = make([]byte, 1024)    //每次读取的缓冲区
	)
	for {
		n, err := tcpConn.Read(buff)
		if err != nil {
			return nil, err
		}
		//读取到数据了
		if n > 0 {
			//append到尾部
			clientData = append(clientData, buff[:n]...)
			//解析数据包
			packet, readErr := common.ReadBillingPacket(clientData)
			if readErr == common.ErrorPacketNotFull {
				//继续读取,直到完整
				continue
			}
			//其它错误
			if readErr != nil {
				return nil, readErr
			}
			return packet, nil
		}
	}
}
