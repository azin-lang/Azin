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
			t.emitStruct(n)
			t.newline()

		case *ast.EnumStmt:
			t.emitEnum(n)
			t.newline()
		}
	}
}

func (t *Transpiler) emitStruct(
	s *ast.StructStmt,
) {
	t.printf(
		"typedef struct %s {\n",
		s.Name.Value,
	)

	t.pushIndent()

	for _, field := range s.Fields {
		t.indentLine()

		t.printf(
			"%s %s;\n",
			emitType(field.Type.Value),
			field.Name.Value,
		)
	}

	t.popIndent()

	t.printf(
		"} %s;\n",
		s.Name.Value,
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
