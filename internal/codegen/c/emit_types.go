package c

import (
	"github.com/azin-lang/Azin/internal/ast"
)

func (t *Transpiler) emitTypes(
	program *ast.Program,
) {
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {

		case *ast.StructStmt:
			if _, ok := t.reachableTypes[n.Name.Value]; ok {
				t.emitStruct(n)
				t.newline()
			}

		case *ast.EnumStmt:
			if _, ok := t.reachableTypes[n.Name.Value]; ok {
				t.emitEnum(n)
				t.newline()
			}
		}
	}
}

func (t *Transpiler) emitStruct(
	s *ast.StructStmt,
) {
	// We already typedef'd this in emitStructDeclarations, so just define it
	t.printf(
		"struct %s {\n",
		s.Name.Value,
	)

	t.pushIndent()

	for _, field := range s.Fields {
		t.indentLine()

		typ := emitType(field.Type.Value)

		// Convert to a pointer if it creates a cycle
		if t.isCyclicField(s.Name.Value, field.Type.Value) {
			typ += "*"
		}

		t.printf(
			"%s %s;\n",
			typ,
			field.Name.Value,
		)
	}

	t.popIndent()

	t.printf(
		"};\n",
	)
}

func (t *Transpiler) emitEnum(
	e *ast.EnumStmt,
) {
	t.printf(
		"typedef enum %s {\n",
		e.Name.Value,
	)

	t.pushIndent()

	for i, variant := range e.Variants {
		t.indentLine()

		t.printf(
			"%s_%s",
			e.Name.Value,
			variant.Value,
		)

		if i != len(e.Variants)-1 {
			t.write(",")
		}

		t.newline()
	}

	t.popIndent()

	t.printf(
		"} %s;\n",
		e.Name.Value,
	)
}
