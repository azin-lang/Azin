package syntax

type Associativity uint8

const (
	Left Associativity = iota
	Right
)

func BinaryPrecedence(kind SyntaxKind) int {
	switch kind {
	case StarToken, SlashToken, AmpersandToken:
		return 6
	case PlusToken, MinusToken, PipeToken, CaretToken:
		return 5
	case LessToken, LessEqualsToken, GreaterToken, GreaterEqualsToken:
		return 4
	case EqualsEqualsToken, BangEqualsToken:
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
	switch kind {
	case PlusToken, MinusToken, BangToken, CaretToken, AmpersandToken, StarToken:
		return 7
	default:
		return 0
	}
}

func AssociativityOf(kind SyntaxKind) Associativity {
	if kind.IsAssignmentOperator() {
		return Right
	}
	return Left
}

func (k SyntaxKind) IsOperator() bool {
	return k.IsArithmeticOperator() || k.IsAssignmentOperator() || k.IsComparisonOperator() || k.IsLogicalOperator() || k.IsBitwiseOperator()
}

func (k SyntaxKind) IsArithmeticOperator() bool {
	switch k {
	case PlusToken, MinusToken, StarToken, SlashToken:
		return true
	default:
		return false
	}
}

func (k SyntaxKind) IsAssignmentOperator() bool {
	switch k {
	case EqualsToken, PlusEqualsToken, MinusEqualsToken, StarEqualsToken, SlashEqualsToken:
		return true
	default:
		return false
	}
}

func (k SyntaxKind) IsComparisonOperator() bool {
	switch k {
	case EqualsEqualsToken, BangEqualsToken, LessToken, LessEqualsToken, GreaterToken, GreaterEqualsToken:
		return true
	default:
		return false
	}
}

func (k SyntaxKind) IsLogicalOperator() bool {
	switch k {
	case AmpersandAmpersandToken, PipePipeToken, BangToken:
		return true
	default:
		return false
	}
}

func (k SyntaxKind) IsBitwiseOperator() bool {
	switch k {
	case AmpersandToken, PipeToken, CaretToken:
		return true
	default:
		return false
	}
}

func (k SyntaxKind) IsPrefixOperator() bool {
	return UnaryPrecedence(k) > 0
}

func (k SyntaxKind) IsPostfixOperator() bool {
	return false
}
