package lexer

import (
	"unicode"

	"github.com/azin-lang/Azin/internal/azin/source"
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
)

func (l *Lexer) scanTrivia(isTrailing bool) green.Node {
	startOffset := l.reader.Offset()

	for !l.reader.EOF() {
		ch, _ := l.reader.Peek()

		if isTrailing && (ch == '\n' || ch == '\r') {
			break
		}

		if unicode.IsSpace(ch) {
			l.reader.Next()
			continue
		}

		if ch == '/' {
			l.reader.Next()
			nextCh, _ := l.reader.Peek()

			if nextCh == '/' {
				l.reader.Next()
				l.scanSingleLineComment()
				continue
			}

			if nextCh == '*' {
				l.reader.Next()
				l.scanMultiLineComment(l.reader.Offset() - 2)
				continue
			}

			l.reader.Backup()
			break
		}
		break
	}

	endOffset := l.reader.Offset()
	if startOffset == endOffset {
		return nil
	}

	text := l.file.Text(source.NewSpan(startOffset, endOffset))
	return green.NewTrivia(syntax.WhitespaceTrivia, text)
}

func (l *Lexer) scanSingleLineComment() {
	for !l.reader.EOF() {
		ch, _ := l.reader.Peek()
		if ch == '\n' || ch == '\r' {
			break
		}
		l.reader.Next()
	}
}

func (l *Lexer) scanMultiLineComment(startOffset uint32) {
	for !l.reader.EOF() {
		ch, _ := l.reader.Next()
		if ch == '*' {
			nextCh, _ := l.reader.Peek()
			if nextCh == '/' {
				l.reader.Next()
				return
			}
		}
	}

	l.diagnostics.ErrorWithNoteAndHelp(
		l.loc(startOffset, l.reader.Offset()),
		"E003",
		"I started reading a comment here...",
		"But I reached the end of the file before I could find the closing '*/'.",
		"Try adding '*/' wherever you'd like this comment to end!",
		"Unfinished comment",
	)
}
