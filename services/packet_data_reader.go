package services

import "golang.org/x/text/encoding/simplifiedchinese"

// PacketDataReader 用于读取[]byte数据的工具
type PacketDataReader struct {
	binaryData []byte
	offset     int
}

// NewPacketDataReader 初始化读取[]byte数据的工具
func NewPacketDataReader(binaryData []byte) *PacketDataReader {
	return &PacketDataReader{binaryData: binaryData}
}

// ReadByteValue 读取一个字节
func (r *PacketDataReader) ReadByteValue() byte {
	r.offset++
	return r.binaryData[r.offset-1]
}

// ReadUint16 读取两个字节,转换为unit16
func (r *PacketDataReader) ReadUint16() uint16 {
	var targetValue = uint16(r.binaryData[r.offset])
	targetValue <<= 8
	r.offset++
	targetValue += uint16(r.binaryData[r.offset])
	r.offset++
	return targetValue
}

// ReadLeUint16 以小端序读取两个字节,转换为unit16
func (r *PacketDataReader) ReadLeUint16() uint16 {
	var targetValue = uint16(r.binaryData[r.offset])
	r.offset++
	targetValue += (uint16(r.binaryData[r.offset]) << 8)
	r.offset++
	return targetValue
}

// ReadInt 读取4个字节,转换为int
func (r *PacketDataReader) ReadInt() int {
	var targetValue int
	for i := range 4 {
		tmpValue := int(r.binaryData[r.offset])
		r.offset++
		if i < 3 {
			tmpValue <<= uint((3 - i) * 8)
		}
		targetValue += tmpValue
	}
	return targetValue
}

// ReadLeInt 以小端序读取一个int值
func (r *PacketDataReader) ReadLeInt() int {
	var targetValue int
	for i := range 4 {
		tmpValue := int(r.binaryData[r.offset])
		r.offset++
		if i > 0 {
			tmpValue <<= uint(i * 8)
		}
		targetValue += tmpValue
	}
	return targetValue
}

// ReadBytes 读取一部分字节
func (r *PacketDataReader) ReadBytes(n int) []byte {
	r.offset += n
	return r.binaryData[r.offset-n : r.offset]
}

// ReadGbkString 读取GBK字符串
func (r *PacketDataReader) ReadGbkString(n int) []byte {
	gbkData := r.ReadBytes(n)
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	strData, err := gbkDecoder.Bytes(gbkData)
	if err != nil {
		//解码失败
		strData = []byte("?")
	}
	return strData
}

// Skip 跳过一部分字节
func (r *PacketDataReader) Skip(n int) {
	r.offset += n
}
