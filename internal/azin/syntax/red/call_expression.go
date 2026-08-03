package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type CallExpression struct {
	*Node
}

func AsCallExpression(n *Node) *CallExpression {
	if n == nil || n.Kind() != syntax.CallExpression {
		return nil
	}
	return &CallExpression{Node: n}
}

func (c *CallExpression) Expression() *Node {
	return c.Child(0)
}

func (c *CallExpression) OpenParen() *Node {
	return c.Child(1)
}

func (c *CallExpression) CloseParen() *Node {
	if c == nil || c.greenNode == nil {
		return nil
	}
	slotCount := c.greenNode.SlotCount()
	if slotCount < 3 {
		return nil
	}
	return c.Child(slotCount - 1)
}

func (c *CallExpression) ArgumentCount() int {
	if c == nil || c.greenNode == nil {
		return 0
	}
	slotCount := c.greenNode.SlotCount()
	if slotCount < 3 {
		return 0
	}
	return slotCount - 3
}

func (c *CallExpression) Argument(index int) *Node {
	if index < 0 || index >= c.ArgumentCount() {
		return nil
	}
	return c.Child(2 + index)
}
