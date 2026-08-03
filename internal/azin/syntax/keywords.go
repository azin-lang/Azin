package syntax

var keywords = map[string]SyntaxKind{
	"defer":   KeywordDefer,
	"do":      KeywordDo,
	"else":    KeywordElse,
	"end":     KeywordEnd,
	"enum":    KeywordEnum,
	"false":   KeywordFalse,
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
	"true":    KeywordTrue,
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
