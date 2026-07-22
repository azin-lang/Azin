package c

import "github.com/azin-lang/Azin/internal/ast"

// New creates a new C transpiler.
func New() *Transpiler {
	return &Transpiler{
		enums: make(map[string]bool),
	}
}

// Transpile converts an Azin AST into C source code.
func (t *Transpiler) Transpile(program *ast.Program) string {
	t.buf.Reset()

	t.verifyResolvedCalls(program)
	t.emitImports(program)

	for _, stmt := range program.Statements {
		switch n := stmt.(type) {
		case *ast.EnumStmt:
			t.emitEnum(n)
			t.newline()

		case *ast.StructStmt:
			t.emitStruct(n)
			t.newline()
		}
	}

	for _, stmt := range program.Statements {
		switch stmt.(type) {
		case *ast.ImportCStmt,
			*ast.StructStmt,
			*ast.EnumStmt:
			continue

		default:
			t.emitStatement(stmt)
			t.newline()
		}
	}

	return t.buf.String()
}
