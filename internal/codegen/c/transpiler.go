package c

import (
	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/codegen/c/analysis"
)

type Transpiler struct {
	writer

	enums    map[string]struct{}
	includes map[string]struct{}

	funcIndices map[string]int

	lateFuncs []*ast.FuncStmt
	lateSet   map[string]struct{}
}

func New() *Transpiler {
	return &Transpiler{
		enums: make(map[string]struct{}),

		includes: make(map[string]struct{}),

		funcIndices: make(map[string]int),

		lateSet: make(map[string]struct{}),
	}
}

func (t *Transpiler) Transpile(
	program *ast.Program,
) string {
	t.reset()

	t.analyze(program)

	t.emit(program)

	return t.String()
}

func (t *Transpiler) analyze(program *ast.Program) {
	a := analysis.New(t)
	a.Analyze(program)

	for enum := range a.Enums {
		t.SetEnum(enum)
	}
}

func (t *Transpiler) reset() {
	t.writer.reset()

	clear(t.enums)
	clear(t.includes)
	clear(t.funcIndices)
	clear(t.lateSet)

	t.lateFuncs = nil
}

func (t *Transpiler) RequireInclude(
	header string,
) {
	t.includes[header] = struct{}{}
}

func (t *Transpiler) SetFunctionIndex(
	name string,
	index int,
) {
	t.funcIndices[name] = index
}

func (t *Transpiler) FunctionIndex(
	name string,
) (int, bool) {
	index, ok := t.funcIndices[name]
	return index, ok
}

func (t *Transpiler) RegisterForwardDeclaration(
	name string,
	fn *ast.FuncStmt,
) {
	if _, ok := t.lateSet[name]; ok {
		return
	}

	t.lateSet[name] = struct{}{}
	t.lateFuncs = append(
		t.lateFuncs,
		fn,
	)
}

func (t *Transpiler) SetEnum(name string) {
	t.enums[name] = struct{}{}
}
