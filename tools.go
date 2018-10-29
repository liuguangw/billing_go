package main

import(
	"fmt"
)

type ServerConfig struct {
	Ip string
	Port int
	Db_host string
	Db_port int
	Db_user string
	Db_password string
    Db_name string
    Allow_ips []string
}
func showErrorInfo(tipText string,err error){
	fmt.Println("[error]",tipText," : ",err.Error())
}