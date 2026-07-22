package analysis

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/ast"
)

func (a *Analyzer) ReportError(format string, args ...any) {
	a.Errors = append(a.Errors, fmt.Errorf(format, args...))
}

func (a *Analyzer) HasErrors() bool {
	return len(a.Errors) > 0
}

func (a *Analyzer) Analyze(program *ast.Program) {
	if program == nil {
		return
	}
	a.CollectFunctions(program)
	a.CollectTypes(program)

	if a.HasErrors() {
		return
	}

	a.BuildCallGraph()
	a.ComputeReachability()
	a.RemoveUnusedFunctions(program)

	a.AnalyzeUsage(program)
	a.RemoveUnusedVariables(program)

	a.CollectReachableTypes(program)
	a.RemoveUnusedTypes(program)
}
