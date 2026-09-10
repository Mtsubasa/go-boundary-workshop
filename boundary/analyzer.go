// Package boundaryは、ワークショップで作成するanalyzerを提供します。
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

	// TODO(Step 1): 関数宣言を探し、関数名を表示する。
	_ = inspectResult

	return nil, nil
}

// containsTestFileとisTestFileはあらかじめ用意した補助関数です。
// analysisの実行側は、通常のパッケージとテストを含むパッケージをそれぞれ扱います。
// 今回は、実装コードとテストコードを両方含む場合だけ解析を続けます。
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
