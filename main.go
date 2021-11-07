package main

import (
	"fmt"
	"github.com/liuguangw/billing_go/config"
	"github.com/liuguangw/billing_go/server"
	"github.com/liuguangw/billing_go/services"
	"github.com/liuguangw/billing_go/tools"
	"log"
	"os"
	"runtime"
)

func main() {
	appConfig, err := services.LoadServerConfig()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", appConfig)
	__main()
}

func __main() {
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
		} else if len(os.Args) > 2 {
			// ./billing up -d
			// 使用上面的命令时, 程序会在后台运行(支持类unix系统, 不支持windows)
			// 在当前shell退出后也能保持在后台运行
			if os.Args[1] == "up" && os.Args[2] == "-d" && runtime.GOOS != "windows" {
				server.RunBillingAtBackground(os.Args[0])
				return
			}
		}
	}
	tools.LogMessage("powered by liuguang @github https://github.com/liuguangw")
	tools.LogMessage("build by " + runtime.Version())
	//启动
	server.RunBilling(serverConfig)
}
