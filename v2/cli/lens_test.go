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

func TestIsComparableType(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "basic type - string",
			code:     "type T struct { F string }",
			expected: true,
		},
		{
			name:     "basic type - int",
			code:     "type T struct { F int }",
			expected: true,
		},
		{
			name:     "basic type - bool",
			code:     "type T struct { F bool }",
			expected: true,
		},
		{
			name:     "pointer type",
			code:     "type T struct { F *string }",
			expected: true,
		},
		{
			name:     "slice type - not comparable",
			code:     "type T struct { F []string }",
			expected: false,
		},
		{
			name:     "map type - not comparable",
			code:     "type T struct { F map[string]int }",
			expected: false,
		},
		{
			name:     "array type - comparable if element is",
			code:     "type T struct { F [5]int }",
			expected: true,
		},
		{
			name:     "interface type",
			code:     "type T struct { F interface{} }",
			expected: true,
		},
		{
			name:     "channel type",
			code:     "type T struct { F chan int }",
			expected: true,
		},
		{
			name:     "function type - not comparable",
			code:     "type T struct { F func() }",
			expected: false,
		},
		{
			name:     "struct literal - conservatively not comparable",
			code:     "type T struct { F struct{ X int } }",
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
			result := isComparableType(fieldType)
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

func TestParseFileWithOmitEmpty(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// fp-go:Lens
type Config struct {
	Name     string
	Value    string  ` + "`json:\"value,omitempty\"`" + `
	Count    int     ` + "`json:\",omitempty\"`" + `
	Optional *string ` + "`json:\"optional,omitempty\"`" + `
	Required int     ` + "`json:\"required\"`" + `
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 1)

	// Check Config struct
	config := structs[0]
	assert.Equal(t, "Config", config.Name)
	assert.Len(t, config.Fields, 5)

	// Name - no tag, not optional
	assert.Equal(t, "Name", config.Fields[0].Name)
	assert.Equal(t, "string", config.Fields[0].TypeName)
	assert.False(t, config.Fields[0].IsOptional)

	// Value - has omitempty, should be optional
	assert.Equal(t, "Value", config.Fields[1].Name)
	assert.Equal(t, "string", config.Fields[1].TypeName)
	assert.True(t, config.Fields[1].IsOptional, "Value field with omitempty should be optional")

	// Count - has omitempty (no field name in tag), should be optional
	assert.Equal(t, "Count", config.Fields[2].Name)
	assert.Equal(t, "int", config.Fields[2].TypeName)
	assert.True(t, config.Fields[2].IsOptional, "Count field with omitempty should be optional")

	// Optional - pointer with omitempty, should be optional
	assert.Equal(t, "Optional", config.Fields[3].Name)
	assert.Equal(t, "*string", config.Fields[3].TypeName)
	assert.True(t, config.Fields[3].IsOptional)

	// Required - has json tag but no omitempty, not optional
	assert.Equal(t, "Required", config.Fields[4].Name)
	assert.Equal(t, "int", config.Fields[4].TypeName)
	assert.False(t, config.Fields[4].IsOptional, "Required field without omitempty should not be optional")
}

func TestParseFileWithComparableTypes(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// fp-go:Lens
type TypeTest struct {
	Name      string
	Age       int
	Pointer   *string
	Slice     []string
	Map       map[string]int
	Channel   chan int
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 1)

	// Check TypeTest struct
	typeTest := structs[0]
	assert.Equal(t, "TypeTest", typeTest.Name)
	assert.Len(t, typeTest.Fields, 6)

	// Name - string is comparable
	assert.Equal(t, "Name", typeTest.Fields[0].Name)
	assert.Equal(t, "string", typeTest.Fields[0].TypeName)
	assert.False(t, typeTest.Fields[0].IsOptional)
	assert.True(t, typeTest.Fields[0].IsComparable, "string should be comparable")

	// Age - int is comparable
	assert.Equal(t, "Age", typeTest.Fields[1].Name)
	assert.Equal(t, "int", typeTest.Fields[1].TypeName)
	assert.False(t, typeTest.Fields[1].IsOptional)
	assert.True(t, typeTest.Fields[1].IsComparable, "int should be comparable")

	// Pointer - pointer is optional, IsComparable not checked for optional fields
	assert.Equal(t, "Pointer", typeTest.Fields[2].Name)
	assert.Equal(t, "*string", typeTest.Fields[2].TypeName)
	assert.True(t, typeTest.Fields[2].IsOptional)

	// Slice - not comparable
	assert.Equal(t, "Slice", typeTest.Fields[3].Name)
	assert.Equal(t, "[]string", typeTest.Fields[3].TypeName)
	assert.False(t, typeTest.Fields[3].IsOptional)
	assert.False(t, typeTest.Fields[3].IsComparable, "slice should not be comparable")

	// Map - not comparable
	assert.Equal(t, "Map", typeTest.Fields[4].Name)
	assert.Equal(t, "map[string]int", typeTest.Fields[4].TypeName)
	assert.False(t, typeTest.Fields[4].IsOptional)
	assert.False(t, typeTest.Fields[4].IsComparable, "map should not be comparable")

	// Channel - comparable (note: getTypeName returns "any" for channel types, but isComparableType correctly identifies them)
	assert.Equal(t, "Channel", typeTest.Fields[5].Name)
	assert.Equal(t, "any", typeTest.Fields[5].TypeName) // getTypeName doesn't handle chan types specifically
	assert.False(t, typeTest.Fields[5].IsOptional)
	assert.True(t, typeTest.Fields[5].IsComparable, "channel should be comparable")
}

func TestLensRefTemplatesWithComparable(t *testing.T) {
	s := structInfo{
		Name: "TestStruct",
		Fields: []fieldInfo{
			{Name: "Name", TypeName: "string", IsOptional: false, IsComparable: true},
			{Name: "Age", TypeName: "int", IsOptional: false, IsComparable: true},
			{Name: "Data", TypeName: "[]byte", IsOptional: false, IsComparable: false},
			{Name: "Pointer", TypeName: "*string", IsOptional: true, IsComparable: false},
		},
	}

	// Test constructor template for RefLenses
	var constructorBuf bytes.Buffer
	err := constructorTmpl.Execute(&constructorBuf, s)
	require.NoError(t, err)

	constructorStr := constructorBuf.String()

	// Check that MakeLensStrict is used for comparable types in RefLenses
	assert.Contains(t, constructorStr, "func MakeTestStructRefLenses() TestStructRefLenses")

	// Name field - comparable, should use MakeLensStrict
	assert.Contains(t, constructorStr, "lensName := L.MakeLensStrict(",
		"comparable field Name should use MakeLensStrict in RefLenses")

	// Age field - comparable, should use MakeLensStrict
	assert.Contains(t, constructorStr, "lensAge := L.MakeLensStrict(",
		"comparable field Age should use MakeLensStrict in RefLenses")

	// Data field - not comparable, should use MakeLensRef
	assert.Contains(t, constructorStr, "lensData := L.MakeLensRef(",
		"non-comparable field Data should use MakeLensRef in RefLenses")

}

func TestGenerateLensHelpersWithComparable(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testCode := `package testpkg

// fp-go:Lens
type TestStruct struct {
	Name  string
	Count int
	Data  []byte
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

	// Check for expected content in RefLenses
	assert.Contains(t, contentStr, "MakeTestStructRefLenses")

	// Name and Count are comparable, should use MakeLensStrict
	assert.Contains(t, contentStr, "L.MakeLensStrict",
		"comparable fields should use MakeLensStrict in RefLenses")

	// Data is not comparable (slice), should use MakeLensRef
	assert.Contains(t, contentStr, "L.MakeLensRef",
		"non-comparable fields should use MakeLensRef in RefLenses")

	// Verify the pattern appears for Name field (comparable)
	namePattern := "lensName := L.MakeLensStrict("
	assert.Contains(t, contentStr, namePattern,
		"Name field should use MakeLensStrict")

	// Verify the pattern appears for Data field (not comparable)
	dataPattern := "lensData := L.MakeLensRef("
	assert.Contains(t, contentStr, dataPattern,
		"Data field should use MakeLensRef")
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
	assert.Contains(t, contentStr, "TestStructLenses")
	assert.Contains(t, contentStr, "MakeTestStructLenses")
	assert.Contains(t, contentStr, "L.Lens[TestStruct, string]")
	assert.Contains(t, contentStr, "LO.LensO[TestStruct, *int]")
	assert.Contains(t, contentStr, "IO.FromZero")
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
	assert.Contains(t, structStr, "NameO LO.LensO[TestStruct, string]")
	assert.Contains(t, structStr, "Value L.Lens[TestStruct, *int]")
	assert.Contains(t, structStr, "ValueO LO.LensO[TestStruct, *int]")

	// Test constructor template
	var constructorBuf bytes.Buffer
	err = constructorTmpl.Execute(&constructorBuf, s)
	require.NoError(t, err)

	constructorStr := constructorBuf.String()
	assert.Contains(t, constructorStr, "func MakeTestStructLenses() TestStructLenses")
	assert.Contains(t, constructorStr, "return TestStructLenses{")
	assert.Contains(t, constructorStr, "Name: lensName,")
	assert.Contains(t, constructorStr, "NameO: lensNameO,")
	assert.Contains(t, constructorStr, "Value: lensValue,")
	assert.Contains(t, constructorStr, "ValueO: lensValueO,")
	assert.Contains(t, constructorStr, "IO.FromZero")
}

func TestLensTemplatesWithOmitEmpty(t *testing.T) {
	s := structInfo{
		Name: "ConfigStruct",
		Fields: []fieldInfo{
			{Name: "Name", TypeName: "string", IsOptional: false},
			{Name: "Value", TypeName: "string", IsOptional: true},    // non-pointer with omitempty
			{Name: "Count", TypeName: "int", IsOptional: true},       // non-pointer with omitempty
			{Name: "Pointer", TypeName: "*string", IsOptional: true}, // pointer
		},
	}

	// Test struct template
	var structBuf bytes.Buffer
	err := structTmpl.Execute(&structBuf, s)
	require.NoError(t, err)

	structStr := structBuf.String()
	assert.Contains(t, structStr, "type ConfigStructLenses struct")
	assert.Contains(t, structStr, "Name L.Lens[ConfigStruct, string]")
	assert.Contains(t, structStr, "NameO LO.LensO[ConfigStruct, string]")
	assert.Contains(t, structStr, "Value L.Lens[ConfigStruct, string]")
	assert.Contains(t, structStr, "ValueO LO.LensO[ConfigStruct, string]", "non-pointer with omitempty should have optional lens")
	assert.Contains(t, structStr, "Count L.Lens[ConfigStruct, int]")
	assert.Contains(t, structStr, "CountO LO.LensO[ConfigStruct, int]", "non-pointer with omitempty should have optional lens")
	assert.Contains(t, structStr, "Pointer L.Lens[ConfigStruct, *string]")
	assert.Contains(t, structStr, "PointerO LO.LensO[ConfigStruct, *string]")

	// Test constructor template
	var constructorBuf bytes.Buffer
	err = constructorTmpl.Execute(&constructorBuf, s)
	require.NoError(t, err)

	constructorStr := constructorBuf.String()
	assert.Contains(t, constructorStr, "func MakeConfigStructLenses() ConfigStructLenses")
	assert.Contains(t, constructorStr, "IO.FromZero[string]()")
	assert.Contains(t, constructorStr, "IO.FromZero[int]()")
	assert.Contains(t, constructorStr, "IO.FromZero[*string]()")
}

func TestLensCommandFlags(t *testing.T) {
	cmd := LensCommand()

	assert.Equal(t, "lens", cmd.Name)
	assert.Equal(t, "generate lens code for annotated structs", cmd.Usage)
	assert.Contains(t, strings.ToLower(cmd.Description), "fp-go:lens")
	assert.Contains(t, strings.ToLower(cmd.Description), "lenso", "Description should mention LensO for optional lenses")

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
