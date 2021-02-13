package cmd

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func BenchmarkRootCommandExecution(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rootCmd.RunE(nil, []string{filepath.Join("..", "testdata", "values.yaml")})
	}
}

func TestUncommentYAML(t *testing.T) {
	input, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "values.yaml"))
	uncommentedInput, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "values.uncommented.yaml"))

	root := yaml.Node{}
	if err := yaml.Unmarshal(input, &root); err != nil {
		t.Fatal(err)
	}
	uncommented := uncommentYAML(&root)

	var actual map[string]interface{}
	if err := uncommented.Decode(&actual); err != nil {
		t.Fatal(err)
	}

	var expected map[string]interface{}
	if err := yaml.Unmarshal(uncommentedInput, &expected); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}
