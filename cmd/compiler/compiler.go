package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/azin-lang/Azin/internal/compiler"
	"github.com/azin-lang/Azin/internal/fs"
)

const VERSION = "0.2.3-dev"

var (
	debug           = flag.Bool("debug", false, "")
	printTokens     = flag.Bool("tokens", false, "")
	output          = flag.String("o", "", "")
	ignoreExtension = flag.Bool("ignore-extension", false, "")
	version         = flag.Bool("version", false, "")
)

// handleHelp prints the help message and exits if the help flag is set
func handleHelp() {
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <filename>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
}

// handleDebug prints debug information if the debug flag is set
func handleDebug() {
	if *debug {
		fmt.Println("Debug:", *debug)
		fmt.Println("Print tokens:", *printTokens)
		fmt.Println("Output:", *output)
	}
}

// if the version flag is set, printVersion prints the version and exits
func printVersion() {
	fmt.Println("Azin compiler version", VERSION)
	os.Exit(0)
}

// checkVersionFlag prints the version if the version flag is set
func checkVersionFlag() {
	if *version {
		printVersion()
	}
}

// loadSource reads the source file from the filesystem
func loadSource(filename string) []byte {
	source, err := fs.ReadSourceFile(filename, *ignoreExtension)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return source
}

// runCompiler compiles the source code and prints any errors
func runCompiler(source []byte) {
	if err := compiler.Compile(source, *printTokens); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// printCompileStats prints the compiled size in KiB if debug mode is enabled
func printCompileStats(source []byte) {
	if *debug {
		fmt.Printf("Compiled %.2f KiB\n", float64(len(source))/1024)
	}
}

// main is the entry point of the compiler
func main() {
	flag.Parse()

	checkVersionFlag()
	handleHelp()

	filename := flag.Arg(0)
	handleDebug()

	source := loadSource(filename)
	runCompiler(source)
	printCompileStats(source)
}
