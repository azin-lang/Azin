package c

import (
	"path/filepath"

	"github.com/azin-lang/Azin/internal/ast"
)

// emitImports handles the import statements in the AST and writes the corresponding C include statements to the buffer.
func (t *Transpiler) emitImports(program *ast.Program) string {
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
