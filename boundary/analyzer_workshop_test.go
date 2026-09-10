//go:build workshop_solution

package boundary

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// TestAnalyzer is excluded from the preflight test by the build tag above.
// Participants enable it as an optional exercise after completing Step 3.
func TestAnalyzer(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "missing_boundary",
			want: []string{"ShippingFee: boundary value 5000 is not tested"},
		},
		{
			name: "complete",
			want: nil,
		},
		{
			name: "missing_left",
			want: []string{"ShippingFee: no test value less than 5000"},
		},
		{
			name: "missing_right",
			want: []string{"ShippingFee: no test value greater than 5000"},
		},
		{
			name: "ignore_want",
			want: []string{"ShippingFee: boundary value 5000 is not tested"},
		},
		{
			name: "separate_functions",
			want: []string{"ShippingFee: boundary value 5000 is not tested"},
		},
		{
			name: "unsupported",
			want: nil,
		},
		{
			name: "named_constant_missing",
			want: []string{"ShippingFee: boundary value 5000 is not tested"},
		},
		{
			name: "named_constant_complete",
			want: nil,
		},
		{
			name: "constant_expression",
			want: []string{"ShippingFee: boundary value 5000 is not tested"},
		},
		{
			name: "runtime_variable",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runOnFixture(t, tt.name)
			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("diagnostics = %q, want %q", got, tt.want)
			}
		})
	}
}

func runOnFixture(t *testing.T, fixtureName string) []string {
	t.Helper()

	directory := filepath.Join("testdata", fixtureName)
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read fixture directory: %v", err)
	}

	fileSet := token.NewFileSet()
	var files []*ast.File
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
			continue
		}

		filename := filepath.Join(directory, entry.Name())
		file, err := parser.ParseFile(fileSet, filename, nil, parser.ParseComments)
		if err != nil {
			t.Fatalf("parse %s: %v", filename, err)
		}
		files = append(files, file)
	}

	typesInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	typesConfig := &types.Config{Importer: importer.Default()}
	typesPackage, err := typesConfig.Check(
		"example.com/shipping",
		fileSet,
		files,
		typesInfo,
	)
	if err != nil {
		t.Fatalf("type-check fixture: %v", err)
	}

	inspectResult := inspector.New(files)
	var diagnostics []string
	pass := &analysis.Pass{
		Analyzer:  Analyzer,
		Fset:      fileSet,
		Files:     files,
		Pkg:       typesPackage,
		TypesInfo: typesInfo,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: inspectResult,
		},
		Report: func(diagnostic analysis.Diagnostic) {
			diagnostics = append(diagnostics, diagnostic.Message)
		},
	}

	if _, err := Analyzer.Run(pass); err != nil {
		t.Fatalf("run analyzer: %v", err)
	}
	return diagnostics
}
