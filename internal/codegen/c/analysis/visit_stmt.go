package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitStmt(
	current string,
	stmt ast.Stmt,
) {
	switch n := stmt.(type) {

	case nil,
		*ast.BadStmt,
		*ast.StopStmt:
		return

	case *ast.ImportCStmt:
		a.requireImport(n.Path.Value)

	case *ast.StructStmt:
		a.visitStruct(n)

	case *ast.EnumStmt:
		a.visitEnum(n)

	case *ast.FuncStmt:
		a.visitFunction(n)

	case *ast.ReturnStmt:
		a.visitExpr(current, n.Value)

	case *ast.VarStmt:
		a.visitVar(current, n)

	case *ast.IfStmt:
		a.visitIf(current, n)

	case *ast.LoopStmt:
		a.visitLoop(current, n)

	case *ast.AssignmentStmt:
		a.visitAssignment(current, n)

	case *ast.ExpressionStmt:
		a.visitExpr(current, n.Expression)
	}
}
