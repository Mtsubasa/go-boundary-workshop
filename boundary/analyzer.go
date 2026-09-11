// Package boundaryは、ワークショップで作成するanalyzerを提供します。
package boundary

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "check whether table-driven tests cover integer boundary values"

// AnalyzerへStepごとに解析処理を追加します。
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

	// TODO(Step 1): ZennのStep 1にあるPreorderをここへ追加する。
	// function := node.(*ast.XXXX)のXXXXを書き換えて完成させる。
	_ = inspectResult

	return nil, nil
}
