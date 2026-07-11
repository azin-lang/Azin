package codegen

import (
	"bytes"
	"fmt"

	"github.com/azin-lang/Azin/internal/ast"
)

// the Transpiler struct is responsible for transpiling the AST to C code.
type Transpiler struct {
	buf bytes.Buffer
}

// create a new Transpiler.
func New() *Transpiler {
	return &Transpiler{}
}

// Transpile transpiles the AST to C code.
func (t *Transpiler) Transpile(program *ast.Program) string {
	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.StructStmt); ok {
			t.compileStruct(s)
		}
	}

	for _, stmt := range program.Statements {
		if _, ok := stmt.(*ast.StructStmt); !ok {
			t.compileStatement(stmt)
		}
	}

	return t.buf.String()
}

func (t *Transpiler) compileStruct(s *ast.StructStmt) {
	fmt.Fprintf(&t.buf, "typedef struct {\n")
	for _, f := range s.Fields {
		typeStr := f.Type.Value
		if typeStr == "int" {
			typeStr = "int"
		}
		fmt.Fprintf(&t.buf, "    %s %s;\n", typeStr, f.Name.Value)
	}
	fmt.Fprintf(&t.buf, "} %s;\n\n", s.Name.Value)
}

func (t *Transpiler) compileStatement(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.FuncStmt:
		t.compileFunc(s)
	case *ast.ReturnStmt:
		t.buf.WriteString("    return ")
		t.compileExpression(s.Value)
		t.buf.WriteString(";\n")
	}
}

func (t *Transpiler) compileFunc(f *ast.FuncStmt) {
	retType := f.ReturnType.Value

	fmt.Fprintf(&t.buf, "%s %s(", retType, f.Name.Value)
	for i, p := range f.Params {
		fmt.Fprintf(&t.buf, "%s %s", p.Type.Value, p.Name.Value)
		if i < len(f.Params)-1 {
			t.buf.WriteString(", ")
		}
	}
	t.buf.WriteString(") {\n")

	for _, stmt := range f.Body {
		t.compileStatement(stmt)
	}

	t.buf.WriteString("}\n")
}

func (t *Transpiler) compileExpression(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.Identifier:
		t.buf.WriteString(e.Value)
	case *ast.MemberExpr:
		t.compileExpression(e.Object)
		t.buf.WriteString(".")
		t.buf.WriteString(e.Property.Value)
	case *ast.BinaryExpr:
		t.compileExpression(e.Left)
		opStr := e.Operator.Kind.String()
		if opStr == "plus" {
			opStr = "+"
		} else if opStr == "minus" {
			opStr = "-"
		}
		fmt.Fprintf(&t.buf, " %s ", opStr)
		t.compileExpression(e.Right)
	case *ast.CallExpr:
		t.buf.WriteString(e.Function.Value)
		t.buf.WriteString("(")
		for i, arg := range e.Args {
			t.compileExpression(arg)
			if i < len(e.Args)-1 {
				t.buf.WriteString(", ")
			}
		}
		t.buf.WriteString(")")
	}
}
