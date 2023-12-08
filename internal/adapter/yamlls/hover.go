package yamlls

import (
	"context"
	"reflect"

	"github.com/mrjosh/helm-ls/internal/util"
	lsp "go.lsp.dev/protocol"
)

// Calls the Completion method of yamlls to get a fitting hover response
// TODO: clarify why the hover method of yamlls can't be used
func (yamllsConnector Connector) CallHover(params lsp.HoverParams, word string) *lsp.Hover {
	if yamllsConnector.Conn == nil {
		return &lsp.Hover{}
	}

	var (
		documentation      string
		hoverResponse      = reflect.New(reflect.TypeOf(lsp.Hover{})).Interface()
		completionResponse = reflect.New(reflect.TypeOf(lsp.CompletionList{})).Interface()
		completionParams   = lsp.CompletionParams{
			TextDocumentPositionParams: params.TextDocumentPositionParams,
		}
	)

	_, err := (*yamllsConnector.Conn).Call(context.Background(), lsp.MethodTextDocumentHover, completionParams, hoverResponse)
	if err != nil {
		logger.Error("Error calling yamlls for hover", err)
		return &lsp.Hover{}
	}

	logger.Debug("Got hover from yamlls", hoverResponse.(*lsp.Hover).Contents.Value)

	if hoverResponse.(*lsp.Hover).Contents.Value != "" {
		return hoverResponse.(*lsp.Hover)
	}

	_, err = (*yamllsConnector.Conn).Call(context.Background(), lsp.MethodTextDocumentCompletion, completionParams, completionResponse)
	if err != nil {
		logger.Error("Error calling yamlls for Completion", err)
		return &lsp.Hover{}
	}

	for _, completionItem := range completionResponse.(*lsp.CompletionList).Items {
		if completionItem.InsertText == word {
			documentation = completionItem.Documentation.(string)
			break
		}
	}

	response := util.BuildHoverResponse(documentation, lsp.Range{})
	return &response
}
