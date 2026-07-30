package codegen

import (
	"github.com/azin-lang/Azin/internal/platform"
)

type Target struct {
	OS           platform.OS
	Arch         platform.Arch
	ABI          platform.ABI
	ObjectFormat platform.Format
	Endian       platform.Endian

	PointerSize uint8
	IntSize     uint8
	LongSize    uint8
}
