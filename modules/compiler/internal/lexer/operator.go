package lexer

import "github.com/azin-lang/azin/compiler/internal/token"

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
	case '"':
		return l.lexString(start)
	default:
		return l.lexUnknown(start)
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

func (l *Lexer) lexUnknown(start token.Position) token.Token {
	l.consumeWhile(func(r rune) bool {
		// Stop on EOF
		if r == 0 {
			return false
		}

		// Stop on whitespace
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return false
		}

		// Stop on any character that could start a valid token
		if isIdentifierStart(r) || isDigit(r) || isPunctuation(r) {
			return false
		}

		// Stop on valid operator characters and quotes
		switch r {
		case '+', '-', '*', '/', '%', '=', '!', '<', '>', '&', '|', '"':
			return false
		}

		// Otherwise, it's more garbage. Keep eating it!
		return true
	})

	length := l.cursor - start.Offset
	text := string(l.file.Slice(start.Offset, l.cursor))

	l.diag.ReportError(start, length, "unexpected characters: %q", text)
	return l.emit(token.Unknown, start)
}
