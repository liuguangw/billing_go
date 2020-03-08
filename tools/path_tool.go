package tools

import (
	"os"
	"path/filepath"
)

//获取程序目录下的文件绝对路径
func getItemFullPath(fname string) (string, error) {
	//当前程序文件的绝对路径
	mainAppPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	//目录
	appDir := filepath.Dir(mainAppPath)
	return filepath.Join(appDir, fname), nil
}

func GetConfigFilePath() (string, error) {
	return getItemFullPath("config.json")
}

func getLogFilePath() (string, error) {
	return getItemFullPath("log.log")
}
