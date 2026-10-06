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

package cli

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fieldTypeOfT extracts the type of the first field of the struct named T from
// the given source.
func fieldTypeOfT(t *testing.T, src string) (ast.Expr, typeDecls) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	require.NoError(t, err)

	decls := make(typeDecls)
	collectTypeDeclsInto(file, decls)

	var fieldType ast.Expr
	ast.Inspect(file, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok && ts.Name.Name == "T" {
			if st, ok := ts.Type.(*ast.StructType); ok {
				fieldType = st.Fields.List[0].Type
			}
			return false
		}
		return true
	})
	require.NotNil(t, fieldType)

	return fieldType, decls
}

// TestGetTypeName_NonComparableShapes verifies that getTypeName renders every
// type shape as valid Go source. Rendering a function, channel or inline struct
// type as "any", or an array as a slice, produces generated code that does not
// compile.
func TestGetTypeName_NonComparableShapes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "function type",
			code:     "type T struct { F func(int) string }",
			expected: "func(int) string",
		},
		{
			name:     "variadic function type with results",
			code:     "type T struct { F func(a int, rest ...string) (int, error) }",
			expected: "func(a int, rest ...string) (int, error)",
		},
		{
			name:     "fixed size array keeps its length",
			code:     "type T struct { F [3]int }",
			expected: "[3]int",
		},
		{
			name:     "array with constant length",
			code:     "const N = 4\ntype T struct { F [N]byte }",
			expected: "[N]byte",
		},
		{
			name:     "bidirectional channel",
			code:     "type T struct { F chan int }",
			expected: "chan int",
		},
		{
			name:     "receive only channel",
			code:     "type T struct { F <-chan int }",
			expected: "<-chan int",
		},
		{
			name:     "send only channel",
			code:     "type T struct { F chan<- int }",
			expected: "chan<- int",
		},
		{
			name:     "slice of functions",
			code:     "type T struct { F []func() error }",
			expected: "[]func() error",
		},
		{
			name:     "map with function values",
			code:     "type T struct { F map[string]func(int) int }",
			expected: "map[string]func(int) int",
		},
		{
			name:     "empty interface renders as any",
			code:     "type T struct { F interface{} }",
			expected: "any",
		},
		{
			name:     "non empty interface keeps its methods",
			code:     "type T struct { F interface{ Close() error } }",
			expected: "interface{Close() error}",
		},
		{
			name:     "inline struct",
			code:     "type T struct { F struct{ X []int } }",
			expected: "struct{X []int}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, _ := fieldTypeOfT(t, "package test\n"+tt.code)
			assert.Equal(t, tt.expected, collapseSpaces(getTypeName(fieldType)))
		})
	}
}

// TestIsComparableType_NamedNonStructTypes verifies that a named type whose
// underlying type is non-comparable (a slice, map or function) is recognized as
// non-comparable. The name alone carries no information, so resolving it
// through the package type map is the only way to get this right.
func TestIsComparableType_NamedNonStructTypes(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected bool
	}{
		{
			name:     "named slice type is not comparable",
			src:      "package test\ntype Names []string\ntype T struct { F Names }",
			expected: false,
		},
		{
			name:     "named map type is not comparable",
			src:      "package test\ntype Lookup map[string]int\ntype T struct { F Lookup }",
			expected: false,
		},
		{
			name:     "named func type is not comparable",
			src:      "package test\ntype Handler func(int) string\ntype T struct { F Handler }",
			expected: false,
		},
		{
			name:     "named type over a named slice type is not comparable",
			src:      "package test\ntype Names []string\ntype Alias Names\ntype T struct { F Alias }",
			expected: false,
		},
		{
			name:     "named array type is comparable",
			src:      "package test\ntype Triple [3]int\ntype T struct { F Triple }",
			expected: true,
		},
		{
			name:     "named array of a named slice type is not comparable",
			src:      "package test\ntype Names []string\ntype Pair [2]Names\ntype T struct { F Pair }",
			expected: false,
		},
		{
			name:     "named basic type is comparable",
			src:      "package test\ntype ID string\ntype T struct { F ID }",
			expected: true,
		},
		{
			name:     "pointer to a named slice type is comparable",
			src:      "package test\ntype Names []string\ntype T struct { F *Names }",
			expected: true,
		},
		{
			name:     "self referential slice type terminates and is not comparable",
			src:      "package test\ntype Tree []Tree\ntype T struct { F Tree }",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, decls := fieldTypeOfT(t, tt.src)
			assert.Equal(t, tt.expected, isComparableType(fieldType, map[string]string{}, decls))
		})
	}
}

