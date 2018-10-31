package main

import (
	"database/sql"
	"fmt"
	"net"
)

func bProcessRequest(billingData *BillingData, db *sql.DB, conn *net.TCPConn, serverConfig *ServerConfig) error {
	var (
		err error
		// 响应的负载数据
		opData []byte
		// 标记是否处理了本次请求
		requestHandled = true
	)
	switch billingData.opType {
	case 0xA0:
		opData, err = bHandleConnect(billingData)
	case 0xA1:
		opData, err = bHandlePing(billingData)
	case 0xA6:
		opData, err = bHandleKeep(billingData)
	case 0xA2:
		opData, err = bHandleLogin(billingData, db, serverConfig.Auto_reg)
	case 0xF1:
		if serverConfig.Auto_reg {
			opData, err = bHandleRegister(billingData, db)
		} else {
			requestHandled = false
		}
	case 0xA3:
		opData, err = bHandleEnterGame(billingData, db)
	case 0xA4:
		opData, err = bHandleLogout(billingData, db)
	case 0xA9:
		opData, err = bHandleKick(billingData)
	case 0xE2:
		opData, err = bHandleCheckPoint(billingData, db)
	default:
		requestHandled = false
	}
	if requestHandled {
		if err != nil {
			// 处理请求出错
			showErrorInfo("process request failed", err)
		} else {
			// 成功获取响应bytes
			var response BillingData
			response.opType = billingData.opType
			response.msgID = billingData.msgID
			response.opData = opData
			responseData := response.PackData()
			_, err := conn.Write(responseData)
			if err != nil {
				return err
			}
			//logMessage("response ok")
			//fmt.Println(response)
		}
	}
	return nil
}

//0xA0
func bHandleConnect(billingData *BillingData) ([]byte, error) {
	var opData = []byte{0x20, 0x00}
	return opData, nil
}

//0xA1
func bHandlePing(billingData *BillingData) ([]byte, error) {
	// ZoneId: 2u
	// WorldId: 2u
	// PlayerCount: 2u
	//
	var opData = []byte{0x01, 0x00}
	return opData, nil
}

//0xA6
func bHandleKeep(billingData *BillingData) ([]byte, error) {
	// username Length: 1u
	// username: *u
	// player level: 2u
	// other : 8u
	//
	usernameLength := billingData.opData[0]
	username := billingData.opData[1 : 1+usernameLength]
	var opData []byte
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	return opData, nil
}

//0xA2
func bHandleLogin(billingData *BillingData, db *sql.DB, allowAutoReg bool) ([]byte, error) {
	var opData []byte
	// username Length: 1u
	// username: *u
	// password Length: 1u
	// password: *u
	// ip Length: 1u
	// ip: *u
	offset := 0
	usernameLength := billingData.opData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := billingData.opData[offset : offset+tmpLength]

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	password := string(billingData.opData[offset : offset+tmpLength])

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	loginIP := string(billingData.opData[offset : offset+tmpLength])
	loginResult := getLoginResult(db, string(username), password)
	// 如果未开启自动注册,当用户不存在时会返回密码错误
	if (!allowAutoReg) && (loginResult == 9) {
		loginResult = 3
	}
	logMessage(fmt.Sprintf("user [%v] try to login from %v : %v", string(username), loginIP, loginResult))
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, loginResult)
	return opData, nil
}

//0xF1
func bHandleRegister(billingData *BillingData, db *sql.DB) ([]byte, error) {
	var opData []byte
	offset := 0
	usernameLength := billingData.opData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := billingData.opData[offset : offset+tmpLength]

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	superPassword := string(billingData.opData[offset : offset+tmpLength])

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	password := string(billingData.opData[offset : offset+tmpLength])

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	registerIP := string(billingData.opData[offset : offset+tmpLength])

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	email := string(billingData.opData[offset : offset+tmpLength])
	//
	regResult := getRegisterResult(db, string(username), password, superPassword, email)
	logMessage(fmt.Sprintf("user [%v](%v) try to register from %v : %v", string(username), email, registerIP, regResult == 1))
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	opData = append(opData, regResult)
	return opData, nil
}

//0xA3
func bHandleEnterGame(billingData *BillingData, db *sql.DB) ([]byte, error) {
	var opData []byte
	offset := 0
	usernameLength := billingData.opData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := billingData.opData[offset : offset+tmpLength]

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	charName := string(billingData.opData[offset : offset+tmpLength])
	// 更新在线状态
	err := updateOnlineStatus(db, string(username), true)
	if err != nil {
		return opData, err
	}
	logMessage("user [" + string(username) + "] " + charName + " entered game")
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	var pResult byte = 1
	opData = append(opData, pResult)
	return opData, nil
}

//0xA4
func bHandleLogout(billingData *BillingData, db *sql.DB) ([]byte, error) {
	var opData []byte
	offset := 0
	usernameLength := billingData.opData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := billingData.opData[offset : offset+tmpLength]

	// 更新在线状态
	err := updateOnlineStatus(db, string(username), false)
	if err != nil {
		return opData, err
	}
	logMessage("user [" + string(username) + "] logout")
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	var pResult byte = 1
	opData = append(opData, pResult)
	return opData, nil
}

//0xA9
func bHandleKick(billingData *BillingData) ([]byte, error) {
	var opData = []byte{0x01}
	return opData, nil
}

//0xE2
func bHandleCheckPoint(billingData *BillingData, db *sql.DB) ([]byte, error) {
	var opData []byte
	offset := 0
	usernameLength := billingData.opData[offset]
	tmpLength := int(usernameLength)
	offset++
	username := billingData.opData[offset : offset+tmpLength]

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	loginIP := string(billingData.opData[offset : offset+tmpLength])

	offset += tmpLength
	tmpLength = int(billingData.opData[offset])
	offset++
	charName := string(billingData.opData[offset : offset+tmpLength])
	// 更新在线状态
	err := updateOnlineStatus(db, string(username), true)
	if err != nil {
		return opData, err
	}
	account, queryOp := getAccountByUsername(db, string(username))
	var accountPoint int32
	if queryOp == 1 {
		accountPoint = (account.point + 1) * 1000
	}
	logMessage("user [" + string(username) + "] " + charName + " check point at " + loginIP)
	opData = append(opData, usernameLength)
	opData = append(opData, username...)
	var tmpByte byte
	tmpByte = byte(accountPoint >> 24)
	opData = append(opData, tmpByte)
	tmpByte = byte((accountPoint >> 16) & 0xff)
	opData = append(opData, tmpByte)
	tmpByte = byte((accountPoint >> 8) & 0xff)
	opData = append(opData, tmpByte)
	tmpByte = byte(accountPoint & 0xff)
	opData = append(opData, tmpByte)
	return opData, nil
}
