package c

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/ast"
)

func (t *Transpiler) emitStatement(stmt ast.Stmt) {
	switch n := stmt.(type) {

	case *ast.StructStmt:
		// already emitted

	case *ast.EnumStmt:
		// already emitted

	case *ast.ImportCStmt:
		// already emitted

	case *ast.IfStmt:
		t.emitIf(n)

	case *ast.LoopStmt:
		t.emitLoop(n)

	case *ast.StopStmt:
		t.writeIndent()
		t.write("break;")
		t.newline()

	case *ast.FuncStmt:
		t.emitFunction(n)

	case *ast.ReturnStmt:
		t.writeIndent()
		t.write("return")

		if n.Value != nil {
			t.write(" ")
			t.emitExpression(n.Value)
		}

		t.write(";")
		t.newline()

	case *ast.ExpressionStmt:
		t.writeIndent()
		t.emitExpression(n.Expression)
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
			t.emitExpression(n.Value)
		}

		t.write(";")
		t.newline()

	case *ast.AssignmentStmt:
		t.writeIndent()
		t.emitExpression(n.Left)
		t.write(" = ")
		t.emitExpression(n.Value)
		t.write(";")
		t.newline()

	default:
		panic(fmt.Sprintf("unsupported statement %T", stmt))
	}
}

func (t *Transpiler) emitIf(n *ast.IfStmt) {
	t.writeIndent()
	t.write("if (")

	t.emitExpression(n.Condition)

	t.write(") {\n")

	t.pushIndent()

	for _, stmt := range n.Then {
		t.emitStatement(stmt)
	}

	t.popIndent()

	t.writeIndent()
	t.write("}")

	if len(n.Else) > 0 {
		t.write(" else {\n")

		t.pushIndent()

		for _, stmt := range n.Else {
			t.emitStatement(stmt)
		}

		t.popIndent()

		t.writeIndent()
		t.write("}")
	}

	t.newline()
}

func (t *Transpiler) emitLoop(n *ast.LoopStmt) {
	t.writeIndent()
	t.write("for (;;) {\n")

	t.pushIndent()

	for _, stmt := range n.Body {
		t.emitStatement(stmt)
	}

	t.popIndent()

	t.writeIndent()
	t.write("}")
	t.newline()
}
