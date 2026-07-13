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
}

// create a new Transpiler.
func New() *Transpiler {
	return &Transpiler{}
}

func handleImports(program *ast.Program) string {
	hasImports := false
	for _, stmt := range program.Statements {
		if imp, ok := stmt.(*ast.ImportCStmt); ok {
			header := imp.Path.Value
			if filepath.Ext(header) == "" {
				header += ".h"
			}
			t.write("#include <" + header + ">\n")
			hasImports = true
		}
	}

	if hasImports {
		t.newline()
	}
}

// Transpile takes an AST and transpiles it to C code.
func (t *Transpiler) Transpile(program *ast.Program) string {
	t.buf.Reset()

	hasImports := false
	for _, stmt := range program.Statements {
		if imp, ok := stmt.(*ast.ImportCStmt); ok {
			header := imp.Path.Value
			if filepath.Ext(header) == "" {
				header += ".h"
			}
			t.write("#include <" + header + ">\n")
			hasImports = true
		}
	}

	if hasImports {
		t.newline()
	}

	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.StructStmt); ok {
			t.compileStruct(s)
			t.newline()
		}
	}

	for _, stmt := range program.Statements {
		switch stmt.(type) {
		case *ast.ImportCStmt, *ast.StructStmt:
			continue
		default:
			t.compileStatement(stmt)
			t.newline()
		}
	}

	return t.buf.String()
}
