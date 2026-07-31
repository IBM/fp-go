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
	"os"
	"path/filepath"
	"testing"

	"github.com/IBM/fp-go/gen/v2/cli/lenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// ParseFile
// ---------------------------------------------------------------------------

func TestParseFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

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
`), 0o644)
	require.NoError(t, err)

	structs, pkg, err := lenses.ParseFile(testFile)
	require.NoError(t, err)

	assert.Equal(t, "testpkg", pkg)
	require.Len(t, structs, 2)

	person := structs[0]
	assert.Equal(t, "Person", person.Name)
	require.Len(t, person.Fields, 3)
	assert.Equal(t, "Name", person.Fields[0].Name)
	assert.Equal(t, "string", person.Fields[0].TypeName)
	assert.False(t, person.Fields[0].IsOptional)
	assert.Equal(t, "Age", person.Fields[1].Name)
	assert.False(t, person.Fields[1].IsOptional)
	assert.Equal(t, "Phone", person.Fields[2].Name)
	assert.Equal(t, "*string", person.Fields[2].TypeName)
	assert.True(t, person.Fields[2].IsOptional)

	address := structs[1]
	assert.Equal(t, "Address", address.Name)
	require.Len(t, address.Fields, 2)
	assert.Equal(t, "Street", address.Fields[0].Name)
	assert.Equal(t, "City", address.Fields[1].Name)
}

func TestParseFileWithOmitEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte("package testpkg\n\n// fp-go:Lens\ntype Config struct {\n\tName     string\n\tValue    string  `json:\"value,omitempty\"`\n\tCount    int     `json:\",omitempty\"`\n\tOptional *string `json:\"optional,omitempty\"`\n\tRequired int     `json:\"required\"`\n}\n"), 0o644)
	require.NoError(t, err)

	structs, pkg, err := lenses.ParseFile(testFile)
	require.NoError(t, err)

	assert.Equal(t, "testpkg", pkg)
	require.Len(t, structs, 1)

	config := structs[0]
	require.Len(t, config.Fields, 5)
	assert.False(t, config.Fields[0].IsOptional, "Name has no omitempty")
	assert.True(t, config.Fields[1].IsOptional, "Value has omitempty")
	assert.True(t, config.Fields[2].IsOptional, "Count has omitempty")
	assert.True(t, config.Fields[3].IsOptional, "Optional is a pointer")
	assert.False(t, config.Fields[4].IsOptional, "Required has no omitempty")
}

func TestParseFileWithComparableTypes(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

// fp-go:Lens
type TypeTest struct {
	Name    string
	Age     int
	Pointer *string
	Slice   []string
	Map     map[string]int
	Channel chan int
}
`), 0o644)
	require.NoError(t, err)

	structs, _, err := lenses.ParseFile(testFile)
	require.NoError(t, err)
	require.Len(t, structs, 1)

	fields := structs[0].Fields
	require.Len(t, fields, 6)
	assert.True(t, fields[0].IsComparable, "string should be comparable")
	assert.True(t, fields[1].IsComparable, "int should be comparable")
	// Pointer — IsComparable is still set to true for pointer types
	assert.True(t, fields[2].IsOptional)
	assert.False(t, fields[3].IsComparable, "slice should not be comparable")
	assert.False(t, fields[4].IsComparable, "map should not be comparable")
	assert.True(t, fields[5].IsComparable, "channel should be comparable")
}

func TestParseFileWithOptionField(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

import "github.com/IBM/fp-go/v2/option"

// fp-go:Lens
type Config struct {
	Name  string
	Value option.Option[string]
	Count int
}
`), 0o644)
	require.NoError(t, err)

	structs, _, err := lenses.ParseFile(testFile)
	require.NoError(t, err)
	require.Len(t, structs, 1)

	cfg := structs[0]
	require.Len(t, cfg.Fields, 3)

	assert.False(t, cfg.Fields[0].IsOption, "string must not be marked IsOption")
	assert.True(t, cfg.Fields[1].IsOption, "option.Option[string] must be marked IsOption")
	assert.True(t, cfg.Fields[1].IsComparable, "option.Option[string] is comparable")
	assert.False(t, cfg.Fields[2].IsOption, "int must not be marked IsOption")
}

