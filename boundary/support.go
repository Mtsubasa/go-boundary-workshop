package boundary

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// containsTestFileとisTestFileは、ワークショップでは変更しない補助関数です。
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

// intParameterNameは、対象とする関数のint型パラメータ名を返します。
func intParameterName(function *ast.FuncDecl) (string, bool) {
	if function.Recv != nil || function.Type.Params == nil {
		return "", false
	}
	if len(function.Type.Params.List) != 1 {
		return "", false
	}

	parameter := function.Type.Params.List[0]
	if len(parameter.Names) != 1 {
		return "", false
	}

	parameterType, ok := parameter.Type.(*ast.Ident)
	if !ok || parameterType.Name != "int" {
		return "", false
	}

	return parameter.Names[0].Name, true
}

// integerLiteralValueは、リテラルの文字列を比較可能なint64へ変換します。
func integerLiteralValue(literal *ast.BasicLit) (int64, bool) {
	value, err := strconv.ParseInt(literal.Value, 0, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}
