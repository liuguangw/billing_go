package config

import (
	"encoding/json"
	"errors"
	"github.com/liuguangw/billing_go/tools"
	"io/ioutil"
)

// 配置对象
type ServerConfig struct {
	Ip               string   `json:"ip"`
	Port             int      `json:"port"`
	DbHost           string   `json:"db_host"`
	DbPort           int      `json:"db_port"`
	DbUser           string   `json:"db_user"`
	DbPassword       string   `json:"db_password"`
	DbName           string   `json:"db_name"`
	AllowOldPassword bool     `json:"allow_old_password"`
	AutoReg          bool     `json:"auto_reg"`
	AllowIps         []string `json:"allow_ips"`
	TransferNumber   int      `json:"transfer_number"`
}

func NewServerConfig() (*ServerConfig, error) {
	//获取配置文件路径
	configFilePath, err := tools.GetConfigFilePath()
	if err != nil {
		return nil, errors.New("Get config file path failed:" + err.Error())
	}
	// 读取配置文件
	fileData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.New("read config file " + configFilePath + " failed: " + err.Error())
	}
	// 初始化字段
	serverConfig := &ServerConfig{
		Ip:             "127.0.0.1",
		Port:           12680,
		DbHost:         "127.0.0.1",
		DbPort:         3306,
		DbUser:         "root",
		DbPassword:     "root",
		DbName:         "web",
		AutoReg:        true,
		AllowIps:       make([]string, 0),
		TransferNumber: 1000,
	}
	// json解析
	err = json.Unmarshal(fileData, serverConfig)
	if err != nil {
		return nil, errors.New("parse config file " + configFilePath + " failed, " + err.Error())
	}
	return serverConfig, nil
}
