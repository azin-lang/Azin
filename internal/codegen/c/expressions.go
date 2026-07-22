package c

import (
	"fmt"
	"strconv"

	"github.com/azin-lang/Azin/internal/ast"
)

func (t *Transpiler) emitExpression(expr ast.Expr) {
	switch n := expr.(type) {

	case *ast.Identifier:
		t.write(n.Value)

	case *ast.BooleanLiteral:
		if n.Value {
			t.write("true")
		} else {
			t.write("false")
		}

	case *ast.IntegerLiteral:
		t.printf("%d", n.Value)

	case *ast.FloatLiteral:
		t.printf("%s", strconv.FormatFloat(n.Value, 'g', -1, 64))

	case *ast.StringLiteral:
		t.printf("%q", n.Value)

	case *ast.CharacterLiteral:
		t.printf("%q", n.Value)

	case *ast.MemberExpr:
		if id, ok := n.Object.(*ast.Identifier); ok && t.enums[id.Value] {
			t.printf("%s_%s", id.Value, n.Property.Value)
			return
		}

		t.emitExpression(n.Object)
		t.write(".")
		t.write(n.Property.Value)

	case *ast.BinaryExpr:
		t.emitExpression(n.Left)
		t.write(" ")
		t.write(emitOperator(n.Operator.Kind))
		t.write(" ")
		t.emitExpression(n.Right)

	case *ast.CallExpr:
		if n.ResolvedName == "" {
			panic("internal compiler error: unresolved function call")
		}
		t.write(n.ResolvedName)
		t.write("(")

		for i, arg := range n.Args {
			if i > 0 {
				t.write(", ")
			}
			t.emitExpression(arg)
		}

		t.write(")")

	default:
		panic(fmt.Sprintf("unsupported expression %T", expr))
	}
}
