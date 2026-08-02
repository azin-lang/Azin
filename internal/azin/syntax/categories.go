package syntax

func (k SyntaxKind) IsTrivia() bool {
	return k >= firstTrivia && k <= lastTrivia
}

func (k SyntaxKind) IsNode() bool {
	return false
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
	return false
}

func (k SyntaxKind) IsStatement() bool {
	return false
}

func (k SyntaxKind) IsDeclaration() bool {
	return false
}

func (k SyntaxKind) IsPrimaryExpression() bool {
	return false
}

func (k SyntaxKind) StartsExpression() bool {
	return false
}

func (k SyntaxKind) StartsStatement() bool {
	return false
}

func (k SyntaxKind) StartsDeclaration() bool {
	return false
}
