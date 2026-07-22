package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) RemoveUnusedTypes(program *ast.Program) {
	out := program.Statements[:0]
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {
		case *ast.StructStmt:
			if n.Name != nil {
				if _, ok := a.ReachableTypes[n.Name.Value]; !ok {
					continue // Prune dead struct
				}
			}
		case *ast.EnumStmt:
			if n.Name != nil {
				if _, ok := a.ReachableTypes[n.Name.Value]; !ok {
					continue // Prune dead enum
				}
			}
		}
		out = append(out, stmt)
	}
	program.Statements = out
}
