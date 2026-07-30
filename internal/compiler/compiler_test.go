package compiler_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/azin-lang/Azin/internal/compiler"
	"github.com/azin-lang/Azin/pkg/source"
)

func writeSource(t *testing.T, dir, name, src string) *source.File {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}
	return source.New(path, []byte(src))
}

func readOutput(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	return string(data)
}

func TestCompileEmitC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		source   string
		contains []string
	}{
		{
			name: "minimal program",
			source: `
fn main: int do
    return 0
end
`,
			contains: []string{"int main"},
		},
		{
			name: "importc",
			source: `
importc "stdio"
fn main: int do
    printf("hello\n")
    return 0
end
`,
			contains: []string{"#include <stdio.h>", "int main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			file := writeSource(t, dir, "test.az", tt.source)
			out := filepath.Join(dir, "output.c")
			opts := compiler.Options{Output: out, EmitC: true}
			if err := compiler.Compile(file, out, opts); err != nil {
				t.Fatalf("Compile() failed: %v", err)
			}
			got := readOutput(t, out)
			if got == "" {
				t.Fatal("generated C source is empty")
			}
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Fatalf("generated output does not contain %q\n\n%s", want, got)
				}
			}
		})
	}
}

func TestCompileDefaultOutputPath(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	file := writeSource(t, dir, "test.az", `
fn main: int do 
    return 0
end 
`)

	opts := compiler.Options{EmitC: true}
	if err := compiler.Compile(file, "", opts); err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}

	output := filepath.Join(dir, "output.c")
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}
}

func TestCompileEmptyProgram(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := writeSource(t, dir, "empty.az", "")
	out := filepath.Join(dir, "empty.c")
	opts := compiler.Options{Output: out, EmitC: true}
	if err := compiler.Compile(file, out, opts); err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}
	got := readOutput(t, out)
	if got != "" {
		t.Fatalf("expected empty output, got:\n%s", got)
	}
}
