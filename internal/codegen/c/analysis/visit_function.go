package analysis

import (
	"github.com/azin-lang/Azin/pkg/ast"
)

func (a *Analyzer) visitFunction(fn *ast.FuncStmt) {
	name := FunctionName(fn)
	if fn.SemaReturnType.IsNominal() {
		a.MarkTypeUsed(fn.SemaReturnType)
	}
	for _, param := range fn.Params {
		if !param.SemaType.IsUseable() {
			continue
		}

		if param.SemaType.IsNominal() {
			a.MarkTypeUsed(param.SemaType)
		}

		a.registerVariable(name, param.Name.Value)
		a.useVariable(name, param.Name.Value)
	}
	for _, stmt := range fn.Body {
		a.visitStmt(name, stmt)
	}
}
