package sema

import (
	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/token"
	types2 "github.com/azin-lang/Azin/pkg/types"
)

func (a *Analyzer) analyzeExpr(expr ast.Expr) {
	switch n := expr.(type) {
	case nil, *ast.BadExpr:
		return

	case *ast.Identifier:
		a.lookup(n.Value)

	case *ast.CallExpr:
		a.analyzeExpr(n.Callee)
		for _, arg := range n.Args {
			a.analyzeExpr(arg)
		}
		a.resolveCallExpr(n)

	case *ast.BinaryExpr:
		a.analyzeExpr(n.Left)
		a.analyzeExpr(n.Right)

	case *ast.MemberExpr:
		a.analyzeExpr(n.Object)
	}
}

func (a *Analyzer) inferExprType(expr ast.Expr) *types2.TypeInfo {
	if expr != nil && expr.Type().IsKnown() {
		// Type has already been inferred, so we can return it directly.
		return expr.Type()
	}

	switch n := expr.(type) {
	case *ast.Identifier:
		resultType, err := a.lookupType(n.Value)
		if err != nil {
			a.errorf(n, "unknown identifier: %s", n.Value)
		}

		n.SemaType = resultType
		return resultType

	case *ast.CallExpr:
		sym := a.resolveCallExpr(n)

		n.SemaReturnType = types2.UnknownType()
		if sym == nil {
			// FIXME: any function returning a C call's result will fail. When are we creating a signature table for the headers?
			// Or better yet, add actual header parsing...
			n.SemaReturnType = types2.ErrorType()
			return types2.ErrorType()
		}

		if sym.Inferring {
			return types2.UnknownType()
		}

		if sym.Type.IsUnknown() {
			// If the function's return type is unknown, we need to infer it now.
			a.inferFunctionReturnType(sym.Function)
		}

		n.SemaReturnType = sym.Type
		return sym.Type

	case *ast.BinaryExpr:
		left := a.inferExprType(n.Left)
		right := a.inferExprType(n.Right)

		if left.IsError() || right.IsError() {
			n.SemaResultType = types2.ErrorType()
			return types2.ErrorType()
		}

		//nolint:exhaustive
		switch n.Operator.Kind {
		case token.Plus, token.Minus, token.Star, token.Slash:
			if !left.IsNumeric() || !right.IsNumeric() {
				a.errorf(
					n,
					"operator '%s' requires numeric operands",
					n.TokenLiteral(),
				)

				n.SemaResultType = types2.ErrorType()
				return types2.ErrorType()
			}

			if left.IsFloat() || right.IsFloat() {
				n.SemaResultType = types2.FloatType()
				return types2.FloatType()
			}

			n.SemaResultType = types2.IntType()
			return types2.IntType()

		case token.EqualEqual, token.BangEqual,
			token.Less, token.LessEqual,
			token.Greater, token.GreaterEqual:

			if !left.Equals(right) {
				a.errorf(
					n,
					"cannot compare %s with %s",
					left.Name,
					right.Name,
				)
				n.SemaResultType = types2.ErrorType()
				return types2.ErrorType()
			}

			n.SemaResultType = types2.BoolType()
			return types2.BoolType()
		}

		n.SemaResultType = types2.ErrorType()
		return types2.ErrorType()

	case *ast.MemberExpr:
		// the object is a type name used as a namespace, so it must be resolved before type inference
		if id, ok := n.Object.(*ast.Identifier); ok {
			if sym := a.lookup(id.Value); sym != nil && sym.Kind == SymbolEnum {
				for _, variant := range sym.Enum.Variants {
					if variant.Value == n.Property.Value {
						n.SemaResultType = types2.NominalType(sym.Enum.Name.Value)
						return n.SemaResultType
					}
				}

				a.errorf(n.Property, "enum '%s' has no variant '%s'", sym.Enum.Name.Value, n.Property.Value)
				n.SemaResultType = types2.ErrorType()
				return n.SemaResultType
			}
		}

		objectType := a.inferExprType(n.Object)
		if objectType.IsError() {
			n.SemaResultType = types2.ErrorType()
			return n.SemaResultType
		}

		if objectType.Kind != types2.Nominal {
			a.errorf(n.Object, "'%s' is not a struct", objectType.Name)
			n.SemaResultType = types2.ErrorType()
			return n.SemaResultType
		}

		strct := a.lookupStruct(objectType.Name)
		if strct == nil {
			a.errorf(n.Object, "'%s' is not a struct", objectType.Name)
			n.SemaResultType = types2.ErrorType()
			return n.SemaResultType
		}

		for _, field := range strct.Fields {
			if field.Name.Value == n.Property.Value {
				n.SemaResultType = field.SemaType
				return n.SemaResultType
			}
		}

		a.errorf(n.Property, "struct '%s' has no field '%s'", strct.Name.Value, n.Property.Value)
		n.SemaResultType = types2.ErrorType()
		return n.SemaResultType
	}

	a.errorf(expr, "internal compiler error: cannot infer type for expression")
	return types2.ErrorType()
}
