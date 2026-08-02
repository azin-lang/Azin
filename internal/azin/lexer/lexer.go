package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/source"
	"github.com/azin-lang/Azin/internal/azin/syntax"
)

type Lexer struct {
	file      *source.File
	reader    *source.Reader
	collector *diagnostics.Collector
}

func New(file *source.File, collector *diagnostics.Collector) *Lexer {
	return &Lexer{
		file:      file,
		reader:    source.NewReader(file),
		collector: collector,
	}
}

func (l *Lexer) NextToken() syntax.Token {
	leading := l.lexLeadingTrivia()

	if l.reader.EOF() {
		end := l.file.Len()
		return syntax.Token{
			Kind:          syntax.EndOfFile,
			Span:          source.NewSpan(end, end),
			LeadingTrivia: leading,
		}
	}

	start := l.reader.Offset()
	ch, _ := l.reader.Next()

	var tokenKind syntax.Kind

	switch {
	case (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_':
		tokenKind = l.lexIdentifierOrKeyword(start)
	case ch >= '0' && ch <= '9':
		tokenKind = syntax.IntegerLiteralToken
		for {
			c, _ := l.reader.Peek()
			if (c >= '0' && c <= '9') || (c >= utf8.RuneSelf && unicode.IsDigit(c)) {
				l.reader.Next()
			} else {
				break
			}
		}
	case ch == '=':
		if l.match('=') {
			tokenKind = syntax.EqualsEqualsToken
		} else {
			tokenKind = syntax.EqualsToken
		}
	case ch == '+':
		if l.match('=') {
			tokenKind = syntax.PlusEqualsToken
		} else {
			tokenKind = syntax.PlusToken
		}
	case ch == '(':
		tokenKind = syntax.OpenParenToken
	case ch == ')':
		tokenKind = syntax.CloseParenToken
	case ch == ';':
		tokenKind = syntax.SemicolonToken
	default:
		for !l.reader.EOF() {
			c, _ := l.reader.Peek()
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' ||
				(c >= '0' && c <= '9') ||
				c == '=' || c == '+' || c == '(' || c == ')' || c == ';' ||
				c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '/' ||
				(c >= utf8.RuneSelf && unicode.IsLetter(c)) {
				break
			}
			l.reader.Next()
		}
		span := source.NewSpan(start, l.reader.Offset())

		l.collector.ErrorWithNoteAndHelp(
			source.Location{File: l.file, Span: span},
			"AZ0001",
			"unexpected curly braces",
			"try replacing with 'end'",
			"Azin uses 'end' to close blocks, not curly braces.",
			"unexpected token sequence",
		)
		tokenKind = syntax.Error
	}

	span := source.NewSpan(start, l.reader.Offset())
	trailing := l.lexTrailingTrivia()

	return syntax.Token{
		Kind:           tokenKind,
		Span:           span,
		LeadingTrivia:  leading,
		TrailingTrivia: trailing,
	}
}

func (l *Lexer) match(ch rune) bool {
	r, _ := l.reader.Peek()
	if r != ch {
		return false
	}
	l.reader.Next()
	return true
}

func (l *Lexer) peekNextRune() rune {
	currOffset := l.reader.Offset()
	_, sz := l.file.Rune(currOffset)
	if sz == 0 {
		return 0
	}
	nextR, _ := l.file.Rune(currOffset + sz)
	return nextR
}

func (l *Lexer) lexLeadingTrivia() []syntax.Trivia {
	trivia := make([]syntax.Trivia, 0, 4)
	for !l.reader.EOF() {
		ch, _ := l.reader.Peek()
		start := l.reader.Offset()

		if ch == ' ' || ch == '\t' {
			// Group contiguous spaces and tabs into a single trivia chunk
			for !l.reader.EOF() {
				c, _ := l.reader.Peek()
				if c == ' ' || c == '\t' {
					l.reader.Next()
				} else {
					break
				}
			}
			trivia = append(trivia, syntax.Trivia{
				Kind: syntax.WhitespaceTrivia,
				Span: source.NewSpan(start, l.reader.Offset()),
			})
		} else if ch == '\n' || ch == '\r' {
			l.reader.Next()
			if ch == '\r' {
				if r, _ := l.reader.Peek(); r == '\n' {
					l.reader.Next()
				}
			}
			trivia = append(trivia, syntax.Trivia{
				Kind: syntax.EndOfLineTrivia,
				Span: source.NewSpan(start, l.reader.Offset()),
			})
		} else if ch == '/' && l.peekNextRune() == '/' {
			l.reader.Next()
			l.reader.Next()
			for !l.reader.EOF() {
				c, _ := l.reader.Next()
				if c == '\n' || c == '\r' {
					break
				}
			}
			trivia = append(trivia, syntax.Trivia{
				Kind: syntax.SingleLineCommentTrivia,
				Span: source.NewSpan(start, l.reader.Offset()),
			})
		} else {
			break
		}
	}
	return trivia
}

func (l *Lexer) lexTrailingTrivia() []syntax.Trivia {
	trivia := make([]syntax.Trivia, 0, 2)
	for !l.reader.EOF() {
		ch, _ := l.reader.Peek()
		start := l.reader.Offset()

		if ch == ' ' || ch == '\t' {
			// Group contiguous spaces and tabs into a single trivia chunk
			for !l.reader.EOF() {
				c, _ := l.reader.Peek()
				if c == ' ' || c == '\t' {
					l.reader.Next()
				} else {
					break
				}
			}
			trivia = append(trivia, syntax.Trivia{
				Kind: syntax.WhitespaceTrivia,
				Span: source.NewSpan(start, l.reader.Offset()),
			})
		} else if ch == '/' && l.peekNextRune() == '/' {
			l.reader.Next()
			l.reader.Next()
			for !l.reader.EOF() {
				c, _ := l.reader.Next()
				if c == '\n' || c == '\r' {
					break
				}
			}
			trivia = append(trivia, syntax.Trivia{
				Kind: syntax.SingleLineCommentTrivia,
				Span: source.NewSpan(start, l.reader.Offset()),
			})
			break
		} else {
			break
		}
	}
	return trivia
}

func (l *Lexer) lexIdentifierOrKeyword(start uint32) syntax.Kind {
	for {
		ch, _ := l.reader.Peek()
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
			l.reader.Next()
		} else if ch >= utf8.RuneSelf && (unicode.IsLetter(ch) || unicode.IsDigit(ch)) {
			l.reader.Next()
		} else {
			break
		}
	}
	text := l.file.Text(source.NewSpan(start, l.reader.Offset()))
	return syntax.LookupKeyword(text)
}
