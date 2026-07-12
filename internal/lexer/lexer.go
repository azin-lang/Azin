package lexer

import (
	"iter"
	"slices"

	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

// Lexer performs lexical analysis on a source file, transforming raw text
// into a stream of syntax tokens.
type Lexer struct {
	file   *source.File
	cursor uint32
	diag   *diagnostics.Engine
}

// New initializes a new Lexer for the given source file.
func New(file *source.File, diag *diagnostics.Engine) *Lexer {
	return &Lexer{
		file: file,
		diag: diag,
	}
}

// Tokenize eagerly scans the entire file and returns a slice of all tokens.
func (l *Lexer) Tokenize() []token.Token {
	return slices.Collect(l.Tokens())
}

// Tokens returns an iterator over the tokens in the source file.
// It yields tokens lazily until the end of the file is reached.
func (l *Lexer) Tokens() iter.Seq[token.Token] {
	return func(yield func(token.Token) bool) {
		for {
			tok := l.nextToken()

			if !yield(tok) {
				return
			}

			if tok.Kind == token.EOF {
				return
			}
		}
	}
}
