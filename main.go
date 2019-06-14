package main

import (
	"billing/config"
	"billing/server"
	"billing/tools"
	"os"
	"runtime"
)

func main() {
	//获取配置文件路径
	configFilePath, err := tools.GetConfigFilePath()
	if err != nil {
		//获取配置文件路径失败
		tools.ShowErrorInfo("Get config file path failed", err)
		return
	}
	//加载配置
	appConfig := new(config.ServerConfig)
	cErr := appConfig.LoadFromFile(configFilePath)
	if cErr != nil {
		tools.ShowErrorInfo("Load config Error", cErr)
		if cErr.ErrorType == config.ErrorParseJson {
			// 如果是解析失败则强制停止
			return
		}
	}
	if len(os.Args) > 1 {
		commandStr := os.Args[1]
		if commandStr == "stop" {
			server.StopBilling(appConfig)
			return
		}
	}
	tools.LogMessage("powered by liuguang @github https://github.com/liuguangw")
	tools.LogMessage("build by " + runtime.Version())
	server.RunBilling(appConfig)
}
