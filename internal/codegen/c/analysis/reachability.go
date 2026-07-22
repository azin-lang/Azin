package analysis

func (a *Analyzer) ComputeReachable() {

	a.Reachable = make(
		map[string]struct{},
	)

	a.reach("main")
}

func (a *Analyzer) reach(
	name string,
) {
	if _, ok := a.Reachable[name]; ok {
		return
	}

	a.Reachable[name] = struct{}{}

	for next := range a.Calls[name] {
		a.reach(next)
	}
}
