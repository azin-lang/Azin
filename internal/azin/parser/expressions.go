package parser

import (
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
)

// ParseExpression is the public entry point for expression parsing.
func (p *Parser) ParseExpression() green.Node {
	return p.parseBinaryExpression(0)
}

// parseBinaryExpression uses precedence climbing for binary operators.
func (p *Parser) parseBinaryExpression(parentPrecedence int) green.Node {
	left := p.parseUnaryExpression()

	for {
		binaryPrec := p.current.Kind().BinaryPrecedence()
		if binaryPrec == 0 || binaryPrec < parentPrecedence {
			break
		}

		opToken := p.current
		p.advance()

		nextPrec := binaryPrec
		if opToken.Kind().Associativity() == syntax.Left {
			nextPrec++
		}

		right := p.parseBinaryExpression(nextPrec)
		left = green.NewBinaryExpression(left, opToken, right)
	}

	return left
}

// parseUnaryExpression handles prefix unary operators (!, -, +, etc.).
func (p *Parser) parseUnaryExpression() green.Node {
	unaryPrec := p.current.Kind().UnaryPrecedence()
	if unaryPrec != 0 {
		opToken := p.current
		p.advance()
		operand := p.parseUnaryExpression()
		return green.NewUnaryExpression(opToken, operand)
	}

	return p.parsePostfixExpression()
}

// parsePostfixExpression handles high-precedence postfix operators like calls and member access.
func (p *Parser) parsePostfixExpression() green.Node {
	expr := p.parsePrimaryExpression()

	for {
		switch p.current.Kind() {
		case syntax.OpenParenToken:
			expr = p.parseCallExpression(expr)
		case syntax.DotToken:
			dot := p.advance()
			member := p.match(syntax.IdentifierToken)
			expr = green.NewMemberAccessExpression(expr, dot, member)
		default:
			return expr
		}
	}
}

// parseCallExpression parses function arguments and constructs a CallExpression node.
func (p *Parser) parseCallExpression(callee green.Node) green.Node {
	openParen := p.match(syntax.OpenParenToken)

	var args []green.Node
	if !p.atAny(syntax.CloseParenToken, syntax.EndOfFileToken) {
		for {
			args = append(args, p.ParseExpression())

			if p.eat(syntax.CommaToken) {
				// Allow optional trailing comma: foo(a, b,)
				if p.at(syntax.CloseParenToken) {
					break
				}
			} else if p.atAny(syntax.CloseParenToken, syntax.EndOfFileToken) {
				break
			} else if p.isAtExpressionStart() {
				// Recover from a missing comma between arguments: foo(a b)
				p.match(syntax.CommaToken)
			} else {
				// Unexpected tokens inside arguments list; synchronize to next comma or ')'
				p.syncTo(syntax.CommaToken, syntax.CloseParenToken)
				if !p.eat(syntax.CommaToken) {
					break
				}
			}
		}
	}

	closeParen := p.match(syntax.CloseParenToken)

	// Wrap the slice in a SyntaxList node so it implements green.Node
	argsList := green.NewSyntaxList(args)

	return green.NewCallExpression(callee, openParen, argsList, closeParen)
}

// parsePrimaryExpression parses atomic elements (literals, identifiers, parenthesized expressions).
func (p *Parser) parsePrimaryExpression() green.Node {
	switch p.current.Kind() {
	case syntax.IntegerLiteralToken, syntax.FloatLiteralToken, syntax.StringLiteralToken, syntax.KeywordTrue, syntax.KeywordFalse:
		tok := p.advance()
		return green.NewLiteralExpression(tok)

	case syntax.IdentifierToken:
		tok := p.advance()
		return green.NewNameExpression(tok)

	case syntax.OpenParenToken:
		openParen := p.match(syntax.OpenParenToken)
		expr := p.ParseExpression()
		closeParen := p.match(syntax.CloseParenToken)
		return green.NewParenthesizedExpression(openParen, expr, closeParen)

	default:
		p.diags.Error(
			p.currentLocation(),
			"E011",
			"Expected expression but got %s",
			p.current.Kind().String(),
		)

		// Skip unexpected token only if it is not a structural delimiter/boundary token
		if !p.isDelimiterOrEOF() {
			p.advance()
		}

		return green.NewToken(syntax.MissingToken, "", nil, nil)
	}
}
