package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type Node interface {
	Kind() syntax.SyntaxKind
	FullWidth() uint32
	Flags() NodeFlags
	SlotCount() int
	Slot(index int) Node
	ContainsDiagnostics() bool
	IsMissing() bool
}
