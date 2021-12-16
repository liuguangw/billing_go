package services

import (
	"errors"
	"os"
	"os/exec"
)

// RunBillingAtBackground 在后台运行程序
func RunBillingAtBackground(billingPath, logFilePath string) error {
	cmd := exec.Command(billingPath)
	if logFilePath != "" {
		cmd.Args = append(cmd.Args, "--log-path", logFilePath)
	}
	// 重定向输出文件
	outFile, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("open output file failed: " + err.Error())
	}
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr
	//设置守护进程模式
	cmd.SysProcAttr = processAttr()
	if err := cmd.Start(); err != nil {
		return errors.New("start billing failed: " + err.Error())
	}
	return nil
}
