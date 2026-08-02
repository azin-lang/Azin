package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type Base struct {
	kind syntax.SyntaxKind

	fullWidth uint32
}

func (b Base) Kind() syntax.SyntaxKind {
	return b.kind
}

func (b Base) FullWidth() uint32 {
	return b.fullWidth
}
