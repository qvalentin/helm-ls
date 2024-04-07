package lsp

import (
	"os"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stretchr/testify/assert"
)

func TestSymbolTableForIncludeDefinitions(t *testing.T) {
	content := `
	{{ define "foo" }}
	{{ .Values.global. }}
	{{ end }}

	{{ define "bar" }}
	{{ .Values.global. }}	
	{{ end }}
	`

	ast := ParseAst(nil, content)

	symbolTable := NewSymbolTable(ast, []byte(content))

	assert.Len(t, symbolTable.includeDefinitions, 2)

	assert.Equal(t, symbolTable.includeDefinitions["bar"], []sitter.Range{{
		StartPoint: sitter.Point{
			Row:    5,
			Column: 0,
		},
		EndPoint: sitter.Point{
			Row:    7,
			Column: 10,
		},
		StartByte: 56,
		EndByte:   110,
	}})
}

func TestSymbolTableForValues(t *testing.T) {
	content := `
{{ with .Values.with.something }}
{{ .test2 }}
{{ end }}

{{ .Test }}
{{ .Values.with.something }}

{{ range .list }}
	{{ . }}
	{{ .listinner }}
	{{ $.dollar }}
	{{ range .nested }}
		{{ .nestedinList }}
	{{ end }}
	{{ range $.Values.dollar }}
		{{ .nestedinList }}
	{{ end }}
{{ end }}

{{ .Test }}
`

	ast := ParseAst(nil, content)

	symbolTable := NewSymbolTable(ast, []byte(content))
	type expectedValue struct {
		path       []string
		startPoint sitter.Point
	}

	expected := []expectedValue{
		{
			path: []string{"Test"},
			startPoint: sitter.Point{
				Row:    5,
				Column: 4,
			},
		},
		{
			path: []string{"Test"},
			startPoint: sitter.Point{
				Row:    20,
				Column: 4,
			},
		},
		{
			path: []string{"Values", "with", "something"},
			startPoint: sitter.Point{
				Row:    1,
				Column: 21,
			},
		},
		{
			path: []string{"Values", "with", "something"},
			startPoint: sitter.Point{
				Row:    6,
				Column: 16,
			},
		},
		{
			path: []string{"list"},
			startPoint: sitter.Point{
				Row:    8,
				Column: 10,
			},
		},
		{
			path: []string{"list[]"},
			startPoint: sitter.Point{
				Row:    9,
				Column: 4,
			},
		},
		{
			path: []string{"list[]", "listinner"},
			startPoint: sitter.Point{
				Row:    10,
				Column: 5,
			},
		},
		{
			path: []string{"dollar"},
			startPoint: sitter.Point{
				Row:    11,
				Column: 6,
			},
		},
		{
			path: []string{"list[]", "nested"},
			startPoint: sitter.Point{
				Row:    12,
				Column: 11,
			},
		},
		{
			path: []string{"list[]", "nested[]", "nestedinList"},
			startPoint: sitter.Point{
				Row:    13,
				Column: 6,
			},
		},
		{
			path: []string{"Values", "dollar"},
			startPoint: sitter.Point{
				Row:    15,
				Column: 19,
			},
		},
		{
			path: []string{"Values", "dollar[]", "nestedinList"},
			startPoint: sitter.Point{
				Row:    16,
				Column: 6,
			},
		},
	}

	for _, v := range expected {
		values := symbolTable.GetTemplateContextRanges(v.path)
		points := []sitter.Point{}
		for _, v := range values {
			points = append(points, v.StartPoint)
		}
		assert.Contains(t, points, v.startPoint)
	}
}

func TestSymbolTableForValuesTestFile(t *testing.T) {
	path := "../../testdata/example/templates/deployment.yaml"

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal("Could not read test file", err)
	}
	ast := ParseAst(nil, string(content))

	symbolTable := NewSymbolTable(ast, []byte(content))
	type expectedValue struct {
		path       []string
		startPoint sitter.Point
	}

	expected := []expectedValue{
		{
			path: []string{"Values", "ingress"},
			startPoint: sitter.Point{
				Row:    0x47,
				Column: 0x18,
			},
		},
	}

	for _, v := range expected {
		values := symbolTable.GetTemplateContextRanges(v.path)
		points := []sitter.Point{}
		for _, v := range values {
			points = append(points, v.StartPoint)
		}
		assert.Contains(t, points, v.startPoint)
	}
}
