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

// LookupKeyword checks if the given identifier text is a reserved keyword.
// It returns IdentifierToken if no match is found.
func LookupKeyword(text string) SyntaxKind {
	if kind, ok := keywords[text]; ok {
		return kind
	}
	return IdentifierToken
}

// IsKeyword reports whether the kind represents a reserved keyword.
func (k SyntaxKind) IsKeyword() bool {
	return k >= KeywordDefer && k <= KeywordWhile
}

// IsLiteral reports whether the kind represents a literal value.
func (k SyntaxKind) IsLiteral() bool {
	return k >= IntegerLiteralToken && k <= StringLiteralToken
}

// IsOperator reports whether the kind represents an operator.
func (k SyntaxKind) IsOperator() bool {
	return k >= PlusToken && k <= StarEqualsToken
}
