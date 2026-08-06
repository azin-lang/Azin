package parser

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
)

// ParseStatement parses a single statement or falls back to an expression statement.
func (p *Parser) ParseStatement() green.Node {
	switch p.current.Kind() {
	case syntax.KeywordVar:
		return p.parseVariableDeclaration()
	case syntax.KeywordIf:
		return p.parseIfStatement()
	case syntax.KeywordReturn:
		return p.parseReturnStatement()
	default:
		return p.ParseExpression()
	}
}

// parseVariableDeclaration parses `var x = expr` or `var x: Type = expr`
func (p *Parser) parseVariableDeclaration() green.Node {
	varKeyword := p.match(syntax.KeywordVar)
	identifier := p.match(syntax.IdentifierToken)

	var colon *green.Token
	var typeAnnotation green.Node
	if p.at(syntax.ColonToken) {
		colon = p.advance()
		typeAnnotation = p.match(syntax.IdentifierToken)
	}

	var equals *green.Token
	var initializer green.Node
	if p.at(syntax.EqualsToken) {
		equals = p.advance()
		initializer = p.ParseExpression()
	}

	return green.NewVariableDeclaration(varKeyword, identifier, colon, typeAnnotation, equals, initializer)
}

// parseIfStatement parses `if condition thenBranch else elseBranch`
func (p *Parser) parseIfStatement() green.Node {
	ifKeyword := p.match(syntax.KeywordIf)
	condition := p.ParseExpression()
	thenBranch := p.ParseStatement()

	var elseKeyword *green.Token
	var elseBranch green.Node
	if p.eat(syntax.KeywordElse) {
		elseKeyword = p.match(syntax.KeywordElse)
		elseBranch = p.ParseStatement()
	}

	return green.NewIfStatement(ifKeyword, condition, thenBranch, elseKeyword, elseBranch)
}

// parseReturnStatement parses `return expr`
func (p *Parser) parseReturnStatement() green.Node {
	returnKeyword := p.match(syntax.KeywordReturn)

	var expr green.Node
	if !p.at(syntax.CloseBraceToken) && !p.at(syntax.EndOfFileToken) {
		expr = p.ParseExpression()
	}

	return green.NewReturnStatement(returnKeyword, expr)
}
