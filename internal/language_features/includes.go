package languagefeatures

import (
	lsp "go.lsp.dev/protocol"

	lsplocal "github.com/mrjosh/helm-ls/internal/lsp"
	"github.com/mrjosh/helm-ls/internal/tree-sitter/gotemplate"
	"github.com/mrjosh/helm-ls/internal/util"
	sitter "github.com/smacker/go-tree-sitter"
)

type IncludesFeature struct {
	*GenericDocumentUseCase
}

type IncludesCallFeature struct {
	*IncludesFeature
}

// should be called on {{ include "name" . }}
func (f *IncludesCallFeature) AppropriateForNode(currentNodeType string, parentNodeType string, node *sitter.Node) bool {
	if parentNodeType != gotemplate.NodeTypeArgumentList {
		return false
	}
	functionCallNode := node.Parent().Parent()
	_, err := lsplocal.ParseIncludeFunctionCall(functionCallNode, []byte(f.GenericDocumentUseCase.Document.Content))
	return err == nil
}

type IncludesDefinitionFeature struct {
	*IncludesFeature
}

// should be called on {{ define "name" }}
func (f *IncludesDefinitionFeature) AppropriateForNode(currentNodeType string, parentNodeType string, node *sitter.Node) bool {
	return parentNodeType == gotemplate.NodeTypeDefineAction && currentNodeType == gotemplate.NodeTypeInterpretedStringLiteral
}

func NewIncludesCallFeature(genericDocumentUseCase *GenericDocumentUseCase) *IncludesCallFeature {
	return &IncludesCallFeature{
		IncludesFeature: &IncludesFeature{genericDocumentUseCase},
	}
}

func NewIncludesDefinitionFeature(genericDocumentUseCase *GenericDocumentUseCase) *IncludesDefinitionFeature {
	return &IncludesDefinitionFeature{
		IncludesFeature: &IncludesFeature{genericDocumentUseCase},
	}
}

func (f *IncludesCallFeature) References() (result []lsp.Location, err error) {
	includeName, err := f.getIncludeName()
	if err != nil {
		return []lsp.Location{}, err
	}

	return f.getReferenceLocations(includeName), nil
}

func (f *IncludesCallFeature) getIncludeName() (string, error) {
	functionCallNode := f.Node.Parent().Parent()
	return lsplocal.ParseIncludeFunctionCall(functionCallNode, []byte(f.GenericDocumentUseCase.Document.Content))
}

func (f *IncludesDefinitionFeature) References() (result []lsp.Location, err error) {
	includeName := util.RemoveQuotes(f.GenericDocumentUseCase.NodeContent())
	return f.getReferenceLocations(includeName), nil
}

func (f *IncludesFeature) getReferenceLocations(includeName string) []lsp.Location {
	locations := []lsp.Location{}
	for _, doc := range f.GenericDocumentUseCase.DocumentStore.GetAllDocs() {
		referenceRanges := doc.SymbolTable.GetIncludeReference(includeName)
		for _, referenceRange := range referenceRanges {
			locations = append(locations, util.RangeToLocation(doc.URI, referenceRange))
		}
	}

	return locations
}

func (f *IncludesFeature) getDefinitionLocations(includeName string) []lsp.Location {
	locations := []lsp.Location{}
	for _, doc := range f.GenericDocumentUseCase.DocumentStore.GetAllDocs() {
		referenceRanges := doc.SymbolTable.GetIncludeDefinitions(includeName)
		for _, referenceRange := range referenceRanges {
			locations = append(locations, util.RangeToLocation(doc.URI, referenceRange))
		}
	}

	return locations
}

func (f *IncludesFeature) getDefinitionsHover(includeName string) util.HoverResultsWithFiles {
	result := util.HoverResultsWithFiles{}
	for _, doc := range f.GenericDocumentUseCase.DocumentStore.GetAllDocs() {
		referenceRanges := doc.SymbolTable.GetIncludeDefinitions(includeName)
		for _, referenceRange := range referenceRanges {
			node := doc.Ast.RootNode().NamedDescendantForPointRange(referenceRange.StartPoint, referenceRange.EndPoint)
			if node != nil {
				result = append(result, util.HoverResultWithFile{
					Value: node.Content([]byte(doc.Content)),
					URI:   doc.URI,
				})
			}
		}
	}

	return result
}

func (f *IncludesCallFeature) Hover() (string, error) {
	includeName, err := f.getIncludeName()
	if err != nil {
		return "", err
	}

	result := f.getDefinitionsHover(includeName)
	return result.Format(f.GenericDocumentUseCase.Document.URI), nil
}

func (f *IncludesCallFeature) Definition() (result []lsp.Location, err error) {
	includeName, err := f.getIncludeName()
	if err != nil {
		return []lsp.Location{}, err
	}
	return f.getDefinitionLocations(includeName), nil
}
