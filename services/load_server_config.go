package services

import (
	"encoding/json"
	"errors"
	"github.com/liuguangw/billing_go/common"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// defaultServerConfig 默认配置
func defaultServerConfig() *common.ServerConfig {
	return &common.ServerConfig{
		IP:               "127.0.0.1",
		Port:             12680,
		DbHost:           "localhost",
		DbPort:           3306,
		DbUser:           "root",
		DbPassword:       "root",
		DbName:           "web",
		AutoReg:          true,
		PointFix:         0,
		MaxClientCount:   500,
		PcMaxClientCount: 3,
	}
}

// LoadServerConfig 加载配置
func LoadServerConfig() (*common.ServerConfig, error) {
	//当前程序文件的绝对路径
	mainAppPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return nil, err
	}
	//程序目录
	appDir := filepath.Dir(mainAppPath)
	//可选配置文件路径
	pathList := []string{
		"./config.yaml",
		"./config.json",
		filepath.Join(appDir, "config.yaml"),
		filepath.Join(appDir, "config.json"),
	}
	var (
		configFile     *os.File
		configFilePath string
	)
	for _, fPath := range pathList {
		configFile, err = os.Open(fPath)
		if err != nil {
			//文件不存在
			if os.IsNotExist(err) {
				continue
			}
			//其它错误
			return nil, errors.New("open config file " + fPath + " error: " + err.Error())
		} else {
			//打开成功
			configFilePath = fPath
			break
		}
	}
	if configFile == nil {
		return nil, errors.New("config file not found")
	}
	defer configFile.Close()
	// 读取配置文件
	fileData, err := io.ReadAll(configFile)
	if err != nil {
		return nil, errors.New("read config file " + configFilePath + " failed: " + err.Error())
	}
	//初始化配置
	serverConfig := defaultServerConfig()
	//解析错误
	var parseError error
	if strings.HasSuffix(configFilePath, ".yaml") {
		parseError = loadServerConfigFromYaml(fileData, serverConfig)
	} else if strings.HasSuffix(configFilePath, ".json") {
		parseError = loadServerConfigFromJSON(fileData, serverConfig)
	} else {
		return nil, errors.New("config file suffix error")
	}
	if parseError != nil {
		return nil, errors.New("parse config file " + configFilePath + " failed: " + parseError.Error())
	}
	return serverConfig, nil
}

// loadServerConfigFromJSON 从json文件中加载配置
func loadServerConfigFromJSON(fileData []byte, serverConfig *common.ServerConfig) error {
	return json.Unmarshal(fileData, serverConfig)
}

// loadServerConfigFromYaml 从yaml文件中加载配置
func loadServerConfigFromYaml(fileData []byte, serverConfig *common.ServerConfig) error {
	return yaml.Unmarshal(fileData, serverConfig)
}
