package lexer

import "github.com/azin-lang/Azin/internal/token"

func (l *Lexer) lexOperator(ch rune, size uint32, start token.Position) token.Token {
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
	case '"':
		return l.lexString(start)
	default:
		l.diag.ReportError(start, size, "unexpected character %q", ch)
		return l.emit(token.Unknown, start)
	}
}

func (l *Lexer) lexPlus(start token.Position) token.Token {
	if l.match('=') {
		return l.emit(token.PlusEqual, start)
	}
	if l.match('+') {
		return l.emit(token.PlusPlus, start)
	}
	return l.emit(token.Plus, start)
}

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
