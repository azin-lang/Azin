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
		binaryPrec := syntax.BinaryPrecedence(p.current.Kind())
		if binaryPrec == 0 || binaryPrec < parentPrecedence {
			break
		}

		opToken := p.current
		p.advance()

		nextPrec := binaryPrec
		if syntax.AssociativityOf(opToken.Kind()) == syntax.Left {
			nextPrec++
		}

		right := p.parseBinaryExpression(nextPrec)
		left = green.NewBinaryExpression(left, opToken, right)
	}

	return left
}

// parseUnaryExpression handles prefix unary operators (!, -, +, etc.).
func (p *Parser) parseUnaryExpression() green.Node {
	unaryPrec := syntax.UnaryPrecedence(p.current.Kind())
	if unaryPrec != 0 {
		opToken := p.current
		p.advance()
		operand := p.parseUnaryExpression()
		return green.NewUnaryExpression(opToken, operand)
	}

	return p.parsePostfixExpression()
}

// parsePostfixExpression handles high-precedence postfix operators like function calls.
func (p *Parser) parsePostfixExpression() green.Node {
	expr := p.parsePrimaryExpression()

	for {
		switch p.current.Kind() {
		case syntax.OpenParenToken:
			expr = p.parseCallExpression(expr)
		default:
			return expr
		}
	}
}

// parseCallExpression parses function arguments and constructs a CallExpression node.
func (p *Parser) parseCallExpression(callee green.Node) green.Node {
	openParen := p.match(syntax.OpenParenToken)

	var args []green.Node
	if p.current.Kind() != syntax.CloseParenToken && p.current.Kind() != syntax.EndOfFileToken {
		for {
			args = append(args, p.ParseExpression())

			if p.current.Kind() == syntax.CommaToken {
				p.advance()
			} else if p.current.Kind() == syntax.CloseParenToken || p.current.Kind() == syntax.EndOfFileToken {
				break
			} else {
				// We encountered unexpected tokens instead of ',' or ')'
				// Synchronize to the next comma or closing parenthesis
				p.syncTo(syntax.CommaToken, syntax.CloseParenToken)

				if p.current.Kind() == syntax.CommaToken {
					p.advance()
				} else {
					break
				}
			}
		}
	}

	closeParen := p.match(syntax.CloseParenToken)
	return green.NewCallExpression(callee, openParen, args, closeParen)
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
			syntax.Display(p.current.Kind()),
		)

		if p.current.Kind() != syntax.EndOfFileToken {
			p.advance()
		}

		return green.NewToken(syntax.MissingToken, "", nil, nil)
	}
}
