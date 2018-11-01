package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

// 日志文件保存目录
var logFilePath = "log.log"

// 是否已停止
var serverStoped = true

func main() {
	//当前程序的绝对路径
	mainAppPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		showErrorInfo("get mainAppPath failed", err)
		return
	}
	// 获取config.json文件绝对路径
	configFilePath := filepath.Join(filepath.Dir(mainAppPath), "config.json")
	// 获取日志文件绝对路径
	logFilePath = filepath.Join(filepath.Dir(mainAppPath), logFilePath)
	// 读取配置文件内容
	jsonBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		showErrorInfo("read config file failed", err)
		return
	}
	// debug输出json内容
	//fmt.Printf("content:%s\n",jsonBytes)
	var serverConfig ServerConfig
	// 默认开启自动注册
	serverConfig.Auto_reg = true
	err = json.Unmarshal(jsonBytes, &serverConfig)
	if err != nil {
		showErrorInfo("parse json failed", err)
		return
	}
	//fmt.Println(serverConfig)
	if len(os.Args) > 1 {
		commandStr := os.Args[1]
		if commandStr == "stop" {
			stopBilling(&serverConfig)
			return
		}
	}
	logMessage("powered by liuguang @github https://github.com/liuguangw")
	runBilling(&serverConfig)
}

func runBilling(config *ServerConfig) {
	//检测MySQL服务器是否可以连接
	db, err := initMysql(config)
	if err != nil {
		showErrorInfo("MySQL error", err)
		return
	}
	//监听端口
	listenAddress := config.Ip + ":" + strconv.Itoa(config.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		showErrorInfo("resolve TCPAddr failed", err)
		return
	}
	ln, err := net.ListenTCP("tcp", serverEndpoint)
	if err != nil {
		// handle error
		showErrorInfo("failed to listen at "+listenAddress, err)
		return
	}
	logMessage("billing server run at " + listenAddress)
	serverStoped = false
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			// handle error
			if !serverStoped {
				showErrorInfo("accept client failed", err)
				continue
			} else {
				// 服务端停止
				logMessage("billing server stoped ok")
				return
			}
		}
		go handleConnection(config, db, conn, ln)
	}
}

func stopBilling(config *ServerConfig) {

	listenAddress := config.Ip + ":" + strconv.Itoa(config.Port)
	serverEndpoint, err := net.ResolveTCPAddr("tcp", listenAddress)
	if err != nil {
		showErrorInfo("resolve TCPAddr failed", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, serverEndpoint)
	if err != nil {
		showErrorInfo("connect to billing error", err)
		return
	}
	defer conn.Close()
	var sendData BillingData
	sendData.msgID = [2]byte{0, 0}
	sendData.opType = 0
	logMessage("stoping billing server ...")
	_, err = conn.Write(sendData.PackData())
	if err != nil {
		showErrorInfo("stop billing failed", err)
		return
	}
}
