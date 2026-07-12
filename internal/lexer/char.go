package lexer

import (
	"unicode"
	"unicode/utf8"
)

// isIdentifierStart reports whether a rune is valid as the first character of an identifier.
func isIdentifierStart(r rune) bool {
	if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_' {
		return true
	}
	return r >= utf8.RuneSelf && unicode.IsLetter(r)
}

// isIdentifierContinue reports whether a rune is valid within an identifier.
func isIdentifierContinue(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r)
}

func isDigit(r rune) bool {
	if '0' <= r && r <= '9' {
		return true
	}
	return r >= utf8.RuneSelf && unicode.IsDigit(r)
}

func isPunctuation(r rune) bool {
	switch r {
	case '(', ')', '{', '}', '[', ']', ',', ';', ':', '.':
		return true
	default:
		return false
	}
}
