package billing

import (
	"errors"
	"github.com/liuguangw/billing_go/common"
	"net"
	"strconv"
)

// Stop 发送停止命令到server
func (s *Server) Stop() error {
	listenAddress := s.Config.IP + ":" + strconv.Itoa(s.Config.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		return errors.New("resolve TCPAddr failed: " + err.Error())
	}
	tcpConn, err := net.DialTCP("tcp", nil, serverEndpoint)
	if err != nil {
		return errors.New("connect to billing failed: " + err.Error())
	}
	defer tcpConn.Close()
	sendData := &common.BillingPacket{
		MsgID:  [2]byte{0, 0},
		OpData: []byte("close billing"),
	}
	_, err = tcpConn.Write(sendData.PackData())
	if err != nil {
		return errors.New("stop billing failed: " + err.Error())
	}
	return nil
}
