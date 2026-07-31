// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lenses_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	S "github.com/IBM/fp-go/v2/string"
	"github.com/IBM/fp-go/gen/v2/cli/lenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// parseFieldType is a test helper that parses a struct definition and returns
// the AST expression for the first named field.
func parseFieldType(t *testing.T, code string) ast.Expr {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", "package test\n"+code, 0)
	require.NoError(t, err)

	var fieldType ast.Expr
	ast.Inspect(file, func(n ast.Node) bool {
		if field, ok := n.(*ast.Field); ok && len(field.Names) > 0 {
			fieldType = field.Type
			return false
		}
		return true
	})
	require.NotNil(t, fieldType, "expected to find a named field in: %s", code)
	return fieldType
}

// ---------------------------------------------------------------------------
// HasLensAnnotation
// ---------------------------------------------------------------------------

func TestHasLensAnnotation(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected bool
	}{
		{"has annotation", "// fp-go:Lens", true},
		{"has annotation with other text", "// This is a struct with fp-go:Lens annotation", true},
		{"no annotation", "// This is just a regular comment", false},
		{"nil comment group", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc *ast.CommentGroup
			if S.IsNonEmpty(tt.comment) {
				doc = &ast.CommentGroup{
					List: []*ast.Comment{{Text: tt.comment}},
				}
			}
			assert.Equal(t, tt.expected, lenses.HasLensAnnotation(doc))
		})
	}
}

// ---------------------------------------------------------------------------
// GetTypeName
// ---------------------------------------------------------------------------

func TestGetTypeName(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"simple type", "type T struct { F string }", "string"},
		{"pointer type", "type T struct { F *string }", "*string"},
		{"slice type", "type T struct { F []int }", "[]int"},
		{"map type", "type T struct { F map[string]int }", "map[string]int"},
		{"qualified type", "type T struct { F pkg.Foo }", "pkg.Foo"},
		{"generic single param", "type T struct { F option.Option[string] }", "option.Option[string]"},
		{"generic multi param", "type T struct { F either.Either[error, string] }", "either.Either[error, string]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseFieldType(t, tt.code)
			assert.Equal(t, tt.expected, lenses.GetTypeName(expr))
		})
	}
}

// ---------------------------------------------------------------------------
// IsPointerType
// ---------------------------------------------------------------------------

func TestIsPointerType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"pointer type", "type T struct { F *string }", true},
		{"non-pointer type", "type T struct { F string }", false},
		{"slice type", "type T struct { F []string }", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseFieldType(t, tt.code)
			assert.Equal(t, tt.expected, lenses.IsPointerType(expr))
		})
	}
}

// ---------------------------------------------------------------------------
// IsOptionType
// ---------------------------------------------------------------------------

func TestIsOptionType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"option.Option[string]", "type T struct { F option.Option[string] }", true},
		{"option.Option[int]", "type T struct { F option.Option[int] }", true},
		{"plain string - not an Option", "type T struct { F string }", false},
		{"pair.Pair[string, int] - not an Option", "type T struct { F pair.Pair[string, int] }", false},
		{"other.Option[string] - any pkg named 'Option' is treated as Option", "type T struct { F other.Option[string] }", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseFieldType(t, tt.code)
			assert.Equal(t, tt.expected, lenses.IsOptionType(expr))
		})
	}
}

// ---------------------------------------------------------------------------
// IsComparableType
// ---------------------------------------------------------------------------

