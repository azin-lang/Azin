package green

// IsNil safely determines whether a green Node interface is nil or holds a typed nil pointer.
// It avoids runtime panics caused by calling methods on typed nil interface values.
func IsNil(n Node) bool {
	return n == nil || n.IsNil()
}

// getProps returns fullWidth and NodeFlags for a node, safely handling nil and typed-nil inputs.
func getProps(n Node) (int, NodeFlags) {
	if IsNil(n) {
		return 0, FlagNone
	}
	return n.FullWidth(), n.Flags()
}

func ComputeProperties1(n1 Node) (int, NodeFlags) {
	return getProps(n1)
}

func ComputeProperties2(n1, n2 Node) (int, NodeFlags) {
	w1, f1 := getProps(n1)
	w2, f2 := getProps(n2)
	return w1 + w2, f1 | f2
}

func ComputeProperties3(n1, n2, n3 Node) (int, NodeFlags) {
	w1, f1 := getProps(n1)
	w2, f2 := getProps(n2)
	w3, f3 := getProps(n3)
	return w1 + w2 + w3, f1 | f2 | f3
}

func ComputeProperties4(n1, n2, n3, n4 Node) (int, NodeFlags) {
	w1, f1 := getProps(n1)
	w2, f2 := getProps(n2)
	w3, f3 := getProps(n3)
	w4, f4 := getProps(n4)
	return w1 + w2 + w3 + w4, f1 | f2 | f3 | f4
}

func ComputeProperties5(n1, n2, n3, n4, n5 Node) (int, NodeFlags) {
	w1, f1 := getProps(n1)
	w2, f2 := getProps(n2)
	w3, f3 := getProps(n3)
	w4, f4 := getProps(n4)
	w5, f5 := getProps(n5)
	return w1 + w2 + w3 + w4 + w5, f1 | f2 | f3 | f4 | f5
}

func ComputeProperties6(n1, n2, n3, n4, n5, n6 Node) (int, NodeFlags) {
	w1, f1 := getProps(n1)
	w2, f2 := getProps(n2)
	w3, f3 := getProps(n3)
	w4, f4 := getProps(n4)
	w5, f5 := getProps(n5)
	w6, f6 := getProps(n6)
	return w1 + w2 + w3 + w4 + w5 + w6, f1 | f2 | f3 | f4 | f5 | f6
}

func ComputePropertiesSlice(nodes []Node) (int, NodeFlags) {
	var totalWidth int
	var totalFlags NodeFlags
	for _, n := range nodes {
		w, f := getProps(n)
		totalWidth += w
		totalFlags |= f
	}
	return totalWidth, totalFlags
}

// GetLeadingTriviaWidth walks down the leftmost descendants to sum leading trivia width.
func GetLeadingTriviaWidth(n Node) int {
	if IsNil(n) {
		return 0
	}
	if t, ok := n.(*Token); ok {
		w, _ := getProps(t.LeadingTrivia())
		return w
	}
	// Note: Short-circuits faster if the leftmost slot is nil but subsequent ones aren't.
	for i := range n.SlotCount() {
		if child := n.Slot(i); !IsNil(child) {
			return GetLeadingTriviaWidth(child)
		}
	}
	return 0
}

// GetTrailingTriviaWidth walks down the rightmost descendants to sum trailing trivia width.
func GetTrailingTriviaWidth(n Node) int {
	if IsNil(n) {
		return 0
	}
	if t, ok := n.(*Token); ok {
		w, _ := getProps(t.TrailingTrivia())
		return w
	}
	for i := n.SlotCount() - 1; i >= 0; i-- {
		if child := n.Slot(i); !IsNil(child) {
			return GetTrailingTriviaWidth(child)
		}
	}
	return 0
}

// SafeNode ensures that typed nils are converted to true, untyped nils.
func SafeNode(n Node) Node {
	if IsNil(n) {
		return nil
	}
	return n
}
