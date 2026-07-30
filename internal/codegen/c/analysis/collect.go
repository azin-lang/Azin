package analysis

import (
	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/types"
)

func (a *Analyzer) CollectTypes(program *ast.Program) {
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {
		case *ast.StructStmt:
			if n.Name == nil {
				continue
			}
			structName := n.Name.Value
			if _, exists := a.Structs[structName]; exists {
				a.ReportError("redeclaration of struct %q", structName)
				continue
			}

			a.Structs[structName] = n
			a.TypeDependencies[n.SemaType] = make(map[*types.TypeInfo]struct{})

			for _, field := range n.Fields {
				if field.SemaType.IsNominal() {
					a.TypeDependencies[n.SemaType][field.SemaType] = struct{}{}
				}
			}

		case *ast.EnumStmt:
			if n.Name == nil {
				continue
			}
			if _, exists := a.Enums[n.Name.Value]; exists {
				a.ReportError("redeclaration of enum %q", n.Name.Value)
				continue
			}
			a.Enums[n.Name.Value] = n
		}
	}
}

func (a *Analyzer) CollectImports(program *ast.Program) {
	for _, stmt := range program.Statements {
		if imp, ok := stmt.(*ast.ImportCStmt); ok {
			a.requireImport(imp.Path.Value)
		}
	}
}

func (a *Analyzer) CollectFunctions(program *ast.Program) {
	for index, stmt := range program.Statements {
		fn, ok := stmt.(*ast.FuncStmt)
		if !ok || fn == nil {
			continue
		}

		name := FunctionName(fn)
		if name == "" {
			continue
		}

		if _, exists := a.Functions[name]; exists {
			a.ReportError("redeclaration of function %q", name)
			continue
		}

		a.Functions[name] = FuncInfo{
			Stmt:  fn,
			Index: index,
		}
		a.Transpiler.SetFunctionIndex(name, index)
	}
}

func FunctionName(fn *ast.FuncStmt) string {
	if fn == nil {
		return ""
	}
	if fn.CName != "" {
		return fn.CName
	}
	if fn.Name != nil {
		return fn.Name.Value
	}
	return ""
}

// CollectReachableTypes scans the AST after dead functions and dead variables
// have been pruned, marking only remaining types (and their field dependencies) as reachable.
func (a *Analyzer) CollectReachableTypes(program *ast.Program) {
	if program == nil {
		return
	}

	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*ast.FuncStmt); ok {
			a.collectTypesFromFunction(fn)
		}
	}
}

func (a *Analyzer) collectTypesFromFunction(fn *ast.FuncStmt) {
	if fn == nil {
		return
	}

	// 1. Mark return type
	if fn.SemaReturnType.IsNominal() {
		a.MarkTypeUsed(fn.SemaReturnType)
	}

	// 2. Mark parameter types
	for _, param := range fn.Params {
		if param.SemaType.IsNominal() {
			a.MarkTypeUsed(param.SemaType)
		}
	}

	// 3. Mark types from live statements inside the function body
	a.collectTypesFromStmts(fn.Body)
}

func (a *Analyzer) collectTypesFromStmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		a.collectTypesFromStmt(stmt)
	}
}

func (a *Analyzer) collectTypesFromStmt(stmt ast.Stmt) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.VarStmt:
		if s.SemaType.IsNominal() {
			a.MarkTypeUsed(s.SemaType)
		}
		if s.Value != nil {
			a.collectTypesFromExpr(s.Value)
		}

	case *ast.IfStmt:
		if s.Condition != nil {
			a.collectTypesFromExpr(s.Condition)
		}
		a.collectTypesFromStmts(s.Then)
		a.collectTypesFromStmts(s.Else)

	case *ast.LoopStmt:
		a.collectTypesFromStmts(s.Body)

	case *ast.ReturnStmt:
		if s.Value != nil {
			a.collectTypesFromExpr(s.Value)
		}

	case *ast.AssignmentStmt:
		if s.Left != nil {
			a.collectTypesFromExpr(s.Left)
		}
		if s.Value != nil {
			a.collectTypesFromExpr(s.Value)
		}

	case *ast.ExpressionStmt:
		if s.Expression != nil {
			a.collectTypesFromExpr(s.Expression)
		}
	}
}

func (a *Analyzer) collectTypesFromExpr(expr ast.Expr) {
	if expr == nil {
		return
	}

	walkExpr(expr, func(e ast.Expr) bool {
		if node, ok := e.(*ast.MemberExpr); ok {
			// Catches static enum or struct member access like `Color.Red`
			if id, ok := node.Object.(*ast.Identifier); ok {
				name := id.Value
				if e, exists := a.Enums[name]; exists {
					a.MarkTypeUsed(e.SemaType)
				}
				if s, exists := a.Structs[name]; exists {
					a.MarkTypeUsed(s.SemaType)
				}
			}
		}
		return true
	})
}
