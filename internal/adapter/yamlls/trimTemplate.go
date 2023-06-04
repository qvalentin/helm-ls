package yamlls

import (

	// lsplocal "github.com/mrjosh/helm-ls/internal/lsp"
	sitter "github.com/smacker/go-tree-sitter"
	// lsp "go.lsp.dev/protocol"
)

func trimTemplateForYamllsFromAst(ast *sitter.Tree, text string) string {

	var result = []byte(text)
	logger.Println(ast.RootNode())
	prettyPrintNode(ast.RootNode(), []byte(text), result)
	return string(result)

}

func prettyPrintNode(node *sitter.Node, previous []byte, result []byte) {

	var childCount = node.ChildCount()

	if childCount == 0 {
		logger.Println("End of recursion", string(node.Content(previous)))
		logger.Println("Type", node.Type())

	}

	switch node.Type() {
	case "block_action":
		earaseTemplate(node.Child(0), previous, result)
		earaseTemplate(node.Child(1), previous, result)
		earaseTemplate(node.Child(2), previous, result)
		earaseTemplate(node.Child(3), previous, result)
		earaseTemplate(node.Child(5), previous, result)
		earaseTemplate(node.Child(6), previous, result)
		earaseTemplate(node.Child(7), previous, result)
	case "if_action":
		earaseTemplate(node.Child(0), previous, result)
		earaseTemplate(node.Child(1), previous, result)
		earaseTemplate(node.Child(2), previous, result)
		earaseTemplate(node.Child(3), previous, result)
		earaseTemplate(node.Child(int(node.ChildCount()-3)), previous, result)
		earaseTemplate(node.Child(int(node.ChildCount()-2)), previous, result)
		earaseTemplate(node.Child(int(node.ChildCount())-1), previous, result)
	default:
		for i := 0; i < int(childCount); i++ {
			prettyPrintNode(node.Child(i), previous, result)
		}
	}

}

func earaseTemplate(node *sitter.Node, previous []byte, result []byte) {
	logger.Println("Content", node.Content(previous))
	for i := range []byte(node.Content(previous)) {
		result[int(node.StartByte())+i] = byte(' ')
	}
}
