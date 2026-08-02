package syntax

type SyntaxKind uint16

const (
	// Unknown or special

	Unknown SyntaxKind = iota
	EndOfFileToken
	BadToken

	// Trivia

	WhitespaceTrivia
	EndOfLineTrivia
	SingleLineCommentTrivia

	// Nodes

	CompilationUnit

	FunctionDeclaration
	Parameter

	BlockStatement
	ExpressionStatement
	ReturnStatement
	StopStatement
	IfStatement
	WhileStatement
	ForStatement

	LiteralExpression
	NameExpression
	UnaryExpression
	BinaryExpression
	ParenthesizedExpression
	AssignmentExpression
	CallExpression

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

	LessToken
	LessEqualsToken

	GreaterToken
	GreaterEqualsToken

	AmpersandToken
	AmpersandAmpersandToken

	PipeToken
	PipePipeToken

	CaretToken

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
