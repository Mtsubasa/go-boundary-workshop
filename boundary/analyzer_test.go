package boundary

import (
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestAnalyzerDefinition(t *testing.T) {
	if err := analysis.Validate([]*analysis.Analyzer{Analyzer}); err != nil {
		t.Fatalf("invalid analyzer: %v", err)
	}
}
