package services

import (
	"go.uber.org/zap"
	"runtime"
)

//定义一些变量,编译时会被替换
var (
	appVersion    = "0.0.0"                                    //应用版本
	appBuildTime  = "2021-01-01 00:00:00"                      //编译时间
	gitCommitHash = "0000000000000000000000000000000000000000" //最后一次git提交的hash
)

// ShowVersionInfo 展示版本信息
func ShowVersionInfo(logger *zap.Logger) {
	logger.Info("project url: https://github.com/liuguangw/billing_go")
	logger.Info("version: "+appVersion,
		zap.String("git-hash", gitCommitHash[:7]),
		zap.String("build-time", appBuildTime),
		zap.String("go-version", runtime.Version()),
		zap.String("arch", runtime.GOARCH),
	)
}
