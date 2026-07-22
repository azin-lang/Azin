package c

import "github.com/azin-lang/Azin/internal/ast"

func (t *Transpiler) emitStruct(s *ast.StructStmt) {
	t.write("typedef struct {\n")
	t.pushIndent()

	for _, field := range s.Fields {
		t.writeIndent()
		t.printf("%s %s;\n", emitType(field.Type.Value), field.Name.Value)
	}

	t.popIndent()
	t.printf("} %s;\n", s.Name.Value)
	t.newline()
}

func (t *Transpiler) emitEnum(e *ast.EnumStmt) {
	t.enums[e.Name.Value] = true

	t.write("typedef enum {\n")
	t.pushIndent()

	for i, variant := range e.Variants {
		t.writeIndent()
		t.printf("%s_%s", e.Name.Value, variant.Value)
		if i < len(e.Variants)-1 {
			t.write(",")
		}
		t.newline()
	}

	t.popIndent()
	t.printf("} %s;\n", e.Name.Value)
}

func (t *Transpiler) emitFunction(fn *ast.FuncStmt) {
	if fn.ReturnType == nil {
		panic("internal compiler error: function has no resolved return type")
	}

	name := fn.Name.Value
	if fn.CName != "" {
		name = fn.CName
	}

	t.printf("%s %s(", emitType(fn.ReturnType.Value), name)

	for i, p := range fn.Params {
		if p.Type == nil {
			panic("internal compiler error: parameter has no resolved type")
		}

		if i > 0 {
			t.write(", ")
		}

		t.printf("%s %s", emitType(p.Type.Value), p.Name.Value)
	}

	t.write(")")
	t.write(" {\n")

	t.pushIndent()

	for _, stmt := range fn.Body {
		t.emitStatement(stmt)
	}

	t.popIndent()

	t.write("}\n")
}