// TestIsStrictlyComparableType verifies that interface typed fields are
// reported as comparable but not strictly comparable: == compiles for them, so
// they can be used with option.FromZero, which only ever compares against the
// zero value, but it panics when two values of the same non-comparable dynamic
// type are compared, so lens.MakeLensStrict must not be used for them.
func TestIsStrictlyComparableType(t *testing.T) {
	tests := []struct {
		name            string
		src             string
		comparable      bool
		strictlyCompare bool
	}{
		{
			name:            "string",
			src:             "package test\ntype T struct { F string }",
			comparable:      true,
			strictlyCompare: true,
		},
		{
			name:            "any",
			src:             "package test\ntype T struct { F any }",
			comparable:      true,
			strictlyCompare: false,
		},
		{
			name:            "interface literal",
			src:             "package test\ntype T struct { F interface{} }",
			comparable:      true,
			strictlyCompare: false,
		},
		{
			name:            "error",
			src:             "package test\ntype T struct { F error }",
			comparable:      true,
			strictlyCompare: false,
		},
		{
			name:            "named interface type",
			src:             "package test\ntype Closer interface{ Close() error }\ntype T struct { F Closer }",
			comparable:      true,
			strictlyCompare: false,
		},
		{
			name:            "struct containing an interface field",
			src:             "package test\ntype Inner struct { Payload any }\ntype T struct { F Inner }",
			comparable:      true,
			strictlyCompare: false,
		},
		{
			name:            "array of any",
			src:             "package test\ntype T struct { F [2]any }",
			comparable:      true,
			strictlyCompare: false,
		},
		{
			name:            "pointer stays strictly comparable",
			src:             "package test\ntype T struct { F *[]int }",
			comparable:      true,
			strictlyCompare: true,
		},
		{
			name:            "slice",
			src:             "package test\ntype T struct { F []string }",
			comparable:      false,
			strictlyCompare: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, decls := fieldTypeOfT(t, tt.src)
			assert.Equal(t, tt.comparable, isComparableType(fieldType, map[string]string{}, decls), "isComparableType")
			assert.Equal(t, tt.strictlyCompare, isStrictlyComparableType(fieldType, map[string]string{}, decls), "isStrictlyComparableType")
		})
	}
}

// TestGenerateLensHelpers_NonComparableFields is an end to end check that a
// struct whose fields are all non-comparable produces a file that compiles: no
// optional lenses, hence no unused imports, and every field type rendered
// verbatim.
func TestGenerateLensHelpers_NonComparableFields(t *testing.T) {
	tmpDir := t.TempDir()

	testCode := "package testpkg\n" +
		"\n" +
		"type Names []string\n" +
		"\n" +
		"// fp-go:Lens\n" +
		"type AllNonComparable struct {\n" +
		"\tTags     []string\n" +
		"\tMeta     map[string]int\n" +
		"\tCallback func(int) string\n" +
		"\tNamed    Names\n" +
		"}\n"

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "types.go"), []byte(testCode), 0o644))
	require.NoError(t, generateLensHelpers(tmpDir, "gen_lens.go", false, false))

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen_lens.go"))
	require.NoError(t, err)
	contentStr := collapseSpaces(string(content))

	// The generated file must be valid Go
	_, err = parser.ParseFile(token.NewFileSet(), "gen_lens.go", content, 0)
	require.NoError(t, err, "generated code must parse")

	// No optional lens for any field, so the option packages must not be imported
	assert.NotContains(t, contentStr, "__lens_option")
	assert.NotContains(t, contentStr, "__iso_option")

	// Field types are rendered verbatim
	assert.Contains(t, contentStr, "Tags __lens.Lens[AllNonComparable, []string]")
	assert.Contains(t, contentStr, "Meta __lens.Lens[AllNonComparable, map[string]int]")
	assert.Contains(t, contentStr, "Callback __lens.Lens[AllNonComparable, func(int) string]")
	assert.Contains(t, contentStr, "Named __lens.Lens[AllNonComparable, Names]")

	// Reference lenses must not compare non-comparable values
	assert.NotContains(t, contentStr, "MakeLensStrictWithName")
	assert.Contains(t, contentStr, "MakeLensRefWithName")

	// Prisms use Some instead of FromNonZero for non-comparable fields
	assert.Contains(t, contentStr, "__option.Some(s.Tags)")
	assert.NotContains(t, contentStr, "__option.FromNonZero")
}

