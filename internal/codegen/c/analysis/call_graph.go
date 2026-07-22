package analysis

func (a *Analyzer) addCall(
	from string,
	to string,
) {
	if from == "" {
		return
	}

	calls, ok := a.Calls[from]

	if !ok {
		calls = make(
			map[string]struct{},
		)

		a.Calls[from] = calls
	}

	calls[to] = struct{}{}
}
