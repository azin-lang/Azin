package lexer

import (
	token "github.com/azin-lang/Azin/pkg/token"
)

// lexNumber processes and returns a token for integer, floating-point, hexadecimal, or binary number literals.
func (l *Lexer) lexNumber(pos token.Position) token.Token {
	if l.src[pos.Offset] == '0' {
		if l.matchAny("xX") {
			if !isHexDigit(l.peek()) {
				l.diag.ReportError(pos, int(l.cursor-pos.Offset), "empty hex literal")
			}
			l.consumeWhile(isHexDigit)
			return l.emit(token.IntegerLiteral, pos)
		}

		if l.matchAny("bB") {
			if !isBinaryDigit(l.peek()) {
				l.diag.ReportError(pos, int(l.cursor-pos.Offset), "empty binary literal")
			}
			l.consumeWhile(isBinaryDigit)
			return l.emit(token.IntegerLiteral, pos)
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
		return l.emit(token.FloatLiteral, pos)
	}

	return l.emit(token.IntegerLiteral, pos)
}

// isHexDigit reports whether r is a valid hexadecimal digit.
func isHexDigit(r rune) bool {
	return ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

// isBinaryDigit reports whether r is a valid binary digit.
func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
}

// lexCharacter scans a character literal starting at pos, validating any escape sequences.
func (l *Lexer) lexCharacter(pos token.Position) token.Token {
	if l.eof() {
		l.diag.ReportError(pos, 1, "unterminated character literal")
		return l.emit(token.CharacterLiteral, pos)
	}

	ch, _ := l.advance()

	// Reject ''
	if ch == '\'' {
		l.diag.ReportError(pos, 2, "empty character literal")
		return l.emit(token.CharacterLiteral, pos)
	}

	if ch == '\\' {
		if l.eof() {
			l.diag.ReportError(pos, 1, "unterminated escape sequence")
			return l.emit(token.CharacterLiteral, pos)
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
		l.diag.ReportError(pos, int(l.cursor-pos.Offset), "unterminated character literal")
		return l.emit(token.CharacterLiteral, pos)
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

	return l.emit(token.CharacterLiteral, pos)
}

// lexString scans a string literal starting at pos, handling multi-line checks and escape sequence validation.
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
