package main
import(
	"net"
	"fmt"
	"strings"
)

//处理客户端连接
func handleConnection(serverConfig *ServerConfig,conn *net.TCPConn){
	var remoteAddr=conn.RemoteAddr().String()
	var remoteIp = remoteAddr[:strings.LastIndex(remoteAddr,":")]
	// 当数组不为空时,只允许指定的ip连接
	if len(serverConfig.Allow_ips)>0 {
		ipAllowed:=false
		for _,allowIp:=range serverConfig.Allow_ips{
			if allowIp == remoteIp{
				ipAllowed = true
				break
			}
		}
		if !ipAllowed {
			fmt.Println("client ip "+remoteIp+" is not allowed !")
			conn.Close()
			return
		}
	}
	fmt.Printf("client ip %v entered\n",remoteIp)
	for {

	}
}