package c

import (
	"github.com/azin-lang/Azin/internal/ast"
)

func (t *Transpiler) emitFunction(
	fn *ast.FuncStmt,
) {
	t.emitFunctionSignature(fn)

	t.write(" {\n")

	t.pushIndent()

	for _, stmt := range fn.Body {
		t.emitStatement(stmt)
	}

	if t.functionName(fn) == "main" &&
		fn.ReturnType != nil &&
		emitType(fn.ReturnType.Value) == "void" {

		t.indentLine()
		t.write("return 0;\n")
	}

	t.popIndent()

	t.write("}\n")
}

func (t *Transpiler) emitFunctionSignature(
	fn *ast.FuncStmt,
) {
	name := t.functionName(fn)

	if fn.ReturnType == nil {
		t.printf("void %s(", name)
		t.write("void")
		t.write(") /* missing return type */")
		return
	}

	ret := emitType(
		fn.ReturnType.Value,
	)

	if name == "main" &&
		ret == "void" {

		ret = "int"
	}

	t.printf(
		"%s %s(",
		ret,
		name,
	)

	if len(fn.Params) == 0 {
		t.write("void")
	} else {
		for i, param := range fn.Params {

			if i > 0 {
				t.write(", ")
			}

			if param.Type == nil {
				t.printf("void %s /* missing type */", param.Name.Value)
				continue
			}

			t.printf(
				"%s %s",
				emitType(param.Type.Value),
				param.Name.Value,
			)
		}
	}

	t.write(")")
}

func (t *Transpiler) functionName(
	fn *ast.FuncStmt,
) string {
	if fn.CName != "" {
		return fn.CName
	}

	return fn.Name.Value
}
