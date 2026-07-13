package codegen

import (
	"bytes"
	"path/filepath"

	"github.com/azin-lang/Azin/internal/ast"
)

// the Transpiler struct is responsible for transpiling the AST to C code.
type Transpiler struct {
	buf    bytes.Buffer
	indent int
	enums  map[string]bool
}

// create a new Transpiler.
func New() *Transpiler {
	return &Transpiler{
		enums: map[string]bool{},
	}
}

// Transpile transpiles the AST to C code.
func (t *Transpiler) Transpile(program *ast.Program) string {
	hasImports := false

	for _, stmt := range program.Statements {
		if imp, ok := stmt.(*ast.ImportCStmt); ok {
			header := imp.Path.Value

			if filepath.Ext(header) == "" {
				header += ".h"
			}

			t.printf("#include <%s>\n", header)
			hasImports = true
		}
	}

	if hasImports {
		t.newline()
	}

	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.EnumStmt:
			t.compileEnum(s)
		case *ast.StructStmt:
			t.compileStruct(s)
		}
	}

	for _, stmt := range program.Statements {
		switch stmt.(type) {
		case *ast.StructStmt, *ast.EnumStmt:
			continue
		}

		t.compileStatement(stmt)
		t.newline()
	}

	return t.buf.String()
}
