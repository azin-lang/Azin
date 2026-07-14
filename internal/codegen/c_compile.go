package codegen

import (
	"fmt"
	"strconv"

	"github.com/azin-lang/Azin/internal/ast"
)

func (t *Transpiler) compileStruct(s *ast.StructStmt) {
	t.write("typedef struct {\n")
	t.pushIndent()

	for _, field := range s.Fields {
		t.writeIndent()
		t.printf("%s %s;\n", emitType(field.Type.Value), field.Name.Value)
	}

	t.popIndent()
	t.printf("} %s;\n", s.Name.Value)
	t.newline()
}

func (t *Transpiler) compileStatement(stmt ast.Stmt) {
	switch n := stmt.(type) {

	case *ast.StructStmt:
		// already emitted

	case *ast.ImportCStmt:
		// already emitted

	case *ast.IfStmt:
		t.compileIf(n)

	case *ast.LoopStmt:
		t.compileLoop(n)

	case *ast.StopStmt:
		t.writeIndent()
		t.write("break;")
		t.newline()

	case *ast.FuncStmt:
		t.compileFunc(n)

	case *ast.ReturnStmt:
		t.writeIndent()
		t.write("return")

		if n.Value != nil {
			t.write(" ")
			t.compileExpression(n.Value)
		}

		t.write(";")
		t.newline()

	case *ast.ExpressionStmt:
		t.writeIndent()
		t.compileExpression(n.Expression)
		t.write(";")
		t.newline()

	case *ast.VarStmt:
		t.writeIndent()

		if n.Type == nil {
			panic("internal compiler error: variable '" + n.Name.Value + "' has no resolved type")
		}

		if !n.Mutable {
			t.write("const ")
		}

		t.printf("%s %s", emitType(n.Type.Value), n.Name.Value)

		if n.Value != nil {
			t.write(" = ")
			t.compileExpression(n.Value)
		}

		t.write(";")
		t.newline()

	case *ast.AssignmentStmt:
		t.writeIndent()
		t.compileExpression(n.Left)
		t.write(" = ")
		t.compileExpression(n.Value)
		t.write(";")
		t.newline()

	default:
		panic(fmt.Sprintf("unsupported statement %T", stmt))
	}
}

func (t *Transpiler) compileFunc(fn *ast.FuncStmt) {
	if fn.ReturnType == nil {
		panic("internal compiler error: function has no resolved return type")
	}

	t.printf("%s %s(", emitType(fn.ReturnType.Value), fn.Name.Value)

	for i, p := range fn.Params {
		if p.Type == nil {
			panic("internal compiler error: parameter has no resolved type")
		}

		if i > 0 {
			t.write(", ")
		}

		t.printf("%s %s", emitType(p.Type.Value), p.Name.Value)
	}

	t.write(")")
	t.write(" {\n")

	t.pushIndent()

	for _, stmt := range fn.Body {
		t.compileStatement(stmt)
	}

	t.popIndent()

	t.write("}\n")
}

func (t *Transpiler) compileExpression(expr ast.Expr) {
	switch n := expr.(type) {

	case *ast.Identifier:
		t.write(n.Value)

	case *ast.BooleanLiteral:
		if n.Value {
			t.write("true")
		} else {
			t.write("false")
		}

	case *ast.IntegerLiteral:
		t.printf("%d", n.Value)

	case *ast.FloatLiteral:
		t.printf("%s", strconv.FormatFloat(n.Value, 'g', -1, 64))

	case *ast.StringLiteral:
		t.printf("%q", n.Value)

	case *ast.CharacterLiteral:
		t.printf("%q", n.Value)

	case *ast.MemberExpr:
		t.compileExpression(n.Object)
		t.write(".")
		t.write(n.Property.Value)

	case *ast.BinaryExpr:
		t.compileExpression(n.Left)
		t.write(" ")
		t.write(emitOperator(n.Operator.Kind))
		t.write(" ")
		t.compileExpression(n.Right)

	case *ast.CallExpr:
		t.compileExpression(n.Callee)
		t.write("(")

		for i, arg := range n.Args {
			if i > 0 {
				t.write(", ")
			}
			t.compileExpression(arg)
		}

		t.write(")")

	default:
		panic(fmt.Sprintf("unsupported expression %T", expr))
	}
}

func (t *Transpiler) compileIf(n *ast.IfStmt) {
	t.writeIndent()
	t.write("if (")

	t.compileExpression(n.Condition)

	t.write(") {\n")

	t.pushIndent()

	for _, stmt := range n.Then {
		t.compileStatement(stmt)
	}

	t.popIndent()

	t.writeIndent()
	t.write("}")

	if len(n.Else) > 0 {
		t.write(" else {\n")

		t.pushIndent()

		for _, stmt := range n.Else {
			t.compileStatement(stmt)
		}

		t.popIndent()

		t.writeIndent()
		t.write("}")
	}

	t.newline()
}

func (t *Transpiler) compileLoop(n *ast.LoopStmt) {
	t.writeIndent()
	t.write("for (;;) {\n")

	t.pushIndent()

	for _, stmt := range n.Body {
		t.compileStatement(stmt)
	}

	t.popIndent()

	t.writeIndent()
	t.write("}")
	t.newline()
}
