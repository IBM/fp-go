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

	S "github.com/IBM/fp-go/v2/string"
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
			if S.IsNonEmpty(tt.comment) {
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
			result := isComparableType(fieldType, map[string]string{})
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
			if S.IsNonEmpty(tt.tag) {
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

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
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

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
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

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
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
	assert.Contains(t, constructorStr, "lensName := __lens.MakeLensStrictWithName(",
		"comparable field Name should use MakeLensStrictWithName in RefLenses")

	// Age field - comparable, should use MakeLensStrict
	assert.Contains(t, constructorStr, "lensAge := __lens.MakeLensStrictWithName(",
		"comparable field Age should use MakeLensStrictWithName in RefLenses")

	// Data field - not comparable, should use MakeLensRef
	assert.Contains(t, constructorStr, "lensData := __lens.MakeLensRefWithName(",
		"non-comparable field Data should use MakeLensRefWithName in RefLenses")

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
	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Generate lens code
	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
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

	// Name and Count are comparable, should use MakeLensStrictWithName
	assert.Contains(t, contentStr, "__lens.MakeLensStrictWithName",
		"comparable fields should use MakeLensStrictWithName in RefLenses")

	// Data is not comparable (slice), should use MakeLensRefWithName
	assert.Contains(t, contentStr, "__lens.MakeLensRefWithName",
		"non-comparable fields should use MakeLensRefWithName in RefLenses")

	// Verify the pattern appears for Name field (comparable)
	namePattern := "lensName := __lens.MakeLensStrictWithName("
	assert.Contains(t, contentStr, namePattern,
		"Name field should use MakeLensStrictWithName")

	// Verify the pattern appears for Data field (not comparable)
	dataPattern := "lensData := __lens.MakeLensRefWithName("
	assert.Contains(t, contentStr, dataPattern,
		"Data field should use MakeLensRefWithName")
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
	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Generate lens code
	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
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
	assert.Contains(t, contentStr, "__lens.Lens[TestStruct, string]")
	assert.Contains(t, contentStr, "__lens_option.LensO[TestStruct, *int]")
	assert.Contains(t, contentStr, "__iso_option.FromZero")
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
	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Generate lens code (should not create file)
	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
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
			{Name: "Name", TypeName: "string", IsOptional: false, IsComparable: true},
			{Name: "Value", TypeName: "*int", IsOptional: true, IsComparable: true},
		},
	}

	// Test struct template
	var structBuf bytes.Buffer
	err := structTmpl.Execute(&structBuf, s)
	require.NoError(t, err)

	structStr := structBuf.String()
	assert.Contains(t, structStr, "type TestStructLenses struct")
	assert.Contains(t, structStr, "Name __lens.Lens[TestStruct, string]")
	assert.Contains(t, structStr, "NameO __lens_option.LensO[TestStruct, string]")
	assert.Contains(t, structStr, "Value __lens.Lens[TestStruct, *int]")
	assert.Contains(t, structStr, "ValueO __lens_option.LensO[TestStruct, *int]")

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
	assert.Contains(t, constructorStr, "__iso_option.FromZero")
}

