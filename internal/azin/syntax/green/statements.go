package green

import "github.com/azin-lang/Azin/internal/azin/syntax"

type BlockStatement struct {
	Base
	children [3]Node
}

func NewBlockStatement(startToken *Token, statements Node, endToken *Token) *BlockStatement {
	children := [3]Node{NilSafe(startToken), NilSafe(statements), NilSafe(endToken)}
	width, flags := ComputeProperties(children[:])
	return &BlockStatement{
		Base:     Base{kind: syntax.BlockStatement, fullWidth: width, flags: flags},
		children: children,
	}
}

func (b *BlockStatement) StartToken() *Token  { t, _ := b.children[0].(*Token); return t }
func (b *BlockStatement) Statements() Node    { return b.children[1] }
func (b *BlockStatement) EndToken() *Token    { t, _ := b.children[2].(*Token); return t }
func (b *BlockStatement) IsNil() bool         { return b == nil }
func (b *BlockStatement) SlotCount() int      { return 3 }
func (b *BlockStatement) Slot(index int) Node { return b.children[index] }

type IfStatement struct {
	Base
	children [5]Node
}

func NewIfStatement(ifKeyword *Token, condition, thenBranch Node, elseKeyword *Token, elseBranch Node) *IfStatement {
	children := [5]Node{NilSafe(ifKeyword), NilSafe(condition), NilSafe(thenBranch), NilSafe(elseKeyword), NilSafe(elseBranch)}
	width, flags := ComputeProperties(children[:])
	return &IfStatement{
		Base:     Base{kind: syntax.IfStatement, fullWidth: width, flags: flags},
		children: children,
	}
}

func (i *IfStatement) IfKeyword() *Token   { t, _ := i.children[0].(*Token); return t }
func (i *IfStatement) Condition() Node     { return i.children[1] }
func (i *IfStatement) ThenBranch() Node    { return i.children[2] }
func (i *IfStatement) ElseKeyword() *Token { t, _ := i.children[3].(*Token); return t }
func (i *IfStatement) ElseBranch() Node    { return i.children[4] }
func (i *IfStatement) IsNil() bool         { return i == nil }
func (i *IfStatement) SlotCount() int      { return 5 }
func (i *IfStatement) Slot(index int) Node { return i.children[index] }

type ReturnStatement struct {
	Base
	children [2]Node
}

func NewReturnStatement(returnKeyword *Token, expression Node) *ReturnStatement {
	children := [2]Node{NilSafe(returnKeyword), NilSafe(expression)}
	width, flags := ComputeProperties(children[:])
	return &ReturnStatement{
		Base:     Base{kind: syntax.ReturnStatement, fullWidth: width, flags: flags},
		children: children,
	}
}

func (r *ReturnStatement) ReturnKeyword() *Token { t, _ := r.children[0].(*Token); return t }
func (r *ReturnStatement) Expression() Node      { return r.children[1] }
func (r *ReturnStatement) IsNil() bool           { return r == nil }
func (r *ReturnStatement) SlotCount() int        { return 2 }
func (r *ReturnStatement) Slot(index int) Node   { return r.children[index] }
