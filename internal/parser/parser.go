package parser

import (
	"slices"
	"strconv"
	"strings"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/token"
)

type Parser struct {
	source  string
	tokens  []token.Token
	current int

	diag *diagnostics.Engine
}

func Parse(source string, tokens []token.Token, diag *diagnostics.Engine) (*ast.Program, error) {
	p := New(source, tokens, diag)
	return p.ParseProgram(), p.diag.Err()
}

func New(source string, tokens []token.Token, diag *diagnostics.Engine) *Parser {
	return &Parser{
		source: source,
		tokens: tokens,
		diag:   diag,
	}
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

		switch p.peek().Kind {
		case token.KwFn, token.KwStruct, token.KwVar, token.KwIf,
			token.KwReturn, token.KwElse, token.KwImportC, token.KwEnd:
			return
		default:
			p.advance()
		}
	}
}

func (p *Parser) lexeme(tok token.Token) string {
	return p.source[tok.Position.Offset : tok.Position.Offset+tok.Length]
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

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Stmt{},
	}

	for !p.isAtEnd() {
		before := p.current
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		// If the parser didn't advance, force it to avoid infinite loops
		if p.current == before {
			p.advance()
			p.synchronize()
		}
	}

	return program
}

func (p *Parser) Err() error {
	return p.diag.Err()
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

func (p *Parser) parseVar() ast.Stmt {
	tok := p.advance()
	mutable := p.match(token.KwMut)

	name := p.parseIdentifier()
	if name == nil {
		return badStmt(tok)
	}

	var typ *ast.Identifier
	if p.match(token.Colon) {
		typ = p.parseType()
	}

	var value ast.Expr
	if p.match(token.Equal) {
		value = p.parseExpression(0)
		if value == nil {
			return badStmt(tok)
		}
	} else if typ == nil {
		p.reportError(p.peek(), "expected '=' or an explicit type in variable declaration")
		return badStmt(tok)
	}

	if !p.statementEnd() {
		return badStmt(tok)
	}

	return &ast.VarStmt{
		Token:   tok,
		Name:    name,
		Type:    typ,
		Value:   value,
		Mutable: mutable,
	}
}

func (p *Parser) parseImportC() ast.Stmt {
	tok := p.advance()

	if !p.check(token.StringLiteral) {
		p.reportError(p.peek(), "expected string literal after 'importC'")
		return badStmt(tok)
	}

	path := p.parseStringLiteral()

	if !p.statementEnd() {
		return badStmt(tok)
	}

	return &ast.ImportCStmt{
		Token: tok,
		Path:  path,
	}
}

func (p *Parser) statementEnd() bool {
	if p.match(token.Semicolon) {
		return true
	}

	if p.match(token.Newline) {
		p.skipNewlines()
		return true
	}

	if p.checkAny(token.EOF, token.KwEnd, token.KwElse) {
		return true
	}

	p.reportError(p.peek(), "expected end of statement (newline or ';')")
	return false
}

func (p *Parser) parseStatement() ast.Stmt {
	p.skipNewlines()

	var stmt ast.Stmt
	switch {
	case p.check(token.KwVar):
		stmt = p.parseVar()
	case p.check(token.KwStruct):
		stmt = p.parseStruct()

	case p.check(token.KwFn):
		stmt = p.parseFunc()

	case p.check(token.KwReturn):
		stmt = p.parseReturn()

	case p.check(token.KwIf):
		stmt = p.parseIf()

	case p.check(token.KwImportC):
		stmt = p.parseImportC()

	default:
		stmt = p.parseExpressionOrAssignment()
	}

	if stmt == nil {
		p.synchronize()
		return badStmt(p.peek())
	}

	if _, ok := stmt.(*ast.BadStmt); ok {
		p.synchronize()
	}
	return stmt
}

func (p *Parser) parseExpressionOrAssignment() ast.Stmt {
	expr := p.parseExpression(0)

	if _, ok := expr.(*ast.BadExpr); ok {
		return badStmt(p.peek())
	}

	if p.check(token.Equal) {
		tok := p.advance()

		switch expr.(type) {
		case *ast.Identifier, *ast.MemberExpr:
			// valid assignment target
		default:
			p.reportError(tok, "left side of assignment is not assignable")
			return badStmt(tok)
		}

		value := p.parseExpression(0)
		if _, ok := value.(*ast.BadExpr); ok {
			return badStmt(tok)
		}

		if !p.statementEnd() {
			return badStmt(tok)
		}

		return &ast.AssignmentStmt{
			Token: tok,
			Left:  expr,
			Value: value,
		}
	}

	if !p.statementEnd() {
		return badStmt(p.previous())
	}

	return &ast.ExpressionStmt{
		Token:      p.previous(),
		Expression: expr,
	}
}

func (p *Parser) parseStruct() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()
	if name == nil {
		return badStmt(tok)
	}
	if _, ok := p.expect(token.KwIs, "after struct name"); !ok {
		return badStmt(tok)
	}
	p.skipNewlines()

	fields := []*ast.FieldDecl{}

	for !p.isAtEnd() {
		p.skipNewlines()

		if p.check(token.KwEnd) {
			break
		}

		mutable := p.match(token.KwMut)
		fName := p.parseIdentifier()
		if fName == nil {
			return badStmt(tok)
		}

		if _, ok := p.expect(token.Colon, "after field name"); !ok {
			return badStmt(tok)
		}

		tName := p.parseType()
		if tName == nil {
			return badStmt(tok)
		}

		if !p.statementEnd() {
			return badStmt(tok)
		}

		fields = append(fields, &ast.FieldDecl{
			Name:    fName,
			Type:    tName,
			Mutable: mutable,
		})
	}

	if _, ok := p.expect(token.KwEnd, "to close struct"); !ok {
		return badStmt(tok)
	}

	return &ast.StructStmt{
		Token:  tok,
		Name:   name,
		Fields: fields,
	}
}