func TestLensTemplatesWithOmitEmpty(t *testing.T) {
	s := structInfo{
		Name: "ConfigStruct",
		Fields: []fieldInfo{
			{Name: "Name", TypeName: "string", IsOptional: false, IsComparable: true},
			{Name: "Value", TypeName: "string", IsOptional: true, IsComparable: true},    // non-pointer with omitempty
			{Name: "Count", TypeName: "int", IsOptional: true, IsComparable: true},       // non-pointer with omitempty
			{Name: "Pointer", TypeName: "*string", IsOptional: true, IsComparable: true}, // pointer
		},
	}

	// Test struct template
	var structBuf bytes.Buffer
	err := structTmpl.Execute(&structBuf, s)
	require.NoError(t, err)

	structStr := structBuf.String()
	assert.Contains(t, structStr, "type ConfigStructLenses struct")
	assert.Contains(t, structStr, "Name __lens.Lens[ConfigStruct, string]")
	assert.Contains(t, structStr, "NameO __lens_option.LensO[ConfigStruct, string]")
	assert.Contains(t, structStr, "Value __lens.Lens[ConfigStruct, string]")
	assert.Contains(t, structStr, "ValueO __lens_option.LensO[ConfigStruct, string]", "comparable non-pointer with omitempty should have optional lens")
	assert.Contains(t, structStr, "Count __lens.Lens[ConfigStruct, int]")
	assert.Contains(t, structStr, "CountO __lens_option.LensO[ConfigStruct, int]", "comparable non-pointer with omitempty should have optional lens")
	assert.Contains(t, structStr, "Pointer __lens.Lens[ConfigStruct, *string]")
	assert.Contains(t, structStr, "PointerO __lens_option.LensO[ConfigStruct, *string]")

	// Test constructor template
	var constructorBuf bytes.Buffer
	err = constructorTmpl.Execute(&constructorBuf, s)
	require.NoError(t, err)

	constructorStr := constructorBuf.String()
	assert.Contains(t, constructorStr, "func MakeConfigStructLenses() ConfigStructLenses")
	assert.Contains(t, constructorStr, "__iso_option.FromZero[string]()")
	assert.Contains(t, constructorStr, "__iso_option.FromZero[int]()")
	assert.Contains(t, constructorStr, "__iso_option.FromZero[*string]()")
}

func TestLensCommandFlags(t *testing.T) {
	cmd := LensCommand()

	assert.Equal(t, "lens", cmd.Name)
	assert.Equal(t, "generate lens code for annotated structs", cmd.Usage)
	assert.Contains(t, strings.ToLower(cmd.Description), "fp-go:lens")
	assert.Contains(t, strings.ToLower(cmd.Description), "lenso", "Description should mention LensO for optional lenses")

	// Check flags
	assert.Len(t, cmd.Flags, 4)

	var hasDir, hasFilename, hasVerbose, hasIncludeTestFiles bool
	for _, flag := range cmd.Flags {
		switch flag.Names()[0] {
		case "dir":
			hasDir = true
		case "filename":
			hasFilename = true
		case "verbose":
			hasVerbose = true
		case "include-test-files":
			hasIncludeTestFiles = true
		}
	}

	assert.True(t, hasDir, "should have dir flag")
	assert.True(t, hasFilename, "should have filename flag")
	assert.True(t, hasVerbose, "should have verbose flag")
	assert.True(t, hasIncludeTestFiles, "should have include-test-files flag")
}

func TestParseFileWithEmbeddedStruct(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// Base struct to be embedded
type Base struct {
	ID   int
	Name string
}

// fp-go:Lens
type Extended struct {
	Base
	Extra string
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 1)

	// Check Extended struct
	extended := structs[0]
	assert.Equal(t, "Extended", extended.Name)
	assert.Len(t, extended.Fields, 3, "Should have 3 fields: ID, Name (from Base), and Extra")

	// Check that embedded fields are promoted
	fieldNames := make(map[string]bool)
	for _, field := range extended.Fields {
		fieldNames[field.Name] = true
	}

	assert.True(t, fieldNames["ID"], "Should have promoted ID field from Base")
	assert.True(t, fieldNames["Name"], "Should have promoted Name field from Base")
	assert.True(t, fieldNames["Extra"], "Should have Extra field")
}

func TestGenerateLensHelpersWithEmbeddedStruct(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testCode := `package testpkg

// Base struct to be embedded
type Address struct {
	Street string
	City   string
}

// fp-go:Lens
type Person struct {
	Address
	Name string
	Age  int
}
`

	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Generate lens code
	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
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
	assert.Contains(t, contentStr, "PersonLenses")
	assert.Contains(t, contentStr, "MakePersonLenses")

	// Check that embedded fields are included
	assert.Contains(t, contentStr, "Street __lens.Lens[Person, string]", "Should have lens for embedded Street field")
	assert.Contains(t, contentStr, "City __lens.Lens[Person, string]", "Should have lens for embedded City field")
	assert.Contains(t, contentStr, "Name __lens.Lens[Person, string]", "Should have lens for Name field")
	assert.Contains(t, contentStr, "Age __lens.Lens[Person, int]", "Should have lens for Age field")

	// Check that optional lenses are also generated for embedded fields
	assert.Contains(t, contentStr, "StreetO __lens_option.LensO[Person, string]")
	assert.Contains(t, contentStr, "CityO __lens_option.LensO[Person, string]")
}

