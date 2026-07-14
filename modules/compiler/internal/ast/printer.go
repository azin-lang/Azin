package ast

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/fatih/color"
)

var (
	cNode    = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	cField   = color.New(color.FgYellow).SprintFunc()
	cValue   = color.New(color.FgHiGreen).SprintFunc()
	cLiteral = color.New(color.FgHiCyan).SprintFunc()
	cLabel   = color.New(color.FgWhite).SprintFunc()
	cBranch  = color.New(color.FgHiBlack).SprintFunc()
)

func PrintTree(node Node) {
	// newNormalPrinter(os.Stdout, true).Print(node)
}

func PrintDebugTree(node Node) {
	newDebugPrinter(os.Stdout, true).Print(node)
}

func ExportTree(node Node, path string) error {
	return export(path, func(f *os.File) {
		// newNormalPrinter(f, false).Print(node)
	})
}

func ExportDebugTree(node Node, path string) error {
	return export(path, func(f *os.File) {
		newDebugPrinter(f, false).Print(node)
	})
}

func export(path string, fn func(*os.File)) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fn(f)

	return nil
}

func unwrap(v reflect.Value) reflect.Value {
	for v.IsValid() {
		switch v.Kind() {
		case reflect.Pointer, reflect.Interface:
			if v.IsNil() {
				return reflect.Value{}
			}

			v = v.Elem()

		default:
			return v
		}
	}

	return v
}

func nodeOf(v reflect.Value) Node {
	v = unwrap(v)

	if !v.IsValid() {
		return nil
	}

	if v.CanInterface() {
		if n, ok := v.Interface().(Node); ok {
			return n
		}
	}

	if v.CanAddr() {
		if n, ok := v.Addr().Interface().(Node); ok {
			return n
		}
	}

	return nil
}

func compact(v reflect.Value) (string, bool) {
	n := nodeOf(v)
	if n == nil {
		return "", false
	}

	switch n := n.(type) {

	case *Identifier:
		return n.Value, true

	case *StringLiteral:
		return strconv.Quote(n.Value), true

	case *CharacterLiteral:
		return strconv.QuoteRune(n.Value), true

	case *IntegerLiteral:
		return strconv.FormatInt(n.Value, 10), true

	case *FloatLiteral:
		return strconv.FormatFloat(n.Value, 'g', -1, 64), true

	case *BooleanLiteral:
		if n.Value {
			return "✅", true
		}
		return "❌", true

	case *FieldDecl:
		if n.Type == nil {
			return n.Name.Value, true
		}

		return fmt.Sprintf(
			"%s: %s",
			n.Name.Value,
			n.Type.Label(),
		), true
	}

	return "", false
}

func shouldSkipField(name string) bool {
	switch name {
	case "Token", "Position":
		return true
	default:
		return false
	}
}

func skipFalseOrEmpty(v reflect.Value) bool {
	v = unwrap(v)

	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Slice:
		return v.Len() == 0

	case reflect.Bool:
		return !v.Bool()
	}

	return false
}

func isInlineField(name string) bool {
	switch name {
	case
		"Name",
		"Type",
		"ReturnType",
		"Operator",
		"Property",
		"Object",
		"Callee",
		"Path",
		"Value":
		return true
	}

	return false
}

func isTransparentSlice(name string) bool {
	switch name {
	case
		"Body",
		"Then",
		"Else",
		"Statements",
		"Params",
		"Fields",
		"Args":
		return true
	}

	return false
}

func meaningfulFields(v reflect.Value) []reflect.StructField {
	v = unwrap(v)

	if !v.IsValid() || v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	fields := make([]reflect.StructField, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)

		if !sf.IsExported() || shouldSkipField(sf.Name) {
			continue
		}

		if skipFalseOrEmpty(v.Field(i)) {
			continue
		}

		fields = append(fields, sf)
	}

	return fields
}

func primitive(v reflect.Value) (string, bool) {
	v = unwrap(v)

	if !v.IsValid() {
		return "", false
	}

	switch v.Kind() {

	case reflect.String:
		return strconv.Quote(v.String()), true

	case reflect.Bool:
		if v.Bool() {
			return "✅", true
		}
		return "❌", true

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), true

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), true

	case reflect.Float32,
		reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, 64), true
	}

	return "", false
}
