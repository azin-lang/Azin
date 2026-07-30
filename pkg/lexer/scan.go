package lexer

import (
	"unicode/utf8"

	token "github.com/azin-lang/Azin/pkg/token"
)

func (l *Lexer) eofToken() token.Token {
	return token.Token{
		Kind:     token.EOF,
		Position: l.pos(),
	}
}

func (l *Lexer) eof() bool {
	return int(l.cursor) >= len(l.src)
}

func (l *Lexer) peek() rune {
	if l.eof() {
		return 0
	}

	// Fast path for ASCII
	c := l.src[l.cursor]
	if c < utf8.RuneSelf {
		return rune(c)
	}

	r, _ := utf8.DecodeRune(l.src[l.cursor:])
	return r
}

func (l *Lexer) peekNext() rune {
	if l.eof() {
		return 0
	}

	c := l.src[l.cursor]
	var size uint32 = 1

	if c >= utf8.RuneSelf {
		_, sz := utf8.DecodeRune(l.src[l.cursor:])
		size = uint32(sz)
	}

	nextCursor := l.cursor + size
	if int(nextCursor) >= len(l.src) {
		return 0
	}

	nextC := l.src[nextCursor]
	if nextC < utf8.RuneSelf {
		return rune(nextC)
	}

	nextRune, _ := utf8.DecodeRune(l.src[nextCursor:])
	return nextRune
}

func (l *Lexer) advance() (r rune, size uint32) {
	if l.eof() {
		return 0, 0
	}

	c := l.src[l.cursor]
	if c < utf8.RuneSelf {
		l.cursor++
		return rune(c), 1
	}

	r, sz := utf8.DecodeRune(l.src[l.cursor:])
	l.cursor += uint32(sz)
	return r, uint32(sz)
}

func (l *Lexer) match(ch rune) bool {
	if l.peek() != ch {
		return false
	}
	_, _ = l.advance()
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

func (l *Lexer) either(ch rune, ifMatch, otherwise token.Kind, start token.Position) token.Token {
	if l.match(ch) {
		return l.emit(ifMatch, start)
	}
	return l.emit(otherwise, start)
}

func (l *Lexer) pos() token.Position {
	return token.Position{Offset: l.cursor}
}
