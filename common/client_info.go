package common

// ClientInfo 登录的用户信息
type ClientInfo struct {
	IP       string //客户端IP
	MacMd5   string //客户端MAC地址MD5
	CharName string //角色名称
}

func (c *ClientInfo) String() string {
	return "{ip=" + c.IP + ", mac_md5=" + c.MacMd5 + ", char_name=" + c.CharName + "}"
}
