package token

// map of keywords to their token kinds
// regen
// ok
var Keywords = map[string]Kind{
	"fn":     KwFn,
	"do":     KwDo,
	"var":    KwVar,
	"return": KwReturn,
	"end":    KwEnd,
	"char":   KwChar,
	"int":    KwInt,
	"str":    KwString,
	"float":  KwFloat,
	"if":     KwIf,
	"then":   KwThen,
}
