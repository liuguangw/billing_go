package server

import (
	"fmt"
	"github.com/liuguangw/billing_go/bhandler"
	"github.com/liuguangw/billing_go/config"
	"github.com/liuguangw/billing_go/tools"
	"net"
)

func StopBilling(c *config.ServerConfig) {
	listenAddress := fmt.Sprintf("%s:%d", c.Ip, c.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		tools.ShowErrorInfo("resolve TCPAddr failed", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, serverEndpoint)
	if err != nil {
		tools.ShowErrorInfo("connect to billing error", err)
		return
	}
	defer conn.Close()
	var sendData bhandler.BillingData
	sendData.MsgID = [2]byte{0, 0}
	sendData.OpType = 0
	tools.LogMessage("stoping billing server ...")
	_, err = conn.Write(sendData.PackData())
	if err != nil {
		tools.ShowErrorInfo("stop billing failed", err)
	}
}
