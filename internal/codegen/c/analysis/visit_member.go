package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitMember(fn string, expr *ast.MemberExpr) {
	if id, ok := expr.Object.(*ast.Identifier); ok {
		name := id.Value
		if _, exists := a.Enums[name]; exists {
			a.MarkTypeUsed(name)
			return
		}
		if _, exists := a.Structs[name]; exists {
			a.MarkTypeUsed(name)
		}
	}
	a.visitExpr(fn, expr.Object)
}
