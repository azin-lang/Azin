package lexer

import "github.com/azin-lang/Azin/internal/token"

// skipTrivia consumes whitespace and comments, advancing the cursor
// to the next meaningful token.
func (l *Lexer) skipTrivia() {
	for {
		l.consumeWhile(func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '\r'
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

	l.diag.ReportError(start, l.cursor-start.Offset, "unterminated block comment")
}
