package ast

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/azin-lang/azin/lsp/internal/token"
)

type debugPrinter struct {
	w        io.Writer
	useColor bool
	visited  map[uintptr]struct{}
	maxDepth int
}

func newDebugPrinter(w io.Writer, color bool) *debugPrinter {
	return &debugPrinter{
		w:        w,
		useColor: color,
		visited:  make(map[uintptr]struct{}),
		maxDepth: 128,
	}
}

func (p *debugPrinter) Print(node Node) {
	p.visit(reflect.ValueOf(node), "", true, 0)
}

func (p *debugPrinter) visit(
	v reflect.Value,
	prefix string,
	last bool,
	depth int,
) {
	if depth > p.maxDepth {
		p.line(prefix, last, p.style(cValue, "..."))
		return
	}

	v = unwrap(v)

	if !v.IsValid() {
		p.line(prefix, last, p.style(cValue, "nil"))
		return
	}

	if leave, cycle := p.enter(v, prefix, last); cycle {
		return
	} else if leave != nil {
		defer leave()
	}

	switch {
	case p.printPrimitive(v, prefix, last):
		return

	case p.printToken(v, prefix, last):
		return

	case p.printCompact(v, prefix, last):
		return

	case v.Kind() == reflect.Slice:
		p.visitSlice(v, prefix, last, depth)
		return
	}

	p.visitStruct(v, prefix, last, depth)
}

func (p *debugPrinter) visitStruct(
	v reflect.Value,
	prefix string,
	last bool,
	depth int,
) {
	title := v.Type().Name()

	if n := nodeOf(v); n != nil {
		title = p.nodeTitle(n)
	}

	p.line(prefix, last, title)

	if v.Kind() != reflect.Struct {
		return
	}

	next := p.childPrefix(prefix, last)
	fields := meaningfulFields(v)

	for i, field := range fields {
		p.visitField(
			field.Name,
			v.FieldByIndex(field.Index),
			next,
			i == len(fields)-1,
			depth+1,
		)
	}
}

func (p *debugPrinter) visitField(
	name string,
	value reflect.Value,
	prefix string,
	last bool,
	depth int,
) {
	value = unwrap(value)

	if text, ok := p.inlineValue(name, value); ok {
		p.line(
			prefix,
			last,
			fmt.Sprintf(
				"%s: %s",
				p.style(cField, name),
				text,
			),
		)
		return
	}

	p.line(prefix, last, p.style(cField, name))

	if value.Kind() == reflect.Slice && isTransparentSlice(name) {
		p.visitSlice(value, prefix, last, depth)
		return
	}

	p.visit(
		value,
		p.childPrefix(prefix, last),
		true,
		depth,
	)
}

func (p *debugPrinter) visitSlice(
	v reflect.Value,
	prefix string,
	last bool,
	depth int,
) {
	next := p.childPrefix(prefix, last)

	for i := 0; i < v.Len(); i++ {
		p.visit(
			v.Index(i),
			next,
			i == v.Len()-1,
			depth+1,
		)
	}
}

func (p *debugPrinter) printPrimitive(
	v reflect.Value,
	prefix string,
	last bool,
) bool {
	text, ok := primitive(v)
	if !ok {
		return false
	}

	p.line(prefix, last, p.style(cValue, text))
	return true
}

func (p *debugPrinter) printCompact(
	v reflect.Value,
	prefix string,
	last bool,
) bool {
	text, ok := compact(v)
	if !ok {
		return false
	}

	p.line(prefix, last, p.style(cLiteral, text))
	return true
}

func (p *debugPrinter) printToken(
	v reflect.Value,
	prefix string,
	last bool,
) bool {
	tok, ok := v.Interface().(token.Token)
	if !ok {
		return false
	}

	p.line(
		prefix,
		last,
		fmt.Sprintf(
			"%s %s",
			p.style(cLiteral, tok.Kind.String()),
			p.style(
				cValue,
				fmt.Sprintf("@%d +%d", tok.Position.Offset, tok.Length),
			),
		),
	)

	return true
}

func (p *debugPrinter) inlineValue(
	name string,
	v reflect.Value,
) (string, bool) {
	if text, ok := compact(v); ok {
		return p.style(cLiteral, text), true
	}

	if text, ok := primitive(v); ok {
		return p.style(cValue, text), true
	}

	if !isInlineField(name) {
		return "", false
	}

	fields := meaningfulFields(v)
	if len(fields) != 1 {
		return "", false
	}

	child := unwrap(v.FieldByIndex(fields[0].Index))

	if text, ok := compact(child); ok {
		return p.style(cLiteral, text), true
	}

	if text, ok := primitive(child); ok {
		return p.style(cValue, text), true
	}

	return "", false
}

func (p *debugPrinter) enter(
	v reflect.Value,
	prefix string,
	last bool,
) (func(), bool) {
	if v.Kind() != reflect.Pointer {
		return nil, false
	}

	ptr := v.Pointer()
	if ptr == 0 {
		return nil, false
	}

	if _, ok := p.visited[ptr]; ok {
		p.line(prefix, last, p.style(cValue, "<cycle>"))
		return nil, true
	}

	p.visited[ptr] = struct{}{}

	return func() {
		delete(p.visited, ptr)
	}, false
}

func (p *debugPrinter) nodeTitle(n Node) string {
	t := reflect.TypeOf(n)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	title := p.style(cNode, t.Name())

	if label := n.Label(); label != "" {
		title += " " + p.style(cLabel, strconv.Quote(label))
	}

	return title
}

func (p *debugPrinter) line(
	prefix string,
	last bool,
	text string,
) {
	branch := "├── "
	if last {
		branch = "╰── "
	}

	fmt.Fprintf(
		p.w,
		"%s%s%s\n",
		prefix,
		p.style(cBranch, branch),
		text,
	)
}

func (p *debugPrinter) childPrefix(
	prefix string,
	last bool,
) string {
	if last {
		return prefix + "    "
	}

	return prefix + p.style(cBranch, "│   ")
}

func (p *debugPrinter) style(
	fn func(...interface{}) string,
	s string,
) string {
	if !p.useColor {
		return s
	}

	return fn(s)
}
