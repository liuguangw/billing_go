package main

import (
	"github.com/liuguangw/billing_go/config"
	"github.com/liuguangw/billing_go/server"
	"github.com/liuguangw/billing_go/tools"
	"os"
	"runtime"
)

func main() {
	//获取配置
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		tools.ShowErrorInfoStr(err.Error())
		return
	}
	//命令行参数
	if len(os.Args) > 1 {
		commandStr := os.Args[1]
		//停止billing
		if commandStr == "stop" {
			server.StopBilling(serverConfig)
			return
		}
	}
	tools.LogMessage("powered by liuguang @github https://github.com/liuguangw")
	tools.LogMessage("build by " + runtime.Version())
	//启动
	server.RunBilling(serverConfig)
}
