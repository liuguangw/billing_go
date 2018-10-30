package main
import(
	"net"
	"io"
	//"fmt"
	"strings"
	"database/sql"
)

//处理客户端连接
func handleConnection(serverConfig *ServerConfig,db *sql.DB,conn *net.TCPConn){
	var remoteAddr=conn.RemoteAddr().String()
	var clientIp = remoteAddr[:strings.LastIndex(remoteAddr,":")]
	// 当数组不为空时,只允许指定的ip连接
	if len(serverConfig.Allow_ips)>0 {
		ipAllowed:=false
		for _,allowIp:=range serverConfig.Allow_ips{
			if allowIp == clientIp{
				ipAllowed = true
				break
			}
		}
		if !ipAllowed {
			showErrorInfoStr("client ip "+clientIp+" is not allowed !")
			conn.Close()
			return
		}
	}
	conn.SetKeepAlive(true)
	logMessage("client ip "+clientIp+" connected")
	// 定义客户端传入的所有数据
	var clientData []byte
	for {
		//定义本次读取的缓冲区大小
		buffLength := 500
		var buff []byte = make([]byte,buffLength,buffLength)
		readBytes,err := conn.Read(buff)
		if err != nil {
			if err == io.EOF{
				// 断开连接
				logMessage("client ip "+clientIp+" disconnected")
			}else{
				// 读取错误
				showErrorInfo("read error",err)
			}
			return
		}
		// 当读取到数据时,将数据append到clientData后面
		if readBytes>0 {
			clientData = append(clientData,buff[:readBytes]...)
		}
		//binary data
		//fmt.Println(clientData)
		billingData,resultMask,packLength := readBillingData(&clientData)
		if resultMask == 2 {
			showErrorInfoStr("billing data struct error")
			return
		}else if resultMask == 0 {
			// 将已经读取到的数据移出
			clientData=clientData[packLength:]
			//logMessage("get billingData ok")
			//fmt.Println(billingData)
			// 处理读取到的请求
			err = bProcessRequest(billingData,db,conn,serverConfig)
			if err!=nil {
				showErrorInfo("process request failed",err)
			}
		}
	}
}