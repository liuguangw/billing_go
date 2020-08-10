package tools

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

var logFilePath = ""

// 基本的显示消息功能
func showMessageBase(str string, isError bool) {
	if runtime.GOOS == "linux" {
		if isError {
			// linux下用红色显示错误信息
			fmt.Fprintf(os.Stderr, "%c[1;0;31m%s%c[0m\n", 0x1B, str, 0x1B)
		} else {
			// linux下用绿色显示
			fmt.Printf("%c[1;0;32m%s%c[0m\n", 0x1B, str, 0x1B)
		}
		return
	}
	// 其他系统直接输出
	if isError {
		fmt.Fprintln(os.Stderr, str)
		return
	}
	fmt.Println(str)
}

//记录日志到文件
func logToFile(str string) {
	if logFilePath == "" {
		// 获取日志文件路径
		tPath, err := getLogFilePath()
		if err != nil {
			showMessageBase(err.Error(), true)
			return
		}
		logFilePath = tPath
	}
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		showMessageBase(err.Error(), true)
		return
	}
	_, err = f.Write([]byte(str + "\n"))
	if err != nil {
		showMessageBase(err.Error(), true)
		return
	}
}

func formatMessageString(str string, isError bool) string {
	msgTag := "[log]"
	if isError {
		msgTag = "[err]"
	}
	return msgTag + "[" + time.Now().Format("2006-01-02 15:04:05") + "] " + str
}

//显示日志
func LogMessage(str string) {
	msg := formatMessageString(str, false)
	showMessageBase(msg, false)
	logToFile(msg)
}

//用于显示错误消息文本
func ShowErrorInfoStr(str string) {
	msg := formatMessageString(str, true)
	showMessageBase(msg, true)
	logToFile(msg)
}

//显示错误消息
func ShowErrorInfo(tipText string, err error) {
	ShowErrorInfoStr(tipText + "," + err.Error())
}
