package services

import "syscall"

const detachedProcess uint32 = 0x8

func processAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: detachedProcess,
	}
}
