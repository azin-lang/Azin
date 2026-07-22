package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitStruct(
	s *ast.StructStmt,
) {
	a.Structs[s.Name.Value] = struct{}{}

	for _, field := range s.Fields {
		if field.Type != nil {
			a.markType(field.Type.Value)
		}
	}
}
