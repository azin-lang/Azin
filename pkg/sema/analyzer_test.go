//nolint:unparam
package sema_test

import (
	"strings"
	"testing"

	"github.com/azin-lang/Azin/pkg/ast"
	diagnostics2 "github.com/azin-lang/Azin/pkg/diagnostics"
	"github.com/azin-lang/Azin/pkg/lexer"
	"github.com/azin-lang/Azin/pkg/parser"
	sema "github.com/azin-lang/Azin/pkg/sema"
	"github.com/azin-lang/Azin/pkg/source"
)

func analyzeProgram(t *testing.T, input string) (*ast.Program, *diagnostics2.Engine) {
	t.Helper()
	file := source.New("test.az", []byte(input))
	diag := diagnostics2.New(file)
	tokens := lexer.New(file, diag).Tokenize()
	program, err := parser.Parse(string(file.Slice(0, file.Len())), tokens, diag)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	analyzer := sema.New(diag)
	if err := analyzer.Analyze(program); err != nil {
		return program, diag
	}
	return program, diag
}

func validProgram(t *testing.T, input string) {
	t.Helper()
	_, diag := analyzeProgram(t, input)
	if diag.HasErrors() {
		t.Errorf("unexpected errors: %v", diag.Err())
	}
}

func mustHaveError(t *testing.T, input string) *diagnostics2.Engine {
	t.Helper()
	_, diag := analyzeProgram(t, input)
	if !diag.HasErrors() {
		t.Error("expected sema error, got none")
	}
	return diag
}

func TestSemanticValidProgram(t *testing.T) {
	input := `fn foo: int do
    return 42;
end`
	validProgram(t, input)
}

func TestSemanticVarDeclWithInit(t *testing.T) {
	input := `fn main: int do
    var x: int = 42;
    return x;
end`
	validProgram(t, input)
}

func TestSemanticTypeMismatchInit(t *testing.T) {
	input := `fn main: int do
    var x: int = "hello";
    return 0;
end`
	_ = mustHaveError(t, input)
}

func TestSemanticTypeMismatchAssign(t *testing.T) {
	input := `fn main: int do
    var mut x: int = 42;
    x = "hello";
    return 0;
end`
	_ = mustHaveError(t, input)
}

func TestSemanticWhileLoopConditionTypeMismatch(t *testing.T) {
	input := `fn main: int do
		var mut x: int = 0;
		while x loop
			x = x + 1;
		end
		return x;
	end`
	_ = mustHaveError(t, input)
}

func TestSemanticImmutableAssign(t *testing.T) {
	input := `fn main: int do
    var x: int = 42;
    x = 99;
    return 0;
end`
	_ = mustHaveError(t, input)
}

func TestSemanticMutableAssign(t *testing.T) {
	input := `fn main: int do
    var mut x: int = 42;
    x = 99;
    return x;
end`
	validProgram(t, input)
}

func TestSemanticUnknownVariable(t *testing.T) {
	input := `fn main: int do
    x = 42;
    return 0;
end`
	_ = mustHaveError(t, input)
}

func TestSemanticUnknownFunction(t *testing.T) {
	input := `fn main: int do
    return foo();
end`
	_ = mustHaveError(t, input)
}

func TestSemanticWrongArgCount(t *testing.T) {
	input := `fn add(a: int, b: int): int do
    return a + b;
end
fn main: int do
    return add(1);
end`
	_ = mustHaveError(t, input)
}

func TestSemanticArgTypeMismatch(t *testing.T) {
	input := `fn greet(name: string): int do
    return 0;
end
fn main: int do
    return greet(42);
end`
	_ = mustHaveError(t, input)
}

func TestSemanticReturnTypeMismatch(t *testing.T) {
	input := `fn main: int do
    return "hello";
end`
	_ = mustHaveError(t, input)
}

func TestSemanticReturnInference(t *testing.T) {
	input := `fn foo do
    return 42;
end
fn main: int do
    return foo();
end`
	validProgram(t, input)
}

func TestSemanticStructAccess(t *testing.T) {
	input := `struct Point is
    x: int;
end
fn main: int do
    var p: Point;
    return p.x;
end`
	validProgram(t, input)
}

func TestSemanticStructFieldNotFound(t *testing.T) {
	input := `struct Point is
    x: int;
end
fn main: int do
    var p: Point;
    return p.z;
end`
	_ = mustHaveError(t, input)
}

func TestSemanticStructAssign(t *testing.T) {
	input := `struct Point is
    mut x: int;
end
fn main: int do
    var mut p: Point;
    p.x = 42;
    return 0;
end`
	validProgram(t, input)
}

func TestSemanticDuplicateFunction(t *testing.T) {
	input := `fn foo: int do
    return 1;
end
fn foo: int do
    return 2;
end`
	_ = mustHaveError(t, input)
}

