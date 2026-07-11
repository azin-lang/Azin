package parser

import (
	"slices"

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

func (p *Parser) parseStatement() ast.Stmt {
	if p.check(token.KwStruct) {
		return p.parseStruct()
	}
	if p.check(token.KwFn) {
		return p.parseFunc()
	}
	if p.check(token.KwReturn) {
		return p.parseReturn()
	}
	return nil
}

func (p *Parser) parseStruct() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()

	p.match(token.KwIs)

	fields := []*ast.FieldDecl{}
	for !p.isAtEnd() && !p.check(token.KwEnd) {
		fName := p.parseIdentifier()
		p.match(token.Colon)

		var tName *ast.Identifier
		if p.check(token.KwInt) {
			tTok := p.advance()
			tName = &ast.Identifier{Token: tTok, Value: "int"}
		} else {
			tName = p.parseIdentifier()
		}
		p.match(token.Semicolon)

		fields = append(fields, &ast.FieldDecl{Name: fName, Type: tName})
	}
	p.match(token.KwEnd)
	return &ast.StructStmt{Token: tok, Name: name, Fields: fields}
}

func (p *Parser) parseFunc() ast.Stmt {
	tok := p.advance()
	name := p.parseIdentifier()

	params := []*ast.FieldDecl{}

	if p.match(token.LeftParen) {
		for !p.isAtEnd() && !p.check(token.RightParen) {
			pName := p.parseIdentifier()
			p.match(token.Colon)
			pType := p.parseIdentifier()

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
	p.match(token.Colon)

	var retType *ast.Identifier
	if p.check(token.KwInt) {
		tTok := p.advance()
		retType = &ast.Identifier{Token: tTok, Value: "int"}
	} else {
		retType = p.parseIdentifier()
	}

	p.match(token.KwDo)

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
	val := p.parseExpression(0)
	p.match(token.Semicolon)
	return &ast.ReturnStmt{Token: tok, Value: val}
}

func (p *Parser) parseExpression(precedence int) ast.Expr {
	var left ast.Expr

	if p.check(token.Identifier) {
		left = p.parseIdentifier()
	} else {
		return nil
	}

	for !p.isAtEnd() {
		if p.check(token.LeftParen) && precedence < 4 {
			p.advance()
			args := []ast.Expr{}
			for !p.check(token.RightParen) {
				args = append(args, p.parseExpression(0))
				if !p.check(token.RightParen) {
					p.match(token.Comma)
				}
			}
			p.advance()
			left = &ast.CallExpr{Function: left.(*ast.Identifier), Args: args}
			continue
		}

		if p.check(token.Dot) && precedence < 3 {
			p.advance()
			prop := p.parseIdentifier()
			left = &ast.MemberExpr{Object: left, Property: prop}
			continue
		}

		if p.check(token.Minus) && precedence < 2 {
			op := p.advance()
			right := p.parseExpression(2)
			left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
			continue
		}

		if p.check(token.Plus) && precedence < 1 {
			op := p.advance()
			right := p.parseExpression(1)
			left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
			continue
		}
		break
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

func (p *Parser) match(kinds ...token.Kind) bool {
	if slices.ContainsFunc(kinds, p.check) {
		p.advance()
		return true
	}
	return false
}
