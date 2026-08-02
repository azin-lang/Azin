package syntax

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

	// Declarations

	FunctionDeclaration
	StructDeclaration
	EnumDeclaration
	TypeDeclaration
	VariableDeclaration
	Parameter

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

	BreakStatement
	ContinueStatement

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

	PercentToken
	PercentEqualsToken

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

	// Delimiters

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
