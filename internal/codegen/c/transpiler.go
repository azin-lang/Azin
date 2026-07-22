package c

import (
	"bytes"
	"path/filepath"

	"github.com/azin-lang/Azin/internal/ast"
)

// Transpiler the Transpiler struct is responsible for transpiling the AST to C code.
type Transpiler struct {
	buf    bytes.Buffer
	indent int
	enums  map[string]bool
}

// New create a new Transpiler.
func New() *Transpiler {
	return &Transpiler{
		enums: map[string]bool{},
	}
}

// handleImports handles the import statements in the AST and writes the corresponding C include statements to the buffer.
func (t *Transpiler) handleImports(program *ast.Program) string {
	for _, stmt := range program.Statements {
		if imp, ok := stmt.(*ast.ImportCStmt); ok {
			header := imp.Path.Value
			if filepath.Ext(header) == "" {
				header += ".h"
			}
			t.write("#include <" + header + ">\n")
		}
	}

	t.write("#include <stdbool.h>") // for bool, true, false because *great* C doesn't have those built-in

	t.newline()
	t.newline()

	return t.buf.String()
}

func (t *Transpiler) verifyResolvedCalls(program *ast.Program) {
	var visitExpr func(ast.Expr)
	var visitStmt func(ast.Stmt)

	visitExpr = func(expr ast.Expr) {
		switch n := expr.(type) {
		case nil, *ast.BadExpr:
			return

		case *ast.CallExpr:
			if n.ResolvedName == "" {
				panic("internal compiler error: unresolved function call reached code generation")
			}

			visitExpr(n.Callee)
			for _, arg := range n.Args {
				visitExpr(arg)
			}

		case *ast.BinaryExpr:
			visitExpr(n.Left)
			visitExpr(n.Right)

		case *ast.MemberExpr:
			visitExpr(n.Object)
		}
	}

	visitStmt = func(stmt ast.Stmt) {
		switch n := stmt.(type) {
		case *ast.BadStmt, *ast.ImportCStmt, *ast.StructStmt, *ast.EnumStmt, *ast.StopStmt:
			return

		case *ast.FuncStmt:
			for _, stmt := range n.Body {
				visitStmt(stmt)
			}

		case *ast.ReturnStmt:
			visitExpr(n.Value)

		case *ast.VarStmt:
			visitExpr(n.Value)

		case *ast.IfStmt:
			visitExpr(n.Condition)
			for _, stmt := range n.Then {
				visitStmt(stmt)
			}
			for _, stmt := range n.Else {
				visitStmt(stmt)
			}

		case *ast.LoopStmt:
			for _, stmt := range n.Body {
				visitStmt(stmt)
			}

		case *ast.AssignmentStmt:
			visitExpr(n.Left)
			visitExpr(n.Value)

		case *ast.ExpressionStmt:
			visitExpr(n.Expression)
		}
	}

	for _, stmt := range program.Statements {
		visitStmt(stmt)
	}
}

// Transpile takes an AST and transpiles it to C code.
func (t *Transpiler) Transpile(program *ast.Program) string {
	t.buf.Reset()
	t.verifyResolvedCalls(program)

	t.handleImports(program)

	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.EnumStmt:
			t.compileEnum(s)
			t.newline()
		case *ast.StructStmt:
			t.compileStruct(s)
			t.newline()
		}
	}

	for _, stmt := range program.Statements {
		switch stmt.(type) {
		case *ast.ImportCStmt, *ast.StructStmt, *ast.EnumStmt:
			continue
		default:
			t.compileStatement(stmt)
			t.newline()
		}
	}

	return t.buf.String()
}
