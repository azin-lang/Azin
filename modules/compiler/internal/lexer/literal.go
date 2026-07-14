package lexer

import "github.com/azin-lang/azin/compiler/internal/token"

func (l *Lexer) lexNumber(start token.Position) token.Token {
	if ch, _ := l.file.Rune(start.Offset); ch == '0' {
		if l.matchAny("xX") {
			l.consumeWhile(func(r rune) bool {
				return isDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
			})
			return l.emit(token.IntegerLiteral, start)
		}

		if l.matchAny("bB") {
			l.consumeWhile(func(r rune) bool { return r == '0' || r == '1' })
			return l.emit(token.IntegerLiteral, start)
		}
	}

	l.consumeWhile(isDigit)

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()             // Consume the '.'
		l.consumeWhile(isDigit) // Consume the fractional digits

		return l.emit(token.FloatLiteral, start)
	}

	return l.emit(token.IntegerLiteral, start)
}

func (l *Lexer) lexCharacter(start token.Position) token.Token {
	if l.eof() {
		l.diag.ReportError(start, 1, "unterminated character literal")
		return l.emit(token.CharacterLiteral, start)
	}

	ch, _ := l.advance()

	// Reject ''
	if ch == '\'' {
		l.diag.ReportError(start, 2, "empty character literal")
		return l.emit(token.CharacterLiteral, start)
	}

	if ch == '\\' {
		if l.eof() {
			l.diag.ReportError(start, 1, "unterminated escape sequence")
			return l.emit(token.CharacterLiteral, start)
		}

		escape, _ := l.advance()

		switch escape {
		case '\'', '"', '\\', 'a', 'b', 'f', 'n', 'r', 't', 'v', '0':
			// valid escape
		default:
			l.diag.ReportError(
				token.Position{Offset: l.cursor - 1},
				1,
				"invalid escape sequence \\%c",
				escape,
			)
		}
	}

	if l.eof() {
		l.diag.ReportError(start, l.cursor-start.Offset, "unterminated character literal")
		return l.emit(token.CharacterLiteral, start)
	}

	if l.peek() != '\'' {
		l.diag.ReportError(
			token.Position{Offset: l.cursor},
			1,
			"character literal may contain exactly one character",
		)

		// Recover by skipping to the closing quote or newline.
		for !l.eof() && l.peek() != '\'' && l.peek() != '\n' && l.peek() != '\r' {
			l.advance()
		}
	}

	if l.peek() == '\'' {
		l.advance()
	}

	return l.emit(token.CharacterLiteral, start)
}

func (l *Lexer) lexString(start token.Position) token.Token {
	for !l.eof() {
		ch, _ := l.advance()

		if ch == '\n' || ch == '\r' {
			l.diag.ReportError(start, l.cursor-start.Offset, "unterminated character literal")
			return l.emit(token.CharacterLiteral, start)
		}

		switch ch {
		case '"':
			return l.emit(token.StringLiteral, start)

		case '\\':
			if l.eof() {
				l.diag.ReportError(token.Position{Offset: l.cursor - 1}, 1, "unterminated escape sequence")
				return l.emit(token.StringLiteral, start)
			}
			escape, _ := l.advance()
			switch escape {
			case '"', '\\', 'n', 'r', 't', '0':
				// Valid escape sequence
			default:
				l.diag.ReportError(token.Position{Offset: l.cursor - 1}, 1, "invalid escape sequence \\%c", escape)
			}

		case '\n', '\r':
			l.diag.ReportError(start, l.cursor-start.Offset, "unterminated string literal")
			return l.emit(token.StringLiteral, start)
		}
	}

	l.diag.ReportError(start, l.cursor-start.Offset, "unterminated string literal")
	return l.emit(token.StringLiteral, start)
}
