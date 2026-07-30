package lexer

import (
	"unicode"
	"unicode/utf8"

	token "github.com/azin-lang/Azin/pkg/token"
)

// lexNumber processes and returns a token for integer, floating-point, hexadecimal, or binary number literals.
func (l *Lexer) lexNumber(pos token.Position) token.Token {
	if l.src[pos.Offset] == '0' {
		p := l.peek()
		if p == 'x' || p == 'X' {
			l.advance()
			if !isHexDigit(l.peek()) {
				l.diag.ReportError(pos, int(l.cursor-pos.Offset), "empty hex literal")
			}
			l.consumeHexDigits()
			return l.emit(token.IntegerLiteral, pos)
		}

		if p == 'b' || p == 'B' {
			l.advance()
			if !isBinaryDigit(l.peek()) {
				l.diag.ReportError(pos, int(l.cursor-pos.Offset), "empty binary literal")
			}
			l.consumeBinaryDigits()
			return l.emit(token.IntegerLiteral, pos)
		}
	}

	l.consumeDecimalDigits()

	isFloat := false
	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance() // Consume '.'
		l.consumeDecimalDigits()
		isFloat = true
	}

	p := l.peek()
	if p == 'e' || p == 'E' {
		l.advance()
		isFloat = true

		nxt := l.peek()
		if nxt == '+' || nxt == '-' {
			l.advance()
		}

		if !isDigit(l.peek()) {
			l.diag.ReportError(l.pos(), 1, "malformed floating-point literal: missing exponent")
		} else {
			l.consumeDecimalDigits()
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

	ch := l.advance()

	if ch == '\'' {
		l.diag.ReportError(pos, 2, "empty character literal")
		return l.emit(token.CharacterLiteral, pos)
	}

	if ch == '\\' {
		if l.eof() {
			l.diag.ReportError(pos, 1, "unterminated escape sequence")
			return l.emit(token.CharacterLiteral, pos)
		}

		escape := l.advance()
		switch escape {
		case '\'', '"', '\\', 'a', 'b', 'f', 'n', 'r', 't', 'v', '0':
			// valid
		default:
			l.diag.ReportError(token.Position{Offset: l.cursor - 1}, 1, "invalid escape sequence \\%c", escape)
		}
	}

	if l.eof() {
		l.diag.ReportError(pos, int(l.cursor-pos.Offset), "unterminated character literal")
		return l.emit(token.CharacterLiteral, pos)
	}

	if l.peek() != '\'' {
		l.diag.ReportError(token.Position{Offset: l.cursor}, 1, "character literal may contain exactly one character")
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
		ch := l.advance()

		switch ch {
		case '"':
			return l.emit(token.StringLiteral, start)

		case '\\':
			if l.eof() {
				l.diag.ReportError(token.Position{Offset: l.cursor - 1}, 1, "unterminated escape sequence")
				return l.emit(token.StringLiteral, start)
			}
			escape := l.advance()
			switch escape {
			case '"', '\\', 'n', 'r', 't', '0':
				// valid
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

func (l *Lexer) consumeDecimalDigits() {
	length := len(l.src)
	for int(l.cursor) < length {
		c := l.src[l.cursor]
		if c >= '0' && c <= '9' {
			l.cursor++
			continue
		}
		if c >= utf8.RuneSelf {
			r, sz := utf8.DecodeRune(l.src[l.cursor:])
			if unicode.IsDigit(r) {
				l.cursor += uint32(sz)
				continue
			}
		}
		break
	}
}

func (l *Lexer) consumeHexDigits() {
	length := len(l.src)
	for int(l.cursor) < length {
		c := l.src[l.cursor]
		if ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F') {
			l.cursor++
			continue
		}
		break
	}
}

func (l *Lexer) consumeBinaryDigits() {
	length := len(l.src)
	for int(l.cursor) < length {
		c := l.src[l.cursor]
		if c == '0' || c == '1' {
			l.cursor++
			continue
		}
		break
	}
}
