package red

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/text"
)

type SyntaxToken struct {
	parent   *Node
	green    *green.Token
	position int
	index    int
}

func (t SyntaxToken) IsZero() bool {
	return t.green == nil
}

func (t SyntaxToken) Parent() *Node {
	return t.parent
}

func (t SyntaxToken) Green() *green.Token {
	return t.green
}

func (t SyntaxToken) Kind() syntax.Kind {
	if t.green == nil {
		return syntax.Unknown
	}

	return t.green.Kind()
}

func (t SyntaxToken) Index() int {
	return t.index
}

func (t SyntaxToken) Position() int {
	return t.position
}

func (t SyntaxToken) Text() string {
	if t.green == nil {
		return ""
	}

	return t.green.Text()
}

func (t SyntaxToken) FullWidth() int {
	if t.green == nil {
		return 0
	}

	return int(t.green.FullWidth())
}

func (t SyntaxToken) FullSpan() text.Span {
	if t.green == nil {
		return text.Span{}
	}

	return text.NewSpan(
		t.position,
		int(t.green.FullWidth()),
	)
}

func (t SyntaxToken) Span() text.Span {
	if t.green == nil {
		return text.Span{}
	}

	leading := int(green.GetLeadingTriviaWidth(t.green))
	trailing := int(green.GetTrailingTriviaWidth(t.green))
	fullWidth := int(t.green.FullWidth())

	start := t.position + leading
	width := fullWidth - leading - trailing
	if width < 0 {
		width = 0
	}

	return text.NewSpan(start, width)
}

func (t SyntaxToken) LeadingTrivia() SyntaxTrivia {
	if t.green == nil {
		return SyntaxTrivia{}
	}

	trivia := t.green.LeadingTrivia()
	if green.IsNil(trivia) {
		return SyntaxTrivia{}
	}

	return SyntaxTrivia{
		parent:   t,
		green:    trivia,
		position: t.position,
	}
}

func (t SyntaxToken) TrailingTrivia() SyntaxTrivia {
	if t.green == nil {
		return SyntaxTrivia{}
	}

	trivia := t.green.TrailingTrivia()
	if green.IsNil(trivia) {
		return SyntaxTrivia{}
	}

	leading := int(green.GetLeadingTriviaWidth(t.green))

	return SyntaxTrivia{
		parent:   t,
		green:    trivia,
		position: t.position + leading + len(t.green.Text()),
	}
}
