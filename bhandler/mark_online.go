package bhandler

import "github.com/liuguangw/billing_go/common"

// markOnline 标记用户为在线状态
func markOnline(loginUsers, onlineUsers map[string]*common.ClientInfo, macCounters map[string]int,
	username string, clientInfo *common.ClientInfo) {
	//已经标记为登录了
	if _, userOnline := onlineUsers[username]; userOnline {
		return
	}
	//从loginUsers中删除
	if loginInfo, userLogin := loginUsers[username]; userLogin {
		delete(loginUsers, username)
		//补充字段信息
		clientInfo.MacMd5 = loginInfo.MacMd5
		if clientInfo.IP == "" {
			clientInfo.IP = loginInfo.IP
		}
	}
	//写入onlineUsers
	onlineUsers[username] = clientInfo
	//mac计数+1
	if clientInfo.MacMd5 != "" {
		counterValue := 0
		if value, valueExists := macCounters[clientInfo.MacMd5]; valueExists {
			counterValue = value
		}
		counterValue++
		macCounters[clientInfo.MacMd5] = counterValue
	}
}
