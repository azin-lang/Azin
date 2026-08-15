package red

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/text"
)

type SyntaxTrivia struct {
	parent   SyntaxToken
	green    green.Node
	position int
}

func (t SyntaxTrivia) IsZero() bool {
	return green.IsNil(t.green)
}

func (t SyntaxTrivia) Parent() SyntaxToken {
	return t.parent
}

func (t SyntaxTrivia) Kind() syntax.Kind {
	if t.IsZero() {
		return syntax.Unknown
	}

	return t.green.Kind()
}

func (t SyntaxTrivia) Position() int {
	return t.position
}

func (t SyntaxTrivia) FullWidth() int {
	if t.IsZero() {
		return 0
	}

	return int(t.green.FullWidth())
}

func (t SyntaxTrivia) FullSpan() text.Span {
	if t.IsZero() {
		return text.Span{}
	}

	return text.NewSpan(
		t.position,
		int(t.green.FullWidth()),
	)
}

func (t SyntaxTrivia) Span() text.Span {
	return t.FullSpan()
}
