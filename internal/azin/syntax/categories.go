package syntax

func (k SyntaxKind) IsTrivia() bool {
	return k >= firstTrivia && k <= lastTrivia
}

func (k SyntaxKind) IsNode() bool {
	return k >= CompilationUnit && k <= CastExpression
}

func (k SyntaxKind) IsToken() bool {
	return k >= firstToken && k <= lastToken
}

func (k SyntaxKind) IsKeyword() bool {
	return k >= firstKeyword && k <= lastKeyword
}

func (k SyntaxKind) IsLiteral() bool {
	return k >= firstLiteral && k <= lastLiteral
}

func (k SyntaxKind) IsExpression() bool {
	return k >= LiteralExpression && k <= CastExpression
}

func (k SyntaxKind) IsStatement() bool {
	return k >= BlockStatement && k <= ForStatement
}

func (k SyntaxKind) IsDeclaration() bool {
	return k >= ImportDeclaration && k <= FieldDeclaration
}

func (k SyntaxKind) IsPrimaryExpression() bool {
	return k == LiteralExpression || k == NameExpression || k == ParenthesizedExpression
}

func (k SyntaxKind) StartsExpression() bool {
	return k.IsLiteral() || k == IdentifierToken || k == OpenParenToken || k.IsPrefixOperator()
}

func (k SyntaxKind) StartsStatement() bool {
	switch k {
	case OpenBraceToken, KeywordIf, KeywordWhile, KeywordFor, KeywordLoop, KeywordDefer, KeywordStop, KeywordReturn, SemicolonToken:
		return true
	default:
		return k.StartsExpression()
	}
}

func (k SyntaxKind) StartsDeclaration() bool {
	switch k {
	case KeywordImport, KeywordImportC, KeywordFn, KeywordStruct, KeywordEnum, KeywordType, KeywordVar:
		return true
	default:
		return false
	}
}
