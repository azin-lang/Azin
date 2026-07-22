package analysis

import (
	"github.com/azin-lang/Azin/internal/ast"
)

func (a *Analyzer) visitCall(current string, call *ast.CallExpr) {
	name := call.ResolvedName
	if name == "" {
		a.ReportError("unresolved function call in %s", current)
		return
	}

	if builtin, ok := LookupBuiltin(name); ok {
		a.Transpiler.RequireInclude(builtin.Include)
	} else if info, ok := a.Functions[name]; ok {
		a.addCall(current, name)
		currentIndex, exists := a.Transpiler.FunctionIndex(current)
		// Ensure forward declaration if calling a function defined later
		if !exists || info.Index > currentIndex {
			a.Transpiler.RegisterForwardDeclaration(name, info.Stmt)
		}
	} else {
		a.ReportError("call to unknown function %q in %s", name, current)
	}

	a.visitExpr(current, call.Callee)
	for _, arg := range call.Args {
		a.visitExpr(current, arg)
	}
}
