package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) Analyze(program *ast.Program) {
	a.CollectFunctions(program)

	// Find reachable code
	a.WalkMain(program)

	a.ComputeReachable()

	// Remove dead functions
	a.RemoveUnusedFunctions(program)

	// Now analyze remaining program fully
	a.Walk(program)

	// Remove unused locals and types
	a.RemoveUnusedVariables(program)
	a.RemoveUnusedTypes(program)
}
