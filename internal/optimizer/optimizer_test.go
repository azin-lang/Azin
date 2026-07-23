package optimizer

import (
	"testing"

	"github.com/azin-lang/Azin/internal/ast"
	"github.com/azin-lang/Azin/internal/token"
)

// ---- helpers ----

func id(name string) *ast.Identifier {
	return &ast.Identifier{Value: name}
}

func str(v string) *ast.StringLiteral {
	return &ast.StringLiteral{Value: v}
}

func tok(kind token.Kind) token.Token {
	return token.Token{Kind: kind}
}

func bin(left ast.Expr, op token.Kind, right ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		Left:     left,
		Operator: tok(op),
		Right:    right,
	}
}

// ---- Constant Folding ----

func TestFoldInteger(t *testing.T) {
	tests := []struct {
		name string
		op   token.Kind
		a, b int64
		want int64
	}{
		{"add", token.Plus, 3, 4, 7},
		{"sub", token.Minus, 10, 3, 7},
		{"mul", token.Star, 6, 7, 42},
		{"div", token.Slash, 10, 2, 5},
		{"mod", token.Modulo, 10, 3, 1},
		{"eq_true", token.EqualEqual, 5, 5, 1},
		{"eq_false", token.EqualEqual, 5, 6, 0},
		{"ne_true", token.BangEqual, 5, 6, 1},
		{"ne_false", token.BangEqual, 5, 5, 0},
		{"lt_true", token.Less, 3, 5, 1},
		{"lt_false", token.Less, 5, 3, 0},
		{"le_true", token.LessEqual, 5, 5, 1},
		{"le_false", token.LessEqual, 6, 5, 0},
		{"gt_true", token.Greater, 5, 3, 1},
		{"gt_false", token.Greater, 3, 5, 0},
		{"ge_true", token.GreaterEqual, 5, 5, 1},
		{"ge_false", token.GreaterEqual, 3, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := foldInteger(intLit(tt.a), tok(tt.op), intLit(tt.b))
			if result == nil {
				t.Fatal("got nil, expected literal")
			}
			switch v := result.(type) {
			case *ast.IntegerLiteral:
				if v.Value != tt.want {
					t.Errorf("got %d, want %d", v.Value, tt.want)
				}
			case *ast.BooleanLiteral:
				wantBool := tt.want != 0
				if v.Value != wantBool {
					t.Errorf("got %t, want %t", v.Value, wantBool)
				}
			default:
				t.Fatalf("got %T, expected IntegerLiteral or BooleanLiteral", result)
			}
		})
	}
}

func TestFoldIntegerDivByZero(t *testing.T) {
	if r := foldInteger(intLit(5), tok(token.Slash), intLit(0)); r != nil {
		t.Errorf("expected nil for division by zero, got %v", r)
	}
	if r := foldInteger(intLit(5), tok(token.Modulo), intLit(0)); r != nil {
		t.Errorf("expected nil for modulo by zero, got %v", r)
	}
}

func TestFoldFloat(t *testing.T) {
	tests := []struct {
		name string
		op   token.Kind
		a, b float64
		want float64
	}{
		{"add", token.Plus, 1.5, 2.5, 4.0},
		{"sub", token.Minus, 5.0, 2.0, 3.0},
		{"mul", token.Star, 3.0, 1.5, 4.5},
		{"div", token.Slash, 10.0, 2.0, 5.0},
		{"eq_true", token.EqualEqual, 3.0, 3.0, 1},
		{"eq_false", token.EqualEqual, 1.0, 2.0, 0},
		{"lt_true", token.Less, 1.0, 2.0, 1},
		{"lt_false", token.Less, 2.0, 1.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := foldFloat(floatLit(tt.a), tok(tt.op), floatLit(tt.b))
			if result == nil {
				t.Fatal("got nil, expected literal")
			}
			switch v := result.(type) {
			case *ast.FloatLiteral:
				if v.Value != tt.want {
					t.Errorf("got %g, want %g", v.Value, tt.want)
				}
			case *ast.BooleanLiteral:
				wantBool := tt.want != 0
				if v.Value != wantBool {
					t.Errorf("got %t, want %t", v.Value, wantBool)
				}
			default:
				t.Fatalf("got %T, expected FloatLiteral or BooleanLiteral", result)
			}
		})
	}
}