// TestGenerateLensHelpers_MixedComparability checks a struct that mixes
// comparable and non-comparable fields: only the comparable ones get an
// optional lens, and only the strictly comparable ones get a strict reference
// lens.
func TestGenerateLensHelpers_MixedComparability(t *testing.T) {
	tmpDir := t.TempDir()

	testCode := "package testpkg\n" +
		"\n" +
		"// fp-go:Lens\n" +
		"type Mixed struct {\n" +
		"\tName   string\n" +
		"\tTags   []string\n" +
		"\tMatrix [3]int\n" +
		"\tCh     chan int\n" +
		"\tErr    error\n" +
		"}\n"

	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "types.go"), []byte(testCode), 0o644))
	require.NoError(t, generateLensHelpers(tmpDir, "gen_lens.go", false, false))

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen_lens.go"))
	require.NoError(t, err)
	contentStr := collapseSpaces(string(content))

	_, err = parser.ParseFile(token.NewFileSet(), "gen_lens.go", content, 0)
	require.NoError(t, err, "generated code must parse")

	// Optional lenses exist for the comparable fields only
	assert.Contains(t, contentStr, "NameO __lens_option.LensO[Mixed, string]")
	assert.Contains(t, contentStr, "MatrixO __lens_option.LensO[Mixed, [3]int]")
	assert.Contains(t, contentStr, "ChO __lens_option.LensO[Mixed, chan int]")
	assert.Contains(t, contentStr, "ErrO __lens_option.LensO[Mixed, error]")
	assert.NotContains(t, contentStr, "TagsO")

	// Strict reference lenses only for the strictly comparable fields
	assert.Contains(t, contentStr, "lensName := __lens.MakeLensStrictWithName(")
	assert.Contains(t, contentStr, "lensMatrix := __lens.MakeLensStrictWithName(")
	assert.Contains(t, contentStr, "lensCh := __lens.MakeLensStrictWithName(")
	// error is an interface: == compiles but can panic
	assert.Contains(t, contentStr, "lensErr := __lens.MakeLensRefWithName(")
	assert.Contains(t, contentStr, "lensTags := __lens.MakeLensRefWithName(")
}

// TestGenerateLensHelpers_Deterministic verifies that generating twice produces
// byte identical output. Collected imports live in a map, so they have to be
// sorted before being written.
func TestGenerateLensHelpers_Deterministic(t *testing.T) {
	testCode := "package testpkg\n" +
		"\n" +
		"import (\n" +
		"\t\"time\"\n" +
		"\n" +
		"\t\"github.com/IBM/fp-go/v2/option\"\n" +
		")\n" +
		"\n" +
		"// fp-go:Lens\n" +
		"type Many struct {\n" +
		"\tWhen time.Time\n" +
		"\tOpt  option.Option[string]\n" +
		"\tDur  time.Duration\n" +
		"}\n"

	var first string
	for i := range 5 {
		tmpDir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "types.go"), []byte(testCode), 0o644))
		require.NoError(t, generateLensHelpers(tmpDir, "gen_lens.go", false, false))

		content, err := os.ReadFile(filepath.Join(tmpDir, "gen_lens.go"))
		require.NoError(t, err)

		if i == 0 {
			first = string(content)
			continue
		}
		assert.Equal(t, first, string(content), "generated output must be deterministic")
	}
}
