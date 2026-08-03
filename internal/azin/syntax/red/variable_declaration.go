package red

import "github.com/azin-lang/Azin/internal/azin/syntax"

type VariableDeclaration struct {
	*Node
}

func AsVariableDeclaration(n *Node) *VariableDeclaration {
	if n == nil || n.Kind() != syntax.VariableDeclaration {
		return nil
	}
	return &VariableDeclaration{Node: n}
}

func (v *VariableDeclaration) VarKeyword() *Node     { return v.Child(0) }
func (v *VariableDeclaration) Identifier() *Node     { return v.Child(1) }
func (v *VariableDeclaration) Colon() *Node          { return v.Child(2) }
func (v *VariableDeclaration) TypeAnnotation() *Node { return v.Child(3) }
func (v *VariableDeclaration) Equals() *Node         { return v.Child(4) }
func (v *VariableDeclaration) Initializer() *Node    { return v.Child(5) }
