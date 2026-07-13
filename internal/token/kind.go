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
	KwMut     // kw_mut
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
	KwEnum    // kw_enum

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

// DisplayName takes a token as a parameter, and returns a string of the token's display name
func (k Kind) DisplayName() string {
	switch k {
	case KwFn:
		return "'fn'"
	case KwDo:
		return "'do'"
	case KwVar:
		return "'var'"
	case KwMut:
		return "'mut'"
	case KwReturn:
		return "'return'"
	case KwEnd:
		return "'end'"
	case KwChar:
		return "'char'"
	case KwInt:
		return "'int'"
	case KwUnit:
		return "'unit'"
	case KwString:
		return "'string'"
	case KwFloat:
		return "'float'"
	case KwIf:
		return "'if'"
	case KwThen:
		return "'then'"
	case KwElse:
		return "'else'"
	case KwStruct:
		return "'struct'"
	case KwIs:
		return "'is'"
	case KwImportC:
		return "'importC'"
	case KwEnum:
		return "'enum'"

	case Identifier:
		return "identifier"
	case IntegerLiteral:
		return "integer literal"
	case FloatLiteral:
		return "float literal"
	case StringLiteral:
		return "string literal"
	case CharacterLiteral:
		return "character literal"

	case LeftParen:
		return "'('"
	case RightParen:
		return "')'"
	case LeftBrace:
		return "'{'"
	case RightBrace:
		return "'}'"
	case LeftBracket:
		return "'['"
	case RightBracket:
		return "']'"
	case Comma:
		return "','"
	case Colon:
		return "':'"
	case Semicolon:
		return "';'"
	case Dot:
		return "'.'"

	case Plus:
		return "'+'"
	case Minus:
		return "'-'"
	case Star:
		return "'*'"
	case Slash:
		return "'/'"
	case Equal:
		return "'='"
	case EqualEqual:
		return "'=='"
	case Bang:
		return "'!'"
	case BangEqual:
		return "'!='"
	case Less:
		return "'<'"
	case LessEqual:
		return "'<='"
	case Greater:
		return "'>'"
	case GreaterEqual:
		return "'>='"

	case Newline:
		return "newline"
	case EOF:
		return "end of file"

	default:
		return k.String()
	}
}
