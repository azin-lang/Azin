package ast

import "github.com/azin-lang/Azin/internal/token"

// Node is the interface implemented by every AST node.
type Node interface {
	TokenLiteral() string
}

// Expr represents an expression node.
type Expr interface {
	Node
	exprNode()
}

// Stmt represents a statement node.
type Stmt interface {
	Node
	stmtNode()
}

// Program is the root of the AST.
type Program struct {
	Statements []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

//
// Statements
//

// StructStmt represents a struct declaration.
type StructStmt struct {
	Token  token.Token // struct
	Name   *Identifier
	Fields []*FieldDecl
}

func (s *StructStmt) TokenLiteral() string {
	return s.Token.Kind.String()
}

func (*StructStmt) stmtNode() {}

// FuncStmt represents a function declaration.
type FuncStmt struct {
	Token      token.Token // fn
	Name       *Identifier
	Params     []*FieldDecl
	ReturnType *Identifier
	Body       []Stmt
}

func (f *FuncStmt) TokenLiteral() string {
	return f.Token.Kind.String()
}

func (*FuncStmt) stmtNode() {}

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Token token.Token // return
	Value Expr
}

func (r *ReturnStmt) TokenLiteral() string {
	return r.Token.Kind.String()
}

func (*ReturnStmt) stmtNode() {}

// IfStmt represents an if/else statement.
type IfStmt struct {
	Token     token.Token // if
	Condition Expr
	Then      []Stmt
	Else      []Stmt
}

func (i *IfStmt) TokenLiteral() string {
	return i.Token.Kind.String()
}

func (*IfStmt) stmtNode() {}

// ImportStmt represents an import statement.
type ImportCStmt struct {
	Token token.Token
	Path  *StringLiteral
}

func (i *ImportCStmt) TokenLiteral() string { return i.Token.Kind.String() }
func (i *ImportCStmt) stmtNode()            {}

// ExpressionStmt represents an expression used as a statement.
//
// Example:
//
//	printf("hello");
//	foo();
type ExpressionStmt struct {
	Token      token.Token
	Expression Expr
}

func (e *ExpressionStmt) TokenLiteral() string {
	return e.Token.Kind.String()
}

func (*ExpressionStmt) stmtNode() {}

//
// Declarations
//

// FieldDecl represents either a parameter declaration
// or a struct field declaration.
type FieldDecl struct {
	Name *Identifier
	Type *Identifier
}

func (f *FieldDecl) TokenLiteral() string {
	return f.Name.TokenLiteral()
}

//
// Expressions
//

// Identifier represents an identifier.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Kind.String()
}

func (*Identifier) exprNode() {}

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Kind.String()
}

func (*IntegerLiteral) exprNode() {}

// FloatLiteral represents a floating point literal.
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (i *FloatLiteral) TokenLiteral() string {
	return i.Token.Kind.String()
}

func (*FloatLiteral) exprNode() {}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Kind.String()
}

func (*StringLiteral) exprNode() {}

// CharacterLiteral represents a character literal.
type CharacterLiteral struct {
	Token token.Token
	Value rune
}

func (c *CharacterLiteral) TokenLiteral() string {
	return c.Token.Kind.String()
}

func (*CharacterLiteral) exprNode() {}

// CallExpr represents a function call.
type CallExpr struct {
	Callee Expr
	Args   []Expr
}

func (c *CallExpr) TokenLiteral() string {
	return c.Callee.TokenLiteral()
}

func (*CallExpr) exprNode() {}

// BinaryExpr represents a binary operator expression.
type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *BinaryExpr) TokenLiteral() string {
	return b.Operator.Kind.String()
}

func (*BinaryExpr) exprNode() {}

// MemberExpr represents a member access expression.
//
// Example:
//
//	person.name
type MemberExpr struct {
	Object   Expr
	Property *Identifier
}

func (m *MemberExpr) TokenLiteral() string {
	return m.Property.TokenLiteral()
}

func (*MemberExpr) exprNode() {}

type VarStmt struct {
	Token token.Token // var
	Name  *Identifier
	Type  *Identifier
	Value Expr
}

func (*VarStmt) stmtNode() {}

func (v *VarStmt) TokenLiteral() string {
	return v.Token.Kind.String()
}
