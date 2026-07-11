package token

import "fmt"

type Token struct {
	Kind   Kind
	Offset uint32
	Length uint32
	Line   uint32
	Column uint32
}

func (t Token) String() string {
	return fmt.Sprintf(
		"%-18s %4d:%4d [%d:%d]",
		t.Kind,
		t.Line,
		t.Column,
		t.Offset,
		t.Length,
	)
}

func (t Token) Format(source []byte) string {
	s := fmt.Sprintf(
		"%-18s %4d:%4d [%d:%d]",
		t.Kind,
		t.Line,
		t.Column,
		t.Offset,
		t.Length,
	)

	switch t.Kind {
	case Identifier, IntegerLiteral, FloatLiteral, StringLiteral, CharacterLiteral:
		if end := t.Offset + t.Length; end <= uint32(len(source)) {
			s += fmt.Sprintf(" %q", source[t.Offset:end])
		}
	}

	return s
}
