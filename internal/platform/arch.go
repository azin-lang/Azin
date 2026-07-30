package platform

type Arch uint8

const (
	ArchUnknown Arch = iota

	X86_64
	AArch64
	RISCV64
)