func TestParseFileWithQualifiedField(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

import "some/llmclient"

// fp-go:Lens
type Config struct {
	Name      string
	LLMConfig llmclient.Config
}
`), 0o644)
	require.NoError(t, err)

	structs, _, err := lenses.ParseFile(testFile)
	require.NoError(t, err)
	require.Len(t, structs, 1)

	cfg := structs[0]
	require.Len(t, cfg.Fields, 2)
	assert.True(t, cfg.Fields[0].IsComparable, "string should be comparable")
	assert.False(t, cfg.Fields[1].IsComparable, "qualified struct field should be treated as non-comparable")
}

func TestParseFileWithEmbeddedStruct(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

type Base struct {
	ID   int
	Name string
}

// fp-go:Lens
type Extended struct {
	Base
	Extra string
}
`), 0o644)
	require.NoError(t, err)

	structs, pkg, err := lenses.ParseFile(testFile)
	require.NoError(t, err)

	assert.Equal(t, "testpkg", pkg)
	require.Len(t, structs, 1)

	extended := structs[0]
	assert.Equal(t, "Extended", extended.Name)
	require.Len(t, extended.Fields, 3, "Should have 3 fields: ID, Name (from Base), and Extra")

	fieldNames := make(map[string]bool)
	for _, f := range extended.Fields {
		fieldNames[f.Name] = true
	}
	assert.True(t, fieldNames["ID"])
	assert.True(t, fieldNames["Name"])
	assert.True(t, fieldNames["Extra"])
}

func TestParseFileWithPointerEmbeddedStruct(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

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
`), 0o644)
	require.NoError(t, err)

	structs, _, err := lenses.ParseFile(testFile)
	require.NoError(t, err)
	require.Len(t, structs, 1)

	doc := structs[0]
	assert.Len(t, doc.Fields, 4, "Should have 4 fields: CreatedAt, UpdatedAt (from *Metadata), Title, Content")

	fieldNames := make(map[string]bool)
	for _, f := range doc.Fields {
		fieldNames[f.Name] = true
	}
	assert.True(t, fieldNames["CreatedAt"])
	assert.True(t, fieldNames["UpdatedAt"])
	assert.True(t, fieldNames["Title"])
	assert.True(t, fieldNames["Content"])
}

func TestParseFileWithGenericStruct(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

// fp-go:Lens
type Container[T any] struct {
	Value T
	Count int
}
`), 0o644)
	require.NoError(t, err)

	structs, pkg, err := lenses.ParseFile(testFile)
	require.NoError(t, err)

	assert.Equal(t, "testpkg", pkg)
	require.Len(t, structs, 1)

	container := structs[0]
	assert.Equal(t, "Container", container.Name)
	assert.Equal(t, "[T any]", container.TypeParams)
	require.Len(t, container.Fields, 2)
	assert.Equal(t, "T", container.Fields[0].TypeName)
	assert.Equal(t, "int", container.Fields[1].TypeName)
}

func TestParseFileWithMultipleTypeParams(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

// fp-go:Lens
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}
`), 0o644)
	require.NoError(t, err)

	structs, _, err := lenses.ParseFile(testFile)
	require.NoError(t, err)
	require.Len(t, structs, 1)

	pair := structs[0]
	assert.Equal(t, "[K comparable, V any]", pair.TypeParams)
	require.Len(t, pair.Fields, 2)
	assert.Equal(t, "K", pair.Fields[0].TypeName)
	assert.Equal(t, "V", pair.Fields[1].TypeName)
}

func TestParseFileWithUnexportedFields(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

// fp-go:Lens
type Config struct {
	PublicName   string
	privateName  string
	PublicValue  int
	privateValue *int
}
`), 0o644)
	require.NoError(t, err)

	structs, _, err := lenses.ParseFile(testFile)
	require.NoError(t, err)
	require.Len(t, structs, 1)

	config := structs[0]
	require.Len(t, config.Fields, 4, "Should include both exported and unexported fields")
	assert.Equal(t, "PublicName", config.Fields[0].Name)
	assert.Equal(t, "privateName", config.Fields[1].Name)
	assert.Equal(t, "PublicValue", config.Fields[2].Name)
	assert.Equal(t, "privateValue", config.Fields[3].Name)
	assert.True(t, config.Fields[3].IsOptional)
}

func TestParseFileWithOnlyUnexportedFields(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(`package testpkg

// fp-go:Lens
type PrivateConfig struct {
	name    string
	value   int
	enabled bool
}
`), 0o644)
	require.NoError(t, err)

	structs, _, err := lenses.ParseFile(testFile)
	require.NoError(t, err)
	require.Len(t, structs, 1)

	config := structs[0]
	assert.Equal(t, "PrivateConfig", config.Name)
	require.Len(t, config.Fields, 3)
	assert.Equal(t, "name", config.Fields[0].Name)
	assert.Equal(t, "value", config.Fields[1].Name)
	assert.Equal(t, "enabled", config.Fields[2].Name)
}
