package lexer

import (
	"iter"
	"slices"

	"github.com/azin-lang/Azin/pkg/diagnostics"
	"github.com/azin-lang/Azin/pkg/source"
	token2 "github.com/azin-lang/Azin/pkg/token"
)

// Lexer performs lexical analysis on a source file, transforming raw text
// into a stream of syntax tokens.
type Lexer struct {
	file   *source.File
	src    []byte
	cursor uint32
	diag   *diagnostics.Engine
}

// New returns a new Lexer configured for the given source file.
func New(file *source.File, diag *diagnostics.Engine) *Lexer {
	return &Lexer{
		file:   file,
		src:    file.Bytes(),
		cursor: 0,
		diag:   diag,
	}
}

// Tokenize eagerly scans the entire file and returns a slice of all parsed tokens.
func (l *Lexer) Tokenize() []token2.Token {
	return slices.Collect(l.Tokens())
}

// Tokens returns an iterator over the tokens in the source file.
// It yields tokens lazily until the end of the file is reached.
func (l *Lexer) Tokens() iter.Seq[token2.Token] {
	return func(yield func(token2.Token) bool) {
		for {
			tok := l.nextToken()

			if !yield(tok) {
				return
			}

			if tok.Kind == token2.EOF {
				return
			}
		}
	}
}
