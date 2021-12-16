//go:build !(linux || windows)

package services

import "syscall"

func processAttr() *syscall.SysProcAttr {
	return nil
}
