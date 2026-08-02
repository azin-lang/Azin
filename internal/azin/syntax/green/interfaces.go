package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type Node interface {
	Kind() syntax.SyntaxKind

	FullWidth() uint32

	SlotCount() int

	Slot(int) Node
}
