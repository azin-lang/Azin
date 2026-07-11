// Package lexer converts Azin source text into lexical tokens.
package lexer

import (
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

// Lexer breaks source code into tokens.
type Lexer struct {
	// file is the source file being lexed.
	file *source.File

	// offset is the current byte offset within file.
	offset uint32

	// diag collects diagnostics produced while lexing.
	diag *diagnostics.Engine
}

// New returns a new Lexer for the given file.
func New(file *source.File, diag *diagnostics.Engine) *Lexer {
	return &Lexer{
		file: file,
		diag: diag,
	}
}

// Tokenize reads the entire file and returns its tokens.
func (l *Lexer) Tokenize() []token.Token {
	tokens := make([]token.Token, 0, 128)

	for {
		tok := l.nextToken()
		tokens = append(tokens, tok)

		if tok.Kind == token.EOF {
			break
		}
	}

	return tokens
}

// nextToken lexes and returns the next token from the input stream.
func (l *Lexer) nextToken() token.Token {
	l.skipWhitespace()

	if l.eof() {
		return l.eofToken()
	}

	start := l.position()
	ch := l.advance()

	switch {
	case isAlpha(ch):
		return l.lexIdentifier(start)

	case isDigit(ch):
		return l.lexInteger(start)

	default:
		return l.lexSymbol(ch, start)
	}
}

// lexSymbol lexes punctuation and operator tokens beginning with ch.
func (l *Lexer) lexSymbol(ch byte, start token.Position) token.Token {
	switch ch {
	case '(':
		return l.token(token.LeftParen, start)
	case ')':
		return l.token(token.RightParen, start)
	case '{':
		return l.token(token.LeftBrace, start)
	case '}':
		return l.token(token.RightBrace, start)
	case '[':
		return l.token(token.LeftBracket, start)
	case ']':
		return l.token(token.RightBracket, start)
	case ',':
		return l.token(token.Comma, start)
	case ';':
		return l.token(token.Semicolon, start)
	case ':':
		return l.token(token.Colon, start)
	case '.':
		return l.token(token.Dot, start)
	case '+':
		return l.lexPlus(start)
	case '-':
		return l.lexMinus(start)
	case '"':
		return l.lexString(start)

	case '=':
		if l.match('=') {
			return l.token(token.EqualEqual, start)
		}
		return l.token(token.Equal, start)
	case '!':
		if l.match('=') {
			return l.token(token.BangEqual, start)
		}
		return l.token(token.Bang, start)
	case '<':
		if l.match('=') {
			return l.token(token.LessEqual, start)
		}
		if l.match('<') {
			return l.token(token.LessLess, start)
		}
		return l.token(token.Less, start)
	case '>':
		if l.match('=') {
			return l.token(token.GreaterEqual, start)
		}
		if l.match('>') {
			return l.token(token.GreaterGreater, start)
		}
		return l.token(token.Greater, start)
	case '*':
		if l.match('=') {
			return l.token(token.StarEqual, start)
		}
		return l.token(token.Star, start)
	case '/':
		if l.match('=') {
			return l.token(token.SlashEqual, start)
		}
		return l.token(token.Slash, start)
	case '%':
		if l.match('=') {
			return l.token(token.ModuloEqual, start)
		}
		return l.token(token.Modulo, start)
	case '&':
		if l.match('&') {
			return l.token(token.LogicalAnd, start)
		}
		if l.match('=') {
			return l.token(token.AmpersandEqual, start)
		}
		return l.token(token.Ampersand, start)
	case '|':
		if l.match('|') {
			return l.token(token.LogicalOr, start)
		}
		if l.match('=') {
			return l.token(token.PipeEqual, start)
		}
		return l.token(token.Pipe, start)

	default:
		l.diag.ReportError(start, 1, "unexpected character %q", ch)
		return l.token(token.Unknown, start)
	}
}

// lexPlus lexes '+', '+=', and '++' tokens.
func (l *Lexer) lexPlus(start token.Position) token.Token {
	if l.match('=') {
		return l.token(token.PlusEqual, start)
	}
	if l.match('+') {
		return l.token(token.PlusPlus, start)
	}
	return l.token(token.Plus, start)
}

// lexMinus lexes '-', '-=', '--', and '->' tokens.
func (l *Lexer) lexMinus(start token.Position) token.Token {
	if l.match('=') {
		return l.token(token.MinusEqual, start)
	}
	if l.match('-') {
		return l.token(token.MinusMinus, start)
	}
	if l.match('>') {
		return l.token(token.Arrow, start)
	}
	return l.token(token.Minus, start)
}

// lexIdentifier lexes an identifier or keyword beginning at start.
func (l *Lexer) lexIdentifier(start token.Position) token.Token {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	text := string(l.file.Slice(start.Offset, l.offset))

	if kind, ok := token.Keywords[text]; ok {
		return l.token(kind, start)
	}

	return l.token(token.Identifier, start)
}

// lexInteger lexes a decimal integer literal.
func (l *Lexer) lexInteger(start token.Position) token.Token {
	for isDigit(l.peek()) {
		l.advance()
	}

	return l.token(token.IntegerLiteral, start)
}

// lexString lexes a double-quoted string literal.
//
// Supported escape sequences are:
//
//   - \\
//   - \"
//   - \n
//   - \r
//   - \t
//   - \0
//
// If the string is unterminated or contains an invalid escape sequence,
// a diagnostic is reported before the token is returned.
func (l *Lexer) lexString(start token.Position) token.Token {
	for !l.eof() {
		switch ch := l.advance(); ch {

		case '"':
			// Closing quote.
			return l.token(token.StringLiteral, start)

		case '\\':
			// '\' cannot be the final character
			if l.eof() {
				l.diag.ReportError(
					token.Position{Offset: l.offset - 1},
					1,
					"unterminated escape sequence",
				)
				return l.token(token.StringLiteral, start)
			}

			escape := l.advance()

			switch escape {
			case '"', '\\', 'n', 'r', 't', '0':
				// Valid escape sequence

			default:
				l.diag.ReportError(
					token.Position{Offset: l.offset - 1},
					1,
					"invalid escape sequence \\%c",
					escape,
				)
			}

		case '\n', '\r':
			l.diag.ReportError(
				start,
				l.offset-start.Offset,
				"unterminated string literal",
			)
			return l.token(token.StringLiteral, start)
		}
	}

	l.diag.ReportError(
		start,
		l.offset-start.Offset,
		"unterminated string literal",
	)

	return l.token(token.StringLiteral, start)
}

// eofToken returns the end-of-file token.
func (l *Lexer) eofToken() token.Token {
	return token.Token{
		Kind:     token.EOF,
		Position: l.position(),
	}
}

// eof reports whether the lexer has reached the end of the source file.
func (l *Lexer) eof() bool {
	return l.file.EOF(l.offset)
}

// peek returns the current byte without advancing the lexer.
// It returns 0 if the end of the file has been reached.
func (l *Lexer) peek() byte {
	if l.eof() {
		return 0
	}

	return l.file.Byte(l.offset)
}

// match consumes ch if it is the next byte and reports whether it matched.
func (l *Lexer) match(ch byte) bool {
	if l.peek() != ch {
		return false
	}

	l.advance()
	return true
}

// advance consumes and returns the next byte.
// It returns 0 if the end of the file has been reached.
func (l *Lexer) advance() byte {
	if l.eof() {
		return 0
	}

	ch := l.file.Byte(l.offset)
	l.offset++

	return ch
}

// skipWhitespace consumes consecutive whitespace characters.
func (l *Lexer) skipWhitespace() {
	for !l.eof() {
		switch l.peek() {
		case ' ', '\t', '\r', '\n':
			l.advance()
		default:
			return
		}
	}
}

// token constructs a token of kind beginning at start.
func (l *Lexer) token(kind token.Kind, start token.Position) token.Token {
	return token.Token{
		Kind:     kind,
		Position: start,
		Length:   l.offset - start.Offset,
	}
}

// position returns the current position within the source file.
func (l *Lexer) position() token.Position {
	return token.Position{
		Offset: l.offset,
	}
}

// isAlpha reports whether ch is an ASCII letter or underscore.
func isAlpha(ch byte) bool {
	return ch == '_' ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z')
}

// isDigit reports whether ch is an ASCII decimal digit.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// isAlphaNumeric reports whether ch is an ASCII letter, digit, or underscore.
func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
