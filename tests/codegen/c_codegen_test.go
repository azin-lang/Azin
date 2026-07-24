package codegen_test

// This turns your Azin code into C code.
// C code then turns into bugs that were written 40 years ago.

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/codegen/c"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/parser"
	"github.com/azin-lang/Azin/internal/sema"
	"github.com/azin-lang/Azin/internal/source"
)

var update = flag.Bool("update", false, "update golden .c.expected files")

func parseProgram(t *testing.T, input string) *ast.Program {
	t.Helper()

	file := source.New("test.az", []byte(input))
	diag := diagnostics.New(file)

	tokens := lexer.New(file, diag).Tokenize()

	program, err := parser.Parse(
		string(file.Slice(0, file.Len())),
		tokens,
		diag,
	)
	if err != nil {
		t.Fatalf("parse failed:\n%v", err)
	}

	analyzer := sema.New(diag)
	if err := analyzer.Analyze(program); err != nil {
		t.Fatalf("semantic analysis failed:\n%v", err)
	}

	return program
}

func transpile(t *testing.T, input string) string {
	t.Helper()

	tx := c.New()

	out, _ := tx.Transpile(parseProgram(t, input))
	return normalize(out)
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	return s
}

func requireContains(t *testing.T, output string, values ...string) {
	t.Helper()

	for _, value := range values {
		if !strings.Contains(output, value) {
			t.Fatalf(
				"generated C does not contain %q\n\nGenerated output:\n%s",
				value,
				output,
			)
		}
	}
}

func requireNotContains(t *testing.T, output string, values ...string) {
	t.Helper()

	for _, value := range values {
		if strings.Contains(output, value) {
			t.Fatalf(
				"generated C unexpectedly contains %q\n\nGenerated output:\n%s",
				value,
				output,
			)
		}
	}
}

func TestVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		absent   []string
	}{
		{
			name: "immutable variable",
			input: `
fn main: int do
	var x: int = 42
	return x
end
`,
			contains: []string{
				"int main",
				"const int x",
				"42",
			},
		},
		{
			name: "mutable variable",
			input: `
fn identity(x: int): int do
    return x
end

fn main: int do
    var mut x: int = 42
    x = 99
    return x
end
`,
			contains: []string{
				"int x",
				"x = 99",
			},
			absent: []string{
				"const int x",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			out := transpile(t, tt.input)

			requireContains(t, out, tt.contains...)
			requireNotContains(t, out, tt.absent...)
		})
	}
}

func TestFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "unit return",
			input: `
fn foo do
    return
end

fn main: int do
    foo()
    return 0
end
`,
			contains: []string{
				"void foo",
				"return",
				"foo()",
			},
		},
		{
			name: "function call",
			input: `
fn greet: int do
	return 42
end

fn main: int do
	return greet()
end
`,
			contains: []string{
				"greet()",
			},
		},
		{
			name: "binary expression",
			input: `
fn add(a: int, b: int): int do
	return a + b
end

fn main do 
	return add(1, 2)
end
`,
			contains: []string{
				"a + b",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			requireContains(t, transpile(t, tt.input), tt.contains...)
		})
	}
}

func TestControlFlow(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "if",
			input: `
fn main: int do
	var mut x: int = 0
	if true then
		x = 1
	end
	return x
end
`,
			contains: []string{
				"if (true)",
			},
		},
		{
			name: "if else",
			input: `
fn main: int do
	var mut x: int = 0
	if true then
		x = 1
	else
		x = 2
	end
	return x
end
`,
			contains: []string{
				"if (true)",
				"else",
			},
		},
		{
			name: "loop",
			input: `
fn foo(x: int) do
	if x < 10 then
		return 0
	else
		return 1
	end
end

fn main: int do
	loop
		var mut x: int = 0
        if foo(10) == 0 then	
			return 0
		else	
			return 20
		end
	end
end
`,
			contains: []string{
				"for (;;)",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			requireContains(t, transpile(t, tt.input), tt.contains...)
		})
	}
}

func TestDefer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		absent   []string
	}{
		{
			name: "defer before return",
			input: `
fn side_effect do end

fn foo: int do
	defer side_effect()
	return 42
end

fn main: int do
	return foo()
end
`,
			contains: []string{
				"side_effect()",
				"return 42",
			},
			absent: []string{
				"/* defer */",
			},
		},
		{
			name: "multiple defers reverse order",
			input: `
fn cleanup do end
fn save do end

fn foo: int do
	defer cleanup()
	defer save()
	return 0
end

fn main: int do
	return foo()
end
`,
			contains: []string{
				"save()",
				"cleanup()",
				"return 0",
			},
		},
		{
			name: "defer at function end",
			input: `
fn done do end

fn foo do
	defer done()
end

fn main: int do
	foo()
	return 0
end
`,
			contains: []string{
				"done()",
			},
			absent: []string{
				"/* defer */",
			},
		},
		{
			name: "defer with importc builtin",
			input: `
fn main: int do
	var mut p: int
	defer free(p)
	return 0
end
`,
			contains: []string{
				"free(p)",
			},
			absent: []string{
				"/* defer */",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			out := transpile(t, tt.input)

			requireContains(t, out, tt.contains...)
			requireNotContains(t, out, tt.absent...)
		})
	}
}

