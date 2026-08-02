package syntax

// TODO: Think about how I'll support "unary" expressions, like x.not, since I do not want to implement actual unary
//       operators.

type Associativity uint8

const (
	Left Associativity = iota
	Right
)

func BinaryPrecedence(kind SyntaxKind) int {
	switch kind {
	case StarToken, SlashToken, AmpersandToken:
		return 5
	case PlusToken, MinusToken, PipeToken, CaretToken:
		return 4
	case EqualsEqualsToken, BangEqualsToken, LessToken, LessEqualsToken, GreaterToken, GreaterEqualsToken:
		return 3
	case AmpersandAmpersandToken:
		return 2
	case PipePipeToken:
		return 1
	default:
		return 0
	}
}

func UnaryPrecedence(kind SyntaxKind) int {
	return 0
}

func AssociativityOf(kind SyntaxKind) Associativity {
	return Left
}

func (k SyntaxKind) IsOperator() bool {
	return false
}

func (k SyntaxKind) IsArithmeticOperator() bool {
	return false
}

func (k SyntaxKind) IsAssignmentOperator() bool {
	return false
}

func (k SyntaxKind) IsComparisonOperator() bool {
	return false
}

func (k SyntaxKind) IsLogicalOperator() bool {
	return false
}

func (k SyntaxKind) IsBitwiseOperator() bool {
	return false
}

func (k SyntaxKind) IsPrefixOperator() bool {
	return false
}

func (k SyntaxKind) IsPostfixOperator() bool {
	return false
}
