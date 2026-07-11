package lexer

import (
	"unicode"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

type Lexer struct {
	file   *source.File
	offset uint32
	diag   *diagnostics.Engine
}

func New(file *source.File, diag *diagnostics.Engine) *Lexer {
	return &Lexer{
		file: file,
		diag: diag,
	}
}

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

func (l *Lexer) skipComment(start token.Position) bool {
	switch l.peek() {

	// //
	case '/':
		l.advance()

		for !l.eof() {
			switch l.peek() {
			case '\n', '\r':
				return true
			default:
				l.advance()
			}
		}

		return true

	// /*
	case '*':
		l.advance()

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
						return true
					}
				}
			}
		}

		l.diag.ReportError(
			start,
			l.offset-start.Offset,
			"unterminated block comment",
		)

		return true
	}

	return false
}

func (l *Lexer) nextToken() token.Token {
	for {
		l.skipWhitespace()

		if l.eof() {
			return l.eofToken()
		}

		start := l.position()
		ch, size := l.advance()

		// comments
		if ch == '/' && l.skipComment(start) {
			continue
		}

		switch {
		case isIdentifierStart(ch):
			return l.lexIdentifier(start)

		case isDigit(ch):
			return l.lexInteger(start)

		default:
			return l.lexSymbol(ch, size, start)
		}
	}
}

func (l *Lexer) lexSymbol(ch rune, size uint32, start token.Position) token.Token {
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
		l.diag.ReportError(
			start,
			size,
			"unexpected character %q",
			ch,
		)

		return l.token(token.Unknown, start)
	}
}

func (l *Lexer) lexPlus(start token.Position) token.Token {
	if l.match('=') {
		return l.token(token.PlusEqual, start)
	}
	if l.match('+') {
		return l.token(token.PlusPlus, start)
	}
	return l.token(token.Plus, start)
}

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

func (l *Lexer) lexIdentifier(start token.Position) token.Token {
	for isIdentifierContinue(l.peek()) {
		_, _ = l.advance()
	}

	text := string(l.file.Slice(start.Offset, l.offset))

	if kind, ok := token.Keywords[text]; ok {
		return l.token(kind, start)
	}

	return l.token(token.Identifier, start)
}

func (l *Lexer) lexInteger(start token.Position) token.Token {
	for isDigit(l.peek()) {
		_, _ = l.advance()
	}

	return l.token(token.IntegerLiteral, start)
}

func (l *Lexer) lexString(start token.Position) token.Token {
	for !l.eof() {
		ch, _ := l.advance()

		switch ch {
		case '"':
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

			escape, _ := l.advance()

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

func (l *Lexer) eofToken() token.Token {
	return token.Token{
		Kind:     token.EOF,
		Position: l.position(),
	}
}

func (l *Lexer) eof() bool {
	return l.file.EOF(l.offset)
}

func (l *Lexer) peek() rune {
	if l.eof() {
		return 0
	}

	r, _ := l.file.Rune(l.offset)
	return r
}

func (l *Lexer) match(ch rune) bool {
	if l.peek() != ch {
		return false
	}

	_, _ = l.advance()
	return true
}

func (l *Lexer) advance() (rune, uint32) {
	if l.eof() {
		return 0, 0
	}

	r, size := l.file.Rune(l.offset)
	l.offset += size

	return r, size
}

func (l *Lexer) skipWhitespace() {
	for !l.eof() {
		switch l.peek() {
		case ' ', '\t', '\r', '\n':
			_, _ = l.advance()
		default:
			return
		}
	}
}

func (l *Lexer) token(kind token.Kind, start token.Position) token.Token {
	return token.Token{
		Kind:     kind,
		Position: start,
		Length:   l.offset - start.Offset,
	}
}

func (l *Lexer) position() token.Position {
	return token.Position{
		Offset: l.offset,
	}
}

func isIdentifierStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentifierContinue(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}
