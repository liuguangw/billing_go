package server

import (
	"github.com/liuguangw/billing_go/bhandler"
	"github.com/liuguangw/billing_go/config"
	"github.com/liuguangw/billing_go/tools"
	"net"
	"strconv"
)

func StopBilling(serverConfig *config.ServerConfig) {
	listenAddress := serverConfig.Ip + ":" + strconv.Itoa(serverConfig.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		tools.ShowErrorInfo("resolve TCPAddr failed", err)
		return
	}
	tcpConn, err := net.DialTCP("tcp", nil, serverEndpoint)
	if err != nil {
		tools.ShowErrorInfo("connect to billing error", err)
		return
	}
	defer tcpConn.Close()
	sendData := &bhandler.BillingData{
		MsgID:  [2]byte{0,0},
	}
	tools.LogMessage("stoping billing server ...")
	_, err = tcpConn.Write(sendData.PackData())
	if err != nil {
		tools.ShowErrorInfo("stop billing failed", err)
	}
	tools.LogMessage("stop command sent successfully")
}
