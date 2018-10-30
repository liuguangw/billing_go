package main

import(
	"os"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"net"
	"strconv"
	//"fmt"
)
// 日志文件保存目录
var logFilePath string = "log.log"
func main()  {
	// 获取config.json文件绝对路径
	mainAppPath,err:=filepath.Abs(os.Args[0])
	if err !=nil {
		showErrorInfo("get mainAppPath failed",err)
		return
	}
	configFilePath:=filepath.Join(filepath.Dir(mainAppPath),"config.json")
	logFilePath=filepath.Join(filepath.Dir(mainAppPath),"log.log")
	// 读取配置文件内容
	jsonBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		showErrorInfo("read config file failed",err)
		return
	}
	// debug输出json内容
	//fmt.Printf("content:%s\n",jsonBytes)
	var serverConfig ServerConfig
	// 默认开启自动注册
	serverConfig.Auto_reg = true
	err = json.Unmarshal(jsonBytes, &serverConfig)
	if err !=nil {
		showErrorInfo("parse json failed",err)
		return
	}
	logMessage("powered by liuguang @github https://github.com/liuguangw")
	//fmt.Println(serverConfig)
	runBilling(&serverConfig)
}

func runBilling(config *ServerConfig){
	//检测MySQL服务器是否可以连接
	db,err:=initMysql(config)
	if err !=nil {
		showErrorInfo("MySQL error",err)
		return
	}
	//监听端口
	listenAddress := config.Ip+":"+strconv.Itoa(config.Port)
	serverEndpoint,err := net.ResolveTCPAddr("tcp",listenAddress)
	if err !=nil {
		showErrorInfo("resolve TCPAddr failed",err)
		return
	}
	ln, err := net.ListenTCP("tcp", serverEndpoint)
	if err != nil {
		// handle error
		showErrorInfo("failed to listen at "+listenAddress,err)
		return
	}
	logMessage("billing server run at "+listenAddress)
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			// handle error
			showErrorInfo("accept client failed",err)
			continue
		}
		go handleConnection(config,db,conn)
	}
}