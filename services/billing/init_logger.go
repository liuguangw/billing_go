package billing

import (
	"errors"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

// initLogger 初始化日志系统
func (s *Server) initLogger(logFilePath string) error {
	if logFilePath == "" {
		//使用默认日志路径
		defaultPath, err := defaultLogFilePath()
		if err != nil {
			return err
		}
		logFilePath = defaultPath
	}
	fileFlag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	logFile, err := os.OpenFile(logFilePath, fileFlag, 0644)
	if err != nil {
		return errors.New("Open log file " + logFilePath + " failed: " + err.Error())
	}
	var (
		stdoutWriteSyncer = zapcore.AddSync(colorable.NewColorableStdout())
		stderrWriteSyncer = zapcore.AddSync(colorable.NewColorableStderr())
		fileWriteSyncer   = zapcore.Lock(logFile)
	)
	//普通日志的级别
	commonPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	//错误以及更高级别
	errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	//所有级别的日志
	allPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})
	//日志格式设置
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.ConsoleSeparator = " "
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05 -0700]")
	//控制台可以显示颜色
	consoleCfg := cfg
	consoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//
	fileEncoder := zapcore.NewConsoleEncoder(cfg)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdoutWriteSyncer, commonPriority),
		zapcore.NewCore(consoleEncoder, stderrWriteSyncer, errorPriority),
		zapcore.NewCore(fileEncoder, fileWriteSyncer, allPriority),
	)
	s.logFile = logFile
	s.logger = zap.New(core, zap.AddStacktrace(zapcore.WarnLevel))
	return nil
}

// 默认的日志文件路径
func defaultLogFilePath() (string, error) {
	//当前程序文件的绝对路径
	mainAppPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	//程序目录
	appDir := filepath.Dir(mainAppPath)
	//日志文件
	logFilePath := filepath.Join(appDir, "billing.log")
	return logFilePath, nil
}
