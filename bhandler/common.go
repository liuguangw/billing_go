package bhandler

// 消息类型定义
const (
	packetTypeCommand      byte = 0
	packetTypeConnect      byte = 0xA0
	packetTypePing         byte = 0xA1
	packetTypeLogin        byte = 0xA2
	packetTypeEnterGame    byte = 0xA3
	packetTypeLogout       byte = 0xA4
	packetTypeKeep         byte = 0xA6
	packetTypeKick         byte = 0xA9
	packetTypePrizeFetch   byte = 0xC3
	packetTypeCostLog      byte = 0xC5
	packetTypePrizeCard    byte = 0xC6
	packetTypeConvertPoint byte = 0xE1
	packetTypeQueryPoint   byte = 0xE2
	packetTypePrize        byte = 0xE4
	packetTypeRegister     byte = 0xF1
)