func TestParseFileWithPointerEmbeddedStruct(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// Base struct to be embedded
type Metadata struct {
	CreatedAt string
	UpdatedAt string
}

// fp-go:Lens
type Document struct {
	*Metadata
	Title   string
	Content string
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 1)

	// Check Document struct
	doc := structs[0]
	assert.Equal(t, "Document", doc.Name)
	assert.Len(t, doc.Fields, 4, "Should have 4 fields: CreatedAt, UpdatedAt (from *Metadata), Title, and Content")

	// Check that embedded fields are promoted
	fieldNames := make(map[string]bool)
	for _, field := range doc.Fields {
		fieldNames[field.Name] = true
	}

	assert.True(t, fieldNames["CreatedAt"], "Should have promoted CreatedAt field from *Metadata")
	assert.True(t, fieldNames["UpdatedAt"], "Should have promoted UpdatedAt field from *Metadata")
	assert.True(t, fieldNames["Title"], "Should have Title field")
	assert.True(t, fieldNames["Content"], "Should have Content field")
}

func TestParseFileWithGenericStruct(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// fp-go:Lens
type Container[T any] struct {
	Value T
	Count int
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 1)

	// Check Container struct
	container := structs[0]
	assert.Equal(t, "Container", container.Name)
	assert.Equal(t, "[T any]", container.TypeParams, "Should have type parameter [T any]")
	assert.Len(t, container.Fields, 2)

	assert.Equal(t, "Value", container.Fields[0].Name)
	assert.Equal(t, "T", container.Fields[0].TypeName)

	assert.Equal(t, "Count", container.Fields[1].Name)
	assert.Equal(t, "int", container.Fields[1].TypeName)
}

func TestParseFileWithMultipleTypeParams(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// fp-go:Lens
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 1)

	// Check Pair struct
	pair := structs[0]
	assert.Equal(t, "Pair", pair.Name)
	assert.Equal(t, "[K comparable, V any]", pair.TypeParams, "Should have type parameters [K comparable, V any]")
	assert.Len(t, pair.Fields, 2)

	assert.Equal(t, "Key", pair.Fields[0].Name)
	assert.Equal(t, "K", pair.Fields[0].TypeName)

	assert.Equal(t, "Value", pair.Fields[1].Name)
	assert.Equal(t, "V", pair.Fields[1].TypeName)
}

