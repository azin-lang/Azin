package lexer

import "github.com/azin-lang/azin/compiler/internal/token"

func (l *Lexer) eofToken() token.Token {
	return token.Token{
		Kind:     token.EOF,
		Position: l.pos(),
	}
}

func (l *Lexer) eof() bool {
	return l.file.EOF(l.cursor)
}

func (l *Lexer) peek() rune {
	if l.eof() {
		return 0
	}
	r, _ := l.file.Rune(l.cursor)
	return r
}

func (l *Lexer) peekNext() rune {
	if l.eof() {
		return 0
	}

	_, size := l.file.Rune(l.cursor)
	nextCursor := l.cursor + size

	if l.file.EOF(nextCursor) {
		return 0
	}

	nextRune, _ := l.file.Rune(nextCursor)
	return nextRune
}

func (l *Lexer) peekString(s string) bool {
	offset := l.cursor
	for _, r := range s {
		if l.file.EOF(offset) {
			return false
		}
		ch, size := l.file.Rune(offset)
		if ch != r {
			return false
		}
		offset += size
	}
	return true
}

func (l *Lexer) advance() (rune, uint32) {
	if l.eof() {
		return 0, 0
	}
	r, size := l.file.Rune(l.cursor)
	l.cursor += size
	return r, size
}

func (l *Lexer) match(ch rune) bool {
	if l.peek() != ch {
		return false
	}
	_, _ = l.advance()
	return true
}

func (l *Lexer) matchString(s string) bool {
	start := l.pos()

	for _, r := range s {
		if l.peek() != r {
			l.rewind(start)
			return false
		}
		_, _ = l.advance()
	}
	return true
}

func (l *Lexer) matchAny(chars string) bool {
	r := l.peek()
	for _, ch := range chars {
		if r == ch {
			_, _ = l.advance()
			return true
		}
	}
	return false
}

func (l *Lexer) rewind(pos token.Position) {
	l.cursor = pos.Offset
}

func (l *Lexer) consumeWhile(pred func(rune) bool) {
	for pred(l.peek()) {
		_, _ = l.advance()
	}
}

func (l *Lexer) emit(kind token.Kind, start token.Position) token.Token {
	return token.Token{
		Kind:     kind,
		Position: start,
		Length:   l.cursor - start.Offset,
	}
}

func (l *Lexer) either(ch rune, ifMatch token.Kind, otherwise token.Kind, start token.Position) token.Token {
	if l.match(ch) {
		return l.emit(ifMatch, start)
	}
	return l.emit(otherwise, start)
}

func (l *Lexer) pos() token.Position {
	return token.Position{Offset: l.cursor}
}
