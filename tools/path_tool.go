package tools

import (
	"os"
	"path/filepath"
)

var appPath = ""

func getItemFullPath(fname string) (fullPath string, err error) {
	if appPath == "" {
		//当前程序的绝对路径
		mainAppPath, tErr := filepath.Abs(os.Args[0])
		if tErr != nil {
			err = tErr
			return
		}
		appPath = filepath.Dir(mainAppPath)
	}
	fullPath = filepath.Join(appPath, fname)
	return
}

func GetConfigFilePath() (string, error) {
	return getItemFullPath("config.json")
}

func getLogFilePath() (string, error) {
	return getItemFullPath("log.log")
}
