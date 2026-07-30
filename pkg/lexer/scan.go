package lexer

import (
	"unicode/utf8"

	"github.com/azin-lang/Azin/pkg/token"
)

// eofToken returns an end-of-file token at the current cursor position.
func (l *Lexer) eofToken() token.Token {
	return token.Token{
		Kind:     token.EOF,
		Position: l.pos(),
	}
}

// eof reports whether the cursor has reached or exceeded the end of the source buffer.
func (l *Lexer) eof() bool {
	return int(l.cursor) >= len(l.src)
}

// peek returns the next rune without advancing the cursor.
func (l *Lexer) peek() rune {
	if l.eof() {
		return 0
	}
	c := l.src[l.cursor]
	if c < utf8.RuneSelf {
		return rune(c)
	}
	r, _ := utf8.DecodeRune(l.src[l.cursor:])
	return r
}

// peekNext returns the rune following the current peeked rune without advancing the cursor.
func (l *Lexer) peekNext() rune {
	if l.eof() {
		return 0
	}
	c := l.src[l.cursor]
	nextCursor := l.cursor + 1
	if c >= utf8.RuneSelf {
		_, sz := utf8.DecodeRune(l.src[l.cursor:])
		nextCursor = l.cursor + uint32(sz)
	}

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

// advance reads the next rune, increments the cursor, and returns the rune and its byte size.
func (l *Lexer) advance() rune {
	if l.eof() {
		return 0
	}
	c := l.src[l.cursor]
	if c < utf8.RuneSelf {
		l.cursor++
		return rune(c)
	}
	r, sz := utf8.DecodeRune(l.src[l.cursor:])
	l.cursor += uint32(sz)
	return r
}

// match reports whether the next rune matches ch, advancing the cursor if it does.
func (l *Lexer) match(ch rune) bool {
	if l.peek() != ch {
		return false
	}
	// We already know we have room, so manually advancing is faster
	if ch < utf8.RuneSelf {
		l.cursor++
	} else {
		_, sz := utf8.DecodeRune(l.src[l.cursor:])
		l.cursor += uint32(sz)
	}
	return true
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
