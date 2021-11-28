package services

import (
	"encoding/json"
	"errors"
	"github.com/liuguangw/billing_go/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

//defaultServerConfig 默认配置
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
	var (
		configFile    *os.File
		filenames     = []string{"config.yaml", "config.json"}
		filenameIndex = -1
	)
	for index, filename := range filenames {
		configFilePath := filepath.Join(appDir, filename)
		configFile, err = os.OpenFile(configFilePath, os.O_RDONLY, 0)
		if err != nil {
			//文件不存在
			if os.IsNotExist(err) {
				continue
			}
			//其它错误
			return nil, errors.New("open config file " + configFilePath + " error: " + err.Error())
		}
		//打开成功,标记index
		filenameIndex = index
		break
	}
	if filenameIndex < 0 {
		return nil, errors.New("config file not found")
	}
	defer configFile.Close()
	// 读取配置文件
	fileData, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, errors.New("read config file failed: " + err.Error())
	}
	//初始化配置
	serverConfig := defaultServerConfig()
	if filenameIndex == 0 {
		if err := loadServerConfigFromYaml(fileData, serverConfig); err != nil {
			return nil, err
		}
		return serverConfig, nil
	}
	//json格式
	if err := loadServerConfigFromJSON(fileData, serverConfig); err != nil {
		return nil, err
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
