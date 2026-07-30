package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/azin-lang/Azin/pkg/token"
)

// lexIdentifier consumes a sequence of identifier continuation characters and resolves it to a keyword or identifier token.
func (l *Lexer) lexIdentifier(start token.Position) token.Token {
	// Inline loop with an ASCII fast-path
	for int(l.cursor) < len(l.src) {
		c := l.src[l.cursor]
		if c < utf8.RuneSelf {
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '_' {
				l.cursor++
				continue
			}
			break
		}

		// Fallback for Unicode identifiers
		r, sz := utf8.DecodeRune(l.src[l.cursor:])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r) {
			l.cursor += uint32(sz)
			continue
		}
		break
	}

	if kind, ok := token.Keywords[string(l.src[start.Offset:l.cursor])]; ok {
		return l.emit(kind, start)
	}

	return l.emit(token.Identifier, start)
}
