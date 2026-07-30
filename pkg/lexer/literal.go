package lexer

import (
	token "github.com/azin-lang/Azin/pkg/token"
)

func (l *Lexer) lexNumber(start token.Position) token.Token {
	if l.src[start.Offset] == '0' {
		if l.matchAny("xX") {
			if !isHexDigit(l.peek()) {
				l.diag.ReportError(start, int(l.cursor-start.Offset), "empty hex literal")
			}
			l.consumeWhile(isHexDigit)
			return l.emit(token.IntegerLiteral, start)
		}

		if l.matchAny("bB") {
			if !isBinaryDigit(l.peek()) {
				l.diag.ReportError(start, int(l.cursor-start.Offset), "empty binary literal")
			}
			l.consumeWhile(isBinaryDigit)
			return l.emit(token.IntegerLiteral, start)
		}
	}

	l.consumeWhile(isDigit)

	isFloat := false

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()             // Consume the '.'
		l.consumeWhile(isDigit) // Consume the fractional digits
		isFloat = true
	}

	// Support scientific notation (e.g., 1e10, 1.23E-4)
	if l.peek() == 'e' || l.peek() == 'E' {
		l.advance() // Consume 'e' or 'E'
		isFloat = true

		if l.peek() == '+' || l.peek() == '-' {
			l.advance()
		}

		if !isDigit(l.peek()) {
			l.diag.ReportError(l.pos(), 1, "malformed floating-point literal: missing exponent")
		} else {
			l.consumeWhile(isDigit)
		}
	}

	if isFloat {
		return l.emit(token.FloatLiteral, start)
	}

	return l.emit(token.IntegerLiteral, start)
}

func isHexDigit(r rune) bool {
	return ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
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
		l.diag.ReportError(start, int(l.cursor-start.Offset), "unterminated character literal")
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
			l.diag.ReportError(start, int(l.cursor-start.Offset), "unterminated string literal")
			return l.emit(token.StringLiteral, start)
		}
	}

	l.diag.ReportError(start, int(l.cursor-start.Offset), "unterminated string literal")
	return l.emit(token.StringLiteral, start)
}
