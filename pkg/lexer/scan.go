package lexer

import (
	"unicode/utf8"

	token "github.com/azin-lang/Azin/pkg/token"
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

	// Fast path for ASCII
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

// advance reads the next rune, increments the cursor, and returns the rune and its byte size.
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

// match reports whether the next rune matches ch, advancing the cursor if it does.
func (l *Lexer) match(ch rune) bool {
	if l.peek() != ch {
		return false
	}
	_, _ = l.advance()
	return true
}

// matchAny reports whether the next rune matches any character in the given string, advancing the cursor if it does.
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

// consumeWhile advances the cursor as long as the given predicate returns true.
func (l *Lexer) consumeWhile(pred func(rune) bool) {
	for pred(l.peek()) {
		_, _ = l.advance()
	}
}

// emit returns a token of the specified kind starting from the given position up to the current cursor.
func (l *Lexer) emit(kind token.Kind, start token.Position) token.Token {
	return token.Token{
		Kind:     kind,
		Position: start,
		Length:   l.cursor - start.Offset,
	}
}

// either returns a token of kind ifMatch if the next character matches ch (advancing the cursor); otherwise, it returns a token of kind otherwise.
func (l *Lexer) either(ch rune, ifMatch, otherwise token.Kind, start token.Position) token.Token {
	if l.match(ch) {
		return l.emit(ifMatch, start)
	}
	return l.emit(otherwise, start)
}

// pos returns the current position offset of the cursor.
func (l *Lexer) pos() token.Position {
	return token.Position{Offset: l.cursor}
}
