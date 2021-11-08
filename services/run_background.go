package services

import (
	"errors"
	"os"
	"os/exec"
)

func RunBillingAtBackground(billingPath string) error {
	cmd := exec.Command(billingPath)
	// 重定向输出文件
	outFile, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("open output file failed: " + err.Error())
	}
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return errors.New("start billing failed: " + err.Error())
	}
	return nil
}