func TestFoldFloatDivByZero(t *testing.T) {
	if r := foldFloat(floatLit(5.0), tok(token.Slash), floatLit(0)); r != nil {
		t.Errorf("expected nil for float division by zero, got %v", r)
	}
}

func TestFoldBoolean(t *testing.T) {
	tests := []struct {
		name string
		op   token.Kind
		a, b bool
		want bool
	}{
		{"and_true", token.LogicalAnd, true, true, true},
		{"and_false", token.LogicalAnd, true, false, false},
		{"or_true", token.LogicalOr, false, true, true},
		{"or_false", token.LogicalOr, false, false, false},
		{"eq_true", token.EqualEqual, true, true, true},
		{"ne_true", token.BangEqual, true, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := foldBoolean(boolLit(tt.a), tok(tt.op), boolLit(tt.b))
			if result == nil {
				t.Fatal("got nil, expected boolean literal")
			}
			got, ok := result.(*ast.BooleanLiteral)
			if !ok {
				t.Fatalf("got %T, expected *ast.BooleanLiteral", result)
			}
			if got.Value != tt.want {
				t.Errorf("got %t, want %t", got.Value, tt.want)
			}
		})
	}
}

func TestFoldBinaryExpr(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want int64
	}{
		{"int_add", bin(intLit(1), token.Plus, intLit(2)), 3},
		{"int_mul", bin(intLit(3), token.Star, intLit(4)), 12},
		{"int_sub", bin(intLit(10), token.Minus, intLit(3)), 7},
		{"char_add", bin(&ast.CharacterLiteral{Value: 'A'}, token.Plus, &ast.CharacterLiteral{Value: 1}), 66},
		{"char_int_add", bin(&ast.CharacterLiteral{Value: 'A'}, token.Plus, intLit(1)), 66},
		{"int_float_add", bin(intLit(3), token.Plus, floatLit(2.5)), -1},
		{"float_int_add", bin(floatLit(1.5), token.Plus, intLit(2)), -1},
		{"not_foldable", bin(id("x"), token.Plus, intLit(1)), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			binExpr := tt.expr.(*ast.BinaryExpr)
			result := foldBinaryExpr(binExpr.Left, binExpr.Operator, binExpr.Right)
			if tt.want >= 0 {
				if result == nil {
					t.Fatal("expected fold, got nil")
				}
				got, ok := result.(*ast.IntegerLiteral)
				if !ok {
					t.Fatalf("got %T, want *ast.IntegerLiteral", result)
				}
				if got.Value != tt.want {
					t.Errorf("got %d, want %d", got.Value, tt.want)
				}
			} else {
				if result != nil {
					// float ops should also produce result, just different type
					switch result.(type) {
					case *ast.FloatLiteral:
						// OK
					default:
						t.Errorf("expected nil or FloatLiteral, got %T", result)
					}
				}
			}
		})
	}
}

// ---- Algebraic Simplification ----

