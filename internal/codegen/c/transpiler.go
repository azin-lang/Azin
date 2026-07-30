package c

import (
	"fmt"
	"slices"
	"strings"

	"github.com/azin-lang/Azin/internal/codegen/c/analysis"
	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/types"
)

type Transpiler struct {
	writer

	enums          map[string]struct{}
	structs        map[string]struct{}
	reachableTypes map[*types.TypeInfo]struct{}
	structDeps     map[string][]string
	includes       map[string]struct{}

	funcIndices map[string]int

	lateFuncs []*ast.FuncStmt
	lateSet   map[string]struct{}

	defers []ast.Expr
}

func New() *Transpiler {
	return &Transpiler{
		enums:          make(map[string]struct{}),
		structs:        make(map[string]struct{}),
		reachableTypes: make(map[*types.TypeInfo]struct{}),
		structDeps:     make(map[string][]string),
		includes:       make(map[string]struct{}),
		funcIndices:    make(map[string]int),
		lateSet:        make(map[string]struct{}),
	}
}

func (t *Transpiler) Transpile(
	program *ast.Program,
) (string, error) {
	t.reset()

	if err := t.analyze(program); err != nil {
		return "", err
	}

	t.emit(program)

	return t.String(), t.err
}

func (t *Transpiler) analyze(program *ast.Program) error {
	a := analysis.New(t)
	a.Analyze(program)

	if a.HasErrors() {
		errs := make([]string, len(a.Errors))
		for i, err := range a.Errors {
			errs[i] = err.Error()
		}
		return fmt.Errorf("analysis: %s", strings.Join(errs, "; "))
	}

	for enum := range a.Enums {
		t.SetEnum(enum)
	}

	// Capture reachable types from the analyzer
	for tname := range a.ReachableTypes {
		t.reachableTypes[tname] = struct{}{}
	}

	// Map all structs and their dependencies to resolve cycles later
	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.StructStmt); ok {
			t.structs[s.Name.Value] = struct{}{}

			var deps []string
			for _, f := range s.Fields {
				if f.SemaType.IsComplete() {
					deps = append(deps, f.SemaType.Name)
				}
			}
			t.structDeps[s.Name.Value] = deps
		}
	}

	return nil
}

func (t *Transpiler) reset() {
	t.writer.reset()

	clear(t.enums)
	clear(t.structs)
	clear(t.reachableTypes)
	clear(t.structDeps)
	clear(t.includes)
	clear(t.funcIndices)
	clear(t.lateSet)

	t.lateFuncs = nil
}

// Detects if 'fieldType' creates a cycle back to 'parentStruct'
func (t *Transpiler) isCyclicField(parentStruct, fieldType string) bool {
	if _, isStruct := t.structs[fieldType]; !isStruct {
		return false
	}

	// Direct self-reference (e.g. LinkedListNode -> LinkedListNode)
	if parentStruct == fieldType {
		return true
	}

	// DFS to see if fieldType eventually contains parentStruct
	visited := make(map[string]bool)
	var dfs func(string) bool
	dfs = func(curr string) bool {
		if curr == parentStruct {
			return true
		}
		if visited[curr] {
			return false
		}
		visited[curr] = true
		return slices.ContainsFunc(t.structDeps[curr], dfs)
	}

	return dfs(fieldType)
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
