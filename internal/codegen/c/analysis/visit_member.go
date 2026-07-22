package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitMember(
	fn string,
	expr *ast.MemberExpr,
) {
	if id, ok := expr.Object.(*ast.Identifier); ok {
		name := id.Value

		// Enum access: Color.Blue
		if _, ok := a.Types[name]; ok {
			return
		}

		// Struct access: point.x
		a.markType(name)
	}

	a.visitExpr(fn, expr.Object)
}
