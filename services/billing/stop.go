package billing

import (
	"errors"
	"github.com/liuguangw/billing_go/common"
	"net"
	"strconv"
)

// sendPacketToServer 发送Packet到server
func (s *Server) sendPacketToServer(packet *common.BillingPacket) error {
	//如果监听的是0.0.0.0, 发送命令时应该发送到127.0.0.1
	serverIP := s.config.IP
	if serverIP == "0.0.0.0" {
		serverIP = "127.0.0.1"
	}
	listenAddress := serverIP + ":" + strconv.Itoa(s.config.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		return errors.New("resolve TCPAddr failed: " + err.Error())
	}
	tcpConn, err := net.DialTCP("tcp", nil, serverEndpoint)
	if err != nil {
		return errors.New("connect to billing failed: " + err.Error())
	}
	defer tcpConn.Close()
	if _, err = tcpConn.Write(packet.PackData()); err != nil {
		return errors.New("send packet failed: " + err.Error())
	}
	return nil
}

// Stop 发送停止命令到server
func (s *Server) Stop() error {
	packet := &common.BillingPacket{
		MsgID:  [2]byte{0, 0},
		OpData: []byte("close"),
	}
	return s.sendPacketToServer(packet)
}
