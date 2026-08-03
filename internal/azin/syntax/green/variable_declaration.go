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

func NewVariableDeclaration(
	varKeyword *Token,
	identifier *Token,
	colon *Token,
	typeAnnotation Node,
	equals *Token,
	initializer Node,
) *VariableDeclaration {
	var width uint32
	nodes := []Node{varKeyword, identifier, colon, typeAnnotation, equals, initializer}
	for _, n := range nodes {
		if n != nil {
			width += n.FullWidth()
		}
	}

	return &VariableDeclaration{
		Base: Base{
			kind:      syntax.VariableDeclaration,
			fullWidth: width,
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
		return v.varKeyword
	case 1:
		return v.identifier
	case 2:
		return v.colon
	case 3:
		return v.typeAnnotation
	case 4:
		return v.equals
	case 5:
		return v.initializer
	default:
		return nil
	}
}
