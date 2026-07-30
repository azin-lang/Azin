package lexer

import (
	"github.com/azin-lang/Azin/pkg/token"
)

// lexOperator evaluates and returns a token for multi-character or single-character operator symbols.
func (l *Lexer) lexOperator(ch rune, start token.Position) token.Token {
	switch ch {
	case '+':
		return l.lexPlus(start)
	case '-':
		return l.lexMinus(start)
	case '*':
		return l.either('=', token.StarEqual, token.Star, start)
	case '/':
		return l.either('=', token.SlashEqual, token.Slash, start)
	case '%':
		return l.either('=', token.ModuloEqual, token.Modulo, start)
	case '=':
		return l.either('=', token.EqualEqual, token.Equal, start)
	case '!':
		return l.either('=', token.BangEqual, token.Bang, start)
	case '<':
		if l.match('=') {
			return l.emit(token.LessEqual, start)
		}
		if l.match('<') {
			return l.emit(token.LessLess, start)
		}
		return l.emit(token.Less, start)
	case '>':
		if l.match('=') {
			return l.emit(token.GreaterEqual, start)
		}
		if l.match('>') {
			return l.emit(token.GreaterGreater, start)
		}
		return l.emit(token.Greater, start)
	case '&':
		if l.match('&') {
			return l.emit(token.LogicalAnd, start)
		}
		if l.match('=') {
			return l.emit(token.AmpersandEqual, start)
		}
		return l.emit(token.Ampersand, start)
	case '|':
		if l.match('|') {
			return l.emit(token.LogicalOr, start)
		}
		if l.match('=') {
			return l.emit(token.PipeEqual, start)
		}
		return l.emit(token.Pipe, start)
	default:
		return l.lexUnknown(start)
	}
}

// lexPlus handles plus, increment, and addition assignment operators.
func (l *Lexer) lexPlus(start token.Position) token.Token {
	if l.match('=') {
		return l.emit(token.PlusEqual, start)
	}
	if l.match('+') {
		return l.emit(token.PlusPlus, start)
	}
	return l.emit(token.Plus, start)
}

// lexMinus handles minus, decrement, subtraction assignment, and arrow operators.
func (l *Lexer) lexMinus(start token.Position) token.Token {
	if l.match('=') {
		return l.emit(token.MinusEqual, start)
	}
	if l.match('-') {
		return l.emit(token.MinusMinus, start)
	}
	if l.match('>') {
		return l.emit(token.Arrow, start)
	}
	return l.emit(token.Minus, start)
}

// lexUnknown handles unexpected or invalid sequences of characters, reporting a diagnostic error and returning an Unknown token.
func (l *Lexer) lexUnknown(start token.Position) token.Token {
	for !l.eof() {
		r := l.peek()
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			break
		}
		if isIdentifierStart(r) || isDigit(r) || isPunctuation(r) {
			break
		}

		stop := false
		switch r {
		case '+', '-', '*', '/', '%', '=', '!', '<', '>', '&', '|', '"':
			stop = true
		}
		if stop {
			break
		}

		l.advance()
	}

	length := l.cursor - start.Offset
	text := string(l.src[start.Offset:l.cursor])

	l.diag.ReportError(start, int(length), "unexpected characters: %q", text)
	return l.emit(token.Unknown, start)
}
