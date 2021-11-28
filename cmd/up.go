package cmd

import (
	"errors"
	"github.com/liuguangw/billing_go/services"
	"github.com/liuguangw/billing_go/services/billing"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

// UpCommand 运行server的命令
func UpCommand() *cli.Command {
	upFlags := append(runFlags, &cli.BoolFlag{
		Name:    "daemon",
		Aliases: []string{"d"},
		Usage:   "daemon mode, run server at background",
	})
	return &cli.Command{
		Name:   "up",
		Usage:  "run the billing server",
		Flags:  upFlags,
		Action: runUpCommand,
	}
}

//runUpCommand 运行billing服务
func runUpCommand(c *cli.Context) error {
	isDaemon := c.IsSet("daemon")
	logPath := c.String("log-path")
	//后台模式
	if isDaemon {
		if runtime.GOOS == "windows" {
			return errors.New("daemon mode is not supported on windows")
		}
		return services.RunBillingAtBackground(os.Args[0], logPath)
	}
	//初始化server
	server, err := billing.NewServer()
	if err != nil {
		return err
	}
	server.Run(logPath)
	return nil
}
