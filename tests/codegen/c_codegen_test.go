package codegen_test

// This turns your Azin code into C code.
// C code then turns into bugs that were written 40 years ago.
import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/azin-lang/Azin/internal/codegen/c"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/parser"
	"github.com/azin-lang/Azin/internal/sema"
	"github.com/azin-lang/Azin/internal/source"
)

var update = flag.Bool("update", false, "update golden .c.expected files")

func transpile(t *testing.T, input string) string {
	t.Helper()
	file := source.New("test.az", []byte(input))
	diag := diagnostics.New(file)
	tokens := lexer.New(file, diag).Tokenize()
	program, err := parser.Parse(string(file.Slice(0, file.Len())), tokens, diag)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	analyzer := sema.New(diag)
	if err := analyzer.Analyze(program); err != nil {
		t.Fatalf("sema error: %v", err)
	}

	tx := c.New()
	return tx.Transpile(program)
}

func TestCodegenVarDecl(t *testing.T) {
	c := transpile(t, `fn main: int do
    var x: int = 42;
    return x;
end
`)
	if !strings.Contains(c, "int main") {
		t.Error("output missing function definition")
	}
	if !strings.Contains(c, "const int x") {
		t.Error("expected const int x for immutable var")
	}
	if !strings.Contains(c, "42") {
		t.Error("output missing initializer")
	}
}

func TestCodegenMutableVar(t *testing.T) {
	c := transpile(t, `fn main: int do
    var mut x: int = 42;
    x = 99;
    return x;
end
`)
	if !strings.Contains(c, "int x") {
		t.Error("expected non-const int x for mutable var")
	}
	if !strings.Contains(c, "x = 99") {
		t.Error("output missing assignment")
	}
}

func TestCodegenUnitReturn(t *testing.T) {
	c := transpile(t, `fn foo do
    return;
end
`)
	if !strings.Contains(c, "void foo") {
		t.Error("expected void return type for unit function")
	}
	if !strings.Contains(c, "return;") {
		t.Error("output missing return statement")
	}
}

func TestCodegenIf(t *testing.T) {
	c := transpile(t, `fn main: int do
    var mut x: int = 0;
    if true then
        x = 1;
    end
    return x;
end
`)
	if !strings.Contains(c, "if (true)") {
		t.Error("output missing if condition")
	}
}

func TestCodegenIfElse(t *testing.T) {
	c := transpile(t, `fn main: int do
    var mut x: int = 0;
    if true then
        x = 1;
    else
        x = 2;
    end
    return x;
end
`)
	if !strings.Contains(c, "else") {
		t.Error("output missing else branch")
	}
}

func TestCodegenLoop(t *testing.T) {
	c := transpile(t, `fn main: int do
    loop
        return 0;
    end
end
`)
	if !strings.Contains(c, "for (;;)") {
		t.Error("expected infinite for loop")
	}
}

func TestCodegenStruct(t *testing.T) {
	c := transpile(t, `struct Point is
    x: int;
    y: int;
end
fn main: int do
    return 0;
end
`)
	if !strings.Contains(c, "typedef struct") {
		t.Error("output missing typedef struct")
	}
	if !strings.Contains(c, "Point") {
		t.Error("output missing struct name")
	}
}

func TestCodegenFunctionCall(t *testing.T) {
	c := transpile(t, `fn greet: int do
    return 42;
end
fn main: int do
    return greet();
end
`)
	if !strings.Contains(c, "greet()") {
		t.Error("output missing function call")
	}
}

func TestCodegenStructFieldAccess(t *testing.T) {
	c := transpile(t, `struct Point is
    x: int;
end
fn main: int do
    var p: Point;
    return 0;
end
`)
	if !strings.Contains(c, "Point p") {
		t.Error("output missing struct declaration")
	}
}

func TestCodegenBinaryExpr(t *testing.T) {
	c := transpile(t, `fn add(a: int, b: int): int do
    return a + b;
end
`)
	if !strings.Contains(c, "a + b") {
		t.Error("output missing binary expression")
	}
}

func TestCodegenTypes(t *testing.T) {
	c := transpile(t, `fn main: int do
    var a: int = 1;
    var b: float = 2.0;
    var c: bool = true;
    var d: char = 'x';
    return 0;
end
`)
	if !strings.Contains(c, "int a") {
		t.Error("expected int type")
	}
	if !strings.Contains(c, "float b") {
		t.Error("expected float type")
	}
	if !strings.Contains(c, "bool c") {
		t.Error("expected bool type")
	}
	if !strings.Contains(c, "char d") {
		t.Error("expected char type")
	}
}

func TestCodegenStringType(t *testing.T) {
	c := transpile(t, `fn main: int do
    var s: string = "hello";
    return 0;
end
`)
	if !strings.Contains(c, "char *s") && !strings.Contains(c, "const char") {
		t.Errorf("expected const char* for string type, got:\n%s", c)
	}
}

func TestCodegenImports(t *testing.T) {
	c := transpile(t, `importc "stdio.h"
fn main: int do
    return 0;
end
`)
	if !strings.Contains(c, "#include <stdio.h>") {
		t.Error("output missing include directive")
	}
}

func TestCodegenStdboolIncluded(t *testing.T) {
	c := transpile(t, `fn main: int do
    return 0;
end
`)
	if !strings.Contains(c, "#include <stdbool.h>") {
		t.Error("output missing stdbool.h include")
	}
}

// Snapshot tests with golden files

func TestCodegenSnapshot(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "*.az"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Skip("no testdata files found")
	}

	for _, azfile := range files {
		name := filepath.Base(azfile)
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(azfile)
			if err != nil {
				t.Fatal(err)
			}

			got := transpile(t, string(input))
			got = strings.TrimSpace(got)

			golden := azfile + ".c.expected"
			if *update {
				if err := os.WriteFile(golden, []byte(got+"\n"), 0644); err != nil {
					t.Fatal(err)
				}
			}

			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("golden file %s: %v", golden, err)
			}

			wantStr := strings.ReplaceAll(string(want), "\r\n", "\n")
			if got != strings.TrimSpace(wantStr) {
				t.Errorf("C output mismatch for %s\n--- got:\n%s\n--- want:\n%s", name, got, wantStr)
			}
		})
	}
}