func (p *Parser) parseType() *ast.Identifier {
	switch p.peek().Kind {
	case token.KwUnit, token.KwInt, token.KwFloat, token.KwString, token.KwChar:
		tok := p.advance()
		return &ast.Identifier{Token: tok, Value: p.lexeme(tok)}
	default:
		return p.parseIdentifier()
	}
}

func (p *Parser) parseFunc() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()
	if name == nil {
		return badStmt(tok)
	}

	params := []*ast.FieldDecl{}

	if p.match(token.LeftParen) {
		for !p.isAtEnd() && !p.check(token.RightParen) {
			pName := p.parseIdentifier()
			if pName == nil {
				return badStmt(tok)
			}
			p.match(token.Colon)
			pType := p.parseType()

			if pType == nil {
				return badStmt(tok)
			}

			params = append(params, &ast.FieldDecl{
				Name: pName,
				Type: pType,
			})

			if !p.check(token.RightParen) {
				if _, ok := p.expect(token.Comma, "between parameters"); !ok {
					return badStmt(tok)
				}
			}
		}

		if _, ok := p.expect(token.RightParen, "after parameter list"); !ok {
			return badStmt(tok)
		}
	}
	var retType *ast.Identifier
	if p.match(token.Colon) {
		retType = p.parseType()
		if retType == nil {
			return badStmt(tok)
		}
	}

	if _, ok := p.expect(token.KwDo, "after function declaration"); !ok {
		return badStmt(tok)
	}

	p.skipNewlines()

	body := []ast.Stmt{}

	for {
		p.skipNewlines()

		if p.isAtEnd() || p.check(token.KwEnd) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}

	if _, ok := p.expect(token.KwEnd, "to close function"); !ok {
		return badStmt(tok)
	}

	return &ast.FuncStmt{Token: tok, Name: name, Params: params, ReturnType: retType, Body: body}
}

func (p *Parser) parseReturn() ast.Stmt {
	tok := p.advance()
	var value ast.Expr

	if !p.checkAny(token.Semicolon, token.Newline, token.KwEnd, token.EOF) {
		value = p.parseExpression(0)
	}

	if !p.statementEnd() {
		return badStmt(tok)
	}

	return &ast.ReturnStmt{
		Token: tok,
		Value: value,
	}
}

