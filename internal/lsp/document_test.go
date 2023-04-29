package lsp

import (
	lsp "go.lsp.dev/protocol"
	"testing"
)

func TestBytePositionAt(t *testing.T) {
	input := `line 1
line 2
line 3`
	expected := 21

	sut := document{
		Content: input,
	}

	position := lsp.Position{
		Character: 5,
		Line:      2,
	}

	result := sut.BytePositionAt(position)

	if expected != result {
		t.Errorf("Expected %d but got %d.", expected, result)
	}
}

func TestBytePositionAtSecond(t *testing.T) {
	input := `line 1
line 2
line 3
line 4`
	expected := 27

	sut := document{
		Content: input,
	}

	position := lsp.Position{
		Character: 5,
		Line:      3,
	}

	result := sut.BytePositionAt(position)

	if expected != result {
		t.Errorf("Expected %d but got %d.", expected, result)
	}
}

func TestApplyChanges(t *testing.T) {

	previous := "{{ $.Values.image.repository }}"
	final := "{{ $.Values.image.tag }}"

	sut := document{
		Content: previous,
		Ast:     ParseAst(previous),
	}

	change := lsp.TextDocumentContentChangeEvent{
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      0,
				Character: 18,
			}, End: lsp.Position{
				Line:      0,
				Character: 28,
			}},
		RangeLength: 10,
	}
	change2 := lsp.TextDocumentContentChangeEvent{
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      0,
				Character: 18,
			}, End: lsp.Position{
				Line:      0,
				Character: 18,
			}},
		RangeLength: 0,
		Text:        "tag",
	}

	sut.ApplyChanges([]lsp.TextDocumentContentChangeEvent{change})
	sut.ApplyChanges([]lsp.TextDocumentContentChangeEvent{change2})

	if final != sut.Content {
		t.Errorf("Expected %s but got %s.", final, sut.Content)
	}

	changedText := sut.Ast.RootNode().Child(1).Child(2).Content([]byte(final))
	if changedText != "tag" {
		t.Errorf("Expected %s but got %s.", "tag", changedText)
	}

}

func TestApplyWithNewLine(t *testing.T) {

	previous := "{{ $.Values.image.tag }}"
	final := `{{ $.Values.image.tag }}
{{ $.Values.image.tag }}
`

	println(len(final) == len([]byte(final)))

	sut := document{
		Content: previous,
		Ast:     ParseAst(previous),
	}
	change := lsp.TextDocumentContentChangeEvent{
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      0,
				Character: 24,
			}, End: lsp.Position{
				Line:      1,
				Character: 0,
			}},
		RangeLength: 1,
		Text:        "\n{{ $.Values.image.tag }}\n",
	}

	sut.ApplyChanges([]lsp.TextDocumentContentChangeEvent{change})

	if final != sut.Content {
		t.Errorf("Expected %s but got %s.", final, sut.Content)
	}

	// changedText := sut.Ast.RootNode().Child(2).Child(2).Content([]byte(final))
	// if changedText != "tag" {
	// 	t.Errorf("Expected %s but got %s.", "tag", changedText)
	// }

}
