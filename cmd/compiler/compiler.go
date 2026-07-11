package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/azin-lang/Azin/internal/compiler"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/fs"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

const Version = "0.2.3-dev"

var (
	debug           = flag.Bool("debug", false, "Enable debug output")
	printTokens     = flag.Bool("tokens", false, "Print lexer tokens")
	output          = flag.String("o", "", "Output file")
	ignoreExtension = flag.Bool("ignore-extension", false, "Ignore source file extension")
	version         = flag.Bool("version", false, "Print compiler version")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <file>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("Azin compiler %s\n", Version)
		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *debug {
		printDebug()
	}

	filename := flag.Arg(0)
	data := mustReadSource(filename)
	file := source.New(filename, data)

	diag := diagnostics.New(file)
	tokens := lexer.New(file, diag).Tokenize()
	if *printTokens {
		for _, tok := range tokens {
			fmt.Println(formatToken(file, tok))
		}
	}

	if err := compiler.Compile(file); err != nil {
		fatal(err)
	}

	if err := diag.Err(); err != nil {
		fatal(err)
	}

	if *debug {
		fmt.Printf("Compiled %.2f KiB\n", float64(file.Len())/1024)
	}
}

func printDebug() {
	fmt.Printf("Debug: %t\n", *debug)
	fmt.Printf("Print tokens: %t\n", *printTokens)
	fmt.Printf("Output: %q\n", *output)
}

func mustReadSource(filename string) []byte {
	source, err := fs.ReadSourceFile(filename, *ignoreExtension)
	if err != nil {
		fatal(err)
	}
	return source
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func formatToken(f *source.File, tok token.Token) string {
	line, column := f.LineColumn(tok.Position.Offset)

	s := fmt.Sprintf(
		"%-18s %4d:%4d [%d:%d]",
		tok.Kind,
		line,
		column,
		tok.Position.Offset,
		tok.Length,
	)

	if tok.Kind.HasText() {
		s += fmt.Sprintf(" %q", f.Text(tok))
	}

	return s
}