func TestSimplifyArithmetic(t *testing.T) {
	tests := []struct {
		name    string
		expr    *ast.BinaryExpr
		wantNil bool
		check   func(t *testing.T, got ast.Expr)
	}{
		{
			name: "x_plus_0",
			expr: bin(id("x"), token.Plus, intLit(0)),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want identifier 'x'", got)
				}
			},
		},
		{
			name: "x_minus_0",
			expr: bin(id("x"), token.Minus, intLit(0)),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want identifier 'x'", got)
				}
			},
		},
		{
			name: "x_minus_x",
			expr: bin(id("x"), token.Minus, id("x")),
			check: func(t *testing.T, got ast.Expr) {
				lit, ok := got.(*ast.IntegerLiteral)
				if !ok || lit.Value != 0 {
					t.Errorf("got %v, want 0", got)
				}
			},
		},
		{
			name: "x_times_1",
			expr: bin(id("x"), token.Star, intLit(1)),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want identifier 'x'", got)
				}
			},
		},
		{
			name: "x_times_0",
			expr: bin(id("x"), token.Star, intLit(0)),
			check: func(t *testing.T, got ast.Expr) {
				lit, ok := got.(*ast.IntegerLiteral)
				if !ok || lit.Value != 0 {
					t.Errorf("got %v, want 0", got)
				}
			},
		},
		{
			name: "x_times_2",
			expr: bin(id("x"), token.Star, intLit(2)),
			check: func(t *testing.T, got ast.Expr) {
				bin, ok := got.(*ast.BinaryExpr)
				if !ok {
					t.Fatalf("got %T, want BinaryExpr", got)
				}
				if bin.Operator.Kind != token.LessLess {
					t.Errorf("expected << operator, got %v", bin.Operator.Kind)
				}
				r, ok := bin.Right.(*ast.IntegerLiteral)
				if !ok || r.Value != 1 {
					t.Errorf("expected shift by 1, got %v", r)
				}
			},
		},
		{
			name: "x_div_1",
			expr: bin(id("x"), token.Slash, intLit(1)),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want identifier 'x'", got)
				}
			},
		},
		{
			name: "x_div_x",
			expr: bin(id("x"), token.Slash, id("x")),
			check: func(t *testing.T, got ast.Expr) {
				lit, ok := got.(*ast.IntegerLiteral)
				if !ok || lit.Value != 1 {
					t.Errorf("got %v, want 1", got)
				}
			},
		},
		{
			name: "x_mod_1",
			expr: bin(id("x"), token.Modulo, intLit(1)),
			check: func(t *testing.T, got ast.Expr) {
				lit, ok := got.(*ast.IntegerLiteral)
				if !ok || lit.Value != 0 {
					t.Errorf("got %v, want 0", got)
				}
			},
		},
		{
			name: "side_effect_not_removed",
			expr: bin(&ast.CallExpr{ResolvedName: "foo"}, token.Star, intLit(0)),
			check: func(t *testing.T, got ast.Expr) {
				if got != nil {
					t.Errorf("expected nil for impure x * 0, got %v", got)
				}
			},
		},
		{
			name: "x_times_8",
			expr: bin(id("x"), token.Star, intLit(8)),
			check: func(t *testing.T, got ast.Expr) {
				bin, ok := got.(*ast.BinaryExpr)
				if !ok {
					t.Fatalf("got %T, want BinaryExpr", got)
				}
				if bin.Operator.Kind != token.LessLess {
					t.Errorf("expected << operator, got %v", bin.Operator.Kind)
				}
				r, ok := bin.Right.(*ast.IntegerLiteral)
				if !ok || r.Value != 3 {
					t.Errorf("expected shift by 3, got %v", r)
				}
			},
		},
		{
			name: "x_times_non_power_of_2",
			expr: bin(id("x"), token.Star, intLit(6)),
			check: func(t *testing.T, got ast.Expr) {
				if got != nil {
					t.Errorf("expected nil for non-power-of-2, got %v", got)
				}
			},
		},
		{
			name: "x_times_impure_power_of_2",
			expr: bin(&ast.CallExpr{ResolvedName: "foo"}, token.Star, intLit(8)),
			check: func(t *testing.T, got ast.Expr) {
				if got != nil {
					t.Errorf("expected nil for impure left, got %v", got)
				}
			},
		},
		{
			name: "x_div_8",
			expr: bin(intLit(100), token.Slash, intLit(8)),
			check: func(t *testing.T, got ast.Expr) {
				bin, ok := got.(*ast.BinaryExpr)
				if !ok {
					t.Fatalf("got %T, want BinaryExpr", got)
				}
				if bin.Operator.Kind != token.GreaterGreater {
					t.Errorf("expected >> operator, got %v", bin.Operator.Kind)
				}
				r, ok := bin.Right.(*ast.IntegerLiteral)
				if !ok || r.Value != 3 {
					t.Errorf("expected shift by 3, got %v", r)
				}
			},
		},
		{
			name: "x_div_signed_skip",
			expr: bin(id("x"), token.Slash, intLit(8)),
			check: func(t *testing.T, got ast.Expr) {
				if got != nil {
					t.Errorf("expected nil for signed variable, got %v", got)
				}
			},
		},
		{
			name: "x_div_non_power_of_2",
			expr: bin(intLit(100), token.Slash, intLit(6)),
			check: func(t *testing.T, got ast.Expr) {
				if got != nil {
					t.Errorf("expected nil for non-power-of-2, got %v", got)
				}
			},
		},
		{
			name: "x_mod_8",
			expr: bin(intLit(100), token.Modulo, intLit(8)),
			check: func(t *testing.T, got ast.Expr) {
				bin, ok := got.(*ast.BinaryExpr)
				if !ok {
					t.Fatalf("got %T, want BinaryExpr", got)
				}
				if bin.Operator.Kind != token.Ampersand {
					t.Errorf("expected & operator, got %v", bin.Operator.Kind)
				}
				r, ok := bin.Right.(*ast.IntegerLiteral)
				if !ok || r.Value != 7 {
					t.Errorf("expected mask 7, got %v", r)
				}
			},
		},
		{
			name: "x_mod_signed_skip",
			expr: bin(id("x"), token.Modulo, intLit(8)),
			check: func(t *testing.T, got ast.Expr) {
				if got != nil {
					t.Errorf("expected nil for signed variable, got %v", got)
				}
			},
		},
		{
			name: "x_mod_non_power_of_2",
			expr: bin(intLit(100), token.Modulo, intLit(6)),
			check: func(t *testing.T, got ast.Expr) {
				if got != nil {
					t.Errorf("expected nil for non-power-of-2, got %v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := simplifyArithmetic(tt.expr)
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

// ---- Boolean Simplification ----

func TestSimplifyBoolean(t *testing.T) {
	tests := []struct {
		name  string
		expr  *ast.BinaryExpr
		check func(t *testing.T, got ast.Expr)
	}{
		{
			name: "true_and_x",
			expr: bin(boolLit(true), token.LogicalAnd, id("x")),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want 'x'", got)
				}
			},
		},
		{
			name: "false_and_x",
			expr: bin(boolLit(false), token.LogicalAnd, id("x")),
			check: func(t *testing.T, got ast.Expr) {
				b, ok := got.(*ast.BooleanLiteral)
				if !ok || b.Value {
					t.Errorf("got %v, want false", got)
				}
			},
		},
		{
			name: "x_and_true",
			expr: bin(id("x"), token.LogicalAnd, boolLit(true)),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want 'x'", got)
				}
			},
		},
		{
			name: "false_or_x",
			expr: bin(boolLit(false), token.LogicalOr, id("x")),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want 'x'", got)
				}
			},
		},
		{
			name: "true_or_x",
			expr: bin(boolLit(true), token.LogicalOr, id("x")),
			check: func(t *testing.T, got ast.Expr) {
				b, ok := got.(*ast.BooleanLiteral)
				if !ok || !b.Value {
					t.Errorf("got %v, want true", got)
				}
			},
		},
		{
			name: "x_or_false",
			expr: bin(id("x"), token.LogicalOr, boolLit(false)),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want 'x'", got)
				}
			},
		},
		{
			name: "x_eq_x",
			expr: bin(intLit(5), token.EqualEqual, intLit(5)),
			check: func(t *testing.T, got ast.Expr) {
				b, ok := got.(*ast.BooleanLiteral)
				if !ok || !b.Value {
					t.Errorf("got %v, want true", got)
				}
			},
		},
		{
			name: "x_lt_x",
			expr: bin(id("x"), token.Less, id("x")),
			check: func(t *testing.T, got ast.Expr) {
				b, ok := got.(*ast.BooleanLiteral)
				if !ok || b.Value {
					t.Errorf("got %v, want false", got)
				}
			},
		},
		{
			name: "x_le_x",
			expr: bin(intLit(3), token.LessEqual, intLit(3)),
			check: func(t *testing.T, got ast.Expr) {
				b, ok := got.(*ast.BooleanLiteral)
				if !ok || !b.Value {
					t.Errorf("got %v, want true", got)
				}
			},
		},
		{
			name: "x_eq_true",
			expr: bin(id("x"), token.EqualEqual, boolLit(true)),
			check: func(t *testing.T, got ast.Expr) {
				id, ok := got.(*ast.Identifier)
				if !ok || id.Value != "x" {
					t.Errorf("got %v, want 'x'", got)
				}
			},
		},
		{
			name: "float_eq_float_skipped",
			expr: bin(floatLit(1.0), token.EqualEqual, floatLit(1.0)),
			check: func(t *testing.T, got ast.Expr) {
				// NaN != NaN, so float x == x must NOT optimize
				if got != nil {
					t.Errorf("expected nil for float == float, got %v", got)
				}
			},
		},
		{
			name: "int_eq_int_allowed",
			expr: bin(intLit(5), token.EqualEqual, intLit(5)),
			check: func(t *testing.T, got ast.Expr) {
				b, ok := got.(*ast.BooleanLiteral)
				if !ok || !b.Value {
					t.Errorf("got %v, want true", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := simplifyBoolean(tt.expr)
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

// ---- Normalize ----

func TestReassociateBinary(t *testing.T) {
	// (x + 1) + 2 → x + 3
	inner := bin(id("x"), token.Plus, intLit(1))
	outer := bin(inner, token.Plus, intLit(2))

	got := reassociateBinary(outer)
	if got == nil {
		t.Fatal("expected reassociation, got nil")
	}
	binExpr, ok := got.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("got %T, want *ast.BinaryExpr", got)
	}
	// Check: left = x, op = +, right = 3
	left, ok := binExpr.Left.(*ast.Identifier)
	if !ok || left.Value != "x" {
		t.Errorf("expected left 'x', got %v", binExpr.Left)
	}
	if binExpr.Operator.Kind != token.Plus {
		t.Errorf("expected +, got %v", binExpr.Operator.Kind)
	}
	right, ok := binExpr.Right.(*ast.IntegerLiteral)
	if !ok || right.Value != 3 {
		t.Errorf("expected right 3, got %v", binExpr.Right)
	}

	// (x * 2) * 3 → x * 6
	inner2 := bin(id("x"), token.Star, intLit(2))
	outer2 := bin(inner2, token.Star, intLit(3))

	got2 := reassociateBinary(outer2)
	if got2 == nil {
		t.Fatal("expected reassociation, got nil")
	}
	r2, ok := got2.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("got %T, want *ast.BinaryExpr", got2)
	}
	r2right, ok := r2.Right.(*ast.IntegerLiteral)
	if !ok || r2right.Value != 6 {
		t.Errorf("expected right 6, got %v", r2.Right)
	}
}

func TestCanonicalizeBinary(t *testing.T) {
	// 1 + x → x + 1 (constant moved to right)
	expr := bin(intLit(1), token.Plus, id("x"))
	got := canonicalizeBinary(expr)
	if got != expr {
		t.Fatal("expected same pointer back")
	}
	if isConstant(expr.Left) {
		t.Error("expected non-constant on left side after canonicalization")
	}
	if !isConstant(expr.Right) {
		t.Error("expected constant on right side after canonicalization")
	}
}

// ---- If Optimization ----

func TestOptimizeIfConstantTrue(t *testing.T) {
	ifStmt := &ast.IfStmt{
		Condition: boolLit(true),
		Then:      []ast.Stmt{&ast.ReturnStmt{Value: intLit(1)}},
		Else:      []ast.Stmt{&ast.ReturnStmt{Value: intLit(2)}},
	}
	result := optimizeIf(ifStmt)
	if len(result) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(result))
	}
	_, ok := result[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("expected ReturnStmt, got %T", result[0])
	}
}

