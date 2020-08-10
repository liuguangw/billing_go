package server

import (
	"github.com/liuguangw/billing_go/tools"
	"os"
	"os/exec"
)

func RunBillingAtBackground(billingPath string) {
	cmd := exec.Command(billingPath)
	// 重定向输出文件
	outFile, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	if err != nil {
		tools.ShowErrorInfo("open output file failed", err)
		return
	}
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		tools.ShowErrorInfo("start billing failed", err)
	}
}
