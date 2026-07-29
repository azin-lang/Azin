package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/azin-lang/Azin/internal/driver"
	"github.com/azin-lang/Azin/internal/fs"
	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/diagnostics"
	"github.com/azin-lang/Azin/pkg/lexer"
	"github.com/azin-lang/Azin/pkg/parser"
	"github.com/azin-lang/Azin/pkg/source"
	"github.com/azin-lang/Azin/pkg/token"
)

const Version = "0.2.2"

type Config struct {
	Debug           bool
	PrintTokens     bool
	PrintAST        bool
	Optimization    string
	Output          string
	IgnoreExtension bool
	Version         bool
	EmitC           bool
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := parseFlags()

	if cfg.Version {
		fmt.Printf("Azin compiler v%s\n", Version)
		return nil
	}

	if flag.NArg() != 1 {
		flag.Usage()
		return fmt.Errorf("expected exactly one source file")
	}

	filename := flag.Arg(0)

	if cfg.Debug {
		printDebugInfo(cfg, filename)
	}

	data, err := fs.ReadSourceFile(filename, cfg.IgnoreExtension)
	if err != nil {
		return fmt.Errorf("reading source: %w", err)
	}
	file := source.New(filename, data)

	if cfg.PrintTokens || cfg.PrintAST {
		return handleFrontendDebug(cfg, file)
	}

	opts := driver.Options{
		Output:       cfg.Output,
		EmitC:        cfg.EmitC,
		Optimization: cfg.Optimization,
		Debug:        cfg.Debug,
	}

	if err := driver.Compile(file, cfg.Output, opts); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	if cfg.Debug {
		fmt.Printf("Compiled %.2f KiB\n", float64(file.Len())/1024)
	}

	return nil
}

func parseFlags() Config {
	var cfg Config

	flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug output")
	flag.BoolVar(&cfg.PrintTokens, "print-tokens", false, "Print lexer tokens and exit")
	flag.BoolVar(&cfg.PrintAST, "print-ast", false, "Print the parsed AST and exit")
	flag.StringVar(&cfg.Optimization, "O", "0", "Optimization level (0,1,2,3,s,z)")
	flag.StringVar(&cfg.Output, "o", "", "Output file path")
	flag.BoolVar(&cfg.IgnoreExtension, "ignore-extension", false, "Ignore source file extension")
	flag.BoolVar(&cfg.Version, "version", false, "Print compiler version")
	flag.BoolVar(&cfg.EmitC, "emit-c", false, "Generate C source instead of compiling")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [flags] <file.az>\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	return cfg
}

func handleFrontendDebug(cfg Config, file *source.File) error {
	diag := diagnostics.New(file)

	l := lexer.New(file, diag)
	tokens := l.Tokenize()
	if err := diag.Err(); err != nil {
		return err
	}

	if cfg.PrintTokens {
		for _, tok := range tokens {
			fmt.Println(formatToken(file, tok))
		}
		return nil
	}

	program, err := parser.Parse(string(file.Slice(0, file.Len())), tokens, diag)
	if err != nil {
		return err
	}
	if err := diag.Err(); err != nil {
		return err
	}

	if cfg.PrintAST {
		if cfg.Output != "" {
			if err := ast.ExportDebugTree(program, cfg.Output); err != nil {
				return fmt.Errorf("exporting AST: %w", err)
			}
		} else {
			ast.PrintDebugTree(program)
		}
	}

	return nil
}

func printDebugInfo(cfg Config, filename string) {
	fmt.Println("=== Compiler Config ===")
	fmt.Printf("File:         %q\n", filename)
	fmt.Printf("Output:       %q\n", cfg.Output)
	fmt.Printf("Optimization: %s\n", cfg.Optimization)
	fmt.Printf("Emit C:       %t\n", cfg.EmitC)
	fmt.Printf("Print Tokens: %t\n", cfg.PrintTokens)
	fmt.Printf("Print AST:    %t\n", cfg.PrintAST)
	fmt.Println("=======================")
}

func formatToken(f *source.File, tok token.Token) string {
	line, column := f.LineColumn(tok.Position.Offset)
	s := fmt.Sprintf("%-18s %4d:%4d [%d:%d]", tok.Kind, line, column, tok.Position.Offset, tok.Length)
	if tok.Kind.HasText() {
		s += fmt.Sprintf(" %q", f.Text(tok))
	}
	return s
}
