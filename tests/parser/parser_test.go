package parser_test

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/diagnostics"
	"github.com/azin-lang/Azin/internal/lexer"
	"github.com/azin-lang/Azin/internal/parser"
	"github.com/azin-lang/Azin/internal/source"
	"github.com/azin-lang/Azin/internal/token"
)

var update = flag.Bool("update", false, "update golden files")

func parseProgram(t *testing.T, input string) (*ast.Program, *diagnostics.Engine) {
	t.Helper()
	file := source.New("test.az", []byte(input))
	diag := diagnostics.New(file)
	tokens := lexer.New(file, diag).Tokenize()
	program, err := parser.Parse(string(file.Slice(0, file.Len())), tokens, diag)
	if err != nil {
		return program, diag
	}
	return program, diag
}

func captureAST(program *ast.Program) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	ast.PrintDebugTree(program)

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old
	return string(out)
}

func TestParserVarDecl(t *testing.T) {
	program, diag := parseProgram(t, "var x: int = 42;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	if _, ok := program.Statements[0].(*ast.VarStmt); !ok {
		t.Fatalf("expected VarStmt, got %T", program.Statements[0])
	}
}

func TestParserVarDeclNoType(t *testing.T) {
	program, diag := parseProgram(t, "var x: int = 42;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestParserVarMut(t *testing.T) {
	program, diag := parseProgram(t, "var mut x: int = 42;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	v := program.Statements[0].(*ast.VarStmt)
	if !v.Mutable {
		t.Error("expected mutable variable")
	}
}

func TestParserVarWithoutInitOrType(t *testing.T) {
	_, diag := parseProgram(t, "var x;\n")
	if !diag.HasErrors() {
		t.Error("expected error for var without type or init")
	}
}

func TestParserFnDecl(t *testing.T) {
	program, diag := parseProgram(t, "fn foo: int do\n    return 0;\nend\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	if _, ok := program.Statements[0].(*ast.FuncStmt); !ok {
		t.Fatalf("expected FuncStmt, got %T", program.Statements[0])
	}
}

func TestParserFnWithParams(t *testing.T) {
	program, diag := parseProgram(t, "fn add(a: int, b: int): int do\n    return a + b;\nend\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	fn := program.Statements[0].(*ast.FuncStmt)
	if len(fn.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Params))
	}
	if fn.Params[0].Name.Value != "a" || fn.Params[1].Name.Value != "b" {
		t.Errorf("param names: %s, %s", fn.Params[0].Name.Value, fn.Params[1].Name.Value)
	}
}

func TestParserFnNoReturnType(t *testing.T) {
	program, diag := parseProgram(t, "fn foo do\n    return;\nend\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	fn := program.Statements[0].(*ast.FuncStmt)
	if fn.ReturnType != nil {
		t.Errorf("expected nil return type, got %v", fn.ReturnType)
	}
}

func TestParserIf(t *testing.T) {
	program, diag := parseProgram(t, "if true then\n    return 1;\nend\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if _, ok := program.Statements[0].(*ast.IfStmt); !ok {
		t.Fatalf("expected IfStmt, got %T", program.Statements[0])
	}
}

func TestParserIfElse(t *testing.T) {
	input := "if true then\n    return 1;\nelse\n    return 2;\nend\n"
	program, diag := parseProgram(t, input)
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	ifstmt := program.Statements[0].(*ast.IfStmt)
	if len(ifstmt.Else) == 0 {
		t.Error("expected else branch to have statements")
	}
}

func TestParserLoop(t *testing.T) {
	program, diag := parseProgram(t, "loop\n    return 1;\nend\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if _, ok := program.Statements[0].(*ast.LoopStmt); !ok {
		t.Fatalf("expected LoopStmt, got %T", program.Statements[0])
	}
}

func TestParserStruct(t *testing.T) {
	input := "struct Point is\n    x: int;\n    y: int;\nend\n"
	program, diag := parseProgram(t, input)
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	s := program.Statements[0].(*ast.StructStmt)
	if s.Name.Value != "Point" {
		t.Errorf("struct name = %q, want %q", s.Name.Value, "Point")
	}
	if len(s.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(s.Fields))
	}
	if s.Fields[0].Name.Value != "x" || s.Fields[1].Name.Value != "y" {
		t.Errorf("field names: %s, %s", s.Fields[0].Name.Value, s.Fields[1].Name.Value)
	}
}

func TestParserStructWithMutable(t *testing.T) {
	input := "struct Point is\n    mut x: int;\n    y: int;\nend\n"
	program, diag := parseProgram(t, input)
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	s := program.Statements[0].(*ast.StructStmt)
	if !s.Fields[0].Mutable {
		t.Error("expected first field to be mutable")
	}
	if s.Fields[1].Mutable {
		t.Error("expected second field to be immutable")
	}
}

func TestParserImportC(t *testing.T) {
	program, diag := parseProgram(t, "importc \"stdio.h\"\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if _, ok := program.Statements[0].(*ast.ImportCStmt); !ok {
		t.Fatalf("expected ImportCStmt, got %T", program.Statements[0])
	}
}

func TestParserAssignment(t *testing.T) {
	program, diag := parseProgram(t, "x = 42;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if _, ok := program.Statements[0].(*ast.AssignmentStmt); !ok {
		t.Fatalf("expected AssignmentStmt, got %T", program.Statements[0])
	}
}

func TestParserMemberAssignment(t *testing.T) {
	program, diag := parseProgram(t, "p.x = 42;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	a := program.Statements[0].(*ast.AssignmentStmt)
	if _, ok := a.Left.(*ast.MemberExpr); !ok {
		t.Fatalf("expected MemberExpr on left, got %T", a.Left)
	}
}

func TestParserExpressionStmt(t *testing.T) {
	program, diag := parseProgram(t, "foo();\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if _, ok := program.Statements[0].(*ast.ExpressionStmt); !ok {
		t.Fatalf("expected ExpressionStmt, got %T", program.Statements[0])
	}
}

func TestParserExpressionPrecedence(t *testing.T) {
	// a + b * c should parse as (a + (b * c))
	program, diag := parseProgram(t, "x = a + b * c;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	assign := program.Statements[0].(*ast.AssignmentStmt)
	bin, ok := assign.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr for value, got %T", assign.Value)
	}
	if bin.Operator.Kind.String() != "plus" {
		t.Errorf("expected root operator +, got %s", bin.Operator.Kind.String())
	}
	// Right side should be (b * c)
	if _, ok := bin.Right.(*ast.BinaryExpr); !ok {
		t.Fatalf("expected BinaryExpr on right side")
	}
}

func TestParserBoolLiterals(t *testing.T) {
	program, diag := parseProgram(t, "var a: bool = true;\nvar b: bool = false;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
}

func TestParserStringLiteral(t *testing.T) {
	program, diag := parseProgram(t, "var s: string = \"hello\";\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if _, ok := program.Statements[0].(*ast.VarStmt); !ok {
		t.Fatalf("expected VarStmt, got %T", program.Statements[0])
	}
}

func TestParserCharLiteral(t *testing.T) {
	program, diag := parseProgram(t, "var c: char = 'x';\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	v := program.Statements[0].(*ast.VarStmt)
	if _, ok := v.Value.(*ast.CharacterLiteral); !ok {
		t.Fatalf("expected CharacterLiteral, got %T", v.Value)
	}
}

func TestParserHexLiteral(t *testing.T) {
	program, diag := parseProgram(t, "var x: int = 0xFF;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	v := program.Statements[0].(*ast.VarStmt)
	if _, ok := v.Value.(*ast.IntegerLiteral); !ok {
		t.Fatalf("expected IntegerLiteral, got %T", v.Value)
	}
}

func TestParserBinaryLiteral(t *testing.T) {
	program, diag := parseProgram(t, "var x: int = 0b1010;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	v := program.Statements[0].(*ast.VarStmt)
	if _, ok := v.Value.(*ast.IntegerLiteral); !ok {
		t.Fatalf("expected IntegerLiteral, got %T", v.Value)
	}
}

func TestParserFloatLiteral(t *testing.T) {
	program, diag := parseProgram(t, "var x: float = 3.14;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	v := program.Statements[0].(*ast.VarStmt)
	if _, ok := v.Value.(*ast.FloatLiteral); !ok {
		t.Fatalf("expected FloatLiteral, got %T", v.Value)
	}
}

func TestParserCallExpr(t *testing.T) {
	program, diag := parseProgram(t, "foo(1, 2, 3);\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	expr := program.Statements[0].(*ast.ExpressionStmt)
	call, ok := expr.Expression.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", expr.Expression)
	}
	if len(call.Args) != 3 {
		t.Errorf("expected 3 args, got %d", len(call.Args))
	}
}

func TestParserMemberExpr(t *testing.T) {
	program, diag := parseProgram(t, "point.x;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	expr := program.Statements[0].(*ast.ExpressionStmt)
	if _, ok := expr.Expression.(*ast.MemberExpr); !ok {
		t.Fatalf("expected MemberExpr, got %T", expr.Expression)
	}
}

func TestParserMissingEndError(t *testing.T) {
	_, diag := parseProgram(t, "fn foo do\n    return 0;\n")
	if !diag.HasErrors() {
		t.Error("expected error for missing 'end'")
	}
}

func TestParserMissingThenError(t *testing.T) {
	_, diag := parseProgram(t, "if true\n    return 1;\nend\n")
	if !diag.HasErrors() {
		t.Error("expected error for missing 'then'")
	}
}

func TestParserStructNoEndError(t *testing.T) {
	_, diag := parseProgram(t, "struct Point is\n    x: int;\n")
	if !diag.HasErrors() {
		t.Error("expected error for unclosed struct")
	}
}

func TestParserTrailingNewline(t *testing.T) {
	_ = flag.Lookup("test.v")
	// An empty line at the end of your code should not break everything.
	// But it did. Now it doesn't. That's progress.
	program, diag := parseProgram(t, "fn main: int do\n    return 0;\nend\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	for _, stmt := range program.Statements {
		if _, ok := stmt.(*ast.BadStmt); ok {
			t.Error("found unexpected BadStmt from trailing newline")
		}
	}
}

func TestParserMultipleNewlines(t *testing.T) {
	program, diag := parseProgram(t, "var a: int = 1;\n\n\nvar b: int = 2;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
}

func TestParserComparisonChaining(t *testing.T) {
	// a < b == c should parse as (a < b) == c
	program, diag := parseProgram(t, "x = a < b == c;\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	assign := program.Statements[0].(*ast.AssignmentStmt)
	bin, ok := assign.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", assign.Value)
	}
	// Root should be ==
	if bin.Operator.Kind != token.EqualEqual {
		t.Errorf("expected root ==, got %s", bin.Operator.Kind)
	}
}

func TestParserDefer(t *testing.T) {
	program, diag := parseProgram(t, "defer foo();\n")
	if diag.HasErrors() {
		t.Fatalf("unexpected errors: %v", diag.Err())
	}
	deferStmt, ok := program.Statements[0].(*ast.DeferStmt)
	if !ok {
		t.Fatalf("expected DeferStmt, got %T", program.Statements[0])
	}
	if _, ok := deferStmt.Call.(*ast.CallExpr); !ok {
		t.Fatalf("expected CallExpr in defer, got %T", deferStmt.Call)
	}
}

// Parser snapshot tests with golden files

func TestParserSnapshot(t *testing.T) {
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

			program, _ := parseProgram(t, string(input))
			got := captureAST(program)

			golden := azfile + ".golden"
			if *update {
				if err := os.WriteFile(golden, []byte(got), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("golden file %s: %v", golden, err)
			}

			wantStr := strings.ReplaceAll(string(want), "\r\n", "\n")
			if strings.TrimSpace(got) != strings.TrimSpace(wantStr) {
				t.Errorf("AST mismatch for %s\n--- got:\n%s\n--- want:\n%s", name, got, wantStr)
			}
		})
	}
}
