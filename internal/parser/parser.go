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

func New(source string, tokens []token.Token, diag *diagnostics.Engine) *Parser {
	return &Parser{
		source: source,
		tokens: tokens,
		diag:   diag,
	}
}

func (p *Parser) synchronize() {
	p.advance() // consume at least a token (no deadlocks I hope)

	for !p.isAtEnd() {
		if p.previous().Kind == token.Newline ||
			p.previous().Kind == token.Semicolon {
			return
		}

		switch p.peek().Kind {
		case token.KwFn,
			token.KwStruct,
			token.KwEnum,
			token.KwVar,
			token.KwIf,
			token.KwReturn,
			token.KwElse,
			token.KwImportC,
			token.KwEnd:
			return
		}

		p.advance()
	}
}

func (p *Parser) error(tok token.Token, format string, args ...any) {
	p.diag.ReportError(tok.Position, tok.Length, format, args...)
}

func (p *Parser) reportError(tok token.Token, format string, args ...any) {
	p.error(tok, format, args...)
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
			continue
		}

		if p.current == before {
			p.advance()
		}

		p.synchronize()
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

	if context == "" {
		p.reportError(
			got,
			"expected %s %s, found %s",
			kind.DisplayName(),
			context,
			got.Kind.DisplayName(),
		)
	} else {
		p.reportError(
			got,
			"expected %s %s, found %s",
			kind.DisplayName(),
			context,
			got.Kind.DisplayName(),
		)
	}

	return got, false
}

func (p *Parser) parseVar() ast.Stmt {
	tok := p.advance()

	mutable := false
	if p.match(token.KwMut) {
		mutable = true
	}

	name := p.parseIdentifier()
	if name == nil {
		return nil
	}

	var typ *ast.Identifier
	if p.match(token.Colon) {
		typ = p.parseType()
	}

	var value ast.Expr

	if p.match(token.Equal) {
		value = p.parseExpression(0)
		if value == nil {
			return nil
		}
	} else if typ == nil {
		p.reportError(
			p.peek(),
			"expected '=' or an explicit type in variable declaration",
		)
	}

	if !p.statementEnd() {
		return nil
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
		return nil
	}

	path := p.parseStringLiteral()

	if !p.statementEnd() {
		return nil
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

	if p.check(token.EOF) ||
		p.check(token.KwEnd) ||
		p.check(token.KwElse) {
		return true
	}

	p.reportError(
		p.peek(),
		"expected end of statement (newline or ';')",
	)

	return false
}

func (p *Parser) parseStatement() ast.Stmt {
	p.skipNewlines()
	switch {
	case p.check(token.KwVar):
		return p.parseVar()
	case p.check(token.KwStruct):
		return p.parseStruct()

	case p.check(token.KwEnum):
		return p.parseEnum()

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

		if p.check(token.Equal) {
			tok := p.advance()

			switch expr.(type) {
			case *ast.Identifier, *ast.MemberExpr:
				// valid assignment target
			default:
				p.reportError(tok, "left side of assignment is not assignable")
				return nil
			}

			value := p.parseExpression(0)
			if value == nil {
				return nil
			}

			if !p.statementEnd() {
				return nil
			}

			return &ast.AssignmentStmt{
				Token: tok,
				Left:  expr,
				Value: value,
			}
		}

		if !p.statementEnd() {
			return nil
		}

		return &ast.ExpressionStmt{
			Token:      p.previous(),
			Expression: expr,
		}
	}
}

func (p *Parser) parseStruct() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()
	if name == nil {
		return nil
	}

	if _, ok := p.expect(token.KwIs, "after struct name"); !ok {
		return nil
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
			return nil
		}

		if _, ok := p.expect(token.Colon, "after field name"); !ok {
			return nil
		}

		tName := p.parseType()
		if tName == nil {
			return nil
		}

		if !p.statementEnd() {
			return nil
		}

		fields = append(fields, &ast.FieldDecl{
			Name:    fName,
			Type:    tName,
			Mutable: mutable,
		})
	}

	if _, ok := p.expect(token.KwEnd, "to close struct"); !ok {
		return nil
	}

	return &ast.StructStmt{
		Token:  tok,
		Name:   name,
		Fields: fields,
	}
}

