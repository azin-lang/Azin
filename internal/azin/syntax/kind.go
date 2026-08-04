package syntax

//go:generate stringer -type=Kind

type Kind uint16

const (
	// Special
	Unknown Kind = iota
	BadToken
	EndOfFileToken
	MissingToken

	// Trivia
	WhitespaceTrivia
	EndOfLineTrivia
	SingleLineCommentTrivia
	MultiLineCommentTrivia
	SkippedTextTrivia

	// Root
	CompilationUnit
	SyntaxList
	SeparatedSyntaxList

	// Declarations
	ImportDeclaration
	FunctionDeclaration
	StructDeclaration
	EnumDeclaration
	TypeDeclaration
	Parameter
	VariableDeclaration
	FieldDeclaration

	// Statements
	BlockStatement
	EmptyStatement
	ExpressionStatement
	ReturnStatement
	StopStatement
	DeferStatement
	IfStatement
	WhileStatement
	ForStatement

	// Expressions
	LiteralExpression
	NameExpression
	ParenthesizedExpression
	UnaryExpression
	BinaryExpression
	AssignmentExpression
	CallExpression
	MemberAccessExpression
	IndexExpression

	// Future
	TupleExpression
	CastExpression

	// Tokens
	IdentifierToken
	IntegerLiteralToken
	FloatLiteralToken
	CharacterLiteralToken
	BooleanLiteralToken
	StringLiteralToken

	// Keywords
	KeywordDefer
	KeywordDo
	KeywordElse
	KeywordEnd
	KeywordEnum
	KeywordFalse
	KeywordFn
	KeywordFor
	KeywordIf
	KeywordImport
	KeywordImportC
	KeywordIs
	KeywordLoop
	KeywordMut
	KeywordReturn
	KeywordStop
	KeywordStruct
	KeywordThen
	KeywordTrue
	KeywordType
	KeywordVar
	KeywordWhile

	// Operators
	PlusToken
	PlusEqualsToken
	MinusToken
	MinusEqualsToken
	StarToken
	StarEqualsToken
	SlashToken
	SlashEqualsToken
	EqualsToken
	EqualsEqualsToken
	BangEqualsToken
	LessToken
	LessEqualsToken
	GreaterToken
	GreaterEqualsToken
	BangToken
	AmpersandToken
	AmpersandAmpersandToken
	PipeToken
	PipePipeToken
	CaretToken
	QuestionToken

	// Punctuation
	OpenParenToken
	CloseParenToken
	OpenBracketToken
	CloseBracketToken
	OpenBraceToken
	CloseBraceToken
	CommaToken
	ColonToken
	SemicolonToken
	DotToken
	NewlineToken
)

const (
	firstTrivia = WhitespaceTrivia
	lastTrivia  = SkippedTextTrivia

	firstToken = IdentifierToken
	lastToken  = NewlineToken

	firstKeyword = KeywordDefer
	lastKeyword  = KeywordWhile

	firstLiteral = IntegerLiteralToken
	lastLiteral  = StringLiteralToken
)

func (k Kind) IsTrivia() bool  { return k >= firstTrivia && k <= lastTrivia }
func (k Kind) IsToken() bool   { return k >= firstToken && k <= lastToken }
func (k Kind) IsKeyword() bool { return k >= firstKeyword && k <= lastKeyword }
func (k Kind) IsLiteral() bool { return k >= firstLiteral && k <= lastLiteral }

func (k Kind) IsNode() bool {
	return k >= CompilationUnit && k <= CastExpression
}

func (k Kind) IsExpression() bool {
	return k >= LiteralExpression && k <= CastExpression
}

func (k Kind) IsStatement() bool {
	return k >= BlockStatement && k <= ForStatement
}

func (k Kind) IsDeclaration() bool {
	return k >= ImportDeclaration && k <= FieldDeclaration
}

func (k Kind) IsPrimaryExpression() bool {
	return k == LiteralExpression || k == NameExpression || k == ParenthesizedExpression
}

func (k Kind) StartsExpression() bool {
	return k.IsLiteral() || k == IdentifierToken || k == OpenParenToken || k.IsPrefixOperator()
}

func (k Kind) StartsStatement() bool {
	switch k {
	case OpenBraceToken, KeywordIf, KeywordWhile, KeywordFor, KeywordLoop, KeywordDefer, KeywordStop, KeywordReturn, SemicolonToken:
		return true
	default:
		return k.StartsExpression()
	}
}

func (k Kind) StartsDeclaration() bool {
	switch k {
	case KeywordImport, KeywordImportC, KeywordFn, KeywordStruct, KeywordEnum, KeywordType, KeywordVar:
		return true
	default:
		return false
	}
}
