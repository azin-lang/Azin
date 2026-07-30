package codegen

import (
	"github.com/azin-lang/Azin/internal/platform"
)

type Target struct {
	OS   platform.OS
	Arch platform.Arch
	ABI  platform.ABI

	PointerSize int

	ObjectFormat platform.Format
}
