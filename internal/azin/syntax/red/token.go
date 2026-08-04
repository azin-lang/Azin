package red

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/syntax/text"
)

// SyntaxToken is a stack-allocated value type (struct, not pointer).
// It acts as a lightweight facade over a green.Token.
type SyntaxToken struct {
	parent   *Node
	green    *green.Token
	position uint32
	index    int
}

func NewSyntaxToken(parent *Node, greenTok *green.Token, position uint32, index int) SyntaxToken {
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

func (t SyntaxToken) Kind() syntax.SyntaxKind {
	if t.IsZero() {
		return syntax.Unknown
	}
	return t.green.Kind()
}

func (t SyntaxToken) Parent() *Node {
	return t.parent
}

func (t SyntaxToken) Position() uint32 {
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

func (t SyntaxToken) FullWidth() uint32 {
	if t.IsZero() {
		return 0
	}
	return t.green.FullWidth()
}

// FullSpan includes whitespace and comments.
func (t SyntaxToken) FullSpan() text.TextSpan {
	if t.IsZero() {
		return text.TextSpan{}
	}
	return text.NewTextSpan(t.position, t.green.FullWidth())
}

// Span excludes leading and trailing trivia.
func (t SyntaxToken) Span() text.TextSpan {
	if t.IsZero() {
		return text.TextSpan{}
	}

	leadingW := green.GetLeadingTriviaWidth(t.green)
	trailingW := green.GetTrailingTriviaWidth(t.green)

	start := t.position + leadingW
	length := t.green.FullWidth() - leadingW - trailingW

	return text.NewTextSpan(start, length)
}
