package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

// MaxChildrenPerNode is the branching threshold before a SyntaxList is balanced into a tree.
const MaxChildrenPerNode = 10

var emptySyntaxList = &SyntaxList{
	Base: Base{
		kind:      syntax.SyntaxList,
		fullWidth: 0,
		flags:     FlagNone,
	},
}

type SyntaxList struct {
	Base
	children []Node
}

func (s *SyntaxList) IsNil() bool { return s == nil }

// NewSyntaxList ensures immutability by cloning the external slice,
// then balances the tree using a zero-copy internal builder.
func NewSyntaxList(children []Node) *SyntaxList {
	if len(children) == 0 {
		return emptySyntaxList
	}

	clone := make([]Node, len(children))
	copy(clone, children)
	return buildSyntaxListUnsafe(clone)
}

// buildSyntaxListUnsafe creates balanced chunks without redundant cloning.
func buildSyntaxListUnsafe(children []Node) *SyntaxList {
	for len(children) > MaxChildrenPerNode {
		numChunks := (len(children) + MaxChildrenPerNode - 1) / MaxChildrenPerNode
		chunks := make([]Node, 0, numChunks)
		for i := 0; i < len(children); i += MaxChildrenPerNode {
			end := min(i+MaxChildrenPerNode, len(children))
			chunks = append(chunks, buildSyntaxListUnsafe(children[i:end]))
		}
		children = chunks
	}

	width, flags := ComputePropertiesSlice(children)
	return &SyntaxList{
		Base: Base{
			kind:      syntax.SyntaxList,
			fullWidth: width,
			flags:     flags,
		},
		children: children,
	}
}

func (s *SyntaxList) SlotCount() int {
	if s == nil {
		return 0
	}
	return len(s.children)
}

func (s *SyntaxList) Slot(index int) Node {
	if s == nil || index < 0 || index >= len(s.children) {
		return nil
	}
	return s.children[index]
}
