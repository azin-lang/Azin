package lexer

import "github.com/azin-lang/azin/compiler/internal/token"

func (l *Lexer) nextToken() token.Token {
	l.skipTrivia()

	if l.eof() {
		return l.eofToken()
	}

	start := l.pos()
	ch, _ := l.advance()

	switch {
	case ch == '\n':
		return l.emit(token.Newline, start)

	case ch == '\r':
		if l.peek() == '\n' {
			l.advance() // consume the LF of CRLF
		}
		return l.emit(token.Newline, start)

	case ch == '"':
		return l.lexString(start)

	case ch == '\'':
		return l.lexCharacter(start)

	case isIdentifierStart(ch):
		return l.lexIdentifier(start)

	case isDigit(ch):
		return l.lexNumber(start)

	case isPunctuation(ch):
		return l.lexPunctuation(ch, start)

	default:
		return l.lexOperator(ch, start)
	}
}
