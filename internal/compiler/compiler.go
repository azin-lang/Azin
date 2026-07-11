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
	"github.com/azin-lang/Azin/internal/source"
)

func runMSVC(cl, sourcePath, exeName string) error {
	cmd := exec.Command(
		cl,
		"/nologo",
		"/O2",
		"/Fe:"+exeName,
		sourcePath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("[MSVC] Compiling...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cl.exe compilation failed: %w", err)
	}

	fmt.Printf("[Success] %s\n", exeName)
	return nil
}

func runClang(clang, sourcePath, exeName string) error {
	cmd := exec.Command(
		clang,
		"-std=c23",
		"-O2",
		sourcePath,
		"-o",
		exeName,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("[Clang] Compiling...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clang compilation failed: %w", err)
	}

	fmt.Printf("[Success] %s\n", exeName)
	return nil
}

func runGCC(gcc, sourcePath, exeName string) error {
	cmd := exec.Command(
		gcc,
		"-std=c23",
		"-O2",
		sourcePath,
		"-o",
		exeName,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("[GCC] Compiling...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gcc compilation failed: %w", err)
	}

	fmt.Printf("[Success] %s\n", exeName)
	return nil
}

func runCompiler(sourcePath, exeName string) error {
	switch runtime.GOOS {
	case "windows":
		if cl, err := exec.LookPath("cl.exe"); err == nil {
			return runMSVC(cl, sourcePath, exeName)
		}
		if gcc, err := exec.LookPath("gcc"); err == nil {
			return runGCC(gcc, sourcePath, exeName)
		}
		if clang, err := exec.LookPath("clang"); err == nil {
			return runClang(clang, sourcePath, exeName)
		}

	case "darwin":
		// Apple's default compiler is Clang.
		if clang, err := exec.LookPath("clang"); err == nil {
			return runClang(clang, sourcePath, exeName)
		}
		if gcc, err := exec.LookPath("gcc"); err == nil {
			return runGCC(gcc, sourcePath, exeName)
		}

	default: // Linux, BSD, etc
		if gcc, err := exec.LookPath("gcc"); err == nil {
			return runGCC(gcc, sourcePath, exeName)
		}
		if clang, err := exec.LookPath("clang"); err == nil {
			return runClang(clang, sourcePath, exeName)
		}
	}

	return fmt.Errorf("no supported C compiler found (searched for gcc, clang, and cl.exe)")
}

func Compile(file *source.File, outputPath string, emitC bool) error {
	program, err := parseSource(file)
	if err != nil {
		return err
	}

	cCode := transpileToC(program)
	if emitC {
		if outputPath == "" {
			outputPath = "output.c"
		}

		if filepath.Ext(outputPath) != ".c" {
			outputPath += ".c"
		}

		if err := os.WriteFile(outputPath, []byte(cCode), 0644); err != nil {
			return fmt.Errorf("failed to write C source: %w", err)
		}

		fmt.Printf("[Success] Generated C source: %s\n", outputPath)
		return nil
	}
	exeName := resolveExeName(outputPath)

	tmpPath, err := writeToTempFile(cCode)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	return runCompiler(tmpPath, exeName)
}

func parseSource(file *source.File) (*ast.Program, error) {
	diag := diagnostics.New(file)
	tokens := lexer.New(file, diag).Tokenize()
	if err := diag.Err(); err != nil {
		return nil, err
	}

	sourceString := string(file.Slice(0, file.Len()))
	p := parser.New(sourceString, tokens)
	program := p.ParseProgram()
	return program, diag.Err()
}

func transpileToC(program *ast.Program) string {
	tx := codegen.New()
	return tx.Transpile(program)
}

func resolveExeName(outputPath string) string {
	if outputPath == "" {
		if runtime.GOOS == "windows" {
			return "output.exe"
		}
		return "output"
	}

	if strings.HasSuffix(outputPath, ".c") {
		outputPath = strings.TrimSuffix(outputPath, ".c")
	}

	if runtime.GOOS == "windows" && filepath.Ext(outputPath) != ".exe" {
		outputPath += ".exe"
	}

	return outputPath
}

func writeToTempFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "azin_*.c")
	if err != nil {
		return "", fmt.Errorf("failed to create translation buffer: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to populate compile buffer: %w", err)
	}
	return tmpFile.Name(), nil
}
