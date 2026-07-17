package optimizer

import "github.com/azin-lang/Azin/internal/ast"

func optimizeStatements(stmts []ast.Stmt) []ast.Stmt {
	out := make([]ast.Stmt, 0, len(stmts))

	for _, stmt := range stmts {
		if stmt == nil {
			continue
		}

		out = append(out, optimizeStatement(stmt)...)
	}

	return out
}

func optimizeStatement(stmt ast.Stmt) []ast.Stmt {
	switch n := stmt.(type) {
	case *ast.IfStmt:
		return optimizeIf(n)

	case *ast.LoopStmt:
		optimizeLoop(n)

	case *ast.FuncStmt:
		optimizeFunction(n)

	case *ast.ReturnStmt:
		optimizeReturn(n)

	case *ast.VarStmt:
		optimizeVariable(n)

	case *ast.AssignmentStmt:
		optimizeAssignment(n)

	case *ast.ExpressionStmt:
		optimizeExpressionStmt(n)
	}

	return []ast.Stmt{stmt}
}

func optimizeLoop(n *ast.LoopStmt) {
	n.Body = optimizeStatements(n.Body)
}

func optimizeFunction(n *ast.FuncStmt) {
	n.Body = optimizeStatements(n.Body)
}

func optimizeReturn(n *ast.ReturnStmt) {
	if n.Value != nil {
		n.Value = optimizeExpr(n.Value)
	}
}

func optimizeVariable(n *ast.VarStmt) {
	if n.Value != nil {
		n.Value = optimizeExpr(n.Value)
	}
}

func optimizeAssignment(n *ast.AssignmentStmt) {
	n.Left = optimizeExpr(n.Left)
	n.Value = optimizeExpr(n.Value)
}

func optimizeExpressionStmt(n *ast.ExpressionStmt) {
	if n.Expression != nil {
		n.Expression = optimizeExpr(n.Expression)
	}
}
