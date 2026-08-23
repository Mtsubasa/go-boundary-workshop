// Package boundary provides the analyzer built during the workshop.
package boundary

import (
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "check whether table-driven tests cover integer boundary values"

// Analyzer is the analyzer completed step by step in the workshop.
var Analyzer = &analysis.Analyzer{
	Name: "boundary",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !containsTestFile(pass) {
		return nil, nil
	}

	// TODO(Ex3): Find function declarations and report their names.
	_ = inspectResult

	return nil, nil
}

// containsTestFile and isTestFile are provided boilerplate. The analysis
// driver handles a package both without and with its test files. Workshop code
// compares production and test code in the variant that contains both.
func containsTestFile(pass *analysis.Pass) bool {
	for _, file := range pass.Files {
		if isTestFile(pass, file.Pos()) {
			return true
		}
	}
	return false
}

func isTestFile(pass *analysis.Pass, pos token.Pos) bool {
	filename := pass.Fset.PositionFor(pos, false).Filename
	return strings.HasSuffix(filepath.Base(filename), "_test.go")
}
