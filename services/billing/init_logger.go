package billing

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

// initLogger 初始化日志系统
func (s *Server) initLogger() error {
	//当前程序文件的绝对路径
	mainAppPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}
	//程序目录
	appDir := filepath.Dir(mainAppPath)
	var (
		fileFlag     = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		logFilePaths = []string{
			filepath.Join(appDir, "common.log"),
			filepath.Join(appDir, "error.log"),
		}
		logFiles []*os.File
	)
	//打开日志文件
	for _, logFilePath := range logFilePaths {
		file, err := os.OpenFile(logFilePath, fileFlag, 0644)
		if err != nil {
			return errors.New("Open log file " + logFilePath + " failed: " + err.Error())
		}
		logFiles = append(logFiles, file)
	}
	//合并输出
	commonWriteSyncer := zapcore.NewMultiWriteSyncer(zapcore.Lock(os.Stdout), zapcore.Lock(logFiles[0]))
	errorWriteSyncer := zapcore.NewMultiWriteSyncer(zapcore.Lock(os.Stderr), zapcore.Lock(logFiles[1]))
	//普通日志的级别
	commonPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	//错误以及更高级别
	errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	//日志格式设置
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.ConsoleSeparator = " "
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05 -0700]")
	consoleEncoder := zapcore.NewConsoleEncoder(cfg)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, commonWriteSyncer, commonPriority),
		zapcore.NewCore(consoleEncoder, errorWriteSyncer, errorPriority),
	)
	s.logFiles = logFiles
	s.Logger = zap.New(core, zap.AddStacktrace(zapcore.WarnLevel))
	return nil
}
