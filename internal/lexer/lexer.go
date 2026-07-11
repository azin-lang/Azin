// the lexer package is a lexical analyzer that tokenizes source code.
package lexer

import (
	"github.com/azin-lang/Azin/internal/diagnostics" // diagnostics package is used for reporting errors
	"github.com/azin-lang/Azin/internal/source"      // source package is used for representing source files
	"github.com/azin-lang/Azin/internal/token"       // token package is used for representing tokens
)

// represents a lexical analyzer that tokenizes source code.
type Lexer struct {
	file   *source.File
	offset uint32
	diag   *diagnostics.Engine
}

// creates a new Lexer for the given source file and diagnostics engine.
func New(file *source.File, diag *diagnostics.Engine) *Lexer {
	return &Lexer{
		file: file,
		diag: diag,
	}
}

// tokenizes the source code and returns a slice of tokens.
func (l *Lexer) Tokenize() []token.Token {
	// Added a small capacity to prevent early reallocations
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

// nextToken is a helper function that returns the next token from the source code.
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

// lexes a plus token, handling plus equal and plus plus cases.
func (l *Lexer) lexPlus(start token.Position) token.Token {
	if l.match('=') {
		return l.token(token.PlusEqual, start)
	}
	if l.match('+') {
		return l.token(token.PlusPlus, start)
	}
	return l.token(token.Plus, start)
}

// lexes a minus token, handling minus equal, minus minus, and arrow cases.
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

// lexes an identifier token, handling keywords and regular identifiers.
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

// lexes an integer token, handling integer literals.
func (l *Lexer) lexInteger(start token.Position) token.Token {
	for isDigit(l.peek()) {
		l.advance()
	}

	return l.token(token.IntegerLiteral, start)
}

// returns the end-of-file token.
func (l *Lexer) eofToken() token.Token {
	return token.Token{
		Kind:     token.EOF,
		Position: l.position(),
	}
}

// returns true if the lexer is at the end of the file.
func (l *Lexer) eof() bool {
	return l.file.EOF(l.offset)
}

// returns the next byte in the file without advancing the offset.
func (l *Lexer) peek() byte {
	if l.eof() {
		return 0
	}

	return l.file.Byte(l.offset)
}

// returns true if the next byte in the file matches the given byte, and advances the offset if it does.
func (l *Lexer) match(ch byte) bool {
	if l.peek() != ch {
		return false
	}

	l.advance()
	return true
}

// advances the lexer's offset by one byte and returns the byte that was advanced.
func (l *Lexer) advance() byte {
	if l.eof() {
		return 0
	}

	ch := l.file.Byte(l.offset)
	l.offset++

	return ch
}

// skips over whitespace characters in the file.
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

// returns a token with the given kind and start position.
func (l *Lexer) token(kind token.Kind, start token.Position) token.Token {
	return token.Token{
		Kind:     kind,
		Position: start,
		Length:   l.offset - start.Offset,
	}
}

// returns the current position in the file.
func (l *Lexer) position() token.Position {
	return token.Position{
		Offset: l.offset,
	}
}

// returns true if the given byte is an alphabetic character (letter or underscore).
func isAlpha(ch byte) bool {
	return ch == '_' ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z')
}

// returns true if the given byte is a digit.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// returns true if the given byte is an alphanumeric character (letter, digit, or underscore).
func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
