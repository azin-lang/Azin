package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type VariableDeclaration struct {
	Base
	varKeyword     *Token
	identifier     *Token
	colon          *Token
	typeAnnotation Node
	equals         *Token
	initializer    Node
}

func NewVariableDeclaration(varKeyword *Token, identifier *Token, colon *Token, typeAnnotation Node, equals *Token, initializer Node) *VariableDeclaration {
	width, flags := ComputeProperties6(varKeyword, identifier, colon, typeAnnotation, equals, initializer)
	return &VariableDeclaration{
		Base: Base{
			kind:      syntax.VariableDeclaration,
			fullWidth: width,
			flags:     flags,
		},
		varKeyword:     varKeyword,
		identifier:     identifier,
		colon:          colon,
		typeAnnotation: typeAnnotation,
		equals:         equals,
		initializer:    initializer,
	}
}

func (v *VariableDeclaration) VarKeyword() *Token   { return v.varKeyword }
func (v *VariableDeclaration) Identifier() *Token   { return v.identifier }
func (v *VariableDeclaration) Colon() *Token        { return v.colon }
func (v *VariableDeclaration) TypeAnnotation() Node { return v.typeAnnotation }
func (v *VariableDeclaration) Equals() *Token       { return v.equals }
func (v *VariableDeclaration) Initializer() Node    { return v.initializer }

func (v *VariableDeclaration) SlotCount() int { return 6 }
func (v *VariableDeclaration) Slot(index int) Node {
	switch index {
	case 0:
		if v.varKeyword == nil {
			return nil
		}
		return v.varKeyword
	case 1:
		if v.identifier == nil {
			return nil
		}
		return v.identifier
	case 2:
		if v.colon == nil {
			return nil
		}
		return v.colon
	case 3:
		if v.typeAnnotation == nil {
			return nil
		}
		return v.typeAnnotation
	case 4:
		if v.equals == nil {
			return nil
		}
		return v.equals
	case 5:
		if v.initializer == nil {
			return nil
		}
		return v.initializer
	default:
		return nil
	}
}
