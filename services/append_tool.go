package services

import (
	"encoding/binary"
)

// AppendDataUint16 追加uint16值到data中
func AppendDataUint16(data []byte, value uint16) []byte {
	return binary.BigEndian.AppendUint16(data, value)
}

// AppendDataUint32 追加uint32值到data中
func AppendDataUint32(data []byte, value uint32) []byte {
	return binary.BigEndian.AppendUint32(data, value)
}

// AppendDataUint16 追加uint16值到data中(小端)
func AppendDataLeUint16(data []byte, value uint16) []byte {
	return binary.LittleEndian.AppendUint16(data, value)
}

// AppendDataUint32 追加uint32值到data中(小端)
func AppendDataLeUint32(data []byte, value uint32) []byte {
	return binary.LittleEndian.AppendUint32(data, value)
}

// AppendDataUint64 追加uint64值到data中(小端)
func AppendDataLeUint64(data []byte, value uint64) []byte {
	return binary.LittleEndian.AppendUint64(data, value)
}
