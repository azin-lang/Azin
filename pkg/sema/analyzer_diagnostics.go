package sema

import (
	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/token"
)

func (a *Analyzer) errorf(node ast.Node, format string, args ...any) {
	pos, length := sourceSpan(node)
	a.diag.ReportError(
		pos,
		int(length),
		format,
		args...,
	)
}

func (a *Analyzer) warningf(node ast.Node, format string, args ...any) {
	pos, length := sourceSpan(node)
	a.diag.ReportWarning(
		pos,
		int(length),
		format,
		args...,
	)
}

func sourceSpan(n ast.Node) (pos token.Position, length uint32) {
	pos = n.Pos()
	switch node := n.(type) {
	case *ast.Identifier:
		return pos, node.Token.Length

	case *ast.IntegerLiteral:
		return pos, node.Token.Length
	case *ast.FloatLiteral:
		return pos, node.Token.Length
	case *ast.StringLiteral:
		return pos, node.Token.Length
	case *ast.CharacterLiteral:
		return pos, node.Token.Length
	case *ast.BooleanLiteral:
		return pos, node.Token.Length
	case *ast.BinaryExpr:
		leftEnd := node.Left.Pos().Offset + spanLen(node.Left)
		rightEnd := node.Right.Pos().Offset + spanLen(node.Right)
		if rightEnd > leftEnd {
			return pos, rightEnd - pos.Offset
		}
		return pos, leftEnd - pos.Offset
	case *ast.MemberExpr:
		objEnd := node.Object.Pos().Offset + spanLen(node.Object)
		propEnd := node.Property.Token.Position.Offset + node.Property.Token.Length
		if objEnd > propEnd {
			return pos, objEnd - pos.Offset
		}
		return pos, propEnd - pos.Offset
	case *ast.CallExpr:
		end := node.Callee.Pos().Offset + spanLen(node.Callee)
		for _, arg := range node.Args {
			argEnd := arg.Pos().Offset + spanLen(arg)
			if argEnd > end {
				end = argEnd
			}
		}
		return pos, end - pos.Offset
	case *ast.VarStmt:
		end := node.Name.Token.Position.Offset + node.Name.Token.Length
		if node.SynType != nil {
			tEnd := node.SynType.Token.Position.Offset + node.SynType.Token.Length
			if tEnd > end {
				end = tEnd
			}
		}
		if node.Value != nil {
			vEnd := node.Value.Pos().Offset + spanLen(node.Value)
			if vEnd > end {
				end = vEnd
			}
		}
		return pos, end - pos.Offset
	case *ast.AssignmentStmt:
		end := node.Left.Pos().Offset + spanLen(node.Left)
		vEnd := node.Value.Pos().Offset + spanLen(node.Value)
		if vEnd > end {
			end = vEnd
		}
		return pos, end - pos.Offset
	}

	l := len(n.TokenLiteral())
	if l > 0 {
		length = uint32(l) //nolint:gosec
	}
	return
}

func spanLen(n ast.Node) uint32 {
	_, length := sourceSpan(n)
	return length
}
