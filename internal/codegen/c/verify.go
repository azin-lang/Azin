package c

import "github.com/azin-lang/Azin/internal/ast"

func (t *Transpiler) verifyResolvedCalls(program *ast.Program) {
	var visitExpr func(ast.Expr)
	var visitStmt func(ast.Stmt)

	visitExpr = func(expr ast.Expr) {
		switch n := expr.(type) {
		case nil, *ast.BadExpr:
			return

		case *ast.CallExpr:
			if n.ResolvedName == "" {
				panic("internal compiler error: unresolved function call reached code generation")
			}

			visitExpr(n.Callee)
			for _, arg := range n.Args {
				visitExpr(arg)
			}

		case *ast.BinaryExpr:
			visitExpr(n.Left)
			visitExpr(n.Right)

		case *ast.MemberExpr:
			visitExpr(n.Object)
		}
	}

	visitStmt = func(stmt ast.Stmt) {
		switch n := stmt.(type) {
		case *ast.BadStmt, *ast.ImportCStmt, *ast.StructStmt, *ast.EnumStmt, *ast.StopStmt:
			return

		case *ast.FuncStmt:
			for _, stmt := range n.Body {
				visitStmt(stmt)
			}

		case *ast.ReturnStmt:
			visitExpr(n.Value)

		case *ast.VarStmt:
			visitExpr(n.Value)

		case *ast.IfStmt:
			visitExpr(n.Condition)
			for _, stmt := range n.Then {
				visitStmt(stmt)
			}
			for _, stmt := range n.Else {
				visitStmt(stmt)
			}

		case *ast.LoopStmt:
			for _, stmt := range n.Body {
				visitStmt(stmt)
			}

		case *ast.AssignmentStmt:
			visitExpr(n.Left)
			visitExpr(n.Value)

		case *ast.ExpressionStmt:
			visitExpr(n.Expression)
		}
	}

	for _, stmt := range program.Statements {
		visitStmt(stmt)
	}
}
