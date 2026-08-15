package main

import (
	"fmt"

	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/syntax/green"
	"github.com/azin-lang/Azin/internal/azin/syntax/lexer"
	"github.com/azin-lang/Azin/internal/azin/syntax/parser"
	"github.com/azin-lang/Azin/internal/azin/syntax/red"
	"github.com/azin-lang/Azin/internal/azin/text"
)

func main() {
	code := []byte(`
var x: int = (10 + 2 * max(a, @b == 42) && isReady.not)
var y = 10
`)

	file := text.NewSourceText("showcase.az", 1, code)
	diags := diagnostics.NewCollector()

	lex := lexer.New(file, diags)
	p := parser.New(lex, diags)

	root := red.NewRoot(p.ParseCompilationUnit())

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

	// Extract the literal source if the underlying green node is a Token
	source := ""
	if tok, ok := greenNode.(*green.Token); ok {
		source = fmt.Sprintf(" %q", tok.Text())
	}

	fmt.Printf("%s%s%v%s\n", indent, marker, kind, source)

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