func TestStructs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "struct declaration",
			input: `
struct Point is
	mut x: int
	y: int
end

fn main: int do
	var mut p: Point
	p.x = 4
	return 0
end
`,
			contains: []string{
				"typedef struct",
				"Point",
			},
		},
		{
			name: "field access",
			input: `
struct Point is
	x: int
end

fn main: int do
	var p: Point
	return p.x
end
`,
			contains: []string{
				"Point p",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			requireContains(t, transpile(t, tt.input), tt.contains...)
		})
	}
}

func TestBuiltinTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "primitive types",
			input: `
fn main: int do
	var a: int = 1
	var b: float = 2.0
	var c: bool = true
	var d: char = 'x'
	printf("%d %.1f %d %c\n", a, b, c, d);
	return 0
end
`,
			contains: []string{
				"int a",
				"float b",
				"bool c",
				"char d",
			},
		},
		{
			name: "string",
			input: `
fn foo(a: string) do
	return a
end

fn main: int do
	var s: string = foo("hello")
	return 0
end
`,
			contains: []string{
				"hello",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			out := transpile(t, tt.input)

			switch tt.name {
			case "string":
				if !strings.Contains(out, "char *") &&
					!strings.Contains(out, "const char") {
					t.Fatalf(
						"expected generated string type\n\n%s",
						out,
					)
				}

			default:
				requireContains(t, out, tt.contains...)
			}
		})
	}
}

func TestImports(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "importc",
			input: `
fn main do
	printf("123")
end
`,
			contains: []string{
				"#include <stdio.h>",
			},
		},
		{
			name: "stdbool",
			input: `
fn main do
	var a = true

	if a then
		return 0
	else
		return 1
	end
end
`,
			contains: []string{
				"#include <stdbool.h>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireContains(t, transpile(t, tt.input), tt.contains...)
		})
	}
}

func TestCodegenSnapshot(t *testing.T) {
	t.Helper()

	files, err := filepath.Glob(filepath.Join("testdata", "*.az"))
	if err != nil {
		t.Fatalf("discovering snapshots: %v", err)
	}

	if len(files) == 0 {
		t.Skip("no snapshot tests found")
	}

	for _, input := range files {

		name := strings.TrimSuffix(filepath.Base(input), ".az")

		t.Run(name, func(t *testing.T) {
			runSnapshotTest(t, input)
		})
	}
}

func runSnapshotTest(t *testing.T, sourceFile string) {
	t.Helper()

	source, err := os.ReadFile(sourceFile)
	if err != nil {
		t.Fatalf("reading %q: %v", sourceFile, err)
	}

	got := transpile(t, string(source))

	expectedFile := sourceFile + ".c.expected"

	if *update {
		writeGoldenFile(t, expectedFile, got)
	}

	want := readGoldenFile(t, expectedFile)

	if got != want {
		t.Fatalf("%s\n\n%s",
			diffMessage(sourceFile),
			formatMismatch(want, got),
		)
	}
}

func readGoldenFile(t *testing.T, filename string) string {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading golden file %q: %v", filename, err)
	}

	return normalize(string(data))
}

func writeGoldenFile(t *testing.T, filename string, output string) {
	t.Helper()

	output = normalize(output) + "\n"

	if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
		t.Fatalf("writing golden file %q: %v", filename, err)
	}
}

func diffMessage(file string) string {
	return fmt.Sprintf(
		"generated C does not match snapshot for %s",
		filepath.Base(file),
	)
}

func formatMismatch(expected, actual string) string {
	var b strings.Builder

	b.WriteString("===== EXPECTED =====\n")
	b.WriteString(expected)

	b.WriteString("\n\n")

	b.WriteString("===== GENERATED =====\n")
	b.WriteString(actual)

	return b.String()
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "abc\r\nxyz\r\n",
			want:  "abc\nxyz",
		},
		{
			input: "\n\nhello\n",
			want:  "hello",
		},
		{
			input: "abc",
			want:  "abc",
		},
	}

	for _, tt := range tests {
		got := normalize(tt.input)

		if got != tt.want {
			t.Fatalf(
				"normalize(%q) = %q, want %q",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}

func TestRequireContains(t *testing.T) {
	output := `
int main(void)
{
	return 0
}
`

	requireContains(
		t,
		output,
		"int main",
		"return 0",
	)
}

func TestRequireNotContains(t *testing.T) {
	output := `
const int x = 42
`

	requireNotContains(
		t,
		output,
		"float",
		"double",
	)
}

func BenchmarkCodegen(b *testing.B) {
	const src = `
struct Point is
    mut x: int
    mut y: int
end

fn add(a: int, b: int): int do
    return a + b
end

fn main: int do
    var mut p: Point

    p.x = add(1, 2)
    p.y = add(3, 4)

    if p.x < p.y then
        return p.y
    end

    loop
        break
    end

    return p.x
end
`

	file := source.New("bench.az", []byte(src))
	diag := diagnostics.New(file)

	tokens := lexer.New(file, diag).Tokenize()

	program, err := parser.Parse(
		string(file.Slice(0, file.Len())),
		tokens,
		diag,
	)
	if err != nil {
		b.Fatal(err)
	}

	analyzer := sema.New(diag)

	if err := analyzer.Analyze(program); err != nil {
		b.Fatal(err)
	}

	tx := c.New()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_, _ = tx.Transpile(program)
	}
}
