package lexer

import "github.com/azin-lang/azin/compiler/internal/token"

func (l *Lexer) lexPunctuation(ch rune, start token.Position) token.Token {
	switch ch {
	case '(':
		return l.emit(token.LeftParen, start)
	case ')':
		return l.emit(token.RightParen, start)
	case '{':
		return l.emit(token.LeftBrace, start)
	case '}':
		return l.emit(token.RightBrace, start)
	case '[':
		return l.emit(token.LeftBracket, start)
	case ']':
		return l.emit(token.RightBracket, start)
	case ',':
		return l.emit(token.Comma, start)
	case ';':
		return l.emit(token.Semicolon, start)
	case ':':
		return l.emit(token.Colon, start)
	case '.':
		return l.emit(token.Dot, start)
	}

	// Should be unreachable as long as isPunctuation guards this method.
	panic("unreachable: lexPunctuation called with non-punctuation rune")
}
