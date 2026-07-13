package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/codegen"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/parser"
	"github.com/azin-lang/Azin/internal/semantic"
	"github.com/azin-lang/Azin/internal/source"
)

func runCompilerCommand(name string, args []string, label, output string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("[%s] Compiling...\n", label)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s compilation failed: %w", label, err)
	}

	fmt.Printf("[Success] %s\n", output)
	return nil
}

func runMSVC(cl, sourcePath, exeName string) error {
	return runCompilerCommand(
		cl,
		[]string{
			"/nologo",
			"/O2",
			"/Fe:" + exeName,
			sourcePath,
		},
		"MSVC",
		exeName,
	)
}

func runClang(clang, sourcePath, exeName string) error {
	return runCompilerCommand(
		clang,
		[]string{
			"-std=c23",
			"-O2",
			sourcePath,
			"-o",
			exeName,
		},
		"Clang",
		exeName,
	)
}

func runGCC(gcc, sourcePath, exeName string) error {
	return runCompilerCommand(
		gcc,
		[]string{
			"-std=c23",
			"-O2",
			sourcePath,
			"-o",
			exeName,
		},
		"GCC",
		exeName,
	)
}

type compiler struct {
	name string
	run  func(string, string, string) error
}

func runCompiler(sourcePath, exeName string) error {
	var compilers []compiler

	switch runtime.GOOS {
	case "windows":
		compilers = []compiler{
			{"cl.exe", runMSVC},
			{"gcc", runGCC},
			{"clang", runClang},
		}

	case "darwin":
		compilers = []compiler{
			{"clang", runClang},
			{"gcc", runGCC},
		}

	default:
		compilers = []compiler{
			{"gcc", runGCC},
			{"clang", runClang},
		}
	}

	for _, c := range compilers {
		if path, err := exec.LookPath(c.name); err == nil {
			return c.run(path, sourcePath, exeName)
		}
	}

	return fmt.Errorf("no supported C compiler found")
}

func writeCOutput(code, output string) error {
	if output == "" {
		output = "output.c"
	}
	if filepath.Ext(output) != ".c" {
		output += ".c"
	}

	if err := os.WriteFile(output, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write C source: %w", err)
	}

	fmt.Printf("[Success] Generated C source: %s\n", output)
	return nil
}

// Compile compiles the given source file to a C executable.
func Compile(file *source.File, outputPath string, emitC bool) error {
	diag := diagnostics.New(file)

	program, err := parseSource(file, diag)
	if err != nil {
		return err
	}

	analyzer := semantic.New(diag)

	if err := analyzer.Analyze(program); err != nil {
		return err
	}

	cCode := transpileToC(program)

	if emitC {
		return writeCOutput(cCode, outputPath)
	}

	exeName := resolveExeName(outputPath)

	tmpPath, err := writeToTempFile(cCode)
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(tmpPath)

	return runCompiler(tmpPath, exeName)
}

func parseSource(file *source.File, diag *diagnostics.Engine) (*ast.Program, error) {
	tokens := lexer.New(file, diag).Tokenize()
	if err := diag.Err(); err != nil {
		return nil, err
	}

	parser := parser.New(string(file.Slice(0, file.Len())), tokens, diag)

	program := parser.ParseProgram()

	if err := parser.Err(); err != nil {
		return nil, err
	}

	return program, diag.Err()
}

func transpileToC(program *ast.Program) string {
	tx := codegen.New()
	return tx.Transpile(program)
}

func resolveExeName(output string) string {
	if output == "" {
		if runtime.GOOS == "windows" {
			return "output.exe"
		}
		return "output"
	}

	output = strings.TrimSuffix(output, ".c")

	if runtime.GOOS == "windows" && filepath.Ext(output) != ".exe" {
		output += ".exe"
	}

	return output
}

func writeToTempFile(content string) (string, error) {
	f, err := os.CreateTemp("", "azin_*.c")
	if err != nil {
		return "", fmt.Errorf("failed to create temp source: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write temp source: %w", err)
	}

	return f.Name(), nil
}
