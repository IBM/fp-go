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
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasLensAnnotation(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected bool
	}{
		{
			name:     "has annotation",
			comment:  "// fp-go:Lens",
			expected: true,
		},
		{
			name:     "has annotation with other text",
			comment:  "// This is a struct with fp-go:Lens annotation",
			expected: true,
		},
		{
			name:     "no annotation",
			comment:  "// This is just a regular comment",
			expected: false,
		},
		{
			name:     "nil comment",
			comment:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc *ast.CommentGroup
			if tt.comment != "" {
				doc = &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: tt.comment},
					},
				}
			}
			result := hasLensAnnotation(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTypeName(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "simple type",
			code:     "type T struct { F string }",
			expected: "string",
		},
		{
			name:     "pointer type",
			code:     "type T struct { F *string }",
			expected: "*string",
		},
		{
			name:     "slice type",
			code:     "type T struct { F []int }",
			expected: "[]int",
		},
		{
			name:     "map type",
			code:     "type T struct { F map[string]int }",
			expected: "map[string]int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			require.NoError(t, err)

			var fieldType ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				if field, ok := n.(*ast.Field); ok && len(field.Names) > 0 {
					fieldType = field.Type
					return false
				}
				return true
			})

			require.NotNil(t, fieldType)
			result := getTypeName(fieldType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPointerType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "pointer type",
			code:     "type T struct { F *string }",
			expected: true,
		},
		{
			name:     "non-pointer type",
			code:     "type T struct { F string }",
			expected: false,
		},
		{
			name:     "slice type",
			code:     "type T struct { F []string }",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", "package test\n"+tt.code, 0)
			require.NoError(t, err)

			var fieldType ast.Expr
			ast.Inspect(file, func(n ast.Node) bool {
				if field, ok := n.(*ast.Field); ok && len(field.Names) > 0 {
					fieldType = field.Type
					return false
				}
				return true
			})

			require.NotNil(t, fieldType)
			result := isPointerType(fieldType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasOmitEmpty(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected bool
	}{
		{
			name:     "has omitempty",
			tag:      "`json:\"field,omitempty\"`",
			expected: true,
		},
		{
			name:     "has omitempty with other options",
			tag:      "`json:\"field,omitempty,string\"`",
			expected: true,
		},
		{
			name:     "no omitempty",
			tag:      "`json:\"field\"`",
			expected: false,
		},
		{
			name:     "no tag",
			tag:      "",
			expected: false,
		},
		{
			name:     "different tag",
			tag:      "`xml:\"field\"`",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tag *ast.BasicLit
			if tt.tag != "" {
				tag = &ast.BasicLit{
					Value: tt.tag,
				}
			}
			result := hasOmitEmpty(tag)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// fp-go:Lens
type Person struct {
	Name  string
	Age   int
	Phone *string
}

// fp-go:Lens
type Address struct {
	Street string
	City   string
}

// Not annotated
type Other struct {
	Field string
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 2)

	// Check Person struct
	person := structs[0]
	assert.Equal(t, "Person", person.Name)
	assert.Len(t, person.Fields, 3)

	assert.Equal(t, "Name", person.Fields[0].Name)
	assert.Equal(t, "string", person.Fields[0].TypeName)
	assert.False(t, person.Fields[0].IsOptional)

	assert.Equal(t, "Age", person.Fields[1].Name)
	assert.Equal(t, "int", person.Fields[1].TypeName)
	assert.False(t, person.Fields[1].IsOptional)

	assert.Equal(t, "Phone", person.Fields[2].Name)
	assert.Equal(t, "*string", person.Fields[2].TypeName)
	assert.True(t, person.Fields[2].IsOptional)

	// Check Address struct
	address := structs[1]
	assert.Equal(t, "Address", address.Name)
	assert.Len(t, address.Fields, 2)

	assert.Equal(t, "Street", address.Fields[0].Name)
	assert.Equal(t, "City", address.Fields[1].Name)
}

func TestGenerateLensHelpers(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testCode := `package testpkg

// fp-go:Lens
type TestStruct struct {
	Name  string
	Value *int
}
`

	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Generate lens code
	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false)
	require.NoError(t, err)

	// Verify the generated file exists
	genPath := filepath.Join(tmpDir, outputFile)
	_, err = os.Stat(genPath)
	require.NoError(t, err)

	// Read and verify the generated content
	content, err := os.ReadFile(genPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Check for expected content
	assert.Contains(t, contentStr, "package testpkg")
	assert.Contains(t, contentStr, "Code generated by go generate")
	assert.Contains(t, contentStr, "TestStructLens")
	assert.Contains(t, contentStr, "MakeTestStructLens")
	assert.Contains(t, contentStr, "L.Lens[TestStruct, string]")
	assert.Contains(t, contentStr, "LO.LensO[TestStruct, *int]")
	assert.Contains(t, contentStr, "O.FromNillable")
	assert.Contains(t, contentStr, "O.GetOrElse")
}

func TestGenerateLensHelpersNoAnnotations(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testCode := `package testpkg

// No annotation
type TestStruct struct {
	Name string
}
`

	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Generate lens code (should not create file)
	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false)
	require.NoError(t, err)

	// Verify the generated file does not exist
	genPath := filepath.Join(tmpDir, outputFile)
	_, err = os.Stat(genPath)
	assert.True(t, os.IsNotExist(err))
}

func TestLensTemplates(t *testing.T) {
	s := structInfo{
		Name: "TestStruct",
		Fields: []fieldInfo{
			{Name: "Name", TypeName: "string", IsOptional: false},
			{Name: "Value", TypeName: "*int", IsOptional: true},
		},
	}

	// Test struct template
	var structBuf bytes.Buffer
	err := structTmpl.Execute(&structBuf, s)
	require.NoError(t, err)

	structStr := structBuf.String()
	assert.Contains(t, structStr, "type TestStructLenses struct")
	assert.Contains(t, structStr, "Name L.Lens[TestStruct, string]")
	assert.Contains(t, structStr, "Value LO.LensO[TestStruct, *int]")

	// Test constructor template
	var constructorBuf bytes.Buffer
	err = constructorTmpl.Execute(&constructorBuf, s)
	require.NoError(t, err)

	constructorStr := constructorBuf.String()
	assert.Contains(t, constructorStr, "func MakeTestStructLenses() TestStructLenses")
	assert.Contains(t, constructorStr, "return TestStructLenses{")
	assert.Contains(t, constructorStr, "Name: L.MakeLens(")
	assert.Contains(t, constructorStr, "Value: L.MakeLens(")
	assert.Contains(t, constructorStr, "O.FromNillable")
	assert.Contains(t, constructorStr, "O.GetOrElse")
}

func TestLensCommandFlags(t *testing.T) {
	cmd := LensCommand()

	assert.Equal(t, "lens", cmd.Name)
	assert.Equal(t, "generate lens code for annotated structs", cmd.Usage)
	assert.Contains(t, strings.ToLower(cmd.Description), "fp-go:lens")
	assert.Contains(t, strings.ToLower(cmd.Description), "lenso")

	// Check flags
	assert.Len(t, cmd.Flags, 3)

	var hasDir, hasFilename, hasVerbose bool
	for _, flag := range cmd.Flags {
		switch flag.Names()[0] {
		case "dir":
			hasDir = true
		case "filename":
			hasFilename = true
		case "verbose":
			hasVerbose = true
		}
	}

	assert.True(t, hasDir, "should have dir flag")
	assert.True(t, hasFilename, "should have filename flag")
	assert.True(t, hasVerbose, "should have verbose flag")
}
