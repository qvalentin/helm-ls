package yamlls

import (
	"context"
	"encoding/json"

	lsplocal "github.com/mrjosh/helm-ls/internal/lsp"
	sitter "github.com/smacker/go-tree-sitter"
	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"
)

func handleDiagnostics(req jsonrpc2.Request, clientConn jsonrpc2.Conn, documents *lsplocal.DocumentStore) {
	var params lsp.PublishDiagnosticsParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		logger.Println("Error handling diagnostic", err)
	}

	doc, ok := documents.Get(params.URI)
	if !ok {
		logger.Println("Error handling diagnostic. Could not get document: " + params.URI.Filename())
	}
	doc.DiagnosticsCache.Yamldiagnostics = filterDiagnostics(params.Diagnostics, doc.Ast)
	params.Diagnostics = doc.DiagnosticsCache.GetMergedDiagnostics()

	clientConn.Notify(context.Background(), lsp.MethodTextDocumentPublishDiagnostics, &params)
}

func filterDiagnostics(diagnostics []lsp.Diagnostic, ast *sitter.Tree) (filtered []lsp.Diagnostic) {
	filtered = []lsp.Diagnostic{}
	for _, diagnostic := range diagnostics {
		node := lsplocal.FindRelevantChildNode(ast.RootNode(), lsplocal.GetSitterPointForLspPos(diagnostic.Range.Start))
		if node.Type() == "text" {
			filtered = append(filtered, diagnostic)
		}
	}
	return filtered
}