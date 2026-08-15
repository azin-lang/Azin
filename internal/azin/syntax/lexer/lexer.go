package lexer

import (
	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/text"
)

// Lexer breaks down source code into syntax tokens.
type Lexer struct {
	file        *text.SourceText
	reader      *text.Reader
	diagnostics *diagnostics.Collector
}

// New creates and initializes a new Lexer.
func New(file *text.SourceText, diags *diagnostics.Collector) *Lexer {
	return &Lexer{
		file:        file,
		reader:      text.NewReader(file),
		diagnostics: diags,
	}
}

// File returns the underlying source text buffer.
func (l *Lexer) File() *text.SourceText {
	return l.file
}

// NextToken is the main entry point called by the Parser.
func (l *Lexer) NextToken() *green.Token {
	leadingTrivia := l.scanTrivia(false)

	if l.reader.EOF() {
		return green.NewToken(syntax.EndOfFileToken, "", leadingTrivia, nil)
	}

	startOffset := l.reader.Offset()
	kind := l.scanSyntaxToken(startOffset)
	endOffset := l.reader.Offset()

	source := l.file.Text(text.SpanFromBounds(startOffset, endOffset))
	trailingTrivia := l.scanTrivia(true)

	return green.NewToken(kind, source, leadingTrivia, trailingTrivia)
}

// match consumes the expected character if it matches the current peek.
func (l *Lexer) match(expected rune) bool {
	if l.reader.EOF() {
		return false
	}
	ch, _ := l.reader.Peek()
	if ch == expected {
		l.reader.Next()
		return true
	}
	return false
}

// loc is a helper to clean up span/location boilerplate.
func (l *Lexer) loc(startOffset, endOffset int) text.Location {
	return text.Location{
		File: l.file,
		Span: text.SpanFromBounds(startOffset, endOffset),
	}
}
