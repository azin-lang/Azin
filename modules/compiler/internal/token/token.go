// Package token defines constants and types representing Azin lexical tokens.
package token

// Token represents a single lexical unit.
type Token struct {
	Kind     Kind
	Position Position
	Length   uint32
}

// HasText reports whether the token kind contains variable text content.
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
