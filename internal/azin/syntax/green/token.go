package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type Token struct {
	Base

	text           string
	leadingTrivia  Node
	trailingTrivia Node
}

// NewToken constructs an immutable green token, calculating its full width
// including any attached trivia.
func NewToken(kind syntax.SyntaxKind, text string, leading Node, trailing Node) *Token {
	width := uint32(len(text))
	if leading != nil {
		width += leading.FullWidth()
	}
	if trailing != nil {
		width += trailing.FullWidth()
	}

	return &Token{
		Base: Base{
			kind:      kind,
			fullWidth: width,
		},
		text:           text,
		leadingTrivia:  leading,
		trailingTrivia: trailing,
	}
}

func (t *Token) Text() string {
	return t.text
}

func (t *Token) LeadingTrivia() Node {
	return t.leadingTrivia
}

func (t *Token) TrailingTrivia() Node {
	return t.trailingTrivia
}

// SlotCount for a Token is always 0 because tokens are leaf nodes in the AST.
func (t *Token) SlotCount() int {
	return 0
}

func (t *Token) Slot(index int) Node {
	return nil
}
