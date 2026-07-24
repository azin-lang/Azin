package optimizer

import (
	"reflect"
	"sort"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

func reassociateBinary(n *ast.BinaryExpr) ast.Expr {
	switch n.Operator.Kind {
	case token.Plus, token.Star:
	default:
		return nil
	}

	if _, ok := n.Left.(*ast.BinaryExpr); ok {
		return reassociate(n)
	}
	if _, ok := n.Right.(*ast.BinaryExpr); ok {
		return reassociate(n)
	}

	if isConstant(n.Left) && isConstant(n.Right) {
		return reassociate(n)
	}

	return nil
}

func reassociate(root *ast.BinaryExpr) ast.Expr {
	var (
		terms    []ast.Expr
		constant ast.Expr
	)

	var collect func(ast.Expr)
	collect = func(expr ast.Expr) {
		if bin, ok := expr.(*ast.BinaryExpr); ok &&
			bin.Operator.Kind == root.Operator.Kind {
			collect(bin.Left)
			collect(bin.Right)
			return
		}

		if isConstant(expr) {
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

	//nolint:exhaustive
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

	sort.SliceStable(terms, func(i, j int) bool {
		return exprLess(terms[i], terms[j])
	})

	var rebuilt ast.Expr
	switch len(terms) {
	case 0:
		rebuilt = constant
	case 1:
		rebuilt = terms[0]
	default:
		rebuilt = buildAssociativeTree(root.Operator, terms)
	}

	if reflect.DeepEqual(root, rebuilt) {
		return nil
	}

	return rebuilt
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
	//nolint:exhaustive
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
