package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitBinary(fn string, expr *ast.BinaryExpr) {
	a.visitExpr(fn, expr.Left)
	a.visitExpr(fn, expr.Right)
}
