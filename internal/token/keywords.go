package token

// Keywords maps string literals to keyword token kinds.
var Keywords = map[string]Kind{
	"fn":      KwFn,
	"do":      KwDo,
	"var":     KwVar,
	"mut":     KwMut,
	"return":  KwReturn,
	"end":     KwEnd,
	"char":    KwChar,
	"int":     KwInt,
	"bool":    KwBool,
	"unit":    KwUnit,
	"string":  KwString,
	"float":   KwFloat,
	"if":      KwIf,
	"then":    KwThen,
	"else":    KwElse,
	"struct":  KwStruct,
	"is":      KwIs,
	"importc": KwImportC,
	"loop":    KwLoop,
	"null":    KwNull,
}
