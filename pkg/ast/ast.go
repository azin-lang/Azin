//nolint:unused
package ast

import (
	"fmt"
	"strconv"
	"strings"

	token2 "github.com/azin-lang/Azin/pkg/token"
	types2 "github.com/azin-lang/Azin/pkg/types"
)

// Node is the interface implemented by every AST node.
type Node interface {
	TokenLiteral() string
	Pos() token2.Position
	Label() string
}

// Expr represents an expression node.
type Expr interface {
	Node
	exprNode()
	Equals(other Expr) bool
	Type() *types2.TypeInfo
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

func (p *Program) Pos() token2.Position {
	if len(p.Statements) == 0 {
		return token2.Position{}
	}
	return p.Statements[0].Pos()
}

func (p *Program) Label() string {
	return "Program"
}

// Bad nodes

type BadExpr struct {
	Token token2.Token
}

func (*BadExpr) exprNode()              {}
func (b *BadExpr) TokenLiteral() string { return b.Token.Kind.String() }
func (b *BadExpr) Pos() token2.Position { return b.Token.Position }
func (*BadExpr) Label() string          { return "BadExpr" }
func (*BadExpr) Equals(other Expr) bool { return false }
func (*BadExpr) Type() *types2.TypeInfo { return types2.ErrorType() }

type BadStmt struct {
	Token token2.Token
}

func (*BadStmt) stmtNode()              {}
func (b *BadStmt) TokenLiteral() string { return b.Token.Kind.String() }
func (b *BadStmt) Pos() token2.Position { return b.Token.Position }
func (*BadStmt) Label() string          { return "BadStmt" }

// Statements

type VarStmt struct {
	Token token2.Token // var
	Name  *Identifier
	// SynType is the AST identifier for the type of the variable. It should not be modified.
	SynType *Identifier
	Value   Expr
	Mutable bool

	// SemaType contains the type information for the variable, set after semantic analysis.
	SemaType *types2.TypeInfo
}

func (*VarStmt) stmtNode()              {}
func (v *VarStmt) TokenLiteral() string { return v.Token.Kind.String() }
func (v *VarStmt) Pos() token2.Position { return v.Token.Position }
func (v *VarStmt) Label() string {
	s := "var "

	if v.Mutable {
		s += "mut "
	}

	s += v.Name.Value

	if v.SynType != nil {
		s += ": " + v.SynType.Value
	}

	return s
}

type AssignmentStmt struct {
	Token token2.Token // =
	Left  Expr
	Value Expr
}

func (*AssignmentStmt) stmtNode()              {}
func (a *AssignmentStmt) TokenLiteral() string { return a.Token.Kind.String() }
func (a *AssignmentStmt) Pos() token2.Position { return a.Left.Pos() }
func (*AssignmentStmt) Label() string {
	return "assign"
}

type StructStmt struct {
	Token  token2.Token // struct
	Name   *Identifier
	Fields []*FieldDecl

	// SemaType contains the type information for the struct, set after semantic analysis.
	SemaType *types2.TypeInfo
}

func (*StructStmt) stmtNode()              {}
func (s *StructStmt) TokenLiteral() string { return s.Token.Kind.String() }
func (s *StructStmt) Pos() token2.Position { return s.Token.Position }
func (s *StructStmt) Label() string {
	return "struct " + s.Name.Value
}

type EnumStmt struct {
	Token    token2.Token // enum
	Name     *Identifier
	Variants []*Identifier

	// SemaType contains the type information for the enum, set after semantic analysis.
	SemaType *types2.TypeInfo
}

func (*EnumStmt) stmtNode()              {}
func (e *EnumStmt) TokenLiteral() string { return e.Token.Kind.String() }
func (e *EnumStmt) Pos() token2.Position { return e.Token.Position }
func (e *EnumStmt) Label() string {
	return "enum " + e.Name.Value
}

type FuncStmt struct {
	Token  token2.Token // fn
	Name   *Identifier
	Params []*FieldDecl
	// SynReturnType is the AST identifier for the return type of the function, can be nil for void functions. It should not be modified.
	SynReturnType *Identifier
	Body          []Stmt
	CName         string

	// SemaReturnType is the return type of the function, set after semantic analysis.
	SemaReturnType *types2.TypeInfo
}

func (*FuncStmt) stmtNode()              {}
func (f *FuncStmt) TokenLiteral() string { return f.Token.Kind.String() }
func (f *FuncStmt) Pos() token2.Position { return f.Token.Position }
func (f *FuncStmt) Label() string {
	var s strings.Builder
	s.WriteString("fn " + f.Name.Value + "(")

	for i, p := range f.Params {
		if i != 0 {
			s.WriteString(", ")
		}

		if p.Mutable {
			s.WriteString("mut ")
		}

		s.WriteString(p.Name.Value)

		if p.SynType != nil {
			s.WriteString(": " + p.SynType.Value)
		}
	}

	s.WriteString(")")

	if f.SynReturnType != nil {
		s.WriteString(": " + f.SynReturnType.Value)
	}

	return s.String()
}

type ReturnStmt struct {
	Token token2.Token // return
	Value Expr
}

func (*ReturnStmt) stmtNode()              {}
func (r *ReturnStmt) TokenLiteral() string { return r.Token.Kind.String() }
func (r *ReturnStmt) Pos() token2.Position { return r.Token.Position }
func (*ReturnStmt) Label() string {
	return "return"
}

type IfStmt struct {
	Token     token2.Token // if
	Condition Expr
	Then      []Stmt
	Else      []Stmt
}

func (*IfStmt) stmtNode()              {}
func (i *IfStmt) TokenLiteral() string { return i.Token.Kind.String() }
func (i *IfStmt) Pos() token2.Position { return i.Token.Position }
func (*IfStmt) Label() string {
	return "if"
}

type WhileStmt struct {
	Token     token2.Token // while
	Condition Expr
	Body      []Stmt
}

func (*WhileStmt) stmtNode()              {}
func (w *WhileStmt) TokenLiteral() string { return w.Token.Kind.String() }
func (w *WhileStmt) Pos() token2.Position { return w.Token.Position }
func (*WhileStmt) Label() string {
	return "while"
}

type LoopStmt struct {
	Token token2.Token // loop
	Body  []Stmt
}

func (*LoopStmt) stmtNode()              {}
func (l *LoopStmt) TokenLiteral() string { return l.Token.Kind.String() }
func (l *LoopStmt) Pos() token2.Position { return l.Token.Position }
func (*LoopStmt) Label() string {
	return "loop"
}

type StopStmt struct {
	Token token2.Token // stop
}

func (*StopStmt) stmtNode()              {}
func (s *StopStmt) TokenLiteral() string { return s.Token.Kind.String() }
func (s *StopStmt) Pos() token2.Position { return s.Token.Position }
func (*StopStmt) Label() string {
	return "stop"
}

type DeferStmt struct {
	Token token2.Token // defer
	Call  Expr
}

func (*DeferStmt) stmtNode()              {}
func (d *DeferStmt) TokenLiteral() string { return d.Token.Kind.String() }
func (d *DeferStmt) Pos() token2.Position { return d.Token.Position }
func (*DeferStmt) Label() string {
	return "defer"
}

type ImportCStmt struct {
	Token token2.Token
	Path  *StringLiteral
}

func (*ImportCStmt) stmtNode()              {}
func (i *ImportCStmt) TokenLiteral() string { return i.Token.Kind.String() }
func (i *ImportCStmt) Pos() token2.Position { return i.Token.Position }
func (i *ImportCStmt) Label() string {
	return `importc "` + i.Path.Value + `"`
}

type ImportStmt struct {
	Token token2.Token
	Path  *StringLiteral
}

func (*ImportStmt) stmtNode()              {}
func (i *ImportStmt) TokenLiteral() string { return i.Token.Kind.String() }
func (i *ImportStmt) Pos() token2.Position { return i.Token.Position }
func (i *ImportStmt) Label() string {
	return `import "` + i.Path.Value + `"`
}

type ExpressionStmt struct {
	Token      token2.Token
	Expression Expr
}

func (*ExpressionStmt) stmtNode()              {}
func (e *ExpressionStmt) TokenLiteral() string { return e.Token.Kind.String() }
func (e *ExpressionStmt) Pos() token2.Position { return e.Expression.Pos() }
func (e *ExpressionStmt) Label() string {
	if e.Expression != nil {
		return e.Expression.Label()
	}
	return "expr"
}

// Declarations

type FieldDecl struct {
	Name *Identifier
	// SynType is the AST identifier for the type of the field. It should not be modified.
	SynType *Identifier
	Mutable bool

	// SemaType contains the type information for the field, set after semantic analysis.
	SemaType *types2.TypeInfo
}

func (f *FieldDecl) TokenLiteral() string { return f.Name.TokenLiteral() }
func (f *FieldDecl) Pos() token2.Position { return f.Name.Pos() }
func (f *FieldDecl) Label() string {
	s := ""

	if f.Mutable {
		s += "mut "
	}

	s += f.Name.Value

	if f.SynType != nil {
		s += ": " + f.SynType.Value
	}

	return s
}

// Expressions

type Identifier struct {
	Token token2.Token
	Value string

	// SemaType contains the type information for the identifier, set after semantic analysis.
	SemaType *types2.TypeInfo
}

func (*Identifier) exprNode()              {}
func (i *Identifier) TokenLiteral() string { return i.Value }
func (i *Identifier) Pos() token2.Position { return i.Token.Position }
func (i *Identifier) Type() *types2.TypeInfo {
	if i.SemaType != nil {
		return i.SemaType
	}

	return types2.UnknownType()
}
func (i *Identifier) Label() string {
	return i.Value
}
func (i *Identifier) Equals(other Expr) bool {
	o, ok := other.(*Identifier)
	return ok && i.Value == o.Value
}

type IntegerLiteral struct {
	Token token2.Token
	Value int64
}

func (*IntegerLiteral) exprNode()                {}
func (i *IntegerLiteral) TokenLiteral() string   { return fmt.Sprintf("%d", i.Value) }
func (i *IntegerLiteral) Pos() token2.Position   { return i.Token.Position }
func (i *IntegerLiteral) Type() *types2.TypeInfo { return types2.IntType() }
func (i *IntegerLiteral) Label() string {
	return strconv.FormatInt(i.Value, 10)
}
func (i *IntegerLiteral) Equals(other Expr) bool {
	o, ok := other.(*IntegerLiteral)
	return ok && i.Value == o.Value
}

type FloatLiteral struct {
	Token token2.Token
	Value float64
}

func (*FloatLiteral) exprNode()                {}
func (f *FloatLiteral) TokenLiteral() string   { return fmt.Sprintf("%f", f.Value) }
func (f *FloatLiteral) Pos() token2.Position   { return f.Token.Position }
func (f *FloatLiteral) Type() *types2.TypeInfo { return types2.FloatType() }
func (f *FloatLiteral) Label() string {
	return strconv.FormatFloat(f.Value, 'g', -1, 64)
}
func (f *FloatLiteral) Equals(other Expr) bool {
	o, ok := other.(*FloatLiteral)
	return ok && f.Value == o.Value
}

type StringLiteral struct {
	Token token2.Token
	Value string
}

func (*StringLiteral) exprNode()                {}
func (s *StringLiteral) TokenLiteral() string   { return s.Value }
func (s *StringLiteral) Pos() token2.Position   { return s.Token.Position }
func (s *StringLiteral) Type() *types2.TypeInfo { return types2.StringType() }
func (s *StringLiteral) Label() string {
	return strconv.Quote(s.Value)
}
func (s *StringLiteral) Equals(other Expr) bool {
	o, ok := other.(*StringLiteral)
	return ok && s.Value == o.Value
}

type CharacterLiteral struct {
	Token token2.Token
	Value rune
}

func (*CharacterLiteral) exprNode()                {}
func (c *CharacterLiteral) TokenLiteral() string   { return string(c.Value) }
func (c *CharacterLiteral) Pos() token2.Position   { return c.Token.Position }
func (c *CharacterLiteral) Type() *types2.TypeInfo { return types2.CharType() }
func (c *CharacterLiteral) Label() string {
	return strconv.QuoteRune(c.Value)
}
func (c *CharacterLiteral) Equals(other Expr) bool {
	o, ok := other.(*CharacterLiteral)
	return ok && c.Value == o.Value
}

type BooleanLiteral struct {
	Token token2.Token
	Value bool
}

func (*BooleanLiteral) exprNode() {}
func (b *BooleanLiteral) TokenLiteral() string {
	return strconv.FormatBool(b.Value)
}
func (b *BooleanLiteral) Pos() token2.Position   { return b.Token.Position }
func (b *BooleanLiteral) Type() *types2.TypeInfo { return types2.BoolType() }
func (b *BooleanLiteral) Label() string {
	return strconv.FormatBool(b.Value)
}
func (b *BooleanLiteral) Equals(other Expr) bool {
	o, ok := other.(*BooleanLiteral)
	return ok && b.Value == o.Value
}

type CallExpr struct {
	Callee       Expr
	Args         []Expr
	ResolvedName string

	// SemaReturnType contains the return type of the function call, set after semantic analysis.
	SemaReturnType *types2.TypeInfo
}

func (*CallExpr) exprNode()              {}
func (c *CallExpr) TokenLiteral() string { return c.Callee.TokenLiteral() }
func (c *CallExpr) Pos() token2.Position { return c.Callee.Pos() }
func (c *CallExpr) Type() *types2.TypeInfo {
	if c.SemaReturnType != nil {
		return c.SemaReturnType
	}

	return types2.UnknownType()
}
func (c *CallExpr) Label() string {
	switch callee := c.Callee.(type) {
	case *Identifier:
		return "call " + callee.Value

	case *MemberExpr:
		return "call " + callee.Label()

	default:
		return "call"
	}
}
func (c *CallExpr) Equals(other Expr) bool {
	o, ok := other.(*CallExpr)
	if !ok || len(c.Args) != len(o.Args) || !c.Callee.Equals(o.Callee) {
		return false
	}
	for i := range c.Args {
		if !c.Args[i].Equals(o.Args[i]) {
			return false
		}
	}
	return true
}

type BinaryExpr struct {
	Left     Expr
	Operator token2.Token
	Right    Expr

	// SemaResultType contains the result type of the binary operation, set after semantic analysis.
	SemaResultType *types2.TypeInfo
}

func (*BinaryExpr) exprNode()                {}
func (b *BinaryExpr) TokenLiteral() string   { return b.Operator.Kind.String() }
func (b *BinaryExpr) Pos() token2.Position   { return b.Left.Pos() }
func (b *BinaryExpr) Type() *types2.TypeInfo { return b.SemaResultType }
func (b *BinaryExpr) Label() string {
	return b.Operator.Kind.String()
}
func (b *BinaryExpr) Equals(other Expr) bool {
	o, ok := other.(*BinaryExpr)
	return ok && b.Operator.Kind == o.Operator.Kind &&
		b.Left.Equals(o.Left) &&
		b.Right.Equals(o.Right)
}

type MemberExpr struct {
	Object   Expr
	Property *Identifier

	// SemaResultType contains the type of the member access, set after semantic analysis.
	SemaResultType *types2.TypeInfo
}

func (*MemberExpr) exprNode()              {}
func (m *MemberExpr) TokenLiteral() string { return m.Property.TokenLiteral() }
func (m *MemberExpr) Pos() token2.Position { return m.Object.Pos() }
func (m *MemberExpr) Type() *types2.TypeInfo {
	if m.SemaResultType != nil {
		return m.SemaResultType
	}

	return types2.UnknownType()
}
func (m *MemberExpr) Label() string {
	if id, ok := m.Object.(*Identifier); ok {
		return id.Value + "." + m.Property.Value
	}

	return "." + m.Property.Value
}
func (m *MemberExpr) Equals(other Expr) bool {
	o, ok := other.(*MemberExpr)
	return ok && m.Object.Equals(o.Object) && m.Property.Equals(o.Property)
}
