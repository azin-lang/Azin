// Package lexer implements lexical analysis for the Azin programming language.
// It transforms raw source text into a stream of structured syntax tokens.
package lexer

import (
	"unicode"
	"unicode/utf8"
)

// isIdentifierStart reports whether a rune can start an identifier (letters and underscores).
func isIdentifierStart(r rune) bool {
	if r < utf8.RuneSelf {
		return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || r == '_'
	}
	return unicode.IsLetter(r)
}

// isIdentifierContinue reports whether a rune can continue an identifier (letters, digits, underscores, and marks).
func isIdentifierContinue(r rune) bool {
	if r < utf8.RuneSelf {
		return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r == '_'
	}
	return unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r)
}

// isDigit reports whether a rune is a decimal digit.
func isDigit(r rune) bool {
	if r < utf8.RuneSelf {
		return '0' <= r && r <= '9'
	}
	return unicode.IsDigit(r)
}

// isPunctuation reports whether a rune is a valid punctuation token separator.
func isPunctuation(r rune) bool {
	switch r {
	case '(', ')', '{', '}', '[', ']', ',', ';', ':', '.':
		return true
	default:
		return false
	}
}