func (p *Parser) parseEnum() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()
	if name == nil {
		return nil
	}

	if _, ok := p.expect(token.KwIs, "after enum name"); !ok {
		return nil
	}
	p.skipNewlines()

	variants := []*ast.Identifier{}

	for !p.isAtEnd() {
		p.skipNewlines()

		if p.check(token.KwEnd) {
			break
		}

		v := p.parseIdentifier()
		if v == nil {
			return nil
		}

		if !p.statementEnd() {
			return nil
		}

		variants = append(variants, v)
	}

	if _, ok := p.expect(token.KwEnd, "to close enum"); !ok {
		return nil
	}

	return &ast.EnumStmt{
		Token:    tok,
		Name:     name,
		Variants: variants,
	}
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
	if name == nil {
		return nil
	}

	params := []*ast.FieldDecl{}

	if p.match(token.LeftParen) {
		for !p.isAtEnd() && !p.check(token.RightParen) {
			pName := p.parseIdentifier()
			if pName == nil {
				return nil
			}
			p.match(token.Colon)
			pType := p.parseType()

			if pType == nil {
				return nil
			}

			params = append(params, &ast.FieldDecl{
				Name: pName,
				Type: pType,
			})

			if !p.check(token.RightParen) {
				if _, ok := p.expect(token.Comma, "between parameters"); !ok {
					return nil
				}
			}
		}

		if _, ok := p.expect(token.RightParen, "after parameter list"); !ok {
			return nil
		}
	}
	var retType *ast.Identifier

	if p.match(token.Colon) {
		retType = p.parseType()
		if retType == nil {
			return nil
		}
	}

	if _, ok := p.expect(token.KwDo, "after function declaration"); !ok {
		return nil
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
		return nil
	}

	return &ast.FuncStmt{Token: tok, Name: name, Params: params, ReturnType: retType, Body: body}
}

func (p *Parser) parseReturn() ast.Stmt {
	tok := p.advance()

	var value ast.Expr
	if !p.checkAny(
		token.Semicolon,
		token.Newline,
		token.KwEnd,
		token.EOF,
	) {
		value = p.parseExpression(0)
		if value == nil {
			return nil
		}
	}

	if !p.statementEnd() {
		return nil
	}

	return &ast.ReturnStmt{
		Token: tok,
		Value: value,
	}
}

func (p *Parser) skipNewlines() {
	for p.match(token.Newline) {
	}
}

func (p *Parser) parseIf() ast.Stmt {
	tok := p.advance() // if

	condition := p.parseExpression(0)
	if condition == nil {
		return nil
	}

	if _, ok := p.expect(token.KwThen, "after if condition"); !ok {
		return nil
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
		return nil
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
	var left ast.Expr

	switch {
	case p.check(token.Identifier):
		left = p.parseIdentifier()
		if left == nil {
			return nil
		}

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

		if p.checkAny(
			token.Newline,
			token.Semicolon,
			token.KwEnd,
			token.KwElse,
			token.KwThen,
			token.EOF,
		) {
			break
		}

		nextPrec := getPrecedence(p.peek().Kind)

		if nextPrec == 0 || precedence >= nextPrec {
			break
		}

		switch {
		case p.check(token.LeftParen):
			if _, ok := p.expect(token.LeftParen, "to start argument list"); !ok {
				return nil
			}
			var args []ast.Expr
			for !p.check(token.RightParen) {
				arg := p.parseExpression(0)
				if arg == nil {
					return nil
				}

				args = append(args, arg)
				if !p.check(token.RightParen) {
					if _, ok := p.expect(token.Comma, "between arguments"); !ok {
						return nil
					}
				}
			}
			if _, ok := p.expect(token.RightParen, "to close argument list"); !ok {
				return nil
			}
			left = &ast.CallExpr{Callee: left, Args: args}

		case p.check(token.Dot):
			if _, ok := p.expect(token.Dot, ""); !ok {
				return nil
			}
			prop := p.parseIdentifier()
			if prop == nil {
				return nil
			}
			left = &ast.MemberExpr{Object: left, Property: prop}

		default:
			// All binary operators are handled by this
			op := p.advance()

			right := p.parseExpression(nextPrec)

			if right == nil {
				return nil
			}

			left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
		}
	}
	return left
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
		p.reportError(tok, "invalid string literal: %v", err)
		return nil
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
		p.reportError(tok, "invalid character literal: %v", err)
		return nil
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
