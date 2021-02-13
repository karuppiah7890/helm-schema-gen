package cmd

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkRootCommandExecution(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rootCmd.RunE(nil, []string{filepath.Join("..", "testdata", "values.yaml")})
	}
}

func TestUncommentYAML_Pattern1(t *testing.T) {
	input := `
foo: {}
  # bar: 1
  # baz: 2
  # qux: 3
`
	expected := `
foo:
  bar: 1
  baz: 2
  qux: 3
`

	root := yaml.Node{}
	if err := yaml.Unmarshal([]byte(input), &root); err != nil {
		t.Fatal(err)
	}
	uncommented := uncommentYAML(&root)

	buf := &bytes.Buffer{}
	enc := yaml.NewEncoder(buf)
	enc.SetIndent(2)
	if err := enc.Encode(uncommented); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(strings.TrimSpace(buf.String()), strings.TrimSpace(expected)); diff != "" {
		t.Error(diff)
	}
}
