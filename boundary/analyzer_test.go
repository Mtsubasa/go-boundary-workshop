package boundary

import (
	"go/ast"
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

func TestClassifyInputs(t *testing.T) {
	tests := []struct {
		name        string
		inputs      []int64
		boundary    int64
		wantLess    bool
		wantEqual   bool
		wantGreater bool
	}{
		{
			name:        "all classes",
			inputs:      []int64{4999, 5000, 5001},
			boundary:    5000,
			wantLess:    true,
			wantEqual:   true,
			wantGreater: true,
		},
		{
			name:        "boundary missing",
			inputs:      []int64{4999, 5001},
			boundary:    5000,
			wantLess:    true,
			wantEqual:   false,
			wantGreater: true,
		},
		{
			name:        "no inputs",
			inputs:      nil,
			boundary:    5000,
			wantLess:    false,
			wantEqual:   false,
			wantGreater: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			less, equal, greater := classifyInputs(tt.inputs, tt.boundary)
			if less != tt.wantLess || equal != tt.wantEqual || greater != tt.wantGreater {
				t.Fatalf(
					"classifyInputs() = (%t, %t, %t), want (%t, %t, %t)",
					less,
					equal,
					greater,
					tt.wantLess,
					tt.wantEqual,
					tt.wantGreater,
				)
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

	inspectResult := inspector.New(files)
	var diagnostics []string
	pass := &analysis.Pass{
		Analyzer: Analyzer,
		Fset:     fileSet,
		Files:    files,
		Pkg:      types.NewPackage("example.com/shipping", "shipping"),
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
