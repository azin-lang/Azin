package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

// Trivia represents whitespace, newlines, or comments.
type Trivia struct {
	Base
	text string
}

func NewTrivia(kind syntax.SyntaxKind, text string) *Trivia {
	return &Trivia{
		Base: Base{
			kind:      kind,
			fullWidth: uint32(len(text)),
		},
		text: text,
	}
}

func (t *Trivia) Text() string {
	return t.text
}

// SlotCount is 0 because trivia cannot have children.
func (t *Trivia) SlotCount() int {
	return 0
}

func (t *Trivia) Slot(index int) Node {
	return nil
}
