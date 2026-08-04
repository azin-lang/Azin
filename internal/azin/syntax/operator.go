package syntax

type Associativity uint8

const (
	Left Associativity = iota
	Right
)

func (k Kind) BinaryPrecedence() int {
	switch k {
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

func (k Kind) UnaryPrecedence() int {
	switch k {
	case PlusToken, MinusToken, BangToken, CaretToken, AmpersandToken, StarToken:
		return 7
	default:
		return 0
	}
}

func (k Kind) Associativity() Associativity {
	if k.IsAssignmentOperator() {
		return Right
	}
	return Left
}

func (k Kind) IsOperator() bool {
	return k.IsArithmeticOperator() ||
		k.IsAssignmentOperator() ||
		k.IsComparisonOperator() ||
		k.IsLogicalOperator() ||
		k.IsBitwiseOperator()
}

func (k Kind) IsArithmeticOperator() bool {
	switch k {
	case PlusToken, MinusToken, StarToken, SlashToken:
		return true
	default:
		return false
	}
}

func (k Kind) IsAssignmentOperator() bool {
	switch k {
	case EqualsToken, PlusEqualsToken, MinusEqualsToken, StarEqualsToken, SlashEqualsToken:
		return true
	default:
		return false
	}
}

func (k Kind) IsComparisonOperator() bool {
	switch k {
	case EqualsEqualsToken, BangEqualsToken, LessToken, LessEqualsToken, GreaterToken, GreaterEqualsToken:
		return true
	default:
		return false
	}
}

func (k Kind) IsLogicalOperator() bool {
	switch k {
	case AmpersandAmpersandToken, PipePipeToken, BangToken:
		return true
	default:
		return false
	}
}

func (k Kind) IsBitwiseOperator() bool {
	switch k {
	case AmpersandToken, PipeToken, CaretToken:
		return true
	default:
		return false
	}
}

func (k Kind) IsPrefixOperator() bool {
	return k.UnaryPrecedence() > 0
}

func (k Kind) IsPostfixOperator() bool {
	return false
}
