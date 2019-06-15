package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	// 读文件失败
	ErrorReadFile = iota
	// json解析失败
	ErrorParseJson
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
	ConvertNumber    int      `json:"convert_number"`
}

// 配置错误信息
type ServerConfigError struct {
	ErrorType    int
	ErrorMessage string
}

func (e *ServerConfigError) Error() string {
	return e.ErrorMessage
}

func (c *ServerConfig) LoadFromFile(filepath string) *ServerConfigError {
	//初始化各字段的默认值
	c.Ip = "127.0.0.1"
	c.Port = 12680
	c.DbHost = c.Ip
	c.DbPort = 3306
	c.DbUser = "root"
	c.DbPassword = "root"
	c.DbName = "web"
	c.AllowOldPassword = false
	c.AutoReg = true
	c.AllowIps = make([]string, 0)
	c.ConvertNumber = 1000
	// 读取文件
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return &ServerConfigError{
			ErrorReadFile,
			"read config file " + filepath + " failed, " + err.Error(),
		}
	}
	// json解析
	err = json.Unmarshal(data, c)
	if err != nil {
		return &ServerConfigError{
			ErrorParseJson,
			"parse config file " + filepath + " failed, " + err.Error(),
		}
	}
	return nil
}