func TestOptimizeIfConstantFalse(t *testing.T) {
	ifStmt := &ast.IfStmt{
		Condition: boolLit(false),
		Then:      []ast.Stmt{&ast.ReturnStmt{Value: intLit(1)}},
		Else:      []ast.Stmt{&ast.ReturnStmt{Value: intLit(2)}},
	}
	result := optimizeIf(ifStmt)
	if len(result) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(result))
	}
	rs, ok := result[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("expected ReturnStmt, got %T", result[0])
	}
	v, ok := rs.Value.(*ast.IntegerLiteral)
	if !ok || v.Value != 2 {
		t.Errorf("expected return 2, got %v", rs.Value)
	}
}

func TestOptimizeIfBothBranchesEmpty(t *testing.T) {
	// Both branches empty: pure if condition gets eliminated entirely
	ifStmt := &ast.IfStmt{
		Condition: boolLit(true),
		Then:      nil,
		Else:      nil,
	}
	result := optimizeIf(ifStmt)
	if len(result) != 0 {
		t.Fatalf("expected 0 stmts (pure condition eliminated), got %d", len(result))
	}

	// Both branches empty with impure condition becomes ExpressionStmt
	ifStmt2 := &ast.IfStmt{
		Condition: &ast.CallExpr{ResolvedName: "side_effect"},
		Then:      nil,
		Else:      nil,
	}
	result2 := optimizeIf(ifStmt2)
	if len(result2) != 1 {
		t.Fatalf("expected 1 stmt (expression retained), got %d", len(result2))
	}
	_, ok := result2[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("expected ExpressionStmt, got %T", result2[0])
	}
}

