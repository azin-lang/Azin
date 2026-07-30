package platform

type Endian uint8

const (
	LittleEndian Endian = iota
	BigEndian
	NativeEndian
)
