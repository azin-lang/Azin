package main

import (
	"github.com/azin-lang/azin/lsp/internal/parser"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "Azin"

// keeps track of the current file's content
var documents = make(map[string]string)

var version = "0.2.0" // um, why???
var handler protocol.Handler

func documentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI

	// Because we set TextDocumentSyncKindFull in initialize, we get the entire file text
	if len(params.ContentChanges) > 0 {
		if change, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole); ok {
			documents[uri] = change.Text
			_, err := parser.Parse(context, uri, change.Text)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func documentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI
	text := params.TextDocument.Text

	documents[uri] = text
	analyze(context, uri, text)
	return nil
}

func main() {
	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		SetTrace:    setTrace,

		TextDocumentHover:     documentDidHover,
		TextDocumentDidOpen:   documentDidOpen,
		TextDocumentDidChange: documentDidChange,
	}

	lsp := server.NewServer(&handler, lsName, false)

	err := lsp.RunStdio()
	if err != nil {
		return
	}
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	capabilities.TextDocumentSync = &protocol.TextDocumentSyncOptions{
		OpenClose: new(true),
		Change:    new(protocol.TextDocumentSyncKindFull),
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func documentDidHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "**Azin LSP:** You hovered over some code!",
		},
	}, nil
}
