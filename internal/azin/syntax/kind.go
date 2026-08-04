package syntax

//go:generate stringer -type=SyntaxKind

type SyntaxKind uint16

const (
	// Special
	Unknown SyntaxKind = iota
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

func (k SyntaxKind) IsTrivia() bool  { return k >= firstTrivia && k <= lastTrivia }
func (k SyntaxKind) IsToken() bool   { return k >= firstToken && k <= lastToken }
func (k SyntaxKind) IsKeyword() bool { return k >= firstKeyword && k <= lastKeyword }
func (k SyntaxKind) IsLiteral() bool { return k >= firstLiteral && k <= lastLiteral }

func (k SyntaxKind) IsNode() bool {
	return k >= CompilationUnit && k <= CastExpression
}

func (k SyntaxKind) IsExpression() bool {
	return k >= LiteralExpression && k <= CastExpression
}

func (k SyntaxKind) IsStatement() bool {
	return k >= BlockStatement && k <= ForStatement
}

func (k SyntaxKind) IsDeclaration() bool {
	return k >= ImportDeclaration && k <= FieldDeclaration
}

func (k SyntaxKind) IsPrimaryExpression() bool {
	return k == LiteralExpression || k == NameExpression || k == ParenthesizedExpression
}

func (k SyntaxKind) StartsExpression() bool {
	return k.IsLiteral() || k == IdentifierToken || k == OpenParenToken || k.IsPrefixOperator()
}

func (k SyntaxKind) StartsStatement() bool {
	switch k {
	case OpenBraceToken, KeywordIf, KeywordWhile, KeywordFor, KeywordLoop, KeywordDefer, KeywordStop, KeywordReturn, SemicolonToken:
		return true
	default:
		return k.StartsExpression()
	}
}

func (k SyntaxKind) StartsDeclaration() bool {
	switch k {
	case KeywordImport, KeywordImportC, KeywordFn, KeywordStruct, KeywordEnum, KeywordType, KeywordVar:
		return true
	default:
		return false
	}
}
