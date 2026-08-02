package main

import (
	"fmt"
	"strings"

	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/syntax/red"
)

func main() {
	fmt.Println("Source code: \"2020 + 1024\"")
	fmt.Println(strings.Repeat("-", 54))

	space := green.NewTrivia(syntax.WhitespaceTrivia, " ")

	t1 := green.NewToken(syntax.IntegerLiteralToken, "2020", nil, space)
	tPlus := green.NewToken(syntax.PlusToken, "+", nil, space)
	t2 := green.NewToken(syntax.IntegerLiteralToken, "1024", nil, nil)

	greenRoot := green.NewBinaryExpression(t1, tPlus, t2)

	redRoot := red.NewRoot(greenRoot)
	expr := red.AsBinaryExpression(redRoot)

	fmt.Printf("%-30s | %-8s | %-10s\n", "NODE KIND", "POS", "WIDTH")
	fmt.Println(strings.Repeat("-", 54))

	fmt.Printf("%-30s | %-8d | %-10d\n",
		expr.Kind().String(),
		expr.Position(),
		expr.FullWidth(),
	)

	left := expr.Left()
	fmt.Printf("  ├── %-24s | %-8d | %-10d\n",
		left.Kind().String(),
		left.Position(),
		left.FullWidth(),
	)

	op := expr.Operator()
	fmt.Printf("  ├── %-24s | %-8d | %-10d\n",
		op.Kind().String(),
		op.Position(),
		op.FullWidth(),
	)

	right := expr.Right()
	fmt.Printf("  └── %-24s | %-8d | %-10d\n",
		right.Kind().String(),
		right.Position(),
		right.FullWidth(),
	)
}
