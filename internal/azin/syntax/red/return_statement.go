package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type ReturnStatement struct {
	*Node
}

func AsReturnStatement(n *Node) *ReturnStatement {
	if n == nil || n.Kind() != syntax.ReturnStatement {
		return nil
	}
	return &ReturnStatement{Node: n}
}

func (r *ReturnStatement) ReturnKeyword() *Node { return r.Child(0) }
func (r *ReturnStatement) Expression() *Node    { return r.Child(1) }
