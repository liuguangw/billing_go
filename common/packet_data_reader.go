package common

// PacketDataReader 用于读取[]byte数据的工具
type PacketDataReader struct {
	binaryData []byte
	offset     int
}

func NewPacketDataReader(binaryData []byte) *PacketDataReader {
	return &PacketDataReader{binaryData: binaryData}
}

// ReadByte 读取一个字节
func (r *PacketDataReader) ReadByte() byte {
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

// ReadInt 读取4个字节,转换为int
func (r *PacketDataReader) ReadInt() int {
	var targetValue int
	for i := 0; i < 4; i++ {
		tmpValue := int(r.binaryData[r.offset])
		r.offset++
		if i < 3 {
			tmpValue <<= uint((3 - i) * 8)
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

// Skip 跳过一部分字节
func (r *PacketDataReader) Skip(n int) {
	r.offset += n
}
