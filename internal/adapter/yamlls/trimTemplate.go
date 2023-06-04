package yamlls

import (

	// lsplocal "github.com/mrjosh/helm-ls/internal/lsp"
	sitter "github.com/smacker/go-tree-sitter"
	// lsp "go.lsp.dev/protocol"
)

func trimTemplateForYamllsFromAst(ast *sitter.Tree, text string) string {

	return prettyPrintNode(ast.RootNode(), text)

}

func prettyPrintNode(node *sitter.Node, text string) string {

	var result = ""

	var childCount = node.ChildCount()

	if childCount == 0 {
		logger.Println("End of recursion", string(node.Content([]byte(text))))
		logger.Println("Type", node.Type())
		result = result + string(node.Content([]byte(text)))
	}

	switch node.Type() {
	case "interpreted_string_literal":
		result = result + string(node.Content([]byte(text)))
		return result

	}

	for i := 0; i < int(childCount); i++ {
		result = result + prettyPrintNode(node.Child(i), text)
	}

	return result

}
