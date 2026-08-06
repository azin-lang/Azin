package lexer

import (
	"unicode"

	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/text"
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
			// Save our position before consuming the slash
			slashCheckpoint := l.reader.Checkpoint()
			l.reader.Next()

			nextCh, _ := l.reader.Peek()

			if nextCh == '/' {
				l.reader.Next() // Consume second '/'
				l.scanSingleLineComment()
				continue
			}

			if nextCh == '*' {
				l.reader.Next() // Consume '*'
				// Pass the checkpoint as the exact start of the "/*"
				l.scanMultiLineComment(slashCheckpoint)
				continue
			}

			// False alarm: it's a division operator, not a comment.
			// Undo the slash consumption and stop scanning trivia.
			l.reader.Restore(slashCheckpoint)
			break
		}

		// If it's neither whitespace nor a comment, we are done with trivia.
		break
	}

	endOffset := l.reader.Offset()
	if startOffset == endOffset {
		return nil
	}

	triviaText := l.file.Text(text.SpanFromBounds(startOffset, endOffset))
	return green.NewTrivia(syntax.WhitespaceTrivia, triviaText)
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
				l.reader.Next() // Consume the closing '/'
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
