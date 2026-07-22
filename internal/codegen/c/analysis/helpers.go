package analysis

import (
	"path/filepath"

	"github.com/azin-lang/Azin/internal/ast"
)

// registerVariable adds a local variable to the scope tracker initialized with zero usages.
func (a *Analyzer) registerVariable(function, name string) {
	if function == "" || name == "" {
		return
	}
	if a.Variables[function] == nil {
		a.Variables[function] = make(map[string]int)
	}
	a.Variables[function][name] = 0
}

// useVariable safely increments the read count of a scoped variable.
func (a *Analyzer) useVariable(function, name string) {
	if vars := a.Variables[function]; vars != nil {
		if _, exists := vars[name]; exists {
			vars[name]++
		}
	}
}

// MarkTypeUsed resolves a struct or enum type and recursively marks all of its dependent types as used.
func (a *Analyzer) MarkTypeUsed(name string) {
	if name == "" {
		return
	}
	if _, ok := a.ReachableTypes[name]; ok {
		return // Break cyclic dependencies
	}

	a.ReachableTypes[name] = struct{}{}

	if name == "bool" {
		a.Transpiler.RequireInclude("stdbool.h")
	}

	for dep := range a.TypeDependencies[name] {
		a.MarkTypeUsed(dep)
	}
}

func (a *Analyzer) requireImport(path string) {
	if path == "" {
		return
	}
	if filepath.Ext(path) == "" {
		path += ".h"
	}
	a.Transpiler.RequireInclude(path)
}

// walkExpr performs a depth-first search on expressions.
// The visitor function should return true to continue descending, or false to halt the current branch.
func walkExpr(expr ast.Expr, visit func(ast.Expr) bool) {
	if expr == nil || !visit(expr) {
		return
	}
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		walkExpr(e.Left, visit)
		walkExpr(e.Right, visit)
	case *ast.CallExpr:
		walkExpr(e.Callee, visit)
		for _, arg := range e.Args {
			walkExpr(arg, visit)
		}
	case *ast.MemberExpr:
		walkExpr(e.Object, visit)
	}
}
