package ast

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

var (
	cNode    = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	cField   = color.New(color.FgYellow).SprintFunc()
	cValue   = color.New(color.FgHiGreen).SprintFunc()
	cLiteral = color.New(color.FgHiCyan).SprintFunc()
	cLabel   = color.New(color.FgWhite).SprintFunc()
	cBranch  = color.New(color.FgHiBlack).SprintFunc()
)

func PrintTree(node Node) {
	//newNormalPrinter(os.Stdout, true).Print(node)
}

func PrintDebugTree(node Node) {
	newDebugPrinter(os.Stdout, true).Print(node)
}

func ExportTree(node Node, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	//newNormalPrinter(f, false).Print(node)
	return nil
}

func ExportDebugTree(node Node, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	newDebugPrinter(f, false).Print(node)
	return nil
}
