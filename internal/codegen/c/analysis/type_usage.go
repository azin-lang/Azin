package analysis

func (a *Analyzer) markTypeName(
	name string,
) {
	if name == "" {
		return
	}

	a.Types[name] = struct{}{}

	switch name {

	case "bool":
		a.Transpiler.RequireInclude(
			"stdbool.h",
		)

	case "string":
		a.Transpiler.RequireInclude(
			"string.h",
		)
	}
}
