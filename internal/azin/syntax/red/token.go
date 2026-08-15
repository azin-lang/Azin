package red

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/text"
)

// SyntaxToken is a stack-allocated value type (struct, not pointer).
// It acts as a lightweight facade over a green.Token.
type SyntaxToken struct {
	parent   *Node
	green    *green.Token
	position int
	index    int
}

func NewSyntaxToken(parent *Node, greenTok *green.Token, position, index int) SyntaxToken {
	return SyntaxToken{
		parent:   parent,
		green:    greenTok,
		position: position,
		index:    index,
	}
}

func (t SyntaxToken) IsZero() bool {
	return t.green == nil
}

func (t SyntaxToken) Kind() syntax.Kind {
	if t.IsZero() {
		return syntax.Unknown
	}
	return t.green.Kind()
}

func (t SyntaxToken) Parent() *Node {
	return t.parent
}

func (t SyntaxToken) Position() int {
	return t.position
}

func (t SyntaxToken) Index() int {
	return t.index
}

func (t SyntaxToken) Text() string {
	if t.IsZero() {
		return ""
	}
	return t.green.Text()
}

func (t SyntaxToken) FullWidth() int {
	if t.IsZero() {
		return 0
	}
	return t.green.FullWidth()
}

// FullSpan includes whitespace and comments.
func (t SyntaxToken) FullSpan() text.Span {
	if t.IsZero() {
		return text.Span{}
	}
	return text.NewSpan(t.position, t.green.FullWidth())
}

// Span excludes leading and trailing trivia.
func (t SyntaxToken) Span() text.Span {
	if t.IsZero() {
		return text.Span{}
	}

	leadingW := green.GetLeadingTriviaWidth(t.green)
	trailingW := green.GetTrailingTriviaWidth(t.green)

	start := t.position + leadingW
	length := t.green.FullWidth() - leadingW - trailingW

	return text.NewSpan(start, length)
}
