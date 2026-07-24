package optimizer

import (
	"sort"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func reassociateBinary(n *ast.BinaryExpr) ast.Expr {
	switch n.Operator.Kind {
	case token.Plus, token.Star:
		return reassociate(n)

	default:
		return nil
	}
}

func reassociate(root *ast.BinaryExpr) ast.Expr {
	var (
		terms    []ast.Expr
		constant ast.Expr
		changed  bool
	)

	var collect func(ast.Expr)

	collect = func(expr ast.Expr) {
		if bin, ok := expr.(*ast.BinaryExpr); ok &&
			bin.Operator.Kind == root.Operator.Kind {

			changed = true
			collect(bin.Left)
			collect(bin.Right)
			return
		}

		if isConstant(expr) {
			changed = true

			if constant == nil {
				constant = expr
			} else {
				constant = foldConstant(root.Operator, constant, expr)
			}

			return
		}

		terms = append(terms, expr)
	}

	collect(root)

	if !changed {
		return nil
	}

	sort.SliceStable(terms, func(i, j int) bool {
		return exprLess(terms[i], terms[j])
	})

	switch root.Operator.Kind {
	case token.Plus:
		if constant != nil && !isZero(constant) {
			terms = append(terms, constant)
		}

	case token.Star:
		if constant != nil {
			if isZero(constant) {
				return constant
			}

			if !isOne(constant) || len(terms) == 0 {
				terms = append(terms, constant)
			}
		}
	}

	if len(terms) == 0 {
		return constant
	}

	return buildAssociativeTree(root.Operator, terms)
}

func foldConstant(op token.Token, left, right ast.Expr) ast.Expr {
	if folded := foldBinaryExpr(left, op, right); folded != nil {
		return folded
	}

	return &ast.BinaryExpr{
		Left:     left,
		Operator: op,
		Right:    right,
	}
}

func buildAssociativeTree(op token.Token, terms []ast.Expr) ast.Expr {
	switch len(terms) {
	case 0:
		return nil

	case 1:
		return terms[0]

	case 2:
		return &ast.BinaryExpr{
			Left:     terms[0],
			Operator: op,
			Right:    terms[1],
		}
	}

	mid := len(terms) / 2

	return &ast.BinaryExpr{
		Left:     buildAssociativeTree(op, terms[:mid]),
		Operator: op,
		Right:    buildAssociativeTree(op, terms[mid:]),
	}
}

func exprLess(a, b ast.Expr) bool {
	ra := exprRank(a)
	rb := exprRank(b)

	if ra != rb {
		return ra < rb
	}

	switch a := a.(type) {
	case *ast.Identifier:
		return a.Value < b.(*ast.Identifier).Value

	case *ast.IntegerLiteral:
		return a.Value < b.(*ast.IntegerLiteral).Value

	case *ast.FloatLiteral:
		return a.Value < b.(*ast.FloatLiteral).Value

	case *ast.CharacterLiteral:
		return a.Value < b.(*ast.CharacterLiteral).Value

	case *ast.BooleanLiteral:
		return !a.Value && b.(*ast.BooleanLiteral).Value

	case *ast.StringLiteral:
		return a.Value < b.(*ast.StringLiteral).Value
	}

	return false
}

func exprRank(expr ast.Expr) int {
	switch expr.(type) {
	case *ast.Identifier:
		return 0
	case *ast.MemberExpr:
		return 1
	case *ast.BinaryExpr:
		return 2
	default:
		return 3
	}
}

func canonicalizeBinary(n *ast.BinaryExpr) ast.Expr {
	switch n.Operator.Kind {

	case token.Plus,
		token.Star,
		token.EqualEqual,
		token.BangEqual:

		if isConstant(n.Left) && !isConstant(n.Right) {
			n.Left, n.Right = n.Right, n.Left
			return n
		}
	}

	return nil
}
