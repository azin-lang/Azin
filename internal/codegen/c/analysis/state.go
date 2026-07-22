package analysis

import (
	"github.com/azin-lang/Azin/internal/ast"
)

type FuncInfo struct {
	Stmt  *ast.FuncStmt
	Index int
}

type Analyzer struct {
	Transpiler Transpiler

	Functions map[string]FuncInfo

	Calls map[string]map[string]struct{}

	Reachable map[string]struct{}

	Variables map[string]map[string]int

	Types map[string]struct{}

	Enums map[string]struct{}

	Structs map[string]struct{}

	TypeDependencies map[string]map[string]struct{}
}

type Transpiler interface {
	RequireInclude(string)

	SetFunctionIndex(string, int)
	FunctionIndex(string) (int, bool)

	RegisterForwardDeclaration(string, *ast.FuncStmt)

	SetEnum(string)
}

func New(t Transpiler) *Analyzer {
	return &Analyzer{
		Transpiler:       t,
		Functions:        make(map[string]FuncInfo),
		Calls:            make(map[string]map[string]struct{}),
		Reachable:        make(map[string]struct{}),
		Variables:        make(map[string]map[string]int),
		Types:            make(map[string]struct{}),
		Enums:            make(map[string]struct{}),
		Structs:          make(map[string]struct{}),
		TypeDependencies: make(map[string]map[string]struct{}),
	}
}

func (a *Analyzer) ResetTypes() {
	clear(a.Types)
}
