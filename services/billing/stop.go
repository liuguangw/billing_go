package billing

import (
	"errors"
	"github.com/liuguangw/billing_go/common"
	"net"
	"strconv"
)

// sendPacketToServer 发送Packet到server
func (s *Server) sendPacketToServer(packet *common.BillingPacket) error {
	listenAddress := s.config.IP + ":" + strconv.Itoa(s.config.Port)
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
