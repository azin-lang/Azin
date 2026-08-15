package parser_test

import (
	"reflect"
	"testing"

	"github.com/azin-lang/Azin/internal/azin/diagnostics"
	"github.com/azin-lang/Azin/internal/azin/syntax"
	"github.com/azin-lang/Azin/internal/azin/syntax/lexer"
	"github.com/azin-lang/Azin/internal/azin/syntax/parser"
	"github.com/azin-lang/Azin/internal/azin/syntax/red"
	"github.com/azin-lang/Azin/internal/azin/text"
)

func parseSource(sourceText string) (*red.Node, *diagnostics.Collector) {
	file := text.NewSourceText("test.az", 0, []byte(sourceText))
	diags := diagnostics.NewCollector()

	lex := lexer.New(file, diags)
	p := parser.New(lex, diags)

	return red.NewRoot(p.ParseExpression()), diags
}

func parse(t *testing.T, sourceText string) (*red.Node, *diagnostics.Collector) {
	t.Helper()
	return parseSource(sourceText)
}

func assertNoDiagnostics(t *testing.T, diags *diagnostics.Collector) {
	t.Helper()
	if diags.HasErrors() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
}

func assertHasDiagnostics(t *testing.T, diags *diagnostics.Collector) {
	t.Helper()
	if !diags.HasErrors() {
		t.Fatal("expected diagnostic errors, but got none")
	}
}

func TestParseLiterals(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedKind syntax.Kind
		expectedText string
	}{
		{
			name:         "Integer Literal",
			input:        "123",
			expectedKind: syntax.IntegerLiteralToken,
			expectedText: "123",
		},
		{
			name:         "Float Literal",
			input:        "3.14159",
			expectedKind: syntax.FloatLiteralToken,
			expectedText: "3.14159",
		},
		{
			name:         "String Literal",
			input:        `"hello world"`,
			expectedKind: syntax.StringLiteralToken,
			expectedText: `"hello world"`,
		},
		{
			name:         "True Keyword Literal",
			input:        "true",
			expectedKind: syntax.KeywordTrue,
			expectedText: "true",
		},
		{
			name:         "False Keyword Literal",
			input:        "false",
			expectedKind: syntax.KeywordFalse,
			expectedText: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, diags := parse(t, tt.input)
			assertNoDiagnostics(t, diags)

			lit := red.AsLiteralExpression(root)
			if lit == nil {
				t.Fatalf("expected literal expression for input %q", tt.input)
			}

			tok := lit.Token()
			if tok.IsZero() {
				t.Fatal("expected non-zero token")
			}

			if tok.Kind() != tt.expectedKind {
				t.Errorf("expected kind %v, got %v", tt.expectedKind, tok.Kind())
			}

			if tok.Text() != tt.expectedText {
				t.Errorf("expected text %q, got %q", tt.expectedText, tok.Text())
			}
		})
	}
}

func TestParseNameExpression(t *testing.T) {
	root, diags := parse(t, "variableName")
	assertNoDiagnostics(t, diags)

	name := red.AsNameExpression(root)
	if name == nil {
		t.Fatal("expected name expression")
	}

	tok := name.Identifier()
	if tok.IsZero() {
		t.Fatal("expected identifier token")
	}

	if tok.Kind() != syntax.IdentifierToken {
		t.Fatalf("expected IdentifierToken, got %v", tok.Kind())
	}

	if tok.Text() != "variableName" {
		t.Fatalf("expected 'variableName', got %q", tok.Text())
	}
}

func TestParseUnaryExpression(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedOperator syntax.Kind
	}{
		{
			name:             "Unary Minus",
			input:            "-10",
			expectedOperator: syntax.MinusToken,
		},
		{
			name:             "Unary Plus",
			input:            "+42",
			expectedOperator: syntax.PlusToken,
		},
		{
			name:             "Logical NOT",
			input:            "!isReady",
			expectedOperator: syntax.BangToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, diags := parse(t, tt.input)
			assertNoDiagnostics(t, diags)

			unary := red.AsUnaryExpression(root)
			if unary == nil {
				t.Fatalf("expected unary expression for %q", tt.input)
			}

			if unary.Operator().Kind() != tt.expectedOperator {
				t.Fatalf("expected operator %v, got %v", tt.expectedOperator, unary.Operator().Kind())
			}

			if unary.Operand() == nil {
				t.Fatal("expected non-nil operand")
			}
		})
	}
}

