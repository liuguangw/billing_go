package main

import (
	"fmt"
	"github.com/liuguangw/billing_go/services"
	"github.com/liuguangw/billing_go/services/billing"
	"log"
	"os"
	"runtime"
)

func main() {
	//加载配置
	serverConfig, err := services.LoadServerConfig()
	if err != nil {
		log.Fatalln(err)
	}
	// ./billing up -d
	// 使用上面的命令时, 程序会在后台运行(支持类unix系统, 不支持windows)
	if len(os.Args) > 2 {
		if os.Args[1] == "up" && os.Args[2] == "-d" && runtime.GOOS != "windows" {
			if err := services.RunBillingAtBackground(os.Args[0]); err != nil {
				log.Fatalln(err)
			}
			return
		}
	}
	server := billing.NewServer(serverConfig)
	//处理./billing stop
	if len(os.Args) > 1 {
		if os.Args[1] == "stop" {
			fmt.Println("stoping billing server ...")
			if err := server.Stop(); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("stop command sent successfully")
			return
		}
	}
	fmt.Println("powered by liuguang @github https://github.com/liuguangw")
	fmt.Println("build by " + runtime.Version())
	//启动
	server.Run()
}
