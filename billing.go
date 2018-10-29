package main

import(
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"net"
	"strconv"
)
var serverEndpoint *net.TCPAddr

func main()  {
	// 获取config.json文件绝对路径
	mainAppPath,err:=filepath.Abs(os.Args[0])
	if err !=nil {
		showErrorInfo("get mainAppPath failed",err)
		return
	}
	configFilePath:=filepath.Join(filepath.Dir(mainAppPath),"config.json")
	// 读取配置文件内容
	var jsonBytes []byte
	jsonBytes, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		showErrorInfo("read config file failed",err)
		return
	}
	// debug输出json内容
	//fmt.Printf("content:%s\n",jsonBytes)
	var serverConfig ServerConfig
	err = json.Unmarshal(jsonBytes, &serverConfig)
	if err !=nil {
		showErrorInfo("parse json failed",err)
		return
	}
	//fmt.Println(serverConfig)
	var listenAddress string = serverConfig.Ip+":"+strconv.Itoa(serverConfig.Port)
	serverEndpoint,err = net.ResolveTCPAddr("tcp",listenAddress)
	if err !=nil {
		showErrorInfo("resolve TCPAddr failed",err)
		return
	}
	ln, listenErr := net.ListenTCP("tcp", serverEndpoint)
	if listenErr != nil {
		// handle error
		showErrorInfo("failed to listen at "+listenAddress,err)
		return
	}
	fmt.Println("billing server run at "+listenAddress)
	for {
		conn, acceptErr := ln.AcceptTCP()
		if acceptErr != nil {
			// handle error
			showErrorInfo("accept client failed",err)
			continue
		}
		go handleConnection(&serverConfig,conn)
	}

}