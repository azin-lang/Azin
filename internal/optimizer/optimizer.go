// package optimizer provides basic optimitzations before transpiling to C
package optimizer

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

// Optimize traverses the program AST and applies compile-time simplifications.
// It currently performs constant folding and conditional dead-code elimination.
func Optimize(program *ast.Program) {
	if program == nil {
		return
	}

	program.Statements = optimizeStatements(program.Statements)
}

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
		n.Body = optimizeStatements(n.Body)
		return []ast.Stmt{n}

	case *ast.FuncStmt:
		n.Body = optimizeStatements(n.Body)
		return []ast.Stmt{n}

	case *ast.ReturnStmt:
		if n.Value != nil {
			n.Value = optimizeExpr(n.Value)
		}
		return []ast.Stmt{n}

	case *ast.VarStmt:
		if n.Value != nil {
			n.Value = optimizeExpr(n.Value)
		}
		return []ast.Stmt{n}

	case *ast.AssignmentStmt:
		n.Left = optimizeExpr(n.Left)
		n.Value = optimizeExpr(n.Value)
		return []ast.Stmt{n}

	case *ast.ExpressionStmt:
		if n.Expression != nil {
			n.Expression = optimizeExpr(n.Expression)
		}
		return []ast.Stmt{n}

	case *ast.StructStmt, *ast.ImportCStmt, *ast.StopStmt, *ast.BadStmt:
		return []ast.Stmt{stmt}

	default:
		return []ast.Stmt{stmt}
	}
}

func optimizeIf(n *ast.IfStmt) []ast.Stmt {
	n.Condition = optimizeExpr(n.Condition)
	n.Then = optimizeStatements(n.Then)
	n.Else = optimizeStatements(n.Else)

	if cond, ok := n.Condition.(*ast.BooleanLiteral); ok {
		if cond.Value {
			return n.Then
		}
		return n.Else
	}

	// Keep the if-statement when the condition cannot be folded.
	return []ast.Stmt{n}
}

func optimizeExpr(expr ast.Expr) ast.Expr {
	switch n := expr.(type) {
	case *ast.BinaryExpr:
		left := optimizeExpr(n.Left)
		right := optimizeExpr(n.Right)

		n.Left = left
		n.Right = right

		if folded := foldBinaryExpr(left, n.Operator, right); folded != nil {
			return folded
		}

		return n

	case *ast.MemberExpr:
		n.Object = optimizeExpr(n.Object)
		return n

	case *ast.CallExpr:
		n.Callee = optimizeExpr(n.Callee)
		for i, arg := range n.Args {
			n.Args[i] = optimizeExpr(arg)
		}
		return n

	default:
		return expr
	}
}

func foldBinaryExpr(left ast.Expr, op token.Token, right ast.Expr) ast.Expr {
	switch l := left.(type) {
	case *ast.IntegerLiteral:
		switch r := right.(type) {
		case *ast.IntegerLiteral:
			return foldIntegerInteger(l, op, r)
		case *ast.FloatLiteral:
			return foldFloatFloat(&ast.FloatLiteral{Value: float64(l.Value)}, op, r)
		case *ast.BooleanLiteral:
			return nil
		case *ast.CharacterLiteral:
			return foldIntegerInteger(l, op, &ast.IntegerLiteral{Value: int64(r.Value)})
		}

	case *ast.FloatLiteral:
		switch r := right.(type) {
		case *ast.IntegerLiteral:
			return foldFloatFloat(l, op, &ast.FloatLiteral{Value: float64(r.Value)})
		case *ast.FloatLiteral:
			return foldFloatFloat(l, op, r)
		}

	case *ast.BooleanLiteral:
		switch r := right.(type) {
		case *ast.BooleanLiteral:
			return foldBooleanBoolean(l, op, r)
		}

	case *ast.CharacterLiteral:
		switch r := right.(type) {
		case *ast.CharacterLiteral:
			return foldIntegerInteger(&ast.IntegerLiteral{Value: int64(nValue(l))}, op, &ast.IntegerLiteral{Value: int64(r.Value)})
		case *ast.IntegerLiteral:
			return foldIntegerInteger(&ast.IntegerLiteral{Value: int64(l.Value)}, op, r)
		}
	}

	return nil
}

func foldIntegerInteger(left *ast.IntegerLiteral, op token.Token, right *ast.IntegerLiteral) ast.Expr {
	switch op.Kind {
	case token.Plus:
		return &ast.IntegerLiteral{Value: left.Value + right.Value}
	case token.Minus:
		return &ast.IntegerLiteral{Value: left.Value - right.Value}
	case token.Star:
		return &ast.IntegerLiteral{Value: left.Value * right.Value}
	case token.Slash:
		if right.Value == 0 {
			return nil
		}
		return &ast.IntegerLiteral{Value: left.Value / right.Value}
	case token.Modulo:
		if right.Value == 0 {
			return nil
		}
		return &ast.IntegerLiteral{Value: left.Value % right.Value}
	case token.EqualEqual:
		return &ast.BooleanLiteral{Value: left.Value == right.Value}
	case token.BangEqual:
		return &ast.BooleanLiteral{Value: left.Value != right.Value}
	case token.Less:
		return &ast.BooleanLiteral{Value: left.Value < right.Value}
	case token.LessEqual:
		return &ast.BooleanLiteral{Value: left.Value <= right.Value}
	case token.Greater:
		return &ast.BooleanLiteral{Value: left.Value > right.Value}
	case token.GreaterEqual:
		return &ast.BooleanLiteral{Value: left.Value >= right.Value}
	default:
		return nil
	}
}

func foldFloatFloat(left *ast.FloatLiteral, op token.Token, right *ast.FloatLiteral) ast.Expr {
	switch op.Kind {
	case token.Plus:
		return &ast.FloatLiteral{Value: left.Value + right.Value}
	case token.Minus:
		return &ast.FloatLiteral{Value: left.Value - right.Value}
	case token.Star:
		return &ast.FloatLiteral{Value: left.Value * right.Value}
	case token.Slash:
		if right.Value == 0 {
			return nil
		}
		return &ast.FloatLiteral{Value: left.Value / right.Value}
	case token.EqualEqual:
		return &ast.BooleanLiteral{Value: left.Value == right.Value}
	case token.BangEqual:
		return &ast.BooleanLiteral{Value: left.Value != right.Value}
	case token.Less:
		return &ast.BooleanLiteral{Value: left.Value < right.Value}
	case token.LessEqual:
		return &ast.BooleanLiteral{Value: left.Value <= right.Value}
	case token.Greater:
		return &ast.BooleanLiteral{Value: left.Value > right.Value}
	case token.GreaterEqual:
		return &ast.BooleanLiteral{Value: left.Value >= right.Value}
	default:
		return nil
	}
}

func foldBooleanBoolean(left *ast.BooleanLiteral, op token.Token, right *ast.BooleanLiteral) ast.Expr {
	switch op.Kind {
	case token.LogicalAnd:
		return &ast.BooleanLiteral{Value: left.Value && right.Value}
	case token.LogicalOr:
		return &ast.BooleanLiteral{Value: left.Value || right.Value}
	case token.EqualEqual:
		return &ast.BooleanLiteral{Value: left.Value == right.Value}
	case token.BangEqual:
		return &ast.BooleanLiteral{Value: left.Value != right.Value}
	default:
		return nil
	}
}

func nValue(c *ast.CharacterLiteral) rune {
	return c.Value
}
