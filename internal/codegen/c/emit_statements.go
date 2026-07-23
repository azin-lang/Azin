package c

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/ast"
)

func (t *Transpiler) emitStatement(
	stmt ast.Stmt,
) {
	switch n := stmt.(type) {

	case *ast.StructStmt,
		*ast.EnumStmt,
		*ast.ImportCStmt:
		// emitted elsewhere

	case *ast.FuncStmt:
		t.emitFunction(n)

	case *ast.ReturnStmt:
		t.indentLine()

		t.write("return")

		if n.Value != nil {
			t.write(" ")
			t.emitExpression(n.Value)
		}

		t.write(";\n")

	case *ast.ExpressionStmt:
		t.indentLine()

		t.emitExpression(
			n.Expression,
		)

		t.write(";\n")

	case *ast.VarStmt:
		t.emitVariable(n)

	case *ast.AssignmentStmt:
		t.indentLine()

		t.emitExpression(n.Left)

		t.write(" = ")

		t.emitExpression(n.Value)

		t.write(";\n")

	case *ast.IfStmt:
		t.emitIf(n)

	case *ast.LoopStmt:
		t.emitLoop(n)

	case *ast.StopStmt:
		t.indentLine()
		t.write("break;\n")

	default:
		t.write(fmt.Sprintf(
			"/* unsupported statement %T */",
			stmt,
		))
	}
}

func (t *Transpiler) emitVariable(
	stmt *ast.VarStmt,
) {
	t.indentLine()

	if stmt.Type == nil {
		t.write("/* variable missing type */")
		t.write(";\n")
		return
	}

	if stmt.Type.Value == "string" {
		if !stmt.Mutable {
			t.write("const char* const ")
		} else {
			t.write("char* ")
		}
		t.write(stmt.Name.Value)
	} else {
		if !stmt.Mutable {
			t.write("const ")
		}

		t.printf(
			"%s %s",
			emitType(stmt.Type.Value),
			stmt.Name.Value,
		)
	}

	if stmt.Value != nil {
		t.write(" = ")
		t.emitExpression(stmt.Value)
	}

	t.write(";\n")
}

func (t *Transpiler) emitBlock(
	body []ast.Stmt,
) {
	t.pushIndent()

	for _, stmt := range body {
		t.emitStatement(stmt)
	}

	t.popIndent()
}

func (t *Transpiler) emitIf(
	stmt *ast.IfStmt,
) {
	t.indentLine()

	t.write("if (")

	t.emitExpression(
		stmt.Condition,
	)

	t.write(") {\n")

	t.emitBlock(
		stmt.Then,
	)

	t.indentLine()
	t.write("}")

	if len(stmt.Else) != 0 {

		t.write(" else {\n")

		t.emitBlock(
			stmt.Else,
		)

		t.indentLine()
		t.write("}")
	}

	t.newline()
}

func (t *Transpiler) emitLoop(
	stmt *ast.LoopStmt,
) {
	t.indentLine()

	t.write(
		"for (;;) {\n",
	)

	t.emitBlock(
		stmt.Body,
	)

	t.indentLine()

	t.write("}\n")
}
