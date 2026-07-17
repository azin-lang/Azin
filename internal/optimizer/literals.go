package optimizer

import "github.com/azin-lang/Azin/internal/ast"

func intLit(v int64) *ast.IntegerLiteral {
	return &ast.IntegerLiteral{Value: v}
}

func floatLit(v float64) *ast.FloatLiteral {
	return &ast.FloatLiteral{Value: v}
}

func boolLit(v bool) *ast.BooleanLiteral {
	return &ast.BooleanLiteral{Value: v}
}

func charAsInt(c *ast.CharacterLiteral) *ast.IntegerLiteral {
	return intLit(int64(c.Value))
}
