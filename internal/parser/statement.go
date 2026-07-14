package parser

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: make([]ast.Stmt, 0, len(p.tokens)/4),
	}

	for !p.isAtEnd() {
		before := p.current

		stmt := p.parseStatement()

		// Reached EOF while skipping trailing newlines
		if stmt == nil && p.isAtEnd() {
			break
		}

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		// Only recover if we're genuinely stuck
		if p.current == before {
			p.synchronize()
		}
	}

	return program
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
		value = p.parseExpression(PrecLowest)
	} else if typ == nil {
		p.reportError(p.peek(), "expected '=' or an explicit type in variable declaration")
		value = badExpr(p.peek())
	}

	p.consumeStatementEnd()

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
	p.consumeStatementEnd()

	return &ast.ImportCStmt{
		Token: tok,
		Path:  path,
	}
}

func (p *Parser) parseStatement() ast.Stmt {
	p.skipNewlines()

	if p.isAtEnd() {
		return nil
	}

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
	case p.check(token.KwLoop):
		stmt = p.parseLoop()
	case p.check(token.KwStop):
		stmt = p.parseStop()
	default:
		stmt = p.parseExpressionOrAssignment()
	}

	if stmt == nil || isBadStmt(stmt) {
		p.synchronize()
		if stmt == nil {
			return badStmt(p.peek())
		}
	}
	return stmt
}

func (p *Parser) parseExpressionOrAssignment() ast.Stmt {
	expr := p.parseExpression(PrecLowest)

	if isBadExpr(expr) {
		return badStmt(p.peek())
	}

	if p.check(token.Equal) {
		tok := p.advance()

		switch expr.(type) {
		case *ast.Identifier, *ast.MemberExpr:
			// Valid assignment target
		default:
			p.reportError(tok, "left side of assignment is not assignable")
			return badStmt(tok)
		}

		value := p.parseExpression(PrecLowest)
		p.consumeStatementEnd()

		return &ast.AssignmentStmt{
			Token: tok,
			Left:  expr,
			Value: value,
		}
	}

	p.consumeStatementEnd()

	return &ast.ExpressionStmt{
		Token:      p.previous(),
		Expression: expr,
	}
}

func (p *Parser) parseStruct() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()
	if name == nil {
		return badStmt(tok) // Without a name, the struct is useless
	}
	p.expect(token.KwIs, "after struct name")

	var fields []*ast.FieldDecl
	for !p.isAtEnd() {
		p.skipNewlines()
		if p.check(token.KwEnd) {
			break
		}

		field := p.parseFieldDecl(true)
		if field != nil {
			fields = append(fields, field)
		} else {
			// If field parsing failed, sync to the next newline to try parsing the next field
			p.synchronize()
		}

		p.consumeStatementEnd()
	}

	p.expect(token.KwEnd, "to close struct")

	return &ast.StructStmt{
		Token:  tok,
		Name:   name,
		Fields: fields,
	}
}

func (p *Parser) parseType() *ast.Identifier {
	if isBuiltinType(p.peek().Kind) {
		tok := p.advance()
		return &ast.Identifier{Token: tok, Value: p.lexeme(tok)}
	}
	return p.parseIdentifier()
}

func (p *Parser) parseFunc() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()
	if name == nil {
		return badStmt(tok)
	}

	var params []*ast.FieldDecl
	if p.match(token.LeftParen) {
		if !p.check(token.RightParen) {
			for {
				param := p.parseFieldDecl(false)
				if param != nil {
					params = append(params, param)
				}
				if !p.match(token.Comma) {
					break
				}
			}
		}
		p.expect(token.RightParen, "after parameter list")
	}

	var retType *ast.Identifier
	if p.match(token.Colon) {
		retType = p.parseType()
	}

	p.expect(token.KwDo, "after function declaration")
	body := p.parseBlock(token.KwEnd)
	p.expect(token.KwEnd, "to close function")

	// Return partial function node even if errors occurred (allows autocomplete to work inside)
	return &ast.FuncStmt{Token: tok, Name: name, Params: params, ReturnType: retType, Body: body}
}

func (p *Parser) parseReturn() ast.Stmt {
	tok := p.advance()
	var value ast.Expr

	if !p.checkAny(token.Semicolon, token.Newline, token.KwEnd, token.EOF) {
		value = p.parseExpression(PrecLowest)
	}

	p.consumeStatementEnd()

	return &ast.ReturnStmt{
		Token: tok,
		Value: value,
	}
}

func (p *Parser) parseIf() ast.Stmt {
	tok := p.advance()
	condition := p.parseExpression(PrecLowest)

	p.expect(token.KwThen, "after if condition")
	thenBody := p.parseBlock(token.KwElse, token.KwEnd)

	var elseBody []ast.Stmt
	if p.match(token.KwElse) {
		elseBody = p.parseBlock(token.KwEnd)
	}

	p.skipNewlines()
	p.expect(token.KwEnd, "to close if statement")

	return &ast.IfStmt{
		Token:     tok,
		Condition: condition,
		Then:      thenBody,
		Else:      elseBody,
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expr {
	left := p.parsePrefix()
	if left == nil {
		errTok := p.peek()
		p.reportError(errTok, "expected expression")
		return badExpr(errTok)
	}

	for !p.isAtEnd() {
		nextPrec := getPrecedence(p.peek().Kind)

		if nextPrec == PrecLowest || precedence >= nextPrec {
			break
		}

		left = p.parseInfix(left, nextPrec)
	}

	return left
}

func (p *Parser) parsePrefix() ast.Expr {
	switch {
	case p.check(token.Identifier):
		tok := p.peek()
		val := p.lexeme(tok)

		if val == "true" || val == "false" {
			p.advance()
			return &ast.BooleanLiteral{
				Token: tok,
				Value: val == "true",
			}
		}
		return p.parseIdentifier()

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
	case p.match(token.LeftParen):
		var args []ast.Expr

		if !p.check(token.RightParen) {
			for {
				args = append(args, p.parseExpression(PrecLowest))
				if !p.match(token.Comma) {
					break
				}
			}
		}
		p.expect(token.RightParen, "to close argument list")
		return &ast.CallExpr{Callee: left, Args: args}

	case p.match(token.Dot):
		prop := p.parseIdentifier()
		if prop == nil {
			return badExpr(p.previous())
		}
		return &ast.MemberExpr{Object: left, Property: prop}

	default:
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
	return &ast.Identifier{
		Token: tok,
		Value: p.lexeme(tok),
	}
}

func (p *Parser) parseFieldDecl(allowMut bool) *ast.FieldDecl {
	var mutable bool
	if allowMut {
		mutable = p.match(token.KwMut)
	}

	name := p.parseIdentifier()
	if name == nil {
		return nil
	}

	if _, ok := p.expect(token.Colon, "after name"); !ok {
		return nil
	}

	tName := p.parseType()
	if tName == nil {
		return nil
	}

	return &ast.FieldDecl{
		Name:    name,
		Type:    tName,
		Mutable: mutable,
	}
}

func (p *Parser) parseStop() ast.Stmt {
	tok := p.advance()
	p.consumeStatementEnd()

	return &ast.StopStmt{
		Token: tok,
	}
}

func (p *Parser) parseLoop() ast.Stmt {
	tok := p.advance()

	body := p.parseBlock(token.KwEnd)

	p.expect(token.KwEnd, "to close loop")

	return &ast.LoopStmt{
		Token: tok,
		Body:  body,
	}
}
