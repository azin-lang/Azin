package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type VariableDeclaration struct {
	Base
	children [6]Node
}

func NewVariableDeclaration(varKeyword, identifier, colon *Token, typeAnnotation Node, equals *Token, initializer Node) *VariableDeclaration {
	children := [6]Node{
		NilSafe(varKeyword), NilSafe(identifier), NilSafe(colon),
		NilSafe(typeAnnotation), NilSafe(equals), NilSafe(initializer),
	}
	width, flags := ComputeProperties(children[:])
	return &VariableDeclaration{
		Base:     Base{kind: syntax.VariableDeclaration, fullWidth: width, flags: flags},
		children: children,
	}
}

func (v *VariableDeclaration) VarKeyword() *Token   { t, _ := v.children[0].(*Token); return t }
func (v *VariableDeclaration) Identifier() *Token   { t, _ := v.children[1].(*Token); return t }
func (v *VariableDeclaration) Colon() *Token        { t, _ := v.children[2].(*Token); return t }
func (v *VariableDeclaration) TypeAnnotation() Node { return v.children[3] }
func (v *VariableDeclaration) Equals() *Token       { t, _ := v.children[4].(*Token); return t }
func (v *VariableDeclaration) Initializer() Node    { return v.children[5] }
func (v *VariableDeclaration) IsNil() bool          { return v == nil }
func (v *VariableDeclaration) SlotCount() int       { return 6 }
func (v *VariableDeclaration) Slot(index int) Node  { return v.children[index] }
