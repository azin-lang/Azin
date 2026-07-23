package parser

import (
	"slices"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

type ErrorReporter interface {
	ReportError(pos token.Position, length uint32, format string, args ...any)
	Err() error
}

const (
	_ int = iota
	PrecLowest
	PrecEquality   // ==, !=
	PrecComparison // <, >, <=, >=
	PrecTerm       // +, -
	PrecFactor     // *, /
	PrecCall       // (
	PrecMember     // .
)

type Parser struct {
	source  string
	tokens  []token.Token
	current int
	diag    ErrorReporter
}

func Parse(source string, tokens []token.Token, diag ErrorReporter) (*ast.Program, error) {
	p := New(source, tokens, diag)
	return p.ParseProgram(), p.diag.Err()
}

func New(source string, tokens []token.Token, diag ErrorReporter) *Parser {
	return &Parser{
		source: source,
		tokens: tokens,
		diag:   diag,
	}
}

func (p *Parser) Err() error {
	return p.diag.Err()
}

func (p *Parser) synchronize() {
	if p.isAtEnd() {
		return
	}
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind == token.Newline || p.previous().Kind == token.Semicolon {
			return
		}

		if isSyncPoint(p.peek().Kind) {
			return
		}
		p.advance()
	}
}

func (p *Parser) lexeme(tok token.Token) string {
	start := int(tok.Position.Offset)
	end := start + int(tok.Length)
	if start < 0 {
		return ""
	}
	if end > len(p.source) {
		end = len(p.source)
	}
	if start > end {
		return ""
	}
	return p.source[start:end]
}

func (p *Parser) reportError(tok token.Token, format string, args ...any) {
	p.diag.ReportError(tok.Position, tok.Length, format, args...)
}

func badStmt(tok token.Token) *ast.BadStmt {
	return &ast.BadStmt{Token: tok}
}

func badExpr(tok token.Token) *ast.BadExpr {
	return &ast.BadExpr{Token: tok}
}

func isBadStmt(stmt ast.Stmt) bool {
	_, ok := stmt.(*ast.BadStmt)
	return ok
}

func isBadExpr(expr ast.Expr) bool {
	_, ok := expr.(*ast.BadExpr)
	return ok
}

func (p *Parser) expect(kind token.Kind, context string) (token.Token, bool) {
	if p.check(kind) {
		return p.advance(), true
	}

	got := p.peek()
	p.reportError(
		got,
		"expected %s %s, found %s",
		kind.DisplayName(),
		context,
		got.Kind.DisplayName(),
	)

	return got, false
}

func (p *Parser) parseBlock(until ...token.Kind) []ast.Stmt {
	var body []ast.Stmt
	for {
		p.skipNewlines()
		if p.isAtEnd() || p.checkAny(until...) {
			break
		}
		if stmt := p.parseStatement(); stmt != nil {
			body = append(body, stmt)
		}
	}
	return body
}

func (p *Parser) skipNewlines() {
	for p.match(token.Newline) {
	}
}

func (p *Parser) peek() token.Token {
	if p.current >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	if p.current == 0 {
		return p.tokens[0]
	}
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	if p.current >= len(p.tokens) {
		return true
	}
	return p.tokens[p.current].Kind == token.EOF
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(kind token.Kind) bool {
	return p.peek().Kind == kind
}

func (p *Parser) checkAny(kinds ...token.Kind) bool {
	if p.isAtEnd() {
		return false
	}
	return slices.Contains(kinds, p.peek().Kind)
}

func (p *Parser) match(kinds ...token.Kind) bool {
	if p.isAtEnd() {
		return false
	}
	if slices.Contains(kinds, p.peek().Kind) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) consumeStatementEnd() bool {
	switch p.peek().Kind {
	case token.Semicolon:
		p.advance()
		return true
	case token.Newline:
		p.skipNewlines()
		return true
	case token.EOF, token.KwEnd, token.KwElse:
		return true
	default:
		p.reportError(p.peek(), "expected end of statement (newline or ';')")
		return false
	}
}

func isBuiltinType(kind token.Kind) bool {
	switch kind {
	case token.KwUnit, token.KwInt, token.KwFloat, token.KwString, token.KwChar, token.KwBool:
		return true
	default:
		return false
	}
}

func isSyncPoint(kind token.Kind) bool {
	switch kind {
	case token.KwFn, token.KwStruct, token.KwEnum, token.KwVar, token.KwIf, token.KwLoop,
		token.KwReturn, token.KwElse, token.KwImportC, token.KwEnd:
		return true
	default:
		return false
	}
}

func getPrecedence(kind token.Kind) int {
	switch kind {
	case token.LeftParen:
		return PrecCall
	case token.Dot:
		return PrecMember
	case token.Star, token.Slash:
		return PrecFactor
	case token.Plus, token.Minus:
		return PrecTerm
	case token.Less, token.LessEqual, token.Greater, token.GreaterEqual:
		return PrecComparison
	case token.EqualEqual, token.BangEqual:
		return PrecEquality
	default:
		return PrecLowest
	}
}
