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
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/IBM/fp-go/gen/v2/cli/lenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	C "github.com/urfave/cli/v3"
)

// assertContainsField checks if content contains a field declaration, ignoring extra whitespace.
func assertContainsField(t *testing.T, content, fieldName, lensType string, msgAndArgs ...interface{}) {
	t.Helper()
	pattern := regexp.QuoteMeta(fieldName) + `\s+` + regexp.QuoteMeta(lensType)
	matched, err := regexp.MatchString(pattern, content)
	require.NoError(t, err, "regex compilation failed")
	assert.True(t, matched, append(msgAndArgs, "Expected to find field: %s %s", fieldName, lensType)...)
}

// ---------------------------------------------------------------------------
// Template smoke tests — verify the lenses package templates are accessible
// ---------------------------------------------------------------------------

func TestLensTemplates(t *testing.T) {
	s := lenses.StructInfo{
		Name:          "TestStruct",
		QualifiedName: "TestStruct",
		Fields: []lenses.FieldInfo{
			{Name: "Name", TypeName: "string", IsOptional: false, IsComparable: true},
			{Name: "Value", TypeName: "*int", IsOptional: true, IsComparable: true},
		},
	}

	var structBuf bytes.Buffer
	err := lenses.StructTmpl.Execute(&structBuf, s)
	require.NoError(t, err)

	structStr := structBuf.String()
	assert.Contains(t, structStr, "type TestStructLenses struct")
	assert.Contains(t, structStr, "Name __lens.Lens[TestStruct, string]")
	assert.Contains(t, structStr, "NameO __lens_option.LensO[TestStruct, string]")
	assert.Contains(t, structStr, "Value __lens.Lens[TestStruct, *int]")
	assert.Contains(t, structStr, "ValueO __lens_option.LensO[TestStruct, *int]")

	var standaloneBuf bytes.Buffer
	err = lenses.StandaloneTmpl.Execute(&standaloneBuf, s)
	require.NoError(t, err)

	standaloneStr := standaloneBuf.String()
	assert.Contains(t, standaloneStr, "func MakeTestStructNameLens() __lens.Lens[TestStruct, string]")
	assert.Contains(t, standaloneStr, "func MakeTestStructNameLensO() __lens_option.LensO[TestStruct, string]")
	assert.Contains(t, standaloneStr, "func MakeTestStructValueLens() __lens.Lens[TestStruct, *int]")
	assert.Contains(t, standaloneStr, "func MakeTestStructValueLensO() __lens_option.LensO[TestStruct, *int]")
	assert.Contains(t, standaloneStr, "__iso_option.FromZero")

	var constructorBuf bytes.Buffer
	err = lenses.ConstructorTmpl.Execute(&constructorBuf, s)
	require.NoError(t, err)

	constructorStr := constructorBuf.String()
	assert.Contains(t, constructorStr, "func MakeTestStructLenses() TestStructLenses")
	assert.Contains(t, constructorStr, "return TestStructLenses{")
	assert.Contains(t, constructorStr, "Name: MakeTestStructNameLens(),")
	assert.Contains(t, constructorStr, "NameO: MakeTestStructNameLensO(),")
	assert.Contains(t, constructorStr, "Value: MakeTestStructValueLens(),")
	assert.Contains(t, constructorStr, "ValueO: MakeTestStructValueLensO(),")
}

func TestLensTemplatesWithOmitEmpty(t *testing.T) {
	s := lenses.StructInfo{
		Name:          "ConfigStruct",
		QualifiedName: "ConfigStruct",
		Fields: []lenses.FieldInfo{
			{Name: "Name", TypeName: "string", IsOptional: false, IsComparable: true},
			{Name: "Value", TypeName: "string", IsOptional: true, IsComparable: true},
			{Name: "Count", TypeName: "int", IsOptional: true, IsComparable: true},
			{Name: "Pointer", TypeName: "*string", IsOptional: true, IsComparable: true},
		},
	}

	var structBuf bytes.Buffer
	err := lenses.StructTmpl.Execute(&structBuf, s)
	require.NoError(t, err)

	structStr := structBuf.String()
	assert.Contains(t, structStr, "NameO __lens_option.LensO[ConfigStruct, string]")
	assert.Contains(t, structStr, "ValueO __lens_option.LensO[ConfigStruct, string]")
	assert.Contains(t, structStr, "CountO __lens_option.LensO[ConfigStruct, int]")
	assert.Contains(t, structStr, "PointerO __lens_option.LensO[ConfigStruct, *string]")

	var standaloneBuf bytes.Buffer
	err = lenses.StandaloneTmpl.Execute(&standaloneBuf, s)
	require.NoError(t, err)

	standaloneStr := standaloneBuf.String()
	assert.Contains(t, standaloneStr, "__iso_option.FromZero[string]()")
	assert.Contains(t, standaloneStr, "__iso_option.FromZero[int]()")
	assert.Contains(t, standaloneStr, "__iso_option.FromZero[*string]()")

	var constructorBuf bytes.Buffer
	err = lenses.ConstructorTmpl.Execute(&constructorBuf, s)
	require.NoError(t, err)

	constructorStr := constructorBuf.String()
	assert.Contains(t, constructorStr, "func MakeConfigStructLenses() ConfigStructLenses")
	assert.Contains(t, constructorStr, "NameO: MakeConfigStructNameLensO(),")
	assert.Contains(t, constructorStr, "CountO: MakeConfigStructCountLensO(),")
	assert.Contains(t, constructorStr, "PointerO: MakeConfigStructPointerLensO(),")
}

