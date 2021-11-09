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
	//初始化server
	server, err := billing.NewServer()
	if err != nil {
		log.Fatalln(err)
	}
	//处理./billing stop
	if len(os.Args) > 1 {
		commandStr := os.Args[1]
		if commandStr == "stop" {
			fmt.Println("stoping billing server ...")
			if err := server.Stop(); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("stop command sent successfully")
			return
		}
		if commandStr == "show_users" {
			fmt.Println("show users log ...")
			if err := server.ShowUsers(); err != nil {
				log.Fatalln(err)
			}
			return
		}
	}
	//启动
	server.Run()
}
