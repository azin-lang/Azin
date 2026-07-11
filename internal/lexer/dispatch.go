package lexer

import "github.com/azin-lang/Azin/internal/token"

// nextToken scans and returns the next token from the source file.
func (l *Lexer) nextToken() token.Token {
	l.skipTrivia()

	if l.eof() {
		return l.eofToken()
	}

	start := l.pos()
	ch, size := l.advance()

	switch {
	case isIdentifierStart(ch):
		return l.lexIdentifier(start)

	case isDigit(ch):
		return l.lexNumber(start)

	case isPunctuation(ch):
		return l.lexPunctuation(ch, start)

	default:
		return l.lexOperator(ch, size, start)
	}
}