func TestLensRefTemplatesWithComparable(t *testing.T) {
	s := lenses.StructInfo{
		Name:          "TestStruct",
		QualifiedName: "TestStruct",
		Fields: []lenses.FieldInfo{
			{Name: "Name", TypeName: "string", IsOptional: false, IsComparable: true},
			{Name: "Age", TypeName: "int", IsOptional: false, IsComparable: true},
			{Name: "Data", TypeName: "[]byte", IsOptional: false, IsComparable: false},
			{Name: "Pointer", TypeName: "*string", IsOptional: true, IsComparable: false},
		},
	}

	var standaloneBuf bytes.Buffer
	err := lenses.StandaloneTmpl.Execute(&standaloneBuf, s)
	require.NoError(t, err)
	standaloneStr := standaloneBuf.String()

	var constructorBuf bytes.Buffer
	err = lenses.ConstructorTmpl.Execute(&constructorBuf, s)
	require.NoError(t, err)
	assert.Contains(t, constructorBuf.String(), "func MakeTestStructRefLenses() TestStructRefLenses")

	assert.Contains(t, standaloneStr, "return __lens.MakeLensStrictWithName(",
		"comparable field Name should use MakeLensStrictWithName in standalone RefLens")
	assert.Contains(t, standaloneStr, "return __lens.MakeLensRefWithName(",
		"non-comparable field Data should use MakeLensRefWithName in standalone RefLens")

	assert.Contains(t, standaloneStr, "func MakeTestStructNameRefLens() __lens.Lens[*TestStruct, string]")
	assert.Contains(t, standaloneStr, "func MakeTestStructAgeRefLens() __lens.Lens[*TestStruct, int]")
	assert.Contains(t, standaloneStr, "func MakeTestStructDataRefLens() __lens.Lens[*TestStruct, []byte]")
}

// ---------------------------------------------------------------------------
// generateLensHelpers integration tests
// ---------------------------------------------------------------------------

func TestGenerateLensHelpers(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

// fp-go:Lens
type TestStruct struct {
	Name  string
	Value *int
}
`), 0o644)
	require.NoError(t, err)

	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, outputFile))
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "package testpkg")
	assert.Contains(t, contentStr, "Code generated by go generate")
	assert.Contains(t, contentStr, "TestStructLenses")
	assert.Contains(t, contentStr, "MakeTestStructLenses")
	assert.Contains(t, contentStr, "__lens.Lens[TestStruct, string]")
	assert.Contains(t, contentStr, "__lens_option.LensO[TestStruct, *int]")
	assert.Contains(t, contentStr, "__iso_option.FromZero")
}

func TestGenerateLensHelpersNoAnnotations(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

// No annotation
type TestStruct struct {
	Name string
}
`), 0o644)
	require.NoError(t, err)

	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(tmpDir, outputFile))
	assert.True(t, os.IsNotExist(err))
}

func TestGenerateLensHelpersWithComparable(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

// fp-go:Lens
type TestStruct struct {
	Name  string
	Count int
	Data  []byte
}
`), 0o644)
	require.NoError(t, err)

	outputFile := "gen.go"
	err = generateLensHelpers(tmpDir, outputFile, false, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, outputFile))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "MakeTestStructRefLenses")
	assert.Contains(t, contentStr, "return __lens.MakeLensStrictWithName(",
		"comparable fields should use MakeLensStrictWithName")
	assert.Contains(t, contentStr, "return __lens.MakeLensRefWithName(",
		"non-comparable fields should use MakeLensRefWithName")
	assert.Contains(t, contentStr, "func MakeTestStructNameRefLens() __lens.Lens[*TestStruct, string]")
	assert.Contains(t, contentStr, "func MakeTestStructDataRefLens() __lens.Lens[*TestStruct, []byte]")
}

func TestGenerateLensHelpersOptionFieldNoLensO(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

import "github.com/IBM/fp-go/v2/option"

// fp-go:Lens
type Config struct {
	Name  string
	Value option.Option[string]
	Count int
}
`), 0o644)
	require.NoError(t, err)

	err = generateLensHelpers(tmpDir, "gen.go", false, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen.go"))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "__lens.Lens[Config, string]")
	assert.Contains(t, contentStr, "__lens.Lens[Config, option.Option[string]]")
	assert.Contains(t, contentStr, "__lens.Lens[Config, int]")
	assert.Contains(t, contentStr, "__lens_option.LensO[Config, string]")
	assert.Contains(t, contentStr, "__lens_option.LensO[Config, int]")
	assert.NotContains(t, contentStr, "__lens_option.LensO[Config, option.Option[string]]",
		"LensO for Value (option.Option[string]) must NOT be generated")
}

