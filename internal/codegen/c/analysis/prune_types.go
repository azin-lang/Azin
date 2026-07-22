package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) RemoveUnusedTypes(
	program *ast.Program,
) {
	a.expandTypeDependencies()

	out := program.Statements[:0]

	for _, stmt := range program.Statements {

		switch n := stmt.(type) {

		case *ast.StructStmt:

			if !a.TypeUsed(n.Name.Value) {
				continue
			}

		case *ast.EnumStmt:

			if !a.TypeUsed(n.Name.Value) {
				continue
			}
		}

		out = append(out, stmt)
	}

	program.Statements = out
}

func (a *Analyzer) TypeUsed(
	name string,
) bool {
	_, ok := a.Types[name]
	return ok
}

func (a *Analyzer) expandTypeDependencies() {
	changed := true

	for changed {

		changed = false

		for used := range a.Types {

			deps := a.TypeDependencies[used]

			for dep := range deps {

				if _, exists := a.Types[dep]; exists {
					continue
				}

				a.Types[dep] = struct{}{}
				changed = true
			}
		}
	}
}
