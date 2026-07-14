package lexer

import (
	"unicode"
	"unicode/utf8"
)

func isIdentifierStart(r rune) bool {
	if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_' {
		return true
	}
	return r >= utf8.RuneSelf && unicode.IsLetter(r)
}

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