func TestSemanticIfElse(t *testing.T) {
	input := `fn main: int do
    var mut x: int = 0;
    if true then
        x = 1;
    else
        x = 2;
    end
    return x;
end`
	validProgram(t, input)
}

func TestSemanticWhileLoop(t *testing.T) {
	input := `
	fn main: int do
		var mut x: int = 0;
		while x < 10 loop
			x = x + 1;
		end
		return x;
	end`
	validProgram(t, input)
}

func TestSemanticWhileLoopBreak(t *testing.T) {
	input := `fn main: int do
		var mut x: int = 0
		while x < 10 loop
			stop
		end

		return x
	end`
	validProgram(t, input)
}

func TestSemanticLoopBreak(t *testing.T) {
	input := `fn main: int do
    loop
        return 0;
    end
end`
	validProgram(t, input)
}

func TestSemanticGlobalVar(t *testing.T) {
	input := `var x: int = 42;
fn main: int do
    return x;
end`
	validProgram(t, input)
}

func TestSemanticCallFunction(t *testing.T) {
	input := `fn double(x: int): int do
    return x * 2;
end
fn main: int do
    return double(21);
end`
	validProgram(t, input)
}

func TestSemanticDefer(t *testing.T) {
	input := `fn main: int do
    defer printf("cleanup");
    return 0;
end`
	validProgram(t, input)
}

func mustHaveWarning(t *testing.T, input, msg string) {
	t.Helper()
	_, diag := analyzeProgram(t, input)
	for _, d := range diag.Diagnostics() {
		if d.Kind == diagnostics2.Warning && strings.Contains(d.Message, msg) {
			return
		}
	}
	t.Errorf("expected warning containing %q, got diagnostics: %v", msg, diag.Diagnostics())
}

func mustNotHaveWarning(t *testing.T, input string) {
	t.Helper()
	_, diag := analyzeProgram(t, input)
	for _, d := range diag.Diagnostics() {
		if d.Kind == diagnostics2.Warning {
			t.Errorf("unexpected warning: %s", d.Message)
		}
	}
}

func TestSemanticUnusedLocalVar(t *testing.T) {
	input := `fn main: int do
    var x: int = 42;
    return 0;
end`
	mustHaveWarning(t, input, "unused variable: x")
}

func TestSemanticUsedLocalVar(t *testing.T) {
	input := `fn main: int do
    var x: int = 42;
    return x;
end`
	mustNotHaveWarning(t, input)
}

func TestSemanticUnusedParameter(t *testing.T) {
	input := `fn foo(x: int): int do
    return 0;
end`
	mustHaveWarning(t, input, "unused variable: x")
}

func TestSemanticUsedParameter(t *testing.T) {
	input := `fn foo(x: int): int do
    return x;
end`
	mustNotHaveWarning(t, input)
}

func TestSemanticUnusedVarInBlock(t *testing.T) {
	input := `fn main: int do
    if true then
        var x: int = 42;
    end
    return 0;
end`
	mustHaveWarning(t, input, "unused variable: x")
}

func TestSemanticUsedVarInBlock(t *testing.T) {
	input := `fn main: int do
    if true then
        var x: int = 42;
        return x;
    end
    return 0;
end`
	mustNotHaveWarning(t, input)
}

func TestSemanticUnusedGlobalVar(t *testing.T) {
	input := `var x: int = 42;
fn main: int do
    return 0;
end`
	mustHaveWarning(t, input, "unused variable: x")
}

func TestSemanticUsedGlobalVar(t *testing.T) {
	input := `var x: int = 42;
fn main: int do
    return x;
end`
	mustNotHaveWarning(t, input)
}

func TestSemanticUnusedBlankIdentifier(t *testing.T) {
	input := `fn main: int do
    var _: int = 42;
    return 0;
end`
	mustNotHaveWarning(t, input)
}

func TestSemanticUnusedVarInWhileLoop(t *testing.T) {
	input := `fn main: int do
    var mut x: int = 0;
    while x < 10 loop
				var y: int = 42;
        x = x + 1;
    end
    return 0;
end`
	mustHaveWarning(t, input, "unused variable: y")
}

func TestSemanticUnusedVarInLoop(t *testing.T) {
	input := `fn main: int do
    loop
        var x: int = 42;
    end
end`
	mustHaveWarning(t, input, "unused variable: x")
}

func TestSemanticUsedVarInWhileLoop(t *testing.T) {
	input := `fn main: int do
		var mut x: int = 0;
		while x < 10 loop
			var mut y: int = 42;
			x = x + y + 1;
		end
		return x;
	end`
	mustNotHaveWarning(t, input)
}

func TestSemanticUsedVarInLoop(t *testing.T) {
	input := `fn main: int do
    loop
        var x: int = 42;
        return x;
    end
end`
	mustNotHaveWarning(t, input)
}
