// Package boundary provides an analyzer that checks whether table-driven tests
// contain values on both sides of an integer boundary and the boundary itself.
package boundary

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "check whether table-driven tests cover integer boundary values"

// Analyzer checks the deliberately narrow syntax used by the workshop's Core
// exercise. See the repository README for the supported syntax.
var Analyzer = &analysis.Analyzer{
	Name: "boundary",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

type boundaryInfo struct {
	functionName string
	value        int64
	pos          token.Pos
}

func run(pass *analysis.Pass) (any, error) {
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	testInputs := collectTestInputs(pass, inspectResult)
	if len(testInputs) == 0 {
		// The analysis driver also runs this analyzer for a package variant that
		// does not contain test files. Only the variant containing both production
		// and test files has enough information to compare their values.
		return nil, nil
	}

	for _, boundary := range collectBoundaries(pass, inspectResult) {
		inputs, hasTest := testInputs[boundary.functionName]
		if !hasTest {
			continue
		}

		hasLess, hasBoundary, hasGreater := classifyInputs(inputs, boundary.value)
		if !hasLess {
			pass.Reportf(
				boundary.pos,
				"%s: no test value less than %d",
				boundary.functionName,
				boundary.value,
			)
		}
		if !hasBoundary {
			pass.Reportf(
				boundary.pos,
				"%s: boundary value %d is not tested",
				boundary.functionName,
				boundary.value,
			)
		}
		if !hasGreater {
			pass.Reportf(
				boundary.pos,
				"%s: no test value greater than %d",
				boundary.functionName,
				boundary.value,
			)
		}
	}

	return nil, nil
}

func collectBoundaries(pass *analysis.Pass, inspectResult *inspector.Inspector) []boundaryInfo {
	var boundaries []boundaryInfo
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	inspectResult.Preorder(nodeFilter, func(node ast.Node) {
		function := node.(*ast.FuncDecl)
		if isTestFile(pass, function.Pos()) || function.Body == nil {
			return
		}

		parameterName, ok := coreParameterName(function)
		if !ok {
			return
		}

		ast.Inspect(function.Body, func(node ast.Node) bool {
			if _, ok := node.(*ast.FuncLit); ok {
				// A condition in a nested function literal does not belong to the
				// surrounding function.
				return false
			}

			ifStatement, ok := node.(*ast.IfStmt)
			if !ok {
				return true
			}

			comparison, ok := ifStatement.Cond.(*ast.BinaryExpr)
			if !ok || comparison.Op != token.LSS {
				return true
			}

			left, ok := comparison.X.(*ast.Ident)
			if !ok || left.Name != parameterName {
				return true
			}

			literal, ok := comparison.Y.(*ast.BasicLit)
			if !ok || literal.Kind != token.INT {
				return true
			}

			value, ok := integerLiteralValue(literal)
			if !ok {
				return true
			}

			boundaries = append(boundaries, boundaryInfo{
				functionName: function.Name.Name,
				value:        value,
				pos:          literal.Pos(),
			})
			return true
		})
	})

	return boundaries
}

func collectTestInputs(pass *analysis.Pass, inspectResult *inspector.Inspector) map[string][]int64 {
	testInputs := make(map[string][]int64)
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	inspectResult.Preorder(nodeFilter, func(node ast.Node) {
		function := node.(*ast.FuncDecl)
		if !isTestFile(pass, function.Pos()) || function.Body == nil {
			return
		}
		if function.Recv != nil || !strings.HasPrefix(function.Name.Name, "Test") {
			return
		}

		productionFunctionName := strings.TrimPrefix(function.Name.Name, "Test")
		if productionFunctionName == "" {
			return
		}

		// Keep an entry even when no supported input literal is found. A matching
		// test function with an empty/unsupported table is still missing all
		// three classes of values.
		if _, ok := testInputs[productionFunctionName]; !ok {
			testInputs[productionFunctionName] = nil
		}

		ast.Inspect(function.Body, func(node ast.Node) bool {
			keyValue, ok := node.(*ast.KeyValueExpr)
			if !ok {
				return true
			}

			key, ok := keyValue.Key.(*ast.Ident)
			if !ok || key.Name != "input" {
				return true
			}

			literal, ok := keyValue.Value.(*ast.BasicLit)
			if !ok || literal.Kind != token.INT {
				return true
			}

			value, ok := integerLiteralValue(literal)
			if !ok {
				return true
			}

			testInputs[productionFunctionName] = append(
				testInputs[productionFunctionName],
				value,
			)
			return true
		})
	})

	return testInputs
}

func coreParameterName(function *ast.FuncDecl) (string, bool) {
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

func integerLiteralValue(literal *ast.BasicLit) (int64, bool) {
	value, err := strconv.ParseInt(literal.Value, 0, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func classifyInputs(inputs []int64, boundary int64) (less, equal, greater bool) {
	for _, input := range inputs {
		switch {
		case input < boundary:
			less = true
		case input == boundary:
			equal = true
		case input > boundary:
			greater = true
		}
	}
	return less, equal, greater
}

func isTestFile(pass *analysis.Pass, pos token.Pos) bool {
	filename := pass.Fset.PositionFor(pos, false).Filename
	return strings.HasSuffix(filepath.Base(filename), "_test.go")
}
