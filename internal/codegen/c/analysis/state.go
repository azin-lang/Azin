package analysis

import "github.com/azin-lang/Azin/internal/ast"

type FuncInfo struct {
	Stmt  *ast.FuncStmt
	Index int
}

type Analyzer struct {
	Transpiler Transpiler
	Errors     []error

	Functions map[string]FuncInfo
	Structs   map[string]*ast.StructStmt
	Enums     map[string]*ast.EnumStmt

	Calls              map[string]map[string]struct{}
	ReachableFunctions map[string]struct{}
	TypeDependencies   map[string]map[string]struct{}
	ReachableTypes     map[string]struct{}

	Variables map[string]map[string]int
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
		Transpiler:         t,
		Errors:             make([]error, 0),
		Functions:          make(map[string]FuncInfo),
		Structs:            make(map[string]*ast.StructStmt),
		Enums:              make(map[string]*ast.EnumStmt),
		Calls:              make(map[string]map[string]struct{}),
		ReachableFunctions: make(map[string]struct{}),
		TypeDependencies:   make(map[string]map[string]struct{}),
		ReachableTypes:     make(map[string]struct{}),
		Variables:          make(map[string]map[string]int),
	}
}

func (a *Analyzer) Reset() {
	a.Errors = a.Errors[:0]
	clear(a.Functions)
	clear(a.Structs)
	clear(a.Enums)
	clear(a.Calls)
	clear(a.ReachableFunctions)
	clear(a.TypeDependencies)
	clear(a.ReachableTypes)
	clear(a.Variables)
}
