package bhandler

import "bytes"

//BillingData 数据包结构
type BillingData struct {
	OpType byte
	MsgID  [2]byte
	OpData []byte
}

const (
	//读取数据包成功
	BillingReadOk byte = 0
	//数据包不完整
	BillingDataNotFull byte = 1
	//数据包格式错误
	BillingDataError byte = 2
)

//PackData 数据打包为byte数组
func (billingData *BillingData) PackData() []byte {
	var result []byte
	maskData := []byte{0xAA, 0x55}
	result = append(result, maskData...)
	lengthP := 3 + len(billingData.OpData)
	var tmpByte byte
	// 高8位
	tmpByte = byte(lengthP >> 8)
	result = append(result, tmpByte)
	// 低8位
	tmpByte = byte(lengthP & 0xFF)
	result = append(result, tmpByte)
	// append data
	result = append(result, billingData.OpType)
	result = append(result, billingData.MsgID[0])
	result = append(result, billingData.MsgID[1])
	if lengthP > 3 {
		result = append(result, billingData.OpData...)
	}
	result = append(result, maskData[1])
	result = append(result, maskData[0])
	return result
}

//为响应包进行预处理
func (billingData *BillingData) PrepareResponse(r *BillingData) {
	billingData.OpType = r.OpType
	billingData.MsgID = r.MsgID
}

// u2   头部标识 0xAA 0x55
// u2   有效负载长度 = 1(类型) + 2(ID) + n(opData)
// u1   类型
// u2   ID
// u(n) opData
// u2   尾部标识 0x55 0xAA

// 第二个返回值表示读取结果状态
// 第三个返回值 表示数据包总长度(仅在读取成功时有意义)
func ReadBillingData(binaryData []byte) (*BillingData, byte, int) {
	var result BillingData
	packLength := 0
	maskData := []byte{0xAA, 0x55}
	binaryDataLength := len(binaryData)
	if binaryDataLength < 9 {
		// 数据包总长度的最小值
		return &result, BillingDataNotFull, packLength
	}
	// 检测标识头部
	if bytes.Compare(binaryData[0:2], maskData) != 0 {
		// 头部数据错误
		return &result, BillingDataError, packLength
	}
	//负载数据长度(u2)
	// 负载数据长度需要减去一字节类型标识、两字节的id
	opDataLength := int(binaryData[2])<<8 + int(binaryData[3]) - 3
	// 数据包的总大小
	packLength = 9 + opDataLength
	if binaryDataLength < packLength {
		// 判断数据包总字节数是否达到
		return &result, BillingDataNotFull, packLength
	}
	//检测标识尾部
	if !(binaryData[packLength-2] == maskData[1] && binaryData[packLength-1] == maskData[0]) {
		// 尾部数据错误
		return &result, BillingDataError, packLength
	}
	// 类型标识(u1)
	result.OpType = binaryData[4]
	// 消息id(u2)
	result.MsgID[0] = binaryData[5]
	result.MsgID[1] = binaryData[6]
	// 负载数据(长度为opDataLength)
	if opDataLength > 0 {
		result.OpData = append(result.OpData, binaryData[7:7+opDataLength]...)
	}
	return &result, BillingReadOk, packLength
}
