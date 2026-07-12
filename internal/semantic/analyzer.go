package semantic

import "github.com/azin-lang/Azin/internal/ast"

type Analyzer struct {
	scopes []*Scope
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
			a.inferFunctionReturnType(n)

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

	return nil
}

func (a *Analyzer) visitStatement(stmt ast.Stmt) {
	switch n := stmt.(type) {

	case *ast.FuncStmt:
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

func (a *Analyzer) inferFunctionReturnType(fn *ast.FuncStmt) {
	if fn.ReturnType != nil {
		return
	}

	for _, stmt := range fn.Body {
		if ret, ok := stmt.(*ast.ReturnStmt); ok {
			if ret.Value == nil {
				fn.ReturnType = &ast.Identifier{Value: "unit"}
			} else {
				fn.ReturnType = a.inferExprType(ret.Value)
			}
			return
		}
	}

	fn.ReturnType = &ast.Identifier{Value: "unit"}
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
		if id, ok := n.Callee.(*ast.Identifier); ok {
			if sym := a.lookup(id.Value); sym != nil && sym.Kind == SymbolFunction {
				return sym.Type
			}
		}

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
