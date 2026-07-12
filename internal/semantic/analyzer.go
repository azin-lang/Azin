package semantic

import "github.com/azin-lang/Azin/internal/ast"

type Analyzer struct {
	scopes []*Scope
}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) inferExprType(expr ast.Expr) *ast.Identifier {
	switch expr.(type) {

	case *ast.IntegerLiteral:
		return &ast.Identifier{Value: "int"}

	case *ast.FloatLiteral:
		return &ast.Identifier{Value: "float"}

	case *ast.CharacterLiteral:
		return &ast.Identifier{Value: "char"}

	case *ast.StringLiteral:
		return &ast.Identifier{Value: "string"}

	default:
		return nil
	}
}

func (a *Analyzer) inferFunctionReturnType(fn *ast.FuncStmt) {
	// Explicit return type always wins.
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

	// No return statement at all.
	fn.ReturnType = &ast.Identifier{Value: "unit"}
}

func (a *Analyzer) visitStatement(stmt ast.Stmt) {
	switch n := stmt.(type) {

	case *ast.FuncStmt:
		for _, stmt := range n.Body {
			a.visitStatement(stmt)
		}

		a.inferFunctionReturnType(n)

	case *ast.VarStmt:
		if n.Type == nil {
			n.Type = a.inferExprType(n.Value)
		}
	}
}

func (a *Analyzer) Analyze(program *ast.Program) error {
	for _, stmt := range program.Statements {
		a.visitStatement(stmt)
	}

	return nil
}
