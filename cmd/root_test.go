package cmd

import (
	"path/filepath"
	"testing"
)

func BenchmarkRootCommandExecution(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rootCmd.RunE(nil, []string{filepath.Join("..", "testdata", "values.yaml")})
	}
}
