package syntax

type SyntaxKind uint16

const (
	Unknown SyntaxKind = iota
	EndOfFile
	Error

	WhitespaceTrivia
	EndOfLineTrivia
	SingleLineCommentTrivia

	CompilationUnit
	FunctionDeclaration
	Parameter
	BlockStatement
	StopStatement
	IfStatement
	WhileStatement

	LiteralExpression
	BinaryExpression
	UnaryExpression
	IdentifierExpression

	IdentifierToken
	IntegerLiteralToken
	FloatLiteralToken
	CharacterLiteralToken
	BooleanLiteralToken
	StringLiteralToken

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

	PlusToken
	PlusEqualsToken
	MinusToken
	MinusEqualsToken
	EqualsToken
	EqualsEqualsToken
	StarToken
	StarEqualsToken

	OpenParenToken
	CloseParenToken
	OpenBracketToken
	CloseBracketToken
	OpenBraceToken
	CloseBraceToken
	CommaToken
	SemicolonToken
	ColonToken
	DotToken
	NewlineToken
)
