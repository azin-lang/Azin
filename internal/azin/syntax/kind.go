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
