package platform

type ABI uint8

const (
	ABIUnknown ABI = iota

	SysV
	MicrosoftX64
	AAPCS64
	RISCV
)