func TestOptimizeIfTailMerge(t *testing.T) {
	ifStmt := &ast.IfStmt{
		Condition: id("x"),
		Then: []ast.Stmt{
			&ast.ExpressionStmt{Expression: id("a")},
			&ast.ReturnStmt{Value: intLit(42)},
		},
		Else: []ast.Stmt{
			&ast.ExpressionStmt{Expression: id("b")},
			&ast.ReturnStmt{Value: intLit(42)},
		},
	}
	result := optimizeIf(ifStmt)
	// Should have if + tail return
	foundReturn := false
	for _, s := range result {
		if _, ok := s.(*ast.ReturnStmt); ok {
			foundReturn = true
		}
	}
	if !foundReturn {
		t.Error("expected tail return to be merged out")
	}
}

func TestOptimizeIfUnnestElse(t *testing.T) {
	ifStmt := &ast.IfStmt{
		Condition: id("x"),
		Then:      []ast.Stmt{&ast.ReturnStmt{Value: intLit(1)}},
		Else:      []ast.Stmt{&ast.ReturnStmt{Value: intLit(2)}},
	}
	result := optimizeIf(ifStmt)
	// Terminal then → else unnested after if
	if len(result) != 2 {
		t.Fatalf("expected 2 stmts (if + else), got %d", len(result))
	}
	// First should be the if (with no else)
	ifStmt2, ok := result[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("expected IfStmt, got %T", result[0])
	}
	if len(ifStmt2.Else) != 0 {
		t.Error("expected if to have no else after unnest")
	}
}

