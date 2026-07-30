package lexer

import (
	"github.com/azin-lang/Azin/pkg/token"
)

// skipTrivia consumes whitespace and comments, advancing the cursor
// to the next meaningful syntax token.
func (l *Lexer) skipTrivia() {
	for {
		l.consumeWhile(func(r rune) bool {
			return r == ' ' || r == '\t'
		})

		if l.peek() == '/' {
			next := l.peekNext()

			if next == '/' {
				l.advance() // Consume first '/'
				l.advance() // Consume second '/'
				l.skipLineComment()
				continue
			}

			if next == '*' {
				start := l.pos()
				l.advance() // Consume '/'
				l.advance() // Consume '*'
				l.skipBlockComment(start)
				continue
			}
		}

		return
	}
}

// skipLineComment consumes characters until the end of the line.
func (l *Lexer) skipLineComment() {
	for !l.eof() {
		switch l.peek() {
		case '\n', '\r':
			return
		default:
			_, _ = l.advance()
		}
	}
}

// skipBlockComment consumes characters until the matching closing block comment delimiter, properly supporting nested block comments.
func (l *Lexer) skipBlockComment(start token.Position) {
	depth := 1

	for !l.eof() {
		ch, _ := l.advance()

		switch ch {
		case '/':
			if l.match('*') {
				depth++
			}
		case '*':
			if l.match('/') {
				depth--
				if depth == 0 {
					return
				}
			}
		}
	}

	l.diag.ReportError(start, int(l.cursor-start.Offset), "unterminated block comment")
}
