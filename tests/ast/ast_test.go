package ast_test

import (
	"strings"
	"testing"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func tok(kind token.Kind, offset, length uint32) token.Token {
	return token.Token{Kind: kind, Position: token.Position{Offset: offset}, Length: length}
}

func ident(value string) *ast.Identifier {
	return &ast.Identifier{
		Token: tok(token.Identifier, 0, uint32(len(value))), //nolint:gosec
		Value: value,
	}
}

func TestProgram(t *testing.T) {
	p := &ast.Program{Statements: nil}
	if p.TokenLiteral() != "" {
		t.Errorf("empty program TokenLiteral = %q", p.TokenLiteral())
	}
	if p.Pos() != (token.Position{}) {
		t.Errorf("empty program Pos = %v", p.Pos())
	}
	if p.Label() != "Program" {
		t.Errorf("Label = %q", p.Label())
	}
}

func TestBadNodes(t *testing.T) {
	be := &ast.BadExpr{Token: tok(token.Error, 0, 1)}
	if be.Label() != "BadExpr" {
		t.Errorf("BadExpr Label = %q", be.Label())
	}

	bs := &ast.BadStmt{Token: tok(token.Error, 1, 2)}
	if bs.Label() != "BadStmt" {
		t.Errorf("BadStmt Label = %q", bs.Label())
	}
	if bs.Pos() != (token.Position{Offset: 1}) {
		t.Errorf("BadStmt Pos = %v", bs.Pos())
	}
}

func TestVarStmt(t *testing.T) {
	v := &ast.VarStmt{
		Token:   tok(token.KwVar, 0, 3),
		Name:    ident("x"),
		Type:    ident("int"),
		Mutable: true,
	}
	if !strings.Contains(v.Label(), "var") {
		t.Errorf("Label missing 'var': %q", v.Label())
	}
	if !strings.Contains(v.Label(), "mut") {
		t.Errorf("Label missing 'mut': %q", v.Label())
	}
	if !strings.Contains(v.Label(), "x") {
		t.Errorf("Label missing 'x': %q", v.Label())
	}
	if !strings.Contains(v.Label(), "int") {
		t.Errorf("Label missing 'int': %q", v.Label())
	}
}

func TestFuncStmt(t *testing.T) {
	f := &ast.FuncStmt{
		Token:      tok(token.KwFn, 0, 2),
		Name:       ident("add"),
		Params:     []*ast.FieldDecl{{Name: ident("a"), Type: ident("int")}},
		ReturnType: ident("int"),
	}
	label := f.Label()
	if !strings.Contains(label, "add") {
		t.Errorf("Label missing 'add': %q", label)
	}
	if !strings.Contains(label, "int") {
		t.Errorf("Label missing return type: %q", label)
	}
	if !strings.Contains(label, "a") {
		t.Errorf("Label missing param: %q", label)
	}
}

func TestIfStmt(t *testing.T) {
	s := &ast.IfStmt{
		Token:     tok(token.KwIf, 0, 2),
		Condition: ident("true"),
	}
	if s.Label() != "if" {
		t.Errorf("Label = %q", s.Label())
	}
}

func TestLoopStmt(t *testing.T) {
	s := &ast.LoopStmt{
		Token: tok(token.KwLoop, 0, 4),
	}
	if s.Label() != "loop" {
		t.Errorf("Label = %q", s.Label())
	}
}

func TestReturnStmt(t *testing.T) {
	s := &ast.ReturnStmt{
		Token: tok(token.KwReturn, 0, 6),
	}
	if s.Label() != "return" {
		t.Errorf("Label = %q", s.Label())
	}
}

func TestStructStmt(t *testing.T) {
	s := &ast.StructStmt{
		Token: tok(token.KwStruct, 0, 6),
		Name:  ident("Point"),
	}
	if s.Label() != "struct Point" {
		t.Errorf("Label = %q, want 'struct Point'", s.Label())
	}
}

func TestIdentExpr(t *testing.T) {
	id := ident("foobar")
	if id.Label() != "foobar" {
		t.Errorf("Label = %q", id.Label())
	}
	if id.TokenLiteral() != "foobar" {
		t.Errorf("TokenLiteral = %q", id.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	lit := &ast.IntegerLiteral{Token: tok(token.IntegerLiteral, 0, 2), Value: 42}
	if lit.Label() != "42" {
		t.Errorf("Label = %q", lit.Label())
	}
}

func TestFloatLiteral(t *testing.T) {
	lit := &ast.FloatLiteral{Token: tok(token.FloatLiteral, 0, 4), Value: 3.14}
	if lit.Label() != "3.14" {
		t.Errorf("Label = %q", lit.Label())
	}
}

func TestStringLiteral(t *testing.T) {
	lit := &ast.StringLiteral{Token: tok(token.StringLiteral, 0, 5), Value: "hello"}
	if !strings.Contains(lit.Label(), "hello") {
		t.Errorf("Label missing 'hello': %q", lit.Label())
	}
}

func TestCharacterLiteral(t *testing.T) {
	lit := &ast.CharacterLiteral{Token: tok(token.CharacterLiteral, 0, 3), Value: 'x'}
	if !strings.Contains(lit.Label(), "x") {
		t.Errorf("Label missing 'x': %q", lit.Label())
	}
}

func TestBooleanLiteral(t *testing.T) {
	tLit := &ast.BooleanLiteral{Token: tok(token.Identifier, 0, 4), Value: true}
	if tLit.Label() != "true" {
		t.Errorf("true literal Label = %q", tLit.Label())
	}
	fLit := &ast.BooleanLiteral{Token: tok(token.Identifier, 0, 5), Value: false}
	if fLit.Label() != "false" {
		t.Errorf("false literal Label = %q", fLit.Label())
	}
}

func TestCallExpr(t *testing.T) {
	call := &ast.CallExpr{
		Callee: ident("foo"),
		Args:   []ast.Expr{},
	}
	if call.Label() != "call foo" {
		t.Errorf("Label = %q", call.Label())
	}

	callMember := &ast.CallExpr{
		Callee: &ast.MemberExpr{
			Object:   ident("obj"),
			Property: ident("method"),
		},
	}
	if !strings.Contains(callMember.Label(), "call") {
		t.Errorf("Label missing 'call': %q", callMember.Label())
	}
}

func TestMemberExpr(t *testing.T) {
	m := &ast.MemberExpr{
		Object:   ident("point"),
		Property: ident("x"),
	}
	if m.Label() != "point.x" {
		t.Errorf("Label = %q, want 'point.x'", m.Label())
	}
}

func TestBinaryExpr(t *testing.T) {
	b := &ast.BinaryExpr{
		Left:     &ast.IntegerLiteral{Value: 1},
		Operator: tok(token.Plus, 0, 1),
		Right:    &ast.IntegerLiteral{Value: 2},
	}
	// Label returns the Kind.String(), e.g. "plus" not "+"
	if b.Label() == "" {
		t.Errorf("Label = empty")
	}
}

func TestFieldDecl(t *testing.T) {
	f := &ast.FieldDecl{
		Name: ident("name"),
		Type: ident("string"),
	}
	if !strings.Contains(f.Label(), "name") {
		t.Errorf("Label missing 'name': %q", f.Label())
	}
	if !strings.Contains(f.Label(), "string") {
		t.Errorf("Label missing 'string': %q", f.Label())
	}
}

func TestAssignmentStmt(t *testing.T) {
	a := &ast.AssignmentStmt{
		Token: tok(token.Equal, 0, 1),
		Left:  ident("x"),
		Value: &ast.IntegerLiteral{Value: 42},
	}
	if a.Label() != "assign" {
		t.Errorf("Label = %q, want 'assign'", a.Label())
	}
}

func TestImportCStmt(t *testing.T) {
	i := &ast.ImportCStmt{
		Token: tok(token.KwImportC, 0, 7),
		Path:  &ast.StringLiteral{Value: "stdio.h"},
	}
	if !strings.Contains(i.Label(), "stdio.h") {
		t.Errorf("Label missing stdio.h: %q", i.Label())
	}
}

func TestExpressionStmt(t *testing.T) {
	e := &ast.ExpressionStmt{
		Token:      tok(token.Identifier, 0, 4),
		Expression: ident("test"),
	}
	if e.Label() != "test" {
		t.Errorf("Label = %q, want 'test'", e.Label())
	}

	nilExpr := &ast.ExpressionStmt{}
	if nilExpr.Label() != "expr" {
		t.Errorf("nil expr Label = %q, want 'expr'", nilExpr.Label())
	}
}
