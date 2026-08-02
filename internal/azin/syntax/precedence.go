package syntax

// TODO: Think about how I'll support "unary" expressions, like x.not, since I do not want to implement actual unary
//       operators.

func BinaryPrecedence(kind SyntaxKind) int {
	switch kind {
	case StarToken, SlashToken, AmpersandToken:
		return 5
	case PlusToken, MinusToken, PipeToken, CaretToken:
		return 4
	case EqualsEqualsToken, LessToken, LessEqualsToken, GreaterToken, GreaterEqualsToken:
		return 3
	case AmpersandAmpersandToken:
		return 2
	case PipePipeToken:
		return 1
	default:
		return 0
	}
}
