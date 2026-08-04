package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

// Trivia represents non-executable source elements like whitespace, newlines, or comments.
type Trivia struct {
	Base
	text string
}

func NewTrivia(kind syntax.Kind, text string) *Trivia {
	return &Trivia{
		Base: Base{
			kind:      kind,
			fullWidth: uint32(len(text)),
			flags:     FlagNone,
		},
		text: text,
	}
}

func (t *Trivia) IsNil() bool { return t == nil }

func (t *Trivia) Text() string {
	return t.text
}

func (t *Trivia) SlotCount() int {
	return 0
}

func (t *Trivia) Slot(index int) Node {
	return nil
}
