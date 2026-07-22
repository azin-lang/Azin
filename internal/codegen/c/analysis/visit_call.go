package analysis

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/ast"
)

func (a *Analyzer) visitCall(
	current string,
	call *ast.CallExpr,
) {
	if call.ResolvedName == "" {
		panic(
			"internal compiler error: unresolved call",
		)
	}

	name := call.ResolvedName

	if builtin, ok := LookupBuiltin(name); ok {

		a.Transpiler.RequireInclude(
			builtin.Include,
		)

	} else if info, ok := a.Functions[name]; ok {

		a.addCall(
			current,
			name,
		)

		currentIndex, exists :=
			a.Transpiler.FunctionIndex(current)

		if !exists ||
			info.Index > currentIndex {

			a.Transpiler.RegisterForwardDeclaration(
				name,
				info.Stmt,
			)
		}

	} else {

		panic(fmt.Sprintf(
			"unknown function %q",
			name,
		))
	}

	a.visitExpr(
		current,
		call.Callee,
	)

	for _, arg := range call.Args {
		a.visitExpr(
			current,
			arg,
		)
	}
}
