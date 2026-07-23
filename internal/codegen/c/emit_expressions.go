package c

import (
	"fmt"
	"strconv"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func (t *Transpiler) emitExpression(
	expr ast.Expr,
) {
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
		t.write(
			strconv.FormatFloat(
				n.Value,
				'g',
				-1,
				64,
			),
		)

	case *ast.StringLiteral:
		t.printf("%q", n.Value)

	case *ast.CharacterLiteral:
		t.printf("%q", n.Value)

	case *ast.MemberExpr:
		if id, ok := n.Object.(*ast.Identifier); ok {
			if _, isEnum := t.enums[id.Value]; isEnum {
				t.printf("%s_%s", id.Value, n.Property.Value)
				return
			}
		}

		t.emitExpression(n.Object)

		t.write(".")

		t.write(
			n.Property.Value,
		)

	case *ast.BinaryExpr:
		t.emitExpression(n.Left)

		t.write(" ")
		t.write(
			emitOperator(
				n.Operator.Kind,
			),
		)
		t.write(" ")

		t.emitExpression(n.Right)

	case *ast.CallExpr:
		t.write(
			n.ResolvedName,
		)

		t.write("(")

		for i, arg := range n.Args {

			if i != 0 {
				t.write(", ")
			}

			t.emitExpression(arg)
		}

		t.write(")")

	default:
		t.write(fmt.Sprintf(
			"/* unsupported expression %T */",
			expr,
		))
	}
}

func emitType(
	name string,
) string {
	switch name {

	case "unit":
		return "void"

	case "int":
		return "int"

	case "float":
		return "float"

	case "char":
		return "char"

	case "string":
		return "char*"

	case "bool":
		return "bool"

	default:
		return name
	}
}

func emitOperator(
	kind token.Kind,
) string {
	switch kind {

	case token.Plus:
		return "+"

	case token.Minus:
		return "-"

	case token.Star:
		return "*"

	case token.Slash:
		return "/"

	case token.EqualEqual:
		return "=="

	case token.BangEqual:
		return "!="

	case token.Less:
		return "<"

	case token.LessEqual:
		return "<="

	case token.Greater:
		return ">"

	case token.GreaterEqual:
		return ">="

	default:
		return "/* unsupported operator */"
	}
}