func (p *Parser) parseIf() ast.Stmt {
	tok := p.advance() // if

	condition := p.parseExpression(0)
	if _, ok := p.expect(token.KwThen, "after if condition"); !ok {
		return badStmt(tok)
	}

	p.skipNewlines()
	var thenBody []ast.Stmt

	for {
		p.skipNewlines()
		if p.isAtEnd() || p.checkAny(token.KwElse, token.KwEnd) {
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			thenBody = append(thenBody, stmt)
		}
	}

	var elseBody []ast.Stmt
	if p.match(token.KwElse) {
		p.skipNewlines()
		for {
			p.skipNewlines()
			if p.isAtEnd() || p.check(token.KwEnd) {
				break
			}
			stmt := p.parseStatement()
			if stmt != nil {
				elseBody = append(elseBody, stmt)
			}
		}
	}

	p.skipNewlines()
	if _, ok := p.expect(token.KwEnd, "to close if statement"); !ok {
		return badStmt(tok)
	}
	p.skipNewlines()

	return &ast.IfStmt{
		Token:     tok,
		Condition: condition,
		Then:      thenBody,
		Else:      elseBody,
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expr {
	// Parse prefix/left hand side
	left := p.parsePrefix()
	if left == nil {
		errTok := p.peek()
		p.reportError(errTok, "expected expression")
		return badExpr(errTok)
	}

	// Parse infix/right hand side based on binding power
	for !p.isAtEnd() {
		if p.checkAny(
			token.Newline, token.Semicolon, token.KwEnd,
			token.KwElse, token.KwThen, token.EOF,
		) {
			break
		}

		nextPrec := getPrecedence(p.peek().Kind)
		if nextPrec == 0 || precedence >= nextPrec {
			break
		}

		left = p.parseInfix(left, nextPrec)
	}

	return left
}

func (p *Parser) parsePrefix() ast.Expr {
	switch {
	case p.check(token.Identifier):
		id := p.parseIdentifier()
		if id == nil {
			return badExpr(p.peek())
		}
		return id
	case p.check(token.IntegerLiteral):
		return p.parseIntegerLiteral()
	case p.check(token.FloatLiteral):
		return p.parseFloatLiteral()
	case p.check(token.StringLiteral):
		return p.parseStringLiteral()
	case p.check(token.CharacterLiteral):
		return p.parseCharacterLiteral()
	default:
		return nil
	}
}

func (p *Parser) parseInfix(left ast.Expr, nextPrec int) ast.Expr {
	switch {
	case p.check(token.LeftParen):
		tok := p.advance() // Consume '('
		var args []ast.Expr
		for !p.check(token.RightParen) {
			args = append(args, p.parseExpression(0))
			if !p.check(token.RightParen) {
				if _, ok := p.expect(token.Comma, "between arguments"); !ok {
					return badExpr(tok)
				}
			}
		}
		if _, ok := p.expect(token.RightParen, "to close argument list"); !ok {
			return badExpr(tok)
		}
		return &ast.CallExpr{Callee: left, Args: args}

	case p.check(token.Dot):
		tok := p.advance() // Consume '.'
		prop := p.parseIdentifier()
		if prop == nil {
			return badExpr(tok)
		}
		return &ast.MemberExpr{Object: left, Property: prop}

	default:
		// Binary operators
		op := p.advance()
		right := p.parseExpression(nextPrec)
		return &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	tok, ok := p.expect(token.Identifier, "")

	if !ok {
		return nil
	}

	start := tok.Position.Offset
	end := start + tok.Length
	realValue := p.source[start:end]

	return &ast.Identifier{
		Token: tok,
		Value: realValue,
	}
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	tok := p.advance()
	text := p.lexeme(tok)

	var value int64
	switch {
	case strings.HasPrefix(text, "0x"):
		value, _ = strconv.ParseInt(text[2:], 16, 64)
	case strings.HasPrefix(text, "0b"):
		value, _ = strconv.ParseInt(text[2:], 2, 64)
	default:
		value, _ = strconv.ParseInt(text, 10, 64)
	}

	return &ast.IntegerLiteral{Token: tok, Value: value}
}

func (p *Parser) parseFloatLiteral() *ast.FloatLiteral {
	tok := p.advance()
	value, _ := strconv.ParseFloat(p.lexeme(tok), 64)
	return &ast.FloatLiteral{Token: tok, Value: value}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	tok := p.advance()
	value, err := strconv.Unquote(p.lexeme(tok))
	if err != nil {
		p.reportError(tok, "invalid string literal: %v", err)
		return &ast.StringLiteral{Token: tok, Value: ""}
	}
	return &ast.StringLiteral{Token: tok, Value: value}
}

func (p *Parser) parseCharacterLiteral() *ast.CharacterLiteral {
	tok := p.advance()
	raw := p.lexeme(tok)
	value, _, _, err := strconv.UnquoteChar(raw[1:len(raw)-1], '\'')
	if err != nil {
		p.reportError(tok, "invalid character literal: %v", err)
		return &ast.CharacterLiteral{Token: tok, Value: 0}
	}
	return &ast.CharacterLiteral{Token: tok, Value: value}
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
	if slices.ContainsFunc(kinds, p.check) {
		p.advance()
		return true
	}
	return false
}

// getPrecedence returns the binding power of a given token kind.
// Returns 0 if the token is not an infix operator.
func getPrecedence(kind token.Kind) int {
	switch kind {
	case token.LeftParen:
		return 6

	case token.Dot:
		return 5

	case token.Star, token.Slash:
		return 4

	case token.Plus, token.Minus:
		return 3

	case token.Less,
		token.LessEqual,
		token.Greater,
		token.GreaterEqual:
		return 2

	case token.EqualEqual,
		token.BangEqual:
		return 1

	default:
		return 0
	}
}
