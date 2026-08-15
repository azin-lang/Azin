package parser

import (
	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/syntax/lexer"
	"github.com/azin-lang/Azin/internal/azin/text"
)

// Parser consumes green tokens and builds a green syntax tree.
type Parser struct {
	lex     *lexer.Lexer
	diags   *diagnostics.Collector
	current *green.Token
	pos     int
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

// ParseCompilationUnit parses an entire source file as a sequence of statements until EOF.
func (p *Parser) ParseCompilationUnit() green.Node {
	var stmts []green.Node
	for !p.at(syntax.EndOfFileToken) {
		stmts = append(stmts, p.ParseStatement())
	}
	p.advance()
	return green.NewSyntaxList(stmts)
}

func (p *Parser) advance() *green.Token {
	previous := p.current
	if p.current != nil {
		p.pos += int(p.current.FullWidth())
	}
	p.current = p.lex.NextToken()
	return previous
}

// at reports whether the current token matches the given kind.
func (p *Parser) at(kind syntax.Kind) bool {
	return p.current.Kind() == kind
}

// atAny reports whether the current token matches any of the given kinds.
func (p *Parser) atAny(kinds ...syntax.Kind) bool {
	for _, k := range kinds {
		if p.current.Kind() == k {
			return true
		}
	}
	return false
}

// eat advances the parser if the current token matches kind and reports true.
func (p *Parser) eat(kind syntax.Kind) bool {
	if p.at(kind) {
		p.advance()
		return true
	}
	return false
}

// isAtExpressionStart reports whether the current token can begin an expression.
func (p *Parser) isAtExpressionStart() bool {
	switch p.current.Kind() {
	case syntax.IntegerLiteralToken,
		syntax.FloatLiteralToken,
		syntax.StringLiteralToken,
		syntax.KeywordTrue,
		syntax.KeywordFalse,
		syntax.IdentifierToken,
		syntax.OpenParenToken:
		return true
	default:
		return p.current.Kind().UnaryPrecedence() != 0
	}
}

// isDelimiterOrEOF reports whether the current token is a structural boundary or EOF.
func (p *Parser) isDelimiterOrEOF() bool {
	return p.atAny(syntax.CloseParenToken, syntax.CommaToken, syntax.EndOfFileToken)
}

// syncTo skips tokens until p.current matches one of the target boundary kinds or EOF.
func (p *Parser) syncTo(kinds ...syntax.Kind) {
	for !p.at(syntax.EndOfFileToken) {
		if p.atAny(kinds...) {
			return
		}
		p.advance()
	}
}

// currentLocation returns the source location of the current token.
func (p *Parser) currentLocation() text.Location {
	width := 0
	if p.current != nil {
		width = int(p.current.FullWidth())
	}

	return text.Location{
		File: p.lex.File(),
		Span: text.NewSpan(p.pos, width),
	}
}

// match consumes and returns the current token if it matches the expected kind.
// If it mismatches, it logs a diagnostic with location and returns a MissingToken.
func (p *Parser) match(kind syntax.Kind) *green.Token {
	if p.at(kind) {
		return p.advance()
	}

	p.diags.Error(
		p.currentLocation(),
		"E010",
		"Expected %s but got %s",
		kind.String(),
		p.current.Kind().String(),
	)

	return green.NewToken(syntax.MissingToken, "", nil, nil)
}
