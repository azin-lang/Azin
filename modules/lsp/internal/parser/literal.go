package parser

import (
	"strconv"
	"strings"

	"github.com/azin-lang/azin/lsp/internal/ast"
)

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	tok := p.advance()
	text := p.lexeme(tok)

	var value int64
	var err error

	switch {
	case strings.HasPrefix(text, "0x"):
		value, err = strconv.ParseInt(text[2:], 16, 64)
	case strings.HasPrefix(text, "0b"):
		value, err = strconv.ParseInt(text[2:], 2, 64)
	default:
		value, err = strconv.ParseInt(text, 10, 64)
	}

	if err != nil {
		p.reportError(tok, "invalid integer literal: %v", err)
		return &ast.IntegerLiteral{Token: tok, Value: 0}
	}

	return &ast.IntegerLiteral{Token: tok, Value: value}
}

func (p *Parser) parseFloatLiteral() *ast.FloatLiteral {
	tok := p.advance()
	value, err := strconv.ParseFloat(p.lexeme(tok), 64)
	if err != nil {
		p.reportError(tok, "invalid float literal: %v", err)
		return &ast.FloatLiteral{Token: tok, Value: 0.0}
	}
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

	if len(raw) < 2 || raw[0] != '\'' || raw[len(raw)-1] != '\'' {
		p.reportError(tok, "malformed character literal")
		return &ast.CharacterLiteral{Token: tok, Value: 0}
	}

	value, _, _, err := strconv.UnquoteChar(raw[1:len(raw)-1], '\'')
	if err != nil {
		p.reportError(tok, "invalid character literal: %v", err)
		return &ast.CharacterLiteral{Token: tok, Value: 0}
	}
	return &ast.CharacterLiteral{Token: tok, Value: value}
}
