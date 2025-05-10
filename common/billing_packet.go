package common

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	// ErrorPacketNotFull 数据包不完整
	ErrorPacketNotFull = errors.New("incomplete data packet structure")
	//ErrorPacketInvalid 数据包格式错误
	ErrorPacketInvalid = errors.New("invalid data packet")
)

// billing包头部标识
var BillingPacketHead = [2]byte{0xAA, 0x55}

// packetMinSize 数据包最短长度
const packetMinSize = 9

// BillingPacket 数据包结构
//
// u2   头部标识 0xAA 0x55
//
// u2   有效负载长度 = 1(类型) + 2(ID) + n(opData)
//
// u1   类型
//
// u2   ID
//
// u(n) opData
//
// u2   尾部标识 0x55 0xAA
type BillingPacket struct {
	OpType byte    //类型
	MsgID  [2]byte //消息ID
	OpData []byte  //负载数据
}

// ReadBillingPacket 读取billing包
func ReadBillingPacket(binaryData []byte) (*BillingPacket, error) {
	binaryDataLength := len(binaryData)
	//数据包长度不足
	if binaryDataLength < packetMinSize {
		return nil, ErrorPacketNotFull
	}
	//检测头部标识
	packetHeader := BillingPacketHead[:]
	if !bytes.Equal(packetHeader, binaryData[:2]) {
		return nil, ErrorPacketInvalid
	}
	offset := 2
	//数据长度(u2)
	dataLength := int(binaryData[offset])<<8 + int(binaryData[offset+1])
	offset += 2
	// OpData数据长度需要减去3: OpType(1字节), MsgID(2字节)
	opDataLength := dataLength - 3
	//计算包的总长度
	packetFullLength := opDataLength + packetMinSize
	//数据包不完整
	if binaryDataLength < packetFullLength {
		return nil, ErrorPacketNotFull
	}
	packet := new(BillingPacket)
	// 类型标识(u1)
	packet.OpType = binaryData[offset]
	offset++
	// 消息id(u2)
	packet.MsgID[0] = binaryData[offset]
	packet.MsgID[1] = binaryData[offset+1]
	offset += 2
	//copy数据
	if opDataLength > 0 {
		packet.OpData = make([]byte, opDataLength)
		copy(packet.OpData, binaryData[offset:offset+opDataLength])
	}
	return packet, nil
}

// PrepareResponse 准备响应包
func (packet *BillingPacket) PrepareResponse() *BillingPacket {
	return &BillingPacket{
		OpType: packet.OpType,
		MsgID:  packet.MsgID,
	}
}

// PackData 数据打包为byte数组
func (packet *BillingPacket) PackData() []byte {
	//分配空间
	OpDataLength := len(packet.OpData)
	binaryData := make([]byte, 0, OpDataLength+packetMinSize)
	//头部标识
	packetHeader := BillingPacketHead[:]
	binaryData = append(binaryData, packetHeader...)
	//写入长度 u2
	lengthP := 3 + OpDataLength
	binaryData = append(binaryData, byte(lengthP>>8), byte(lengthP&0xFF))
	// append type, msgID
	binaryData = append(binaryData, packet.OpType, packet.MsgID[0], packet.MsgID[1])
	//append opData
	if OpDataLength > 0 {
		binaryData = append(binaryData, packet.OpData...)
	}
	//尾部
	binaryData = append(binaryData, packetHeader[1], packetHeader[0])
	return binaryData
}

// FullLength 数据包的总长度
func (packet *BillingPacket) FullLength() int {
	return len(packet.OpData) + packetMinSize
}

func (packet *BillingPacket) String() string {
	return fmt.Sprintf("OpType: %x\n"+
		"MsgID: [%x, %x]\n"+
		"Data: \n%s", packet.OpType, packet.MsgID[0], packet.MsgID[1], hex.Dump(packet.OpData))
}

// InitBillingPacketHead 初始化头部标识
func InitBillingPacketHead(billType int) {
	if billType == BillTypeCommon {
		BillingPacketHead = [2]byte{0xAA, 0x55}
	} else if billType == BillTypeHuaiJiu {
		BillingPacketHead = [2]byte{0x55, 0xAA}
	}
}
