package c

import (
	"sort"

	"github.com/azin-lang/Azin/internal/ast"
)

func (t *Transpiler) emit(
	program *ast.Program,
) {
	t.emitImports()
	t.emitStructDeclarations(program)
	t.emitTypes(program)
	t.emitForwardDeclarations()
	t.emitFunctions(program)
}

func (t *Transpiler) emitImports() {
	if len(t.includes) == 0 {
		return
	}

	headers := make([]string, 0, len(t.includes))

	for header := range t.includes {
		headers = append(headers, header)
	}

	sort.Strings(headers)

	for _, header := range headers {
		t.printf(
			"#include <%s>\n",
			header,
		)
	}

	t.newline()
}

func (t *Transpiler) emitStructDeclarations(program *ast.Program) {
	hasStructs := false
	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.StructStmt); ok {
			name := s.Name.Value
			if _, ok := t.reachableTypes[name]; ok {
				t.printf("typedef struct %s %s;\n", name, name)
				hasStructs = true
			}
		}
	}
	if hasStructs {
		t.newline()
	}
}

func (t *Transpiler) emitForwardDeclarations() {
	if len(t.lateFuncs) == 0 {
		return
	}

	for _, fn := range t.lateFuncs {
		t.emitFunctionSignature(fn)
		t.write(";\n")
	}

	t.newline()
}

func (t *Transpiler) emitFunctions(
	program *ast.Program,
) {
	for _, stmt := range program.Statements {
		fn, ok := stmt.(*ast.FuncStmt)

		if !ok {
			continue
		}

		t.emitFunction(fn)
		t.newline()
	}
}
