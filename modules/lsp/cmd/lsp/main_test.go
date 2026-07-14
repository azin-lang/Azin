package main

import (
	"testing"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// wyd bro
// fym wyd
func TestHoverHandler(t *testing.T) {
	ctx := &glsp.Context{}

	params := &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.az",
			},
			Position: protocol.Position{
				Line:      5,
				Character: 10,
			},
		},
	}

	resp, err := documentDidHover(ctx, params)
	if err != nil {
		t.Fatalf("hover returned an unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected hover response, got nil")
	}

	expectedText := "**Azin LSP:** You hovered over some code!"
	if resp.Contents.(protocol.MarkupContent).Value != expectedText {
		t.Errorf("expected hover content %q, got %q", expectedText, resp.Contents.(protocol.MarkupContent).Value)
	}

	println("tests ran successfully!")
}
