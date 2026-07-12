package parser

import (
	"slices"
	"strconv"
	"strings"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

type Parser struct {
	source  string
	tokens  []token.Token
	current int
}

func New(source string, tokens []token.Token) *Parser {
	return &Parser{source: source, tokens: tokens, current: 0}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Stmt{}}

	for !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else {
			p.advance()
		}
	}

	return program
}

func (p *Parser) expect(kind token.Kind, message string) token.Token {
	if p.check(kind) {
		return p.advance()
	}

	panic(message)
}

func (p *Parser) parseVar() ast.Stmt {
	tok := p.advance()

	name := p.parseIdentifier()

	var typ *ast.Identifier
	if p.match(token.Colon) {
		typ = p.parseType()
	}

	p.expect(token.Equal, "expected '=' after variable declaration")

	value := p.parseExpression(0)

	p.match(token.Semicolon)

	return &ast.VarStmt{
		Token: tok,
		Name:  name,
		Type:  typ,
		Value: value,
	}
}

func (p *Parser) parseImportC() ast.Stmt {
	tok := p.advance() // consume 'import'

	p.match(token.LeftParen)

	// Parse the path string literal
	strTok := p.advance()
	pathExpr := &ast.StringLiteral{
		Token: strTok,
		Value: p.source[strTok.Position.Offset : strTok.Position.Offset+strTok.Length], // Slices the path clean
	}

	p.match(token.RightParen)
	p.match(token.Semicolon) // Match optional or explicit trailing semicolon if required

	return &ast.ImportCStmt{Token: tok, Path: pathExpr}
}

func (p *Parser) parseStatement() ast.Stmt {
	switch {
	case p.check(token.KwVar):
		return p.parseVar()
	case p.check(token.KwStruct):
		return p.parseStruct()

	case p.check(token.KwFn):
		return p.parseFunc()

	case p.check(token.KwReturn):
		return p.parseReturn()

	case p.check(token.KwIf):
		return p.parseIf()

	case p.check(token.KwImportC):
		return p.parseImportC()

	default:
		expr := p.parseExpression(0)
		if expr == nil {
			return nil
		}

		p.match(token.Semicolon)

		return &ast.ExpressionStmt{
			Token:      p.previous(),
			Expression: expr,
		}
	}
}

func (p *Parser) parseStruct() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()

	p.match(token.KwIs)

	fields := []*ast.FieldDecl{}
	for !p.isAtEnd() && !p.check(token.KwEnd) {
		fName := p.parseIdentifier()
		p.match(token.Colon)

		tName := p.parseType()
		p.match(token.Semicolon)

		fields = append(fields, &ast.FieldDecl{Name: fName, Type: tName})
	}
	p.match(token.KwEnd)
	return &ast.StructStmt{Token: tok, Name: name, Fields: fields}
}

func (p *Parser) parseType() *ast.Identifier {
	switch p.peek().Kind {

	case token.KwUnit:
		tok := p.advance()
		return &ast.Identifier{Token: tok, Value: "unit"}

	case token.KwInt:
		tok := p.advance()
		return &ast.Identifier{Token: tok, Value: "int"}

	case token.KwFloat:
		tok := p.advance()
		return &ast.Identifier{Token: tok, Value: "float"}

	case token.KwString:
		tok := p.advance()
		return &ast.Identifier{Token: tok, Value: "string"}

	case token.KwChar:
		tok := p.advance()
		return &ast.Identifier{Token: tok, Value: "char"}

	default:
		return p.parseIdentifier()
	}
}

func (p *Parser) parseFunc() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()

	params := []*ast.FieldDecl{}

	if p.match(token.LeftParen) {
		for !p.isAtEnd() && !p.check(token.RightParen) {
			pName := p.parseIdentifier()
			p.match(token.Colon)
			pType := p.parseType()

			params = append(params, &ast.FieldDecl{
				Name: pName,
				Type: pType,
			})

			if !p.check(token.RightParen) {
				p.match(token.Comma)
			}
		}

		p.match(token.RightParen)
	}
	var retType *ast.Identifier

	if p.match(token.Colon) {
		retType = p.parseType()
	}

	p.expect(token.KwDo, "expected 'do' after function declaration")

	body := []ast.Stmt{}
	for !p.isAtEnd() && !p.check(token.KwEnd) {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		} else {
			p.advance()
		}
	}
	p.match(token.KwEnd)

	return &ast.FuncStmt{Token: tok, Name: name, Params: params, ReturnType: retType, Body: body}
}

