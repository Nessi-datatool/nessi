package visualization

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIRenderer_RenderTable(t *testing.T) {
	buf := &bytes.Buffer{}
	renderer := NewCLIRenderer(buf)

	headers := []string{"Name", "Type", "Description"}
	rows := [][]string{
		{"table1", "delta", "First table"},
		{"table2", "parquet", "Second table"},
	}

	err := renderer.RenderTable(headers, rows)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "Type")
	assert.Contains(t, output, "Description")
	assert.Contains(t, output, "table1")
	assert.Contains(t, output, "delta")
	assert.Contains(t, output, "First table")
	assert.Contains(t, output, "table2")
	assert.Contains(t, output, "parquet")
	assert.Contains(t, output, "Second table")
}

func TestCLIRenderer_RenderGraph(t *testing.T) {
	buf := &bytes.Buffer{}
	renderer := NewCLIRenderer(buf)

	nodes := []string{"node1", "node2", "node3"}
	edges := [][2]string{
		{"node1", "node2"},
		{"node2", "node3"},
	}

	err := renderer.RenderGraph(nodes, edges)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Graph Nodes:")
	assert.Contains(t, output, "- node1")
	assert.Contains(t, output, "- node2")
	assert.Contains(t, output, "- node3")
	assert.Contains(t, output, "Graph Edges:")
	assert.Contains(t, output, "- node1 -> node2")
	assert.Contains(t, output, "- node2 -> node3")
}

func TestCLIRenderer_RenderJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	renderer := NewCLIRenderer(buf)

	jsonData := []byte(`{"name":"test","value":123}`)

	err := renderer.RenderJSON(jsonData)
	assert.NoError(t, err)

	output := buf.String()
	assert.Equal(t, string(jsonData), output)
}

func TestNewCLIRenderer_DefaultOutput(t *testing.T) {
	// Test that NewCLIRenderer uses os.Stdout when nil is passed
	renderer := NewCLIRenderer(nil)
	assert.NotNil(t, renderer.output)
}
