package syntax

func Display(kind SyntaxKind) string {
	switch kind {
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
		if text := Text(kind); text != "" {
			return "'" + text + "'"
		}

		return kind.String()
	}
}