func TestChainedUnaryOperators(t *testing.T) {
	root, diags := parse(t, "!--x")
	assertNoDiagnostics(t, diags)

	notUnary := red.AsUnaryExpression(root)
	if notUnary == nil || notUnary.Operator().Kind() != syntax.BangToken {
		t.Fatal("expected outer unary '!'")
	}

	minus1 := red.AsUnaryExpression(notUnary.Operand())
	if minus1 == nil || minus1.Operator().Kind() != syntax.MinusToken {
		t.Fatal("expected middle unary '-'")
	}

	minus2 := red.AsUnaryExpression(minus1.Operand())
	if minus2 == nil || minus2.Operator().Kind() != syntax.MinusToken {
		t.Fatal("expected inner unary '-'")
	}

	if red.AsNameExpression(minus2.Operand()) == nil {
		t.Fatal("expected inner operand to be name expression 'x'")
	}
}

func TestBinaryOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		topOperator   syntax.Kind
		leftOperator  syntax.Kind
		rightOperator syntax.Kind
	}{
		{
			name:          "Multiplication over Addition",
			input:         "1 + 2 * 3",
			topOperator:   syntax.PlusToken,
			rightOperator: syntax.StarToken,
		},
		{
			name:          "Division over Subtraction",
			input:         "10 - 8 / 2",
			topOperator:   syntax.MinusToken,
			rightOperator: syntax.SlashToken,
		},
		{
			name:          "Addition over Relational Comparison",
			input:         "a + b < c + d",
			topOperator:   syntax.LessToken,
			leftOperator:  syntax.PlusToken,
			rightOperator: syntax.PlusToken,
		},
		{
			name:          "Relational over Equality",
			input:         "a < b == c > d",
			topOperator:   syntax.EqualsEqualsToken,
			leftOperator:  syntax.LessToken,
			rightOperator: syntax.GreaterToken,
		},
		{
			name:          "Equality over Logical AND",
			input:         "x == 1 && y != 2",
			topOperator:   syntax.AmpersandAmpersandToken,
			leftOperator:  syntax.EqualsEqualsToken,
			rightOperator: syntax.BangEqualsToken,
		},
		{
			name:          "Logical AND over Logical OR",
			input:         "a || b && c",
			topOperator:   syntax.PipePipeToken,
			rightOperator: syntax.AmpersandAmpersandToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, diags := parse(t, tt.input)
			assertNoDiagnostics(t, diags)

			top := red.AsBinaryExpression(root)
			if top == nil {
				t.Fatalf("expected top binary expression for %q", tt.input)
			}

			if top.Operator().Kind() != tt.topOperator {
				t.Errorf("expected top operator %v, got %v", tt.topOperator, top.Operator().Kind())
			}

			if tt.leftOperator != 0 {
				left := red.AsBinaryExpression(top.Left())
				if left == nil || left.Operator().Kind() != tt.leftOperator {
					t.Errorf("expected left operator %v", tt.leftOperator)
				}
			}

			if tt.rightOperator != 0 {
				right := red.AsBinaryExpression(top.Right())
				if right == nil || right.Operator().Kind() != tt.rightOperator {
					t.Errorf("expected right operator %v", tt.rightOperator)
				}
			}
		})
	}
}

func TestLeftAssociativity(t *testing.T) {
	root, diags := parse(t, "10 - 5 - 2")
	assertNoDiagnostics(t, diags)

	outer := red.AsBinaryExpression(root)
	if outer == nil || outer.Operator().Kind() != syntax.MinusToken {
		t.Fatal("expected outer binary '-'")
	}

	inner := red.AsBinaryExpression(outer.Left())
	if inner == nil || inner.Operator().Kind() != syntax.MinusToken {
		t.Fatal("expected left-associated inner binary '-'")
	}

	if red.AsLiteralExpression(outer.Right()) == nil {
		t.Fatal("expected right operand to be literal '2'")
	}
}

