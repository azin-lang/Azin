// internal/azin/syntax/green/base.go
package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type NodeFlags uint8

const (
	FlagNone                NodeFlags = 0
	FlagContainsDiagnostics NodeFlags = 1 << iota
	FlagIsMissing
)

type Base struct {
	fullWidth int
	kind      syntax.Kind
	flags     NodeFlags
}

func (b Base) Kind() syntax.Kind { return b.kind }
func (b Base) FullWidth() int    { return b.fullWidth }
func (b Base) Flags() NodeFlags  { return b.flags }

func (b Base) ContainsDiagnostics() bool {
	return b.flags&FlagContainsDiagnostics != 0
}

func (b Base) IsMissing() bool {
	return b.flags&FlagIsMissing != 0
}
