package yamlls

import (
	lsplocal "github.com/mrjosh/helm-ls/internal/lsp"
	"testing"
)

func TestTrimTemplate(t *testing.T) {

	var documentText = `
{{ .Values.global. }}
yaml: test

{{block "name"}} T1 {{end}}
`

	var trimmedText = `
{{ .Values.global. }}
yaml: test

                 T1        
`
	doc := &lsplocal.Document{
		Content: documentText,
		Ast:     lsplocal.ParseAst(documentText),
	}

	var trimmed = trimTemplateForYamlls(doc.Ast, documentText)

	var result = trimmed == trimmedText

	if !result {
		t.Errorf("Trimmed templated was not as expected but was %s ", trimmed)
	} else {
		t.Log("Trimmed templated was as expected")
	}

}
func TestTrimTemplateFromAst(t *testing.T) {

	var documentText = `
{{ .Values.global. }}
yaml: test

{{block "name"}} T1 {{end}}
`

	doc := &lsplocal.Document{
		Content: documentText,
		Ast:     lsplocal.ParseAst(documentText),
	}

	var trimmed = trimTemplateForYamllsFromAst(doc.Ast, documentText)

	var result = trimmed == documentText

	if !result {
		t.Errorf("Trimmed templated was not as expected but was %s ", trimmed)
	} else {
		t.Log("Trimmed templated was as expected")
	}

}
