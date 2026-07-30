package codegen

import (
	"github.com/azin-lang/Azin/internal/object"
	"github.com/azin-lang/Azin/pkg/diagnostics"
)

type Context struct {
	Target *Target

	Object *object.Object

	Diagnostics *diagnostics.Engine

	LabelID uint64
}
