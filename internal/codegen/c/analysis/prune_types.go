package analysis

import (
	"github.com/azin-lang/Azin/pkg/ast"
)

func (a *Analyzer) RemoveUnusedTypes(program *ast.Program) {
	out := program.Statements[:0]
	for _, stmt := range program.Statements {
		switch n := stmt.(type) {
		case *ast.StructStmt:
			if n.Name != nil {
				if _, ok := a.ReachableTypes[n.SemaType]; !ok {
					continue // Prune dead struct
				}
			}
		case *ast.EnumStmt:
			if n.Name != nil {
				if _, ok := a.ReachableTypes[n.SemaType]; !ok {
					continue // Prune dead enum
				}
			}
		}
		out = append(out, stmt)
	}
	program.Statements = out
}
