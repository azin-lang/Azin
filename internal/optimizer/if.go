package optimizer

import "github.com/azin-lang/Azin/internal/ast"

func optimizeIf(n *ast.IfStmt) []ast.Stmt {
	n.Condition = optimizeExpr(n.Condition)
	n.Then = optimizeStatements(n.Then)
	n.Else = optimizeStatements(n.Else)

	if b, ok := n.Condition.(*ast.BooleanLiteral); ok {
		if b.Value {
			return n.Then
		}
		return n.Else
	}

	return []ast.Stmt{n}
}
