package token

//// go:generate go run golang.org/x/tools/cmd/stringer -type=Kind -linecomment
//go:generate go tool stringer -type=Kind -linecomment

// Kind is the set of lexical token types.
type Kind uint8

const (
	Unknown Kind = iota // unknown

	Identifier       // identifier
	IntegerLiteral   // integer_literal
	StringLiteral    // string_literal
	FloatLiteral     // float_literal
	CharacterLiteral // character_literal

	KwFn      // kw_fn
	KwDo      // kw_do
	KwVar     // kw_var
	KwReturn  // kw_return
	KwEnd     // kw_end
	KwChar    // kw_char
	KwInt     // kw_int
	KwUnit    // kw_unit
	KwString  // kw_string
	KwFloat   // kw_float
	KwIf      // kw_if
	KwThen    // kw_then
	KwElse    // kw_else
	KwStruct  // kw_struct
	KwIs      // kw_is
	KwImportC // kw_import

	Plus           // plus
	Minus          // minus
	Star           // star
	Slash          // slash
	Equal          // equal
	EqualEqual     // equal_equal
	Bang           // bang
	BangEqual      // bang_equal
	Less           // less
	LessEqual      // less_equal
	Greater        // greater
	GreaterEqual   // greater_equal
	Arrow          // arrow
	Modulo         // modulo
	Pipe           // pipe
	LogicalOr      // logical_or
	LogicalAnd     // logical_and
	Ampersand      // ampersand
	Caret          // caret
	Tilde          // tilde
	PlusEqual      // plus_equal
	MinusEqual     // minus_equal
	StarEqual      // star_equal
	SlashEqual     // slash_equal
	ModuloEqual    // modulo_equal
	CaretEqual     // caret_equal
	PipeEqual      // pipe_equal
	AmpersandEqual // ampersand_equal
	PlusPlus       // plus_plus
	MinusMinus     // minus_minus
	LessLess       // less_less
	GreaterGreater // greater_greater

	LeftParen    // left_paren
	RightParen   // right_paren
	LeftBrace    // left_brace
	RightBrace   // right_brace
	Comma        // comma
	Semicolon    // semicolon
	Colon        // colon
	Dot          // dot
	LeftBracket  // left_bracket
	RightBracket // right_bracket
	Newline      // newline

	EOF // eof

	Error // error
)
