package main

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/lexer"
	"github.com/azin-lang/Azin/internal/azin/parser"
	"github.com/azin-lang/Azin/internal/azin/source"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/syntax/red"
)

func main() {
	code := []byte(`(10 + 2 * max(a, @b == 42) && !isReady)`)

	file := source.NewSourceText("showcase.az", 1, code)
	diags := diagnostics.NewCollector()

	lex := lexer.New(file, diags)
	p := parser.New(lex, diags)

	root := red.NewRoot(p.ParseExpression())

	printTree(root, "", true)
	fmt.Println()
	if diags.HasErrors() {
		fmt.Println(diags.PrintAll())
	} else {
		fmt.Println("Compilation successful with zero errors.")
	}
}

// printTree recursively pretty-prints the red.Node AST.
func printTree(node *red.Node, indent string, isLast bool) {
	if node == nil {
		return
	}

	marker := "├── "
	if isLast {
		marker = "└── "
	}

	greenNode := node.Green()
	kind := greenNode.Kind()

	// Extract the literal text if the underlying green node is a Token
	text := ""
	if tok, ok := greenNode.(*green.Token); ok {
		text = fmt.Sprintf(" %q", tok.Text())
	}

	fmt.Printf("%s%s%v%s\n", indent, marker, kind, text)

	childIndent := indent
	if isLast {
		childIndent += "    "
	} else {
		childIndent += "│   "
	}

	// Dynamically find children using the Green tree's slot count.
	// This works for any node type without needing manual casting.
	type slotter interface {
		SlotCount() int
	}

	var children []*red.Node
	if s, ok := greenNode.(slotter); ok {
		count := s.SlotCount()
		for i := range count {
			if child := node.Child(i); child != nil {
				children = append(children, child)
			}
		}
	}

	for i, child := range children {
		printTree(child, childIndent, i == len(children)-1)
	}
}
