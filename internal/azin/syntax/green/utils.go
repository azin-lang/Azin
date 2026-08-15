package green

// IsNil safely determines whether a green Node interface is nil or holds a typed nil pointer.
// It avoids runtime panics caused by calling methods on typed nil interface values.
func IsNil(n Node) bool {
	return n == nil || n.IsNil()
}

// NilSafe ensures that typed nils are converted to true, untyped nils
// before they enter a node's children array.
func NilSafe(n Node) Node {
	if n == nil || n.IsNil() {
		return nil
	}
	return n
}

// ComputeProperties calculates the combined width and flags for a slice/array of nodes.
func ComputeProperties(nodes []Node) (width uint32, flags NodeFlags) {
	for _, node := range nodes {
		if node == nil {
			continue
		}

		width += node.FullWidth()
		flags |= node.Flags()
	}

	return
}

// GetLeadingTriviaWidth walks down the leftmost descendants to sum leading trivia width.
func GetLeadingTriviaWidth(n Node) uint32 {
	if IsNil(n) {
		return 0
	}
	if t, ok := n.(*Token); ok {
		if t.LeadingTrivia() != nil {
			return t.LeadingTrivia().FullWidth()
		}
		return 0
	}

	count := n.SlotCount()
	for i := range count {
		if child := n.Slot(i); child != nil {
			return GetLeadingTriviaWidth(child)
		}
	}
	return 0
}

// GetTrailingTriviaWidth walks down the rightmost descendants to sum trailing trivia width.
func GetTrailingTriviaWidth(n Node) uint32 {
	if IsNil(n) {
		return 0
	}
	if t, ok := n.(*Token); ok {
		if t.TrailingTrivia() != nil {
			return t.TrailingTrivia().FullWidth()
		}
		return 0
	}

	count := n.SlotCount()
	for i := count - 1; i >= 0; i-- {
		if child := n.Slot(i); child != nil {
			return GetTrailingTriviaWidth(child)
		}
	}
	return 0
}
