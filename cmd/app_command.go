package cmd

import (
	"github.com/urfave/cli/v2"
)

var runFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "log-path",
		Usage: "billing log file path",
	},
}

// AppCommand 项目命令根节点
func AppCommand() *cli.App {
	return &cli.App{
		Name:                 "billing",
		Usage:                "billing server command line tool",
		EnableBashCompletion: true,
		Flags:                runFlags,
		Action:               runUpCommand,
		Commands: []*cli.Command{
			UpCommand(),
			StopCommand(),
			ShowUsersCommand(),
			VersionCommand(),
		},
	}
}
