package analysis

import (
	"path/filepath"
)

func (a *Analyzer) registerVariable(
	function string,
	name string,
) {
	if function == "" {
		return
	}

	vars, ok := a.Variables[function]

	if !ok {
		vars = make(map[string]int)
		a.Variables[function] = vars
	}

	vars[name] = 0
}

func (a *Analyzer) useVariable(
	function string,
	name string,
) {
	vars := a.Variables[function]

	if vars == nil {
		return
	}

	if _, ok := vars[name]; ok {
		vars[name]++
	}
}

func (a *Analyzer) markType(name string) {
	if name == "" {
		return
	}

	a.Types[name] = struct{}{}

	switch name {
	case "bool":
		a.Transpiler.RequireInclude("stdbool.h")
	}
}

func (a *Analyzer) requireImport(path string) {
	if filepath.Ext(path) == "" {
		path += ".h"
	}

	a.Transpiler.RequireInclude(path)
}