func (p *Parser) parseReturn() ast.Stmt {
	tok := p.advance()

	var value ast.Expr
	if !p.check(token.Semicolon) {
		value = p.parseExpression(0)
	}

	p.match(token.Semicolon)

	return &ast.ReturnStmt{
		Token: tok,
		Value: value,
	}
}

func (p *Parser) parseIf() ast.Stmt {
	tok := p.advance() // if

	condition := p.parseExpression(0)

	p.match(token.KwThen)

	var thenBody []ast.Stmt

	for !(p.isAtEnd() || p.checkAny(token.KwElse, token.KwEnd)) {
		if stmt := p.parseStatement(); stmt != nil {
			thenBody = append(thenBody, stmt)
		} else {
			p.advance()
		}
	}

	var elseBody []ast.Stmt

	if p.match(token.KwElse) {
		for !p.isAtEnd() && !p.check(token.KwEnd) {
			if stmt := p.parseStatement(); stmt != nil {
				elseBody = append(elseBody, stmt)
			} else {
				p.advance()
			}
		}
	}

	p.match(token.KwEnd)

	return &ast.IfStmt{
		Token:     tok,
		Condition: condition,
		Then:      thenBody,
		Else:      elseBody,
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expr {
	var left ast.Expr

	switch {
	case p.check(token.Identifier):
		left = p.parseIdentifier()

	case p.check(token.IntegerLiteral):
		left = p.parseIntegerLiteral()

	case p.check(token.FloatLiteral):
		left = p.parseFloatLiteral()

	case p.check(token.StringLiteral):
		left = p.parseStringLiteral()

	case p.check(token.CharacterLiteral):
		left = p.parseCharacterLiteral()

	default:
		return nil
	}

	for !p.isAtEnd() {
		nextPrec := getPrecedence(p.peek().Kind)

		// If the next token isn't an operator, or its precedence is lower/equal, stop.
		if precedence >= nextPrec || nextPrec == 0 {
			break
		}

		switch {
		case p.check(token.LeftParen):
			p.advance()
			var args []ast.Expr
			for !p.check(token.RightParen) {
				args = append(args, p.parseExpression(0))
				if !p.check(token.RightParen) {
					p.match(token.Comma)
				}
			}
			p.advance()
			left = &ast.CallExpr{Callee: left, Args: args}

		case p.check(token.Dot):
			p.advance()
			prop := p.parseIdentifier()
			left = &ast.MemberExpr{Object: left, Property: prop}

		default:
			// All binary operators are handled by this
			op := p.advance()
			right := p.parseExpression(nextPrec)
			left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
		}
	}
	return left
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	tok := p.advance()

	start := tok.Position.Offset
	end := start + tok.Length
	realValue := p.source[start:end]

	return &ast.Identifier{Token: tok, Value: realValue}
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	tok := p.advance()

	start := tok.Position.Offset
	end := start + tok.Length

	text := p.source[start:end]

	var value int64
	switch {
	case strings.HasPrefix(text, "0x"):
		value, _ = strconv.ParseInt(text[2:], 16, 64)
	case strings.HasPrefix(text, "0b"):
		value, _ = strconv.ParseInt(text[2:], 2, 64)
	default:
		value, _ = strconv.ParseInt(text, 10, 64)
	}

	return &ast.IntegerLiteral{
		Token: tok,
		Value: value,
	}
}

func (p *Parser) parseFloatLiteral() *ast.FloatLiteral {
	tok := p.advance()

	start := tok.Position.Offset
	end := start + tok.Length

	text := p.source[start:end]

	value, _ := strconv.ParseFloat(text, 64)

	return &ast.FloatLiteral{
		Token: tok,
		Value: value,
	}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	tok := p.advance()

	start := tok.Position.Offset
	end := start + tok.Length

	raw := p.source[start:end]

	value, err := strconv.Unquote(raw)
	if err != nil {
		panic(err)
	}

	return &ast.StringLiteral{
		Token: tok,
		Value: value,
	}
}

func (p *Parser) parseCharacterLiteral() *ast.CharacterLiteral {
	tok := p.advance()

	start := tok.Position.Offset
	end := start + tok.Length

	raw := p.source[start:end]

	value, _, _, err := strconv.UnquoteChar(raw[1:len(raw)-1], '\'')
	if err != nil {
		panic(err)
	}

	return &ast.CharacterLiteral{
		Token: tok,
		Value: value,
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
	if p.isAtEnd() {
		return false
	}
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
		return 4
	case token.Dot:
		return 3
	case token.Plus, token.Minus:
		return 2
	case token.Greater, token.GreaterEqual, token.Less, token.LessEqual:
		return 1
	default:
		return 0
	}
}
