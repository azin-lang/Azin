package ast

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

/*
 * prints the AST to the console or exports it to a file.
 * @param node the AST node to print
 * @param export whether to export the AST to a file
 * @param outputPath the path to export the AST to
 */
func Print(node Node, export bool, outputPath string) {
	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	if export {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			fmt.Println(err)
			return
		}

		filepath := filepath.Join(outputPath, "ast.json")
		if err := os.WriteFile(filepath, data, 0644); err != nil {
			fmt.Println(err)
		}

		fmt.Println("AST saved to", filepath)
	} else {
		fmt.Println(string(data))
	}
}
