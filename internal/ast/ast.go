package ast

import "github.com/azin-lang/Azin/internal/token"

// Node is the interface for all AST nodes.
type Node interface {
	TokenLiteral() string
}

// Expr is the interface for all expression nodes.
type Expr interface {
	Node
	exprNode()
}

// Stmt is the interface for all statement nodes.
type Stmt interface {
	Node
	stmtNode()
}

// Program is the root node of the AST.
type Program struct {
	Statements []Stmt
}

// TokenLiteral returns the token literal of the program.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// StructStmt represents a struct statement node.
type StructStmt struct {
	Token  token.Token // kw_struct
	Name   *Identifier
	Fields []*FieldDecl
}

// TokenLiteral returns the token literal of the struct statement.
func (s *StructStmt) TokenLiteral() string {
	return s.Token.Kind.String()
}

func (s *StructStmt) stmtNode() {}

// FieldDecl represents a field declaration.
type FieldDecl struct {
	Name *Identifier
	Type *Identifier
}

// TokenLiteral returns the token literal of the field declaration.
func (f *FieldDecl) TokenLiteral() string {
	return f.Name.TokenLiteral()
}

func (f *FieldDecl) exprNode() {}

// FuncStmt represents a function statement node.
type FuncStmt struct {
	Token      token.Token // kw_fn
	Name       *Identifier
	Params     []*FieldDecl
	ReturnType *Identifier
	Body       []Stmt
}

// TokenLiteral returns the token kind of the function statement.
func (f *FuncStmt) TokenLiteral() string {
	return f.Token.Kind.String()
}

func (f *FuncStmt) stmtNode() {}

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Token token.Token // kw_return
	Value Expr
}

// TokenLiteral returns the token kind of the return statement.
func (r *ReturnStmt) TokenLiteral() string {
	return r.Token.Kind.String()
}

func (r *ReturnStmt) stmtNode() {}

// Identifier represents an identifier.
type Identifier struct {
	Token token.Token
	Value string
}

// TokenLiteral returns the token literal of the identifier.
func (i *Identifier) TokenLiteral() string {
	return i.Token.Kind.String()
}

func (i *Identifier) exprNode() {}

// CallExpr represents a call expression.
type CallExpr struct {
	Function *Identifier
	Args     []Expr
}

// TokenLiteral returns the token literal of the call expression.
func (c *CallExpr) TokenLiteral() string {
	return c.Function.TokenLiteral()
}

func (c *CallExpr) exprNode() {}

// BinaryExpr represents a binary expression.
type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

// TokenLiteral returns the token literal of the binary expression.
func (b *BinaryExpr) TokenLiteral() string {
	return b.Operator.Kind.String()
}

func (b *BinaryExpr) exprNode() {}

// MemberExpr represents a member expression.
type MemberExpr struct {
	Object   Expr
	Property *Identifier
}

// TokenLiteral returns the token literal of the member expression.
func (m *MemberExpr) TokenLiteral() string {
	return m.Property.TokenLiteral()
}

func (m *MemberExpr) exprNode() {}
