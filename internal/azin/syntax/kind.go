package syntax

type Kind uint16

const (
	Unknown Kind = iota
	EndOfFile
	Error

	WhitespaceTrivia
	EndOfLineTrivia
	SingleLineCommentTrivia

	IdentifierToken
	IntegerLiteralToken
	FloatLiteralToken
	CharacterLiteralToken
	BooleanLiteralToken
	StringLiteralToken

	KeywordDefer
	KeywordDo
	KeywordElse
	KeywordEnd
	KeywordEnum
	KeywordFn
	KeywordFor
	KeywordIf
	KeywordImport
	KeywordImportC
	KeywordIs
	KeywordLoop
	KeywordMut
	KeywordReturn
	KeywordStop
	KeywordStruct
	KeywordThen
	KeywordType
	KeywordVar
	KeywordWhile

	PlusToken
	PlusEqualsToken
	MinusToken
	MinusEqualsToken
	EqualsToken
	EqualsEqualsToken
	StarToken
	StarEqualsToken

	OpenParenToken
	CloseParenToken
	OpenBracketToken
	CloseBracketToken
	OpenBraceToken
	CloseBraceToken
	CommaToken
	SemicolonToken
	ColonToken
	DotToken
	NewlineToken
)

var keywords = map[string]Kind{
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
func LookupKeyword(text string) Kind {
	if kind, ok := keywords[text]; ok {
		return kind
	}
	return IdentifierToken
}

// IsKeyword reports whether the kind represents a reserved keyword.
func (k Kind) IsKeyword() bool {
	return k >= KeywordDefer && k <= KeywordWhile
}

// IsLiteral reports whether the kind represents a literal value.
func (k Kind) IsLiteral() bool {
	return k >= IntegerLiteralToken && k <= StringLiteralToken
}

// IsOperator reports whether the kind represents an operator.
func (k Kind) IsOperator() bool {
	return k >= PlusToken && k <= StarEqualsToken
}
