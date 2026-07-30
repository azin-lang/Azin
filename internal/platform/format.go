package platform

type Format uint8

const (
	FormatUnknown Format = iota

	ELF
	COFF
	MachO
)
