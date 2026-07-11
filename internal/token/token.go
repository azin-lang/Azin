package token

// represents a token in the source code.
type Token struct {
	Kind     Kind
	Position Position
	Length   uint32
}

// returns whether the token has text associated with it.
func (k Kind) HasText() bool {
	switch k {
	case Identifier,
		IntegerLiteral,
		FloatLiteral,
		StringLiteral,
		CharacterLiteral:
		return true
	default:
		return false
	}
}
