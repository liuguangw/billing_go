package cmd

import (
	"fmt"
	"github.com/liuguangw/billing_go/services/billing"
	"github.com/urfave/cli/v2"
)

// StopCommand 停止server的命令
func StopCommand() *cli.Command {
	return &cli.Command{
		Name:   "stop",
		Usage:  "stop the billing server",
		Action: runStopCommand,
	}
}

// runStopCommand 停止billing服务
func runStopCommand(c *cli.Context) error {
	//初始化server
	server, err := billing.NewServer()
	if err != nil {
		return err
	}
	fmt.Println("stoping billing server ...")
	if err := server.Stop(); err != nil {
		return err
	}
	fmt.Println("stop command sent successfully")
	return nil
}