func TestIsComparableType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"basic type - string", "type T struct { F string }", true},
		{"basic type - int", "type T struct { F int }", true},
		{"basic type - bool", "type T struct { F bool }", true},
		{"pointer type", "type T struct { F *string }", true},
		{"slice type - not comparable", "type T struct { F []string }", false},
		{"map type - not comparable", "type T struct { F map[string]int }", false},
		{"array type - comparable if element is", "type T struct { F [5]int }", true},
		{"interface type", "type T struct { F any }", true},
		{"channel type", "type T struct { F chan int }", true},
		{"function type - not comparable", "type T struct { F func() }", false},
		{"struct literal - conservatively not comparable", "type T struct { F struct{ X int } }", false},
		{"qualified type (pkg.Type) - conservatively not comparable", "type T struct { F pkg.SomeStruct }", false},
		{"context.Context - interface, comparable", "type T struct { F context.Context }", true},
		{"option.Option[string] - comparable type arg", "type T struct { F option.Option[string] }", true},
		{"option.Option[[]byte] - non-comparable type arg", "type T struct { F option.Option[[]byte] }", false},
		{"either.Either[error, string] - both args comparable", "type T struct { F either.Either[error, string] }", true},
		{"either.Either[error, []string] - non-comparable second arg", "type T struct { F either.Either[error, []string] }", false},
		{"pair.Pair[string, int] - both args comparable", "type T struct { F pair.Pair[string, int] }", true},
		{"pair.Pair[string, []int] - non-comparable second arg", "type T struct { F pair.Pair[string, []int] }", false},
		{"tuple.Tuple2[string, int] - both args comparable", "type T struct { F tuple.Tuple2[string, int] }", true},
		{"tuple.Tuple2[string, []int] - non-comparable second arg", "type T struct { F tuple.Tuple2[string, []int] }", false},
		{"result.Result[string] - comparable type arg", "type T struct { F result.Result[string] }", true},
		{"result.Result[[]byte] - non-comparable type arg", "type T struct { F result.Result[[]byte] }", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseFieldType(t, tt.code)
			assert.Equal(t, tt.expected, lenses.IsComparableType(expr, map[string]string{}))
		})
	}
}

// TestIsComparableTypeWithTypeParams verifies that generic type parameter
// constraints are respected.
func TestIsComparableTypeWithTypeParams(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		typeParams  map[string]string
		expected    bool
	}{
		{
			name:       "T any - not comparable",
			code:       "type T struct { F T }",
			typeParams: map[string]string{"T": "any"},
			expected:   false,
		},
		{
			name:       "T comparable - comparable",
			code:       "type T struct { F T }",
			typeParams: map[string]string{"T": "comparable"},
			expected:   true,
		},
		{
			name:       "K comparable, V any - K comparable",
			code:       "type T struct { F K }",
			typeParams: map[string]string{"K": "comparable", "V": "any"},
			expected:   true,
		},
		{
			name:       "K comparable, V any - V not comparable",
			code:       "type T struct { F V }",
			typeParams: map[string]string{"K": "comparable", "V": "any"},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := parseFieldType(t, tt.code)
			assert.Equal(t, tt.expected, lenses.IsComparableType(expr, tt.typeParams))
		})
	}
}

// ---------------------------------------------------------------------------
// HasOmitEmpty
// ---------------------------------------------------------------------------

func TestHasOmitEmpty(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected bool
	}{
		{"has omitempty", "`json:\"field,omitempty\"`", true},
		{"has omitempty with other options", "`json:\"field,omitempty,string\"`", true},
		{"no omitempty", "`json:\"field\"`", false},
		{"no tag", "", false},
		{"different tag", "`xml:\"field\"`", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tag *ast.BasicLit
			if S.IsNonEmpty(tt.tag) {
				tag = &ast.BasicLit{Value: tt.tag}
			}
			assert.Equal(t, tt.expected, lenses.HasOmitEmpty(tag))
		})
	}
}

// ---------------------------------------------------------------------------
// IsFpGoValueStructPkg
// ---------------------------------------------------------------------------

func TestIsFpGoValueStructPkg(t *testing.T) {
	tests := []struct {
		pkgName  string
		typeName string
		expected bool
	}{
		{"option", "Option", true},
		{"either", "Either", true},
		{"pair", "Pair", true},
		{"result", "Result", true},
		{"tuple", "Tuple2", true},
		{"tuple", "Tuple15", true},
		{"option", "Either", false},
		{"unknown", "Option", false},
		{"option", "option", false},
	}

	for _, tt := range tests {
		t.Run(tt.pkgName+"."+tt.typeName, func(t *testing.T) {
			assert.Equal(t, tt.expected, lenses.IsFpGoValueStructPkg(tt.pkgName, tt.typeName))
		})
	}
}
