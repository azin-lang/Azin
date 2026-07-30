package analysis

import (
	"github.com/azin-lang/Azin/pkg/ast"
)

func (a *Analyzer) visitMember(fn string, expr *ast.MemberExpr) {
	if id, ok := expr.Object.(*ast.Identifier); ok {
		name := id.Value
		if e, exists := a.Enums[name]; exists {
			a.MarkTypeUsed(e.SemaType)
			return
		}
		if s, exists := a.Structs[name]; exists {
			a.MarkTypeUsed(s.SemaType)
		}
	}
	a.visitExpr(fn, expr.Object)
}