func TestGenerateLensHelpersWithEmbeddedStruct(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

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
`), 0o644)
	require.NoError(t, err)

	err = generateLensHelpers(tmpDir, "gen.go", false, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen.go"))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "PersonLenses")
	assert.Contains(t, contentStr, "MakePersonLenses")
	assertContainsField(t, contentStr, "Street", "__lens.Lens[Person, string]")
	assertContainsField(t, contentStr, "City", "__lens.Lens[Person, string]")
	assertContainsField(t, contentStr, "Name", "__lens.Lens[Person, string]")
	assertContainsField(t, contentStr, "Age", "__lens.Lens[Person, int]")
	assertContainsField(t, contentStr, "StreetO", "__lens_option.LensO[Person, string]")
	assertContainsField(t, contentStr, "CityO", "__lens_option.LensO[Person, string]")
}

func TestGenerateLensHelpersWithGenericStruct(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

// fp-go:Lens
type Box[T any] struct {
	Content T
	Label   string
}
`), 0o644)
	require.NoError(t, err)

	err = generateLensHelpers(tmpDir, "gen.go", false, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen.go"))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "type BoxLenses[T any] struct")
	assert.Contains(t, contentStr, "type BoxRefLenses[T any] struct")
	assert.Contains(t, contentStr, "func MakeBoxLenses[T any]() BoxLenses[T]")
	assert.Contains(t, contentStr, "func MakeBoxRefLenses[T any]() BoxRefLenses[T]")
	assertContainsField(t, contentStr, "Content", "__lens.Lens[Box[T], T]")
	assertContainsField(t, contentStr, "Label", "__lens.Lens[Box[T], string]")
	assert.NotContains(t, contentStr, "ContentO __lens_option.LensO[Box[T], T]",
		"T any is not comparable, should not have optional lens")
	assert.Contains(t, contentStr, "LabelO __lens_option.LensO[Box[T], string]")
}

func TestGenerateLensHelpersWithComparableTypeParam(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

// fp-go:Lens
type ComparableBox[T comparable] struct {
	Key   T
	Value string
}
`), 0o644)
	require.NoError(t, err)

	err = generateLensHelpers(tmpDir, "gen.go", false, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen.go"))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "type ComparableBoxLenses[T comparable] struct")
	assert.Contains(t, contentStr, "type ComparableBoxRefLenses[T comparable] struct")
	assert.Contains(t, contentStr, "func MakeComparableBoxKeyRefLens[T comparable]() __lens.Lens[*ComparableBox[T], T]")
	assert.Contains(t, contentStr, "return __lens.MakeLensStrictWithName(")
	assert.Contains(t, contentStr, "func MakeComparableBoxValueRefLens[T comparable]() __lens.Lens[*ComparableBox[T], string]")
	assert.NotContains(t, contentStr, "return __lens.MakeLensRefWithName(")
}

func TestGenerateLensHelpersWithUnexportedFields(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

// fp-go:Lens
type MixedStruct struct {
	PublicField     string
	privateField    int
	OptionalPrivate *string
}
`), 0o644)
	require.NoError(t, err)

	err = generateLensHelpers(tmpDir, "gen_lens.go", false, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen_lens.go"))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "MixedStructLenses")
	assert.Contains(t, contentStr, "MakeMixedStructLenses")
	assertContainsField(t, contentStr, "PublicField", "__lens.Lens[MixedStruct, string]")
	assertContainsField(t, contentStr, "privateField", "__lens.Lens[MixedStruct, int]")
	assertContainsField(t, contentStr, "OptionalPrivate", "__lens.Lens[MixedStruct, *string]")
	assert.Contains(t, contentStr, "func(s MixedStruct) string { return s.PublicField }")
	assert.Contains(t, contentStr, "func(s MixedStruct) int { return s.privateField }")
	assert.Contains(t, contentStr, "func(s MixedStruct) *string { return s.OptionalPrivate }")
}

