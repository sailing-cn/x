package convert

import "encoding/binary"

// Uint16ToBinaryString get the string of a uint16 number in binary format.
func (bin *Binary) Uint16ToBinaryString(i uint16) string {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, i)
	return bin.BytesToBinaryString(bs)
}

func (bin *Binary) BinaryStringToUint16(str string) uint16 {
	bs := bin.BinaryStringToBytes(str)
	return binary.BigEndian.Uint16(bs)
}

// Uint32ToBinaryString get the string of a uint32 number in binary format.
func (bin *Binary) Uint32ToBinaryString(i uint32) string {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, i)
	return bin.BytesToBinaryString(bs)
}

// Uint64ToBinaryString get the string of a uint64 number in binary format.
func (bin *Binary) Uint64ToBinaryString(i uint64) string {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, i)
	return bin.BytesToBinaryString(bs)
}
