package cmd

import (
	"fmt"
	"github.com/liuguangw/billing_go/services/billing"
	"github.com/urfave/cli/v2"
)

// ShowUsersCommand 打印用户列表状态的命令
func ShowUsersCommand() *cli.Command {
	return &cli.Command{
		Name:   "show_users",
		Usage:  "show users list status",
		Action: runShowUsersCommand,
	}
}

// runShowUsersCommand 打印用户列表状态
func runShowUsersCommand(c *cli.Context) error {
	//初始化server
	server, err := billing.NewServer()
	if err != nil {
		return err
	}
	fmt.Println("show users log ...")
	return server.ShowUsers()
}
