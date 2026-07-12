package semantic

import (
	"strconv"

	"github.com/azin-lang/Azin/internal/ast"
)

type Analyzer struct {
	scopes []*Scope

	currentFunction *ast.FuncStmt
}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(program *ast.Program) error {
	a.pushScope()
	defer a.popScope()

	// Register every top level symbol first
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {

		case *ast.FuncStmt:
			a.declare(&Symbol{
				Name:     n.Name.Value,
				Type:     n.ReturnType,
				Kind:     SymbolFunction,
				Function: n,
			})

		case *ast.StructStmt:
			a.declare(&Symbol{
				Name:   n.Name.Value,
				Kind:   SymbolStruct,
				Struct: n,
			})
		}
	}

	// Analyze every statement.
	for _, stmt := range program.Statements {
		a.visitStatement(stmt)
	}

	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FuncStmt); ok {
			a.inferFunctionReturnType(fn)

			sym := a.lookup(fn.Name.Value)
			sym.Type = fn.ReturnType
		}
	}

	return nil
}

func (a *Analyzer) visitStatement(stmt ast.Stmt) {
	switch n := stmt.(type) {

	case *ast.FuncStmt:
		old := a.currentFunction
		a.currentFunction = n
		defer func() {
			a.currentFunction = old
		}()

		a.pushScope()

		// Register parameters.
		for _, param := range n.Params {
			a.declare(&Symbol{
				Name: param.Name.Value,
				Type: param.Type,
				Kind: SymbolVariable,
			})
		}

		for _, stmt := range n.Body {
			a.visitStatement(stmt)
		}

		a.popScope()

	case *ast.ReturnStmt:
		if a.currentFunction == nil {
			return
		}

		actual := &ast.Identifier{Value: "unit"}
		if n.Value != nil {
			actual = a.inferExprType(n.Value)
		}

		expected := a.currentFunction.ReturnType

		if expected != nil && actual != nil && expected.Value != actual.Value {
			panic(
				"return type mismatch: expected " +
					expected.Value +
					", got " +
					actual.Value,
			)
		}

	case *ast.VarStmt:
		if n.Type == nil {
			n.Type = a.inferExprType(n.Value)
		}

		a.declare(&Symbol{
			Name: n.Name.Value,
			Type: n.Type,
			Kind: SymbolVariable,
		})

	case *ast.IfStmt:
		a.pushScope()

		for _, stmt := range n.Then {
			a.visitStatement(stmt)
		}

		a.popScope()

		a.pushScope()

		for _, stmt := range n.Else {
			a.visitStatement(stmt)
		}

		a.popScope()
	}
}

func (a *Analyzer) findReturnType(stmt ast.Stmt) *ast.Identifier {
	switch n := stmt.(type) {

	case *ast.ReturnStmt:
		if n.Value == nil {
			return &ast.Identifier{Value: "unit"}
		}
		return a.inferExprType(n.Value)

	case *ast.IfStmt:
		for _, s := range n.Then {
			if t := a.findReturnType(s); t != nil {
				return t
			}
		}

		for _, s := range n.Else {
			if t := a.findReturnType(s); t != nil {
				return t
			}
		}
	}

	return nil
}

func (a *Analyzer) inferFunctionReturnType(fn *ast.FuncStmt) {

	if fn.ReturnType != nil {
		return
	}

	sym := a.lookup(fn.Name.Value)
	if sym != nil {
		if sym.Inferring {
			return
		}
		sym.Inferring = true
		defer func() {
			sym.Inferring = false
		}()
	}

	for _, stmt := range fn.Body {
		if typ := a.findReturnType(stmt); typ != nil {
			fn.ReturnType = typ

			if sym != nil {
				sym.Type = typ
			}

			return
		}
	}

	fn.ReturnType = &ast.Identifier{Value: "unit"}

	if sym != nil {
		sym.Type = fn.ReturnType
	}
}

func (a *Analyzer) inferExprType(expr ast.Expr) *ast.Identifier {
	switch n := expr.(type) {

	case *ast.IntegerLiteral:
		return &ast.Identifier{Value: "int"}

	case *ast.FloatLiteral:
		return &ast.Identifier{Value: "float"}

	case *ast.CharacterLiteral:
		return &ast.Identifier{Value: "char"}

	case *ast.StringLiteral:
		return &ast.Identifier{Value: "string"}

	case *ast.Identifier:
		if sym := a.lookup(n.Value); sym != nil {
			return sym.Type
		}

	case *ast.CallExpr:
		id, ok := n.Callee.(*ast.Identifier)
		if !ok {
			return nil
		}

		sym := a.lookup(id.Value)
		if sym == nil || sym.Kind != SymbolFunction {
			panic("unknown function: " + id.Value)
		}

		if sym.Inferring {
			return nil
		}

		if sym.Type == nil {
			a.inferFunctionReturnType(sym.Function)
		}

		if len(n.Args) != len(sym.Function.Params) {
			panic("wrong number of arguments to " + id.Value)
		}

		for i, arg := range n.Args {
			got := a.inferExprType(arg)
			want := sym.Function.Params[i].Type

			if got == nil {
				continue
			}

			if got.Value != want.Value {
				panic(
					"argument " +
						strconv.Itoa(i+1) +
						" of " +
						id.Value +
						": expected " +
						want.Value +
						", got " +
						got.Value,
				)
			}
		}

		return sym.Type

	case *ast.BinaryExpr:
		left := a.inferExprType(n.Left)
		right := a.inferExprType(n.Right)

		if left == nil || right == nil {
			return nil
		}

		if left.Value == "float" || right.Value == "float" {
			return &ast.Identifier{Value: "float"}
		}

		return left

	case *ast.MemberExpr:
		// TODO: struct field lookup
		return nil
	}

	return nil
}
