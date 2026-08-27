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
	// option is an infrastructure package — TypeName uses the __option alias.
	assert.Contains(t, contentStr, "__lens.Lens[Config, __option.Option[string]]")
	assert.Contains(t, contentStr, "__lens.Lens[Config, int]")
	assert.Contains(t, contentStr, "__lens_option.LensO[Config, string]")
	assert.Contains(t, contentStr, "__lens_option.LensO[Config, int]")
	assert.NotContains(t, contentStr, "__lens_option.LensO[Config, __option.Option[string]]",
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

	assertContainsField(t, contentStr, "Addr", "__lens.Lens[net_http.Server, string]")
	assertContainsField(t, contentStr, "AddrO", "__lens_option.LensO[net_http.Server, string]")
	assertContainsField(t, contentStr, "ReadTimeout", "__lens.Lens[net_http.Server, time.Duration]")
	assertContainsField(t, contentStr, "ReadTimeoutO", "__lens_option.LensO[net_http.Server, time.Duration]")
	assertContainsField(t, contentStr, "TLSConfig", "__lens.Lens[net_http.Server, *crypto_tls.Config]")
	assertContainsField(t, contentStr, "TLSConfigO", "__lens_option.LensO[net_http.Server, *crypto_tls.Config]")
	assertContainsField(t, contentStr, "MaxHeaderBytes", "__lens.Lens[net_http.Server, int]")
	assertContainsField(t, contentStr, "MaxHeaderBytesO", "__lens_option.LensO[net_http.Server, int]")

	assert.NotContains(t, contentStr, "TLSNextProtoO __lens_option.LensO")
	assert.NotContains(t, contentStr, "ConnStateO __lens_option.LensO")
	assert.NotContains(t, contentStr, "BaseContextO __lens_option.LensO")
	assert.NotContains(t, contentStr, "ConnContextO __lens_option.LensO")
}

// ---------------------------------------------------------------------------
// TestLensCommandSourcePackageSameNameAsFieldTypePackage — regression test
// for the non-deterministic alias bug where the source package's short name
// (e.g. "v1") collides with a field-type package's short name, causing
// ResolveImportAliases to swap prefixes randomly across runs.
// ---------------------------------------------------------------------------

// TestLensCommandSourcePackageSameNameAsFieldTypePackage builds a tiny
// self-contained Go module with two packages both named "v1":
//
//	example.com/apps/v1    — defines Widget with a field of type meta/v1.Timestamp
//	example.com/meta/v1    — defines Timestamp
//	example.com/lenses     — the output package (different from source)
//
// It then runs lens --type Widget targeting example.com/lenses so that the
// source package must be imported, triggering the QualifiedName path and the
// conflict that caused non-deterministic alias swapping.
func TestLensCommandSourcePackageSameNameAsFieldTypePackage(t *testing.T) {
	// Build a self-contained module in a temp directory.
	root := t.TempDir()

	goMod := "module example.com\n\ngo 1.21\n"
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte(goMod), 0o644))

	// meta/v1 package: defines Timestamp
	metaDir := filepath.Join(root, "meta", "v1")
	require.NoError(t, os.MkdirAll(metaDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(metaDir, "types.go"), []byte(
		"package v1\n\ntype Timestamp struct{ Seconds int64 }\n",
	), 0o644))

	// apps/v1 package: defines Widget with a Timestamp field from meta/v1
	appsDir := filepath.Join(root, "apps", "v1")
	require.NoError(t, os.MkdirAll(appsDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(appsDir, "types.go"), []byte(
		"package v1\n\nimport metav1 \"example.com/meta/v1\"\n\ntype Widget struct {\n\tName      string\n\tCreatedAt metav1.Timestamp\n}\n",
	), 0o644))

	// lenses output package: a separate package within the same module so that
	// go/packages can locate example.com/apps/v1 while the target package differs
	// from the source, causing QualifiedName to be prefixed.
	lensesDir := filepath.Join(root, "lenses")
	require.NoError(t, os.MkdirAll(lensesDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(lensesDir, "doc.go"), []byte(
		"// Package lenses holds generated lens code.\npackage lenses\n",
	), 0o644))

	app := &C.Command{
		Commands: []*C.Command{LensCommand()},
	}

	err := app.Run(context.Background(), []string{
		"app", "lens",
		"--type", "Widget",
		"--dir", lensesDir,
		"--filename", "gen_lens.go",
		"example.com/apps/v1",
	})
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(lensesDir, "gen_lens.go"))
	require.NoError(t, err)
	contentStr := string(content)

	metaAlias := "example_com_meta_v1"

	// The meta/v1 import must appear with its full-path alias (it was disambiguated
	// from apps/v1 which also declares "package v1").
	assert.Contains(t, contentStr, metaAlias+` "example.com/meta/v1"`,
		"import block must use full-path alias for meta/v1")

	// apps/v1 (the source package) may use either the short alias "v1" (acceptable
	// because meta/v1 was given a distinct full-path alias) or the full-path alias.
	// Either way the import line must be present.
	assert.True(t,
		strings.Contains(contentStr, `"example.com/apps/v1"`) || strings.Contains(contentStr, `example_com_apps_v1 "example.com/apps/v1"`),
		"import block must contain an import for example.com/apps/v1")

	// The Widget struct type must appear in function signatures with SOME qualifier
	// (either v1.Widget or example_com_apps_v1.Widget — never without a prefix since
	// the generated package differs from the source package).
	assert.True(t,
		strings.Contains(contentStr, "v1.Widget") || strings.Contains(contentStr, "example_com_apps_v1.Widget"),
		"Widget must appear with its package qualifier")

	// The Timestamp field type must be qualified with the meta alias — never the apps alias.
	// This is the core regression assertion: before the fix, apps/v1's alias was
	// randomly applied to meta/v1 types half of the time.
	assert.Contains(t, contentStr, metaAlias+".Timestamp",
		"Timestamp field must be qualified with the meta/v1 full-path alias")
	assert.NotContains(t, contentStr, "example_com_apps_v1.Timestamp",
		"Timestamp must NOT carry the apps/v1 alias (regression: swapped prefix)")
	// Also verify that "v1.Timestamp" does not appear — both v1 aliases must have been
	// resolved before any TypeName was written.
	assert.NotContains(t, contentStr, "\tv1.Timestamp",
		"Timestamp must not use the ambiguous short alias v1 (should be fully qualified)")
}