// ---- Statement Optimization ----

func TestOptimizeStatementsDeadCode(t *testing.T) {
	stmts := []ast.Stmt{
		&ast.ReturnStmt{Value: intLit(1)},
		&ast.ReturnStmt{Value: intLit(2)}, // dead code
		&ast.ReturnStmt{Value: intLit(3)}, // dead code
	}
	result := optimizeStatements(stmts)
	if len(result) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(result))
	}
}

func TestOptimizeStatementsExpressionElimination(t *testing.T) {
	stmts := []ast.Stmt{
		&ast.ExpressionStmt{Expression: intLit(42)}, // pure, should be eliminated
		&ast.ReturnStmt{Value: intLit(0)},
	}
	result := optimizeStatements(stmts)
	if len(result) != 1 {
		t.Fatalf("expected 1 stmt (only return), got %d", len(result))
	}
}

func TestOptimizeLoopUnwrapReturn(t *testing.T) {
	// Use non-pure expression so it's not eliminated
	loop := &ast.LoopStmt{
		Body: []ast.Stmt{
			&ast.AssignmentStmt{Left: id("x"), Value: intLit(1)},
			&ast.ReturnStmt{Value: intLit(0)},
		},
	}
	result := optimizeStatement(loop)
	if len(result) != 2 {
		t.Fatalf("expected 2 stmts (unwrapped), got %d", len(result))
	}
	_, ok := result[1].(*ast.ReturnStmt)
	if !ok {
		t.Errorf("expected return as last stmt, got %T", result[1])
	}
}

func TestOptimizeLoopUnwrapStop(t *testing.T) {
	// Use non-pure expression so it's not eliminated
	loop := &ast.LoopStmt{
		Body: []ast.Stmt{
			&ast.AssignmentStmt{Left: id("x"), Value: intLit(1)},
			&ast.StopStmt{},
		},
	}
	result := optimizeStatement(loop)
	if len(result) != 1 {
		t.Fatalf("expected 1 stmt (stop removed), got %d", len(result))
	}
	_, ok := result[0].(*ast.AssignmentStmt)
	if !ok {
		t.Errorf("expected AssignmentStmt, got %T", result[0])
	}
}

