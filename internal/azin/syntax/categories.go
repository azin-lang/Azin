package syntax

func (k SyntaxKind) IsTrivia() bool {
	return false
}

func (k SyntaxKind) IsNode() bool {
	return false
}

func (k SyntaxKind) IsToken() bool {
	return false
}

func (k SyntaxKind) IsKeyword() bool {
	return false
}

func (k SyntaxKind) IsLiteral() bool {
	return false
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
