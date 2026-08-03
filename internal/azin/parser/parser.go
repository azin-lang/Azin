package parser

import (
	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/lexer"
	"github.com/azin-lang/Azin/internal/azin/source"
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
)

// Parser consumes green tokens and builds a green syntax tree.
type Parser struct {
	lex     *lexer.Lexer
	diags   *diagnostics.Collector
	current *green.Token
	pos     uint32
}

func New(lex *lexer.Lexer, diags *diagnostics.Collector) *Parser {
	p := &Parser{
		lex:   lex,
		diags: diags,
		pos:   0,
	}
	// Prime the pump by fetching the first token
	p.advance()
	return p
}

func (p *Parser) advance() *green.Token {
	previous := p.current
	if p.current != nil {
		p.pos += p.current.FullWidth()
	}
	p.current = p.lex.NextToken()
	return previous
}

// syncTo skips tokens until p.current matches one of the target boundary kinds or EOF.
func (p *Parser) syncTo(kinds ...syntax.SyntaxKind) {
	for p.current.Kind() != syntax.EndOfFileToken {
		for _, k := range kinds {
			if p.current.Kind() == k {
				return
			}
		}
		p.advance()
	}
}

// currentLocation returns the source location of the current token.
func (p *Parser) currentLocation() source.Location {
	width := uint32(0)
	if p.current != nil {
		width = p.current.FullWidth()
	}

	span := source.NewSpan(p.pos, p.pos+width)

	return source.Location{
		File: p.lex.File(),
		Span: span,
	}
}

// match consumes and returns the current token if it matches the expected kind.
// If it mismatches, it logs a diagnostic with location and returns a MissingToken.
func (p *Parser) match(kind syntax.SyntaxKind) *green.Token {
	if p.current.Kind() == kind {
		return p.advance()
	}

	p.diags.Error(
		p.currentLocation(),
		"E010",
		"Expected %s but got %s",
		syntax.Display(kind),
		syntax.Display(p.current.Kind()),
	)

	return green.NewToken(syntax.MissingToken, "", nil, nil)
}