func TestParenthesesOverridePrecedence(t *testing.T) {
	root, diags := parse(t, "(1 + 2) * 3")
	assertNoDiagnostics(t, diags)

	mul := red.AsBinaryExpression(root)
	if mul == nil || mul.Operator().Kind() != syntax.StarToken {
		t.Fatal("expected top level '*'")
	}

	paren := red.AsParenthesizedExpression(mul.Left())
	if paren == nil {
		t.Fatal("expected parenthesized expression on the left")
	}

	add := red.AsBinaryExpression(paren.Expression())
	if add == nil || add.Operator().Kind() != syntax.PlusToken {
		t.Fatal("expected '+' inside parentheses on the left")
	}
}

func TestDeeplyNestedParentheses(t *testing.T) {
	root, diags := parse(t, "(((42)))")
	assertNoDiagnostics(t, diags)

	curr := root
	for {
		if paren := red.AsParenthesizedExpression(curr); paren != nil {
			curr = paren.Expression()
		} else {
			break
		}
	}

	lit := red.AsLiteralExpression(curr)
	if lit == nil {
		t.Fatal("expected nested parentheses to unwrap to literal '42'")
	}
}

func TestParseCallExpressionNoArgs(t *testing.T) {
	root, diags := parse(t, "foo()")
	assertNoDiagnostics(t, diags)

	call := red.AsCallExpression(root)
	if call == nil {
		t.Fatal("expected call expression")
	}

	if red.AsNameExpression(call.Expression()) == nil {
		t.Fatal("expected function name expression")
	}

	if call.OpenParen().IsZero() || call.OpenParen().Kind() != syntax.OpenParenToken {
		t.Fatal("expected '(' token")
	}

	if call.CloseParen().IsZero() || call.CloseParen().Kind() != syntax.CloseParenToken {
		t.Fatal("expected ')' token")
	}
}

func TestParseCallExpressionWithArgs(t *testing.T) {
	root, diags := parse(t, "add(x, 10 + y)")
	assertNoDiagnostics(t, diags)

	call := red.AsCallExpression(root)
	if call == nil {
		t.Fatal("expected call expression")
	}

	args := call.Arguments()
	if args == nil || args.Count() < 2 {
		t.Fatal("expected at least 2 arguments")
	}

	arg0 := red.AsNameExpression(args.Item(0))
	if arg0 == nil {
		t.Fatal("expected first argument to be name expression 'x'")
	}

	arg1 := red.AsBinaryExpression(args.Item(1))
	if arg1 == nil {
		t.Fatal("expected second argument to be binary expression '10 + y'")
	}
}

func TestCallExpressionTrailingComma(t *testing.T) {
	root, diags := parse(t, "foo(x, y,)")
	assertNoDiagnostics(t, diags)

	call := red.AsCallExpression(root)
	if call == nil {
		t.Fatal("expected call expression")
	}

	args := call.Arguments()
	if args == nil || args.Count() < 2 {
		t.Fatal("expected at least 2 arguments")
	}

	if red.AsNameExpression(args.Item(0)) == nil || red.AsNameExpression(args.Item(1)) == nil {
		t.Fatal("expected arguments x and y to be parsed successfully")
	}
}

func TestCallExpressionMissingCommaRecovery(t *testing.T) {
	root, diags := parse(t, "foo(x y)")
	assertHasDiagnostics(t, diags)

	call := red.AsCallExpression(root)
	if call == nil {
		t.Fatal("expected call expression to be recovered")
	}

	args := call.Arguments()
	if args == nil || args.Count() < 2 {
		t.Fatal("expected at least 2 arguments in syntax tree")
	}

	if red.AsNameExpression(args.Item(0)) == nil || red.AsNameExpression(args.Item(1)) == nil {
		t.Fatal("expected both arguments 'x' and 'y' to be preserved in the syntax tree")
	}
}