func TestGenerateLensHelpersWithGenericStruct(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testCode := `package testpkg

// fp-go:Lens
type Box[T any] struct {
	Content T
	Label   string
}
`

	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Generate lens code
	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
	require.NoError(t, err)

	// Verify the generated file exists
	genPath := filepath.Join(tmpDir, outputFile)
	_, err = os.Stat(genPath)
	require.NoError(t, err)

	// Read and verify the generated content
	content, err := os.ReadFile(genPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Check for expected content with type parameters
	assert.Contains(t, contentStr, "package testpkg")
	assert.Contains(t, contentStr, "type BoxLenses[T any] struct", "Should have generic BoxLenses type")
	assert.Contains(t, contentStr, "type BoxRefLenses[T any] struct", "Should have generic BoxRefLenses type")
	assert.Contains(t, contentStr, "func MakeBoxLenses[T any]() BoxLenses[T]", "Should have generic constructor")
	assert.Contains(t, contentStr, "func MakeBoxRefLenses[T any]() BoxRefLenses[T]", "Should have generic ref constructor")

	// Check that fields use the generic type parameter
	assert.Contains(t, contentStr, "Content __lens.Lens[Box[T], T]", "Should have lens for generic Content field")
	assert.Contains(t, contentStr, "Label __lens.Lens[Box[T], string]", "Should have lens for Label field")

	// Check optional lenses - only for comparable types
	// T any is not comparable, so ContentO should NOT be generated
	assert.NotContains(t, contentStr, "ContentO __lens_option.LensO[Box[T], T]", "T any is not comparable, should not have optional lens")
	// string is comparable, so LabelO should be generated
	assert.Contains(t, contentStr, "LabelO __lens_option.LensO[Box[T], string]", "string is comparable, should have optional lens")
}

func TestGenerateLensHelpersWithComparableTypeParam(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testCode := `package testpkg

// fp-go:Lens
type ComparableBox[T comparable] struct {
	Key   T
	Value string
}
`

	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Generate lens code
	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
	require.NoError(t, err)

	// Verify the generated file exists
	genPath := filepath.Join(tmpDir, outputFile)
	_, err = os.Stat(genPath)
	require.NoError(t, err)

	// Read and verify the generated content
	content, err := os.ReadFile(genPath)
	require.NoError(t, err)

	contentStr := string(content)

	// Check for expected content with type parameters
	assert.Contains(t, contentStr, "package testpkg")
	assert.Contains(t, contentStr, "type ComparableBoxLenses[T comparable] struct", "Should have generic ComparableBoxLenses type")
	assert.Contains(t, contentStr, "type ComparableBoxRefLenses[T comparable] struct", "Should have generic ComparableBoxRefLenses type")

	// Check that Key field (with comparable constraint) uses MakeLensStrict in RefLenses
	assert.Contains(t, contentStr, "lensKey := __lens.MakeLensStrictWithName(", "Key field with comparable constraint should use MakeLensStrictWithName")

	// Check that Value field (string, always comparable) also uses MakeLensStrict
	assert.Contains(t, contentStr, "lensValue := __lens.MakeLensStrictWithName(", "Value field (string) should use MakeLensStrictWithName")

	// Verify that MakeLensRef is NOT used (since both fields are comparable)
	assert.NotContains(t, contentStr, "__lens.MakeLensRefWithName(", "Should not use MakeLensRefWithName when all fields are comparable")
}

func TestParseFileWithUnexportedFields(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// fp-go:Lens
type Config struct {
	PublicName  string
	privateName string
	PublicValue int
	privateValue *int
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
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
	assert.Len(t, config.Fields, 4, "Should include both exported and unexported fields")

	// Check exported field
	assert.Equal(t, "PublicName", config.Fields[0].Name)
	assert.Equal(t, "string", config.Fields[0].TypeName)
	assert.False(t, config.Fields[0].IsOptional)

	// Check unexported field
	assert.Equal(t, "privateName", config.Fields[1].Name)
	assert.Equal(t, "string", config.Fields[1].TypeName)
	assert.False(t, config.Fields[1].IsOptional)

	// Check exported int field
	assert.Equal(t, "PublicValue", config.Fields[2].Name)
	assert.Equal(t, "int", config.Fields[2].TypeName)
	assert.False(t, config.Fields[2].IsOptional)

	// Check unexported pointer field
	assert.Equal(t, "privateValue", config.Fields[3].Name)
	assert.Equal(t, "*int", config.Fields[3].TypeName)
	assert.True(t, config.Fields[3].IsOptional)
}

func TestGenerateLensHelpersWithUnexportedFields(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testCode := `package testpkg

// fp-go:Lens
type MixedStruct struct {
	PublicField  string
	privateField int
	OptionalPrivate *string
}
`

	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Generate lens code
	outputFile := "gen_lens.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
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
	assert.Contains(t, contentStr, "MixedStructLenses")
	assert.Contains(t, contentStr, "MakeMixedStructLenses")

	// Check that lenses are generated for all fields (exported and unexported)
	assert.Contains(t, contentStr, "PublicField __lens.Lens[MixedStruct, string]")
	assert.Contains(t, contentStr, "privateField __lens.Lens[MixedStruct, int]")
	assert.Contains(t, contentStr, "OptionalPrivate __lens.Lens[MixedStruct, *string]")

	// Check lens constructors
	assert.Contains(t, contentStr, "func(s MixedStruct) string { return s.PublicField }")
	assert.Contains(t, contentStr, "func(s MixedStruct) int { return s.privateField }")
	assert.Contains(t, contentStr, "func(s MixedStruct) *string { return s.OptionalPrivate }")

	// Check setters
	assert.Contains(t, contentStr, "func(s MixedStruct, v string) MixedStruct { s.PublicField = v; return s }")
	assert.Contains(t, contentStr, "func(s MixedStruct, v int) MixedStruct { s.privateField = v; return s }")
	assert.Contains(t, contentStr, "func(s MixedStruct, v *string) MixedStruct { s.OptionalPrivate = v; return s }")
}

func TestParseFileWithOnlyUnexportedFields(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// fp-go:Lens
type PrivateConfig struct {
	name    string
	value   int
	enabled bool
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 1)

	// Check PrivateConfig struct
	config := structs[0]
	assert.Equal(t, "PrivateConfig", config.Name)
	assert.Len(t, config.Fields, 3, "Should include all unexported fields")

	// Check all fields are unexported
	assert.Equal(t, "name", config.Fields[0].Name)
	assert.Equal(t, "value", config.Fields[1].Name)
	assert.Equal(t, "enabled", config.Fields[2].Name)
}

func TestGenerateLensHelpersWithUnexportedEmbeddedFields(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	testCode := `package testpkg

type BaseConfig struct {
	publicBase  string
	privateBase int
}

// fp-go:Lens
type ExtendedConfig struct {
	BaseConfig
	PublicField  string
	privateField bool
}
`

	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Generate lens code
	outputFile := "gen_lens.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
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
	assert.Contains(t, contentStr, "ExtendedConfigLenses")

	// Check that lenses are generated for embedded unexported fields
	assert.Contains(t, contentStr, "publicBase __lens.Lens[ExtendedConfig, string]")
	assert.Contains(t, contentStr, "privateBase __lens.Lens[ExtendedConfig, int]")

	// Check that lenses are generated for direct fields (both exported and unexported)
	assert.Contains(t, contentStr, "PublicField __lens.Lens[ExtendedConfig, string]")
	assert.Contains(t, contentStr, "privateField __lens.Lens[ExtendedConfig, bool]")
}

func TestParseFileWithMixedFieldVisibility(t *testing.T) {
	// Create a temporary test file with various field visibility patterns
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package testpkg

// fp-go:Lens
type ComplexStruct struct {
	// Exported fields
	Name        string
	Age         int
	Email       *string
	
	// Unexported fields
	password    string
	secretKey   []byte
	internalID  *int
	
	// Mixed with tags
	PublicWithTag  string ` + "`json:\"public,omitempty\"`" + `
	privateWithTag int    ` + "`json:\"private,omitempty\"`" + `
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0o644)
	require.NoError(t, err)

	// Parse the file
	structs, pkg, err := parseFile(testFile)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", pkg)
	assert.Len(t, structs, 1)

	// Check ComplexStruct
	complex := structs[0]
	assert.Equal(t, "ComplexStruct", complex.Name)
	assert.Len(t, complex.Fields, 8, "Should include all fields regardless of visibility")

	// Verify field names and types
	fieldNames := []string{"Name", "Age", "Email", "password", "secretKey", "internalID", "PublicWithTag", "privateWithTag"}
	for i, expectedName := range fieldNames {
		assert.Equal(t, expectedName, complex.Fields[i].Name, "Field %d should be %s", i, expectedName)
	}

	// Check optional fields
	assert.False(t, complex.Fields[0].IsOptional, "Name should not be optional")
	assert.True(t, complex.Fields[2].IsOptional, "Email (pointer) should be optional")
	assert.True(t, complex.Fields[5].IsOptional, "internalID (pointer) should be optional")
	assert.True(t, complex.Fields[6].IsOptional, "PublicWithTag (with omitempty) should be optional")
	assert.True(t, complex.Fields[7].IsOptional, "privateWithTag (with omitempty) should be optional")
}
