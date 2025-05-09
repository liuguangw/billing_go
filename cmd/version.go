package cmd

import (
	"github.com/liuguangw/billing_go/services"
	"github.com/mattn/go-colorable"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// VersionCommand 打印版本信息的命令
func VersionCommand() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "show version information",
		Action: runVersionCommand,
	}
}

// runVersionCommand 打印版本信息
func runVersionCommand(c *cli.Context) error {
	//日志格式设置
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.ConsoleSeparator = " "
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05 -0700]")
	//控制台可以显示颜色
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//
	consoleEncoder := zapcore.NewConsoleEncoder(cfg)
	stdoutWriteSyncer := zapcore.AddSync(colorable.NewColorableStdout())
	core := zapcore.NewCore(consoleEncoder, stdoutWriteSyncer, zap.InfoLevel)
	logger := zap.New(core)
	defer logger.Sync()
	services.ShowVersionInfo(logger)
	services.ShowBuilderInfo(logger)
	return nil
}
