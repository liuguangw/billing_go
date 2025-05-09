package handle

import (
	"strings"
)

// allowAddr 判断是否允许此IP连接
func (h *ConnHandle) allowAddr(clientAddr string) bool {
	ipAddr := clientAddr[:strings.LastIndex(clientAddr, ":")]
	// 当配置的白名单为空时,表示允许所有ip连接
	if len(h.allowIps) == 0 {
		return true
	}
	// 当数组不为空时,只允许指定的ip连接
	var ipAllowed bool
	for _, allowIP := range h.allowIps {
		if allowIP == ipAddr {
			ipAllowed = true
			break
		}
	}
	return ipAllowed
}
