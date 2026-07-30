package lexer

import (
	"github.com/azin-lang/Azin/pkg/token"
)

// skipTrivia consumes whitespace and comments, advancing the cursor
// to the next meaningful syntax token.
func (l *Lexer) skipTrivia() {
	length := len(l.src)
	for int(l.cursor) < length {
		c := l.src[l.cursor]

		if c == ' ' || c == '\t' {
			l.cursor++
			continue
		}

		if c == '/' && int(l.cursor)+1 < length {
			next := l.src[l.cursor+1]
			if next == '/' {
				l.cursor += 2
				l.skipLineComment()
				continue
			}
			if next == '*' {
				start := l.pos()
				l.cursor += 2
				l.skipBlockComment(start)
				continue
			}
		}
		return
	}
}

// skipLineComment consumes characters until the end of the line.
func (l *Lexer) skipLineComment() {
	length := len(l.src)
	for int(l.cursor) < length {
		c := l.src[l.cursor]
		if c == '\n' || c == '\r' {
			return
		}
		l.cursor++
	}
}

// skipBlockComment consumes characters until the matching closing delimiter.
func (l *Lexer) skipBlockComment(start token.Position) {
	depth := 1
	length := len(l.src)

	for int(l.cursor)+1 < length {
		c := l.src[l.cursor]
		if c == '*' && l.src[l.cursor+1] == '/' {
			depth--
			l.cursor += 2
			if depth == 0 {
				return
			}
			continue
		}
		if c == '/' && l.src[l.cursor+1] == '*' {
			depth++
			l.cursor += 2
			continue
		}
		l.cursor++
	}

	l.cursor = uint32(length)
	l.diag.ReportError(start, int(l.cursor-start.Offset), "unterminated block comment")
}