func TestGenerateLensHelpersWithUnexportedEmbeddedFields(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

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
`), 0o644)
	require.NoError(t, err)

	err = generateLensHelpers(tmpDir, "gen_lens.go", false, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen_lens.go"))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "ExtendedConfigLenses")
	assertContainsField(t, contentStr, "publicBase", "__lens.Lens[ExtendedConfig, string]")
	assertContainsField(t, contentStr, "privateBase", "__lens.Lens[ExtendedConfig, int]")
	assertContainsField(t, contentStr, "PublicField", "__lens.Lens[ExtendedConfig, string]")
	assertContainsField(t, contentStr, "privateField", "__lens.Lens[ExtendedConfig, bool]")
}

func TestGenerateLensHelpersWithQualifiedField(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(`package testpkg

import "some/llmclient"

// fp-go:Lens
type Config struct {
	Name      string
	LLMConfig llmclient.Config
}
`), 0o644)
	require.NoError(t, err)

	structs, _, err := lenses.ParseFile(filepath.Join(tmpDir, "test.go"))
	require.NoError(t, err)
	require.Len(t, structs, 1)

	cfg := structs[0]
	require.Len(t, cfg.Fields, 2)
	assert.True(t, cfg.Fields[0].IsComparable, "string field should be comparable")
	assert.False(t, cfg.Fields[1].IsComparable, "qualified struct field should be treated as non-comparable")
}

// ---------------------------------------------------------------------------
// LensCommand flags
// ---------------------------------------------------------------------------

func TestLensCommandFlags(t *testing.T) {
	cmd := LensCommand()

	assert.Equal(t, "lens", cmd.Name)
	assert.Contains(t, cmd.Usage, "lens")
	assert.Contains(t, strings.ToLower(cmd.Description), "fp-go:lens")
	assert.Contains(t, strings.ToLower(cmd.Description), "lenso", "Description should mention LensO")

	var hasDir, hasFilename, hasVerbose, hasIncludeTestFiles, hasType bool
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
		case "type":
			hasType = true
		}
	}

	assert.True(t, hasDir, "should have dir flag")
	assert.True(t, hasFilename, "should have filename flag")
	assert.True(t, hasVerbose, "should have verbose flag")
	assert.True(t, hasIncludeTestFiles, "should have include-test-files flag")
	assert.True(t, hasType, "should have type flag")
}

// ---------------------------------------------------------------------------
// LensCommand --type mode (net/http.Server integration test)
// ---------------------------------------------------------------------------

func TestLensCommandHttpServer(t *testing.T) {
	tmpDir := t.TempDir()

	app := &C.Command{
		Commands: []*C.Command{LensCommand()},
	}

	err := app.Run(context.Background(), []string{
		"app", "lens",
		"--type", "Server",
		"--dir", tmpDir,
		"--filename", "gen_lens.go",
		"net/http",
	})
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen_lens.go"))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "package http")
	assert.Contains(t, contentStr, "Code generated by go generate")
	assert.Contains(t, contentStr, "type ServerLenses struct")
	assert.Contains(t, contentStr, "type ServerRefLenses struct")
	assert.Contains(t, contentStr, "func MakeServerLenses() ServerLenses")
	assert.Contains(t, contentStr, "func MakeServerRefLenses() ServerRefLenses")

	assertContainsField(t, contentStr, "Addr", "__lens.Lens[http.Server, string]")
	assertContainsField(t, contentStr, "AddrO", "__lens_option.LensO[http.Server, string]")
	assertContainsField(t, contentStr, "ReadTimeout", "__lens.Lens[http.Server, time.Duration]")
	assertContainsField(t, contentStr, "ReadTimeoutO", "__lens_option.LensO[http.Server, time.Duration]")
	assertContainsField(t, contentStr, "TLSConfig", "__lens.Lens[http.Server, *tls.Config]")
	assertContainsField(t, contentStr, "TLSConfigO", "__lens_option.LensO[http.Server, *tls.Config]")
	assertContainsField(t, contentStr, "MaxHeaderBytes", "__lens.Lens[http.Server, int]")
	assertContainsField(t, contentStr, "MaxHeaderBytesO", "__lens_option.LensO[http.Server, int]")

	assert.NotContains(t, contentStr, "TLSNextProtoO __lens_option.LensO")
	assert.NotContains(t, contentStr, "ConnStateO __lens_option.LensO")
	assert.NotContains(t, contentStr, "BaseContextO __lens_option.LensO")
	assert.NotContains(t, contentStr, "ConnContextO __lens_option.LensO")
}
