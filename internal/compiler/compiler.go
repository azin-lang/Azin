package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/codegen/c"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/optimizer"
	"github.com/azin-lang/Azin/internal/parser"
	"github.com/azin-lang/Azin/internal/sema"
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

func msvcOptimization(opt string) string {
	switch opt {
	case "1", "s", "z":
		return "/O1"
	case "2":
		return "/O2"
	case "3":
		return "/Ox"
	default:
		return "/Od"
	}
}

func runMSVC(cl, sourcePath, exeName string, opts Options) error {
	opt := msvcOptimization(opts.Optimization)

	return runCompilerCommand(
		cl,
		[]string{
			"/nologo",
			opt,
			"/Fe:" + exeName,
			sourcePath,
		},
		"MSVC",
		exeName,
	)
}

func runClang(clang, sourcePath, exeName string, opts Options) error {
	return runCompilerCommand(
		clang,
		[]string{
			"-std=c23",
			"-O" + opts.Optimization,
			sourcePath,
			"-o",
			exeName,
		},
		"Clang",
		exeName,
	)
}

func runGCC(gcc, sourcePath, exeName string, opts Options) error {
	return runCompilerCommand(
		gcc,
		[]string{
			"-std=c23",
			"-O" + opts.Optimization,
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
	run  func(string, string, string, Options) error
}

func runCompiler(sourcePath, exeName string, opts Options) error {
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
			return c.run(path, sourcePath, exeName, opts)
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
func Compile(file *source.File, outputPath string, opts Options) error {
	diag := diagnostics.New(file)

	program, err := parseSource(file, diag)
	if err != nil {
		return err
	}

	analyzer := sema.New(diag)

	if err := analyzer.Analyze(program); err != nil {
		return err
	}

	optimizer.Optimize(program)
	cCode, err := transpileToC(program)
	if err != nil {
		return err
	}

	if opts.EmitC {
		return writeCOutput(cCode, outputPath)
	}

	exeName := resolveExeName(outputPath)

	tmpPath, err := writeToTempFile(cCode)
	if err != nil {
		return err
	}
	defer func(name string) {
		if err := os.Remove(name); err != nil {
			fmt.Printf("warning: failed to remove temp file %s: %v\n", name, err)
		}
	}(tmpPath)

	return runCompiler(tmpPath, exeName, opts)
}

func parseSource(file *source.File, diag *diagnostics.Engine) (*ast.Program, error) {
	tokens := lexer.New(file, diag).Tokenize()
	if err := diag.Err(); err != nil {
		return nil, err
	}

	sourceParser := parser.New(string(file.Slice(0, file.Len())), tokens, diag)

	program := sourceParser.ParseProgram()

	if err := sourceParser.Err(); err != nil {
		return nil, err
	}

	return program, diag.Err()
}

func transpileToC(program *ast.Program) (string, error) {
	tx := c.New()
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

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning: failed to close temp source: %v\n", err)
		}
	}()

	if _, err := f.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write temp source: %w", err)
	}

	return f.Name(), nil
}
