package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//ServerConfig 配置文件结构
type ServerConfig struct {
	Ip          string
	Port        int
	Db_host     string
	Db_port     int
	Db_user     string
	Db_password string
	Db_name     string
	Auto_reg    bool
	Allow_ips   []string
}

//BillingData 数据包结构
type BillingData struct {
	opType byte
	msgID  [2]byte
	opData []byte
}

//PackData 数据打包为byte数组
func (billingData *BillingData) PackData() []byte {
	var result []byte
	maskData := []byte{0xAA, 0x55}
	result = append(result, maskData...)
	lengthP := 3 + len(billingData.opData)
	var tmpByte byte
	// 高8位
	tmpByte = byte(lengthP >> 8)
	result = append(result, tmpByte)
	// 低8位
	tmpByte = byte(lengthP & 0xFF)
	result = append(result, tmpByte)
	// append data
	result = append(result, billingData.opType)
	result = append(result, billingData.msgID[0])
	result = append(result, billingData.msgID[1])
	if lengthP > 3 {
		result = append(result, billingData.opData...)
	}
	result = append(result, maskData[1])
	result = append(result, maskData[0])
	return result
}

//记录日志到文件
func logToFile(str string) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	if _, err := f.Write([]byte(str + "\n")); err != nil {
		return
	}
}

//显示日志
func logMessage(str string) {
	str = "[log][" + time.Now().Format("2006-01-02 15:04:05") + "] " + str
	if runtime.GOOS == "linux" {
		// linux下用绿色显示
		fmt.Printf("%c[1;0;32m%s%c[0m\n", 0x1B, str, 0x1B)
	} else {
		fmt.Println(str)
	}
	logToFile(str)
}

//用于显示错误消息文本
func showErrorInfoStr(str string) {
	str = "[error][" + time.Now().Format("2006-01-02 15:04:05") + "] " + str
	if runtime.GOOS == "linux" {
		// linux下用红色显示错误信息
		fmt.Printf("%c[1;0;31m%s%c[0m\n", 0x1B, str, 0x1B)
	} else {
		fmt.Println(str)
	}
	logToFile(str)
}

//显示错误消息
func showErrorInfo(tipText string, err error) {
	showErrorInfoStr(tipText + "," + err.Error())
}

//initMysql MySQL状态检测和字段初始化
func initMysql(config *ServerConfig) (*sql.DB, error) {
	//user:password@tcp(localhost:3306)/dbname?charset=utf8
	db, err := sql.Open("mysql", config.Db_user+":"+config.Db_password+"@tcp("+config.Db_host+":"+strconv.Itoa(config.Db_port)+")/"+config.Db_name+"?charset=utf8")
	if err != nil {
		return db, err
	}
	// 最大100个连接，最多闲置10个
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	// 判断连接状态
	err = db.Ping()
	if err != nil {
		return db, err
	}
	rows, err := db.Query("SELECT VERSION() as v")
	if err != nil {
		return db, err
	}
	var dbVersion string
	rows.Next()
	err = rows.Scan(&dbVersion)
	if err != nil {
		return db, err
	}
	rows.Close()
	logMessage("mysql version: " + dbVersion)
	// 额外字段
	var (
		hasOnlineField = false
		hasLockField   = false

		onlineFieldName = "is_online"
		lockFieldName   = "is_lock"
	)
	rows, err = db.Query("SHOW COLUMNS FROM account")
	for rows.Next() {
		var (
			fieldName string
			tType     string
			tNull     string
			tKey      string
			tDefault  sql.NullString
			tExtra    string
		)
		err = rows.Scan(&fieldName, &tType, &tNull, &tKey, &tDefault, &tExtra)
		if err != nil {
			return db, err
		}
		// 标记已存在的额外字段
		if fieldName == onlineFieldName {
			hasOnlineField = true
		} else if fieldName == lockFieldName {
			hasLockField = true
		}
		//fmt.Printf("fieldName: %v\n",fieldName)
	}
	rows.Close()
	var extraFields []string
	if !hasOnlineField {
		extraFields = append(extraFields, onlineFieldName)
	}
	if !hasLockField {
		extraFields = append(extraFields, lockFieldName)
	}
	if len(extraFields) > 0 {
		// 添加额外字段
		for _, fName := range extraFields {
			stmt, err := db.Prepare("ALTER TABLE account ADD COLUMN " + fName + " smallint(1) UNSIGNED NOT NULL DEFAULT 0")
			if err != nil {
				return db, err
			}
			_, err = stmt.Exec()
			if err != nil {
				return db, err
			}
			stmt.Close()
		}
	}
	return db, nil
}

// 第二个返回值 0表示读取成功 1表示数据不完整 2表示数据格式错误
// 第三个返回值 表示数据包总长度(仅在读取成功时有意义)
func readBillingData(data *[]byte) (*BillingData, byte, int) {
	binaryData := *data
	var result BillingData
	maskData := []byte{0xAA, 0x55}
	binaryDataLength := len(binaryData)
	if binaryDataLength < 9 {
		// 数据包总长度的最小值
		return &result, 1, 0
	}
	// 检测标识头部
	if bytes.Compare(binaryData[0:2], maskData) != 0 {
		// 头部数据错误
		return &result, 2, 0
	}
	//负载数据长度(u2)
	// 负载数据长度需要减去一字节类型标识、两字节的id
	opDataLength := int(binaryData[2])<<8 + int(binaryData[3]) - 3
	// 计算数据包的大小
	packLength := 2 + 5 + opDataLength + 2
	if binaryDataLength < packLength {
		// 判断数据包总字节数是否达到
		return &result, 1, 0
	}
	//检测标识尾部
	if !(binaryData[packLength-2] == maskData[1] && binaryData[packLength-1] == maskData[0]) {
		// 尾部数据错误
		return &result, 2, 0
	}
	// 类型标识(u1)
	result.opType = binaryData[4]
	// 消息id(u2)
	result.msgID[0] = binaryData[5]
	result.msgID[1] = binaryData[6]
	// 负载数据(长度为opDataLength)
	result.opData = append(result.opData, binaryData[7:7+opDataLength]...)
	return &result, 0, packLength
}
