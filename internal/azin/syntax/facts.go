package syntax

var keywords = map[string]SyntaxKind{
	"defer":   KeywordDefer,
	"do":      KeywordDo,
	"else":    KeywordElse,
	"end":     KeywordEnd,
	"enum":    KeywordEnum,
	"fn":      KeywordFn,
	"for":     KeywordFor,
	"if":      KeywordIf,
	"import":  KeywordImport,
	"importc": KeywordImportC,
	"is":      KeywordIs,
	"loop":    KeywordLoop,
	"mut":     KeywordMut,
	"return":  KeywordReturn,
	"stop":    KeywordStop,
	"struct":  KeywordStruct,
	"then":    KeywordThen,
	"type":    KeywordType,
	"var":     KeywordVar,
	"while":   KeywordWhile,
}

func LookupKeyword(text string) SyntaxKind {
	if kind, ok := keywords[text]; ok {
		return kind
	}

	return IdentifierToken
}

func (k SyntaxKind) IsTrivia() bool {
	return k >= WhitespaceTrivia && k <= SingleLineCommentTrivia
}

func (k SyntaxKind) IsNode() bool {
	return k >= CompilationUnit && k <= CallExpression
}

func (k SyntaxKind) IsToken() bool {
	return k >= IdentifierToken && k <= NewlineToken
}

func (k SyntaxKind) IsKeyword() bool {
	return k >= KeywordDefer && k <= KeywordWhile
}

func (k SyntaxKind) IsLiteral() bool {
	return k >= IntegerLiteralToken && k <= StringLiteralToken
}

func (k SyntaxKind) IsOperator() bool {
	return k >= PlusToken && k <= CaretToken
}

func (k SyntaxKind) IsStatement() bool {
	return k >= BlockStatement && k <= ForStatement
}

func (k SyntaxKind) IsExpression() bool {
	return k >= LiteralExpression && k <= CallExpression
}
