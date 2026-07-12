package codegen

import (
	"bytes"
	"strings"

	"github.com/azin-lang/Azin/internal/ast"
)

// the Transpiler struct is responsible for transpiling the AST to C code.
type Transpiler struct {
	buf    bytes.Buffer
	indent int
}

// create a new Transpiler.
func New() *Transpiler {
	return &Transpiler{}
}

// Transpile transpiles the AST to C code.
func (t *Transpiler) Transpile(program *ast.Program) string {
	for _, stmt := range program.Statements {
		if imp, ok := stmt.(*ast.ImportCStmt); ok {
			// Strip out any surrounding quotes captured by the parser
			cleanPath := strings.Trim(imp.Path.Value, "\"")

			// Emit as a standard C header include string
			t.printf("#include \"%s\"\n", cleanPath)
		}
	}
	t.newline()

	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.StructStmt); ok {
			t.compileStruct(s)
		}
	}

	for _, stmt := range program.Statements {
		if _, ok := stmt.(*ast.StructStmt); ok {
			continue
		}

		t.compileStatement(stmt)
		t.newline()
	}

	return t.buf.String()
}
