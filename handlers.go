package main
import(
	"database/sql"
	"net"
	"fmt"
)

func bProcessRequest(billingData *BillingData,db *sql.DB,conn *net.TCPConn)error{
	var (
		err error
		// 响应的负载数据
		opData []byte
		// 标记是否处理了本次请求
		requestHandled bool = true
	)
	switch billingData.opType {
		case 0xA0:
			opData,err = bHandleConnect(billingData,conn)
		case 0xA1:
			opData,err = bHandlePing(billingData,conn)
		case 0xA6:
			opData,err = bHandleKeep(billingData,conn)
		case 0xA2:
			opData,err = bHandleLogin(billingData,db,conn)
		case 0xF1:
			opData,err = bHandleRegister(billingData,db,conn)
		default :
			requestHandled = false
	}
	if requestHandled {
		if err != nil {
			// 处理请求出错
			showErrorInfo("process request failed",err)
		}else{
			// 成功获取响应bytes
			var response BillingData
			response.opType = billingData.opType
			response.msgId = billingData.msgId
			response.opData = opData
			responseData := response.PackData()
			_,err := conn.Write(responseData)
			if err != nil {
				return err
			}
			logMessage("response ok")
			fmt.Println(response)
		}
	}
	return nil
}
//0xA0
func bHandleConnect(billingData *BillingData,conn *net.TCPConn)([]byte,error){
	var opData=[]byte{0x20,0x00}
	return opData,nil
}

//0xA1
func bHandlePing(billingData *BillingData,conn *net.TCPConn)([]byte,error){
	// ZoneId: 2u
	// WorldId: 2u
	// PlayerCount: 2u 
	//
	var opData=[]byte{0x01,0x00}
	return opData,nil
}

//0xA6
func bHandleKeep(billingData *BillingData,conn *net.TCPConn)([]byte,error){
	// username Length: 1u
	// username: *u
	// player level: 2u
	// other : 8u
	//
	usernameLength := billingData.opData[0]
	username := billingData.opData[1:1+usernameLength]
	var opData []byte
	opData = append(opData,usernameLength)
	opData = append(opData,username...)
	return opData,nil
}

//0xA2
func bHandleLogin(billingData *BillingData,db *sql.DB,conn *net.TCPConn)([]byte,error){
	var opData []byte
	// username Length: 1u
	// username: *u
	// password Length: 1u
	// password: *u
	// ip Length: 1u
	// ip: *u
	offset :=0
	usernameLength:=billingData.opData[offset]
	tmpLength:=int(usernameLength)
	offset++
	username:=billingData.opData[offset:offset+tmpLength]

	offset+=tmpLength
	tmpLength=int(billingData.opData[offset])
	offset++
	password:=string(billingData.opData[offset:offset+tmpLength])
	
	offset+=tmpLength
	tmpLength=int(billingData.opData[offset])
	offset++
	loginIp:=string(billingData.opData[offset:offset+tmpLength])
	var loginResult byte = getLoginResult(db,string(username),password)
	logMessage(fmt.Sprintf("user [%v] try to login from %v : %v",string(username),loginIp,loginResult))
	opData = append(opData,usernameLength)
	opData = append(opData,username...)
	opData = append(opData,loginResult)
	return opData,nil
}

//0xF1
func bHandleRegister(billingData *BillingData,db *sql.DB,conn *net.TCPConn)([]byte,error){
	var opData []byte
	offset :=0
	usernameLength:=billingData.opData[offset]
	tmpLength:=int(usernameLength)
	offset++
	username:=billingData.opData[offset:offset+tmpLength]

	offset+=tmpLength
	tmpLength=int(billingData.opData[offset])
	offset++
	superPassword:=string(billingData.opData[offset:offset+tmpLength])

	offset+=tmpLength
	tmpLength=int(billingData.opData[offset])
	offset++
	password:=string(billingData.opData[offset:offset+tmpLength])

	offset+=tmpLength
	tmpLength=int(billingData.opData[offset])
	offset++
	registerIp:=string(billingData.opData[offset:offset+tmpLength])

	offset+=tmpLength
	tmpLength=int(billingData.opData[offset])
	offset++
	email:=string(billingData.opData[offset:offset+tmpLength])
	// 
	var regResult byte = getRegisterResult(db,string(username),password,superPassword, email)
	logMessage(fmt.Sprintf("user [%v](%v) try to register from %v : %v",string(username),email,registerIp,regResult==1))
	opData = append(opData,usernameLength)
	opData = append(opData,username...)
	opData = append(opData,regResult)
	return opData,nil
}
