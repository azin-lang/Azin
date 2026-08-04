package syntax

var tokenText = [...]string{
	// Keywords
	KeywordDefer:   "defer",
	KeywordDo:      "do",
	KeywordElse:    "else",
	KeywordEnd:     "end",
	KeywordEnum:    "enum",
	KeywordFalse:   "false",
	KeywordFn:      "fn",
	KeywordFor:     "for",
	KeywordIf:      "if",
	KeywordImport:  "import",
	KeywordImportC: "importc",
	KeywordIs:      "is",
	KeywordLoop:    "loop",
	KeywordMut:     "mut",
	KeywordReturn:  "return",
	KeywordStop:    "stop",
	KeywordStruct:  "struct",
	KeywordThen:    "then",
	KeywordTrue:    "true",
	KeywordType:    "type",
	KeywordVar:     "var",
	KeywordWhile:   "while",

	// Operators and punctuation
	PlusToken:               "+",
	PlusEqualsToken:         "+=",
	MinusToken:              "-",
	MinusEqualsToken:        "-=",
	StarToken:               "*",
	StarEqualsToken:         "*=",
	SlashToken:              "/",
	SlashEqualsToken:        "/=",
	EqualsToken:             "=",
	EqualsEqualsToken:       "==",
	BangEqualsToken:         "!=",
	LessToken:               "<",
	LessEqualsToken:         "<=",
	GreaterToken:            ">",
	GreaterEqualsToken:      ">=",
	BangToken:               "!",
	AmpersandToken:          "&",
	AmpersandAmpersandToken: "&&",
	PipeToken:               "|",
	PipePipeToken:           "||",
	CaretToken:              "^",
	OpenParenToken:          "(",
	CloseParenToken:         ")",
	OpenBracketToken:        "[",
	CloseBracketToken:       "]",
	OpenBraceToken:          "{",
	CloseBraceToken:         "}",
	CommaToken:              ",",
	ColonToken:              ":",
	SemicolonToken:          ";",
	DotToken:                ".",
	NewlineToken:            "\n",
}

var keywords map[string]Kind

func init() {
	// Dynamically build the keyword map from the array to prevent duplication
	keywords = make(map[string]Kind, lastKeyword-firstKeyword+1)
	for i := firstKeyword; i <= lastKeyword; i++ {
		if txt := tokenText[i]; txt != "" {
			keywords[txt] = i
		}
	}
}

// LookupKeyword checks if an identifier is a reserved keyword.
func LookupKeyword(text string) Kind {
	if kind, ok := keywords[text]; ok {
		return kind
	}
	return IdentifierToken
}

// Text returns the static string representation of a token.
func (k Kind) Text() string {
	if int(k) < len(tokenText) {
		return tokenText[k]
	}
	return ""
}

// Display returns a human-readable representation for errors and diagnostics.
func (k Kind) Display() string {
	switch k {
	case BadToken:
		return "unexpected token"
	case IdentifierToken:
		return "identifier"
	case IntegerLiteralToken:
		return "integer literal"
	case FloatLiteralToken:
		return "float literal"
	case CharacterLiteralToken:
		return "character literal"
	case BooleanLiteralToken:
		return "boolean literal"
	case StringLiteralToken:
		return "string literal"
	case EndOfFileToken:
		return "end of file"
	default:
		if text := k.Text(); text != "" {
			return "'" + text + "'"
		}

		// Falls back to the auto-generated stringer method (e.g., "KeywordDefer")
		return k.String()
	}
}