func TestHigherOrderFunctionCall(t *testing.T) {
	root, diags := parse(t, "getHandler()(event)")
	assertNoDiagnostics(t, diags)

	outerCall := red.AsCallExpression(root)
	if outerCall == nil {
		t.Fatal("expected outer call expression")
	}

	innerCall := red.AsCallExpression(outerCall.Expression())
	if innerCall == nil {
		t.Fatal("expected inner call expression")
	}
}

func TestParentPointersAndPositions(t *testing.T) {
	source := "   10   +   20  "
	root, diags := parse(t, source)
	assertNoDiagnostics(t, diags)

	bin := red.AsBinaryExpression(root)
	if bin == nil {
		t.Fatal("expected binary expression")
	}

	if bin.Left().Parent() != bin.Node {
		t.Error("incorrect left child parent pointer")
	}
	if bin.Operator().Parent() != bin.Node {
		t.Error("incorrect operator child parent pointer")
	}
	if bin.Right().Parent() != bin.Node {
		t.Error("incorrect right child parent pointer")
	}

	if bin.Position() != 0 {
		t.Errorf("expected root position 0, got %d", bin.Position())
	}

	if bin.Left().Position() != 0 {
		t.Errorf("expected left operand position 0 (including leading whitespace), got %d", bin.Left().Position())
	}

	length := len(source)
	if root.FullWidth() != length {
		t.Errorf("expected full width %d, got %d", length, root.FullWidth())
	}
}

func TestSafeRedNodeCastsOnNilOrMismatch(t *testing.T) {
	var nilNode *red.Node

	if red.AsLiteralExpression(nilNode) != nil {
		t.Error("AsLiteralExpression on nil node should return nil")
	}
	if red.AsBinaryExpression(nilNode) != nil {
		t.Error("AsBinaryExpression on nil node should return nil")
	}
	if red.AsUnaryExpression(nilNode) != nil {
		t.Error("AsUnaryExpression on nil node should return nil")
	}
	if red.AsCallExpression(nilNode) != nil {
		t.Error("AsCallExpression on nil node should return nil")
	}
	if red.AsNameExpression(nilNode) != nil {
		t.Error("AsNameExpression on nil node should return nil")
	}

	root, _ := parse(t, "123")
	if red.AsBinaryExpression(root) != nil {
		t.Error("AsBinaryExpression on Literal node should return nil")
	}
}

func TestMissingClosingParenthesisDiagnostic(t *testing.T) {
	root, diags := parse(t, "(1 + 2")
	assertHasDiagnostics(t, diags)

	if root == nil {
		t.Fatal("parser should produce a resilient syntax tree despite missing ')'")
	}
}

func TestCallMissingClosingParenthesisDiagnostic(t *testing.T) {
	root, diags := parse(t, "foo(x, y")
	assertHasDiagnostics(t, diags)

	call := red.AsCallExpression(root)
	if call == nil {
		t.Fatal("expected parser to recover and create CallExpression node")
	}
}

func BenchmarkParseExpression(b *testing.B) {
	source := "foo(123, 456, 789) + bar(42) * 2"

	b.ResetTimer()

	for b.Loop() {
		parseSource(source)
	}
}

func Test_parseSource(t *testing.T) {
	type args struct {
		sourceText string
	}
	tests := []struct {
		name  string
		args  args
		want  *red.Node
		want1 *diagnostics.Collector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseSource(tt.args.sourceText)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSource() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseSource() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	type args struct {
		t          *testing.T
		sourceText string
	}
	tests := []struct {
		name  string
		args  args
		want  *red.Node
		want1 *diagnostics.Collector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parse(tt.args.t, tt.args.sourceText)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_assertNoDiagnostics(t *testing.T) {
	type args struct {
		t     *testing.T
		diags *diagnostics.Collector
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertNoDiagnostics(tt.args.t, tt.args.diags)
		})
	}
}

func Test_assertHasDiagnostics(t *testing.T) {
	type args struct {
		t     *testing.T
		diags *diagnostics.Collector
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertHasDiagnostics(tt.args.t, tt.args.diags)
		})
	}
}

func TestBenchmarkParseExpression(t *testing.T) {
	type args struct {
		b *testing.B
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BenchmarkParseExpression(tt.args.b)
		})
	}
}
