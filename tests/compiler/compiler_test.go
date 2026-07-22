package compiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/azin-lang/Azin/internal/compiler"
	"github.com/azin-lang/Azin/internal/source"
)

func TestCompileEmitC(t *testing.T) {
	dir := t.TempDir()
	azPath := filepath.Join(dir, "test.az")
	if err := os.WriteFile(azPath, []byte("fn main: int do\n    return 0;\nend\n"), 0644); err != nil {
		t.Fatal(err)
	}

	file := source.New(azPath, []byte("fn main: int do\n    return 0;\nend\n"))

	outPath := filepath.Join(dir, "output.c")
	opts := compiler.Options{
		Output: outPath,
		EmitC:  true,
	}

	err := compiler.Compile(file, outPath, opts)
	if err != nil {
		t.Fatalf("Compile with EmitC failed: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if len(data) == 0 {
		t.Error("output file is empty")
	}
}

func TestCompileEmitCWithImport(t *testing.T) {
	dir := t.TempDir()
	azPath := filepath.Join(dir, "test.az")
	input := []byte("importc \"stdio.h\"\nfn main: int do\n    return 0;\nend\n")
	if err := os.WriteFile(azPath, input, 0644); err != nil {
		t.Fatal(err)
	}

	file := source.New(azPath, input)

	outPath := filepath.Join(dir, "output.c")
	opts := compiler.Options{
		Output: outPath,
		EmitC:  true,
	}

	err := compiler.Compile(file, outPath, opts)
	if err != nil {
		t.Fatalf("Compile with EmitC failed: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if len(data) == 0 {
		t.Error("output file is empty")
	}
}

func TestCompileDefaultOutputPath(t *testing.T) {
	file := source.New("test.az", []byte("fn main: int do\n    return 0;\nend\n"))

	opts := compiler.Options{
		EmitC: true,
	}

	err := compiler.Compile(file, "", opts)
	if err != nil {
		t.Fatalf("Compile with empty output path failed: %v", err)
	}

	// Clean up the default output file
	os.Remove("output.c")
}

func TestCompileEmptyProgram(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "empty_out.c")
	file := source.New("empty.az", []byte{})
	opts := compiler.Options{
		EmitC:  true,
		Output: outPath,
	}

	err := compiler.Compile(file, outPath, opts)
	if err != nil {
		t.Fatalf("expected empty program to succeed, got: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Error("output should be empty")
	}
}
