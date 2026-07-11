package lexer

import (
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/token"
)

type Lexer struct {
	input  []byte
	offset uint32
	line   uint32
	column uint32
	diag   *diagnostics.Engine
}

func New(input []byte, diag *diagnostics.Engine) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 1,
		diag:   diag,
	}
}

func (l *Lexer) Tokenize() []token.Token {
	tokens := make([]token.Token, 0)

	for {
		tok := l.next()
		tokens = append(tokens, tok)

		if tok.Kind == token.EOF {
			break
		}
	}

	return tokens
}

// wyd
func (l *Lexer) next() token.Token {
	l.skipWhitespace()

	if l.eof() {
		return token.Token{
			Kind:   token.EOF,
			Offset: l.offset,
			Length: 0,
			Line:   l.line,
			Column: l.column,
		}
	}

	start := l.offset
	line := l.line
	column := l.column

	ch := l.advance()

	switch {
	case isAlpha(ch):
		return l.lexIdentifier(start, line, column)
	case isDigit(ch):
		return l.lexInteger(start, line, column)
	}

	switch ch {
	case '(':
		return l.token(token.LeftParen, start, line, column)

	case ')':
		return l.token(token.RightParen, start, line, column)

	case '{':
		return l.token(token.LeftBrace, start, line, column)

	case '}':
		return l.token(token.RightBrace, start, line, column)

	case ',':
		return l.token(token.Comma, start, line, column)

	case ';':
		return l.token(token.Semicolon, start, line, column)

	case ':':
		return l.token(token.Colon, start, line, column)

	case '.':
		return l.token(token.Dot, start, line, column)

	case '+':
		return l.lexPlus(start, line, column)

	case '-':
		return l.lexMinus(start, line, column)

	default:
		l.diag.ReportError(
			"unexpected character %q at %d:%d",
			ch,
			line,
			column,
		)
		return l.token(token.Unknown, start, line, column)
	}
}

func (l *Lexer) lexPlus(start, line, column uint32) token.Token {
	switch {
	case l.match('='):
		return l.token(token.PlusEqual, start, line, column)
	default:
		return l.token(token.Plus, start, line, column)
	}
}

func (l *Lexer) lexMinus(start, line, column uint32) token.Token {
	switch {
	case l.match('='):
		return l.token(token.MinusEqual, start, line, column)
	default:
		return l.token(token.Minus, start, line, column)
	}
}

func (l *Lexer) lexIdentifier(start, line, column uint32) token.Token {

	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	text := string(l.input[start:l.offset])
	kind := token.Identifier
	if kw, ok := token.Keywords[text]; ok {
		kind = kw
	}
	return l.token(kind, start, line, column)
}

func (l *Lexer) lexInteger(start, line, column uint32) token.Token {
	for isDigit(l.peek()) {
		l.advance()
	}
	return l.token(token.IntegerLiteral, start, line, column)
}

func (l *Lexer) eof() bool {
	return l.offset >= uint32(len(l.input))
}

func (l *Lexer) peek() byte {
	if l.eof() {
		return 0
	}
	return l.input[l.offset]
}

func (l *Lexer) match(ch byte) bool {
	if l.peek() != ch {
		return false
	}
	l.advance()
	return true
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.peek() {
		case ' ', '\t', '\r', '\n':
			l.advance()
		default:
			return
		}
	}
}

func isAlpha(ch byte) bool {
	return ch == '_' ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func (l *Lexer) advance() byte {
	if l.offset >= uint32(len(l.input)) {
		return 0
	}

	ch := l.input[l.offset]
	l.offset++

	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	return ch
}

func (l *Lexer) token(kind token.Kind, start, line, column uint32) token.Token {
	return token.Token{
		Kind:   kind,
		Offset: start,
		Length: l.offset - start,
		Line:   line,
		Column: column,
	}
}