func TestOptimizeLoopNotUnwrapped(t *testing.T) {
	// Loop without terminal body should not be unwrapped
	loop := &ast.LoopStmt{
		Body: []ast.Stmt{
			&ast.ExpressionStmt{Expression: id("x")},
		},
	}
	result := optimizeStatement(loop)
	if len(result) != 1 {
		t.Fatalf("expected 1 stmt (loop kept), got %d", len(result))
	}
	_, ok := result[0].(*ast.LoopStmt)
	if !ok {
		t.Errorf("expected LoopStmt, got %T", result[0])
	}
}

// ---- Helpers ----

func TestIsTerminal(t *testing.T) {
	if !isTerminal(&ast.ReturnStmt{}) {
		t.Error("ReturnStmt should be terminal")
	}
	if !isTerminal(&ast.StopStmt{}) {
		t.Error("StopStmt should be terminal")
	}
	if isTerminal(&ast.ExpressionStmt{}) {
		t.Error("ExpressionStmt should not be terminal")
	}
}

func TestBlockIsTerminal(t *testing.T) {
	if blockIsTerminal(nil) {
		t.Error("nil block should not be terminal")
	}
	if blockIsTerminal([]ast.Stmt{}) {
		t.Error("empty block should not be terminal")
	}
	if blockIsTerminal([]ast.Stmt{&ast.ExpressionStmt{}}) {
		t.Error("non-terminal block should not be terminal")
	}
	if !blockIsTerminal([]ast.Stmt{&ast.ReturnStmt{}}) {
		t.Error("block ending with return should be terminal")
	}
	if blockIsTerminal([]ast.Stmt{&ast.ReturnStmt{}, &ast.ExpressionStmt{Expression: intLit(1)}}) {
		t.Error("block with dead code after return should still report terminal from optimizer's perspective")
	}
}

func TestIsNotFloat(t *testing.T) {
	if !isNotFloat(intLit(1)) {
		t.Error("int should not be float")
	}
	if !isNotFloat(boolLit(true)) {
		t.Error("bool should not be float")
	}
	if !isNotFloat(&ast.CharacterLiteral{Value: 'a'}) {
		t.Error("char should not be float")
	}
	if !isNotFloat(str("hello")) {
		t.Error("string should not be float")
	}
	if isNotFloat(floatLit(3.14)) {
		t.Error("float should be float")
	}
	if isNotFloat(id("x")) {
		t.Error("identifier should conservatively be considered possibly float")
	}
}

func TestIsZero(t *testing.T) {
	if !isZero(intLit(0)) {
		t.Error("int 0 should be zero")
	}
	if isZero(intLit(1)) {
		t.Error("int 1 should not be zero")
	}
	if !isZero(floatLit(0)) {
		t.Error("float 0 should be zero")
	}
	if !isZero(&ast.CharacterLiteral{Value: 0}) {
		t.Error("char 0 should be zero")
	}
}

func TestIsOne(t *testing.T) {
	if !isOne(intLit(1)) {
		t.Error("int 1 should be one")
	}
	if isOne(intLit(0)) {
		t.Error("int 0 should not be one")
	}
	if !isOne(floatLit(1)) {
		t.Error("float 1 should be one")
	}
}

func TestIsPure(t *testing.T) {
	if !isPure(intLit(42)) {
		t.Error("literal should be pure")
	}
	if !isPure(id("x")) {
		t.Error("identifier should be pure")
	}
	if !isPure(bin(intLit(1), token.Plus, intLit(2))) {
		t.Error("pure binary expr should be pure")
	}
	if isPure(&ast.CallExpr{ResolvedName: "foo"}) {
		t.Error("call expression should not be pure")
	}
}

// ---- Full Optimize ----

func TestOptimizeFullProgram(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Stmt{
			&ast.ExpressionStmt{Expression: bin(intLit(1), token.Plus, intLit(2))},
			&ast.ReturnStmt{Value: bin(intLit(3), token.Minus, intLit(1))},
		},
	}
	Optimize(program)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 stmt after optimization, got %d", len(program.Statements))
	}
	rs, ok := program.Statements[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("expected ReturnStmt, got %T", program.Statements[0])
	}
	v, ok := rs.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected IntegerLiteral, got %T", rs.Value)
	}
	if v.Value != 2 {
		t.Errorf("expected 2, got %d", v.Value)
	}
}

func TestOptimizeNilProgram(t *testing.T) {
	// Should not panic
	Optimize(nil)
}
