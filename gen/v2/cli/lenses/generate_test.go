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
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/IBM/fp-go/gen/v2/cli/lenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Template execution
// ---------------------------------------------------------------------------

func TestStructTemplate(t *testing.T) {
	s := lenses.StructInfo{
		Name:          "TestStruct",
		QualifiedName: "TestStruct",
		Fields: []lenses.FieldInfo{
			{Name: "Name", TypeName: "string", IsComparable: true},
			{Name: "Value", TypeName: "*int", IsOptional: true, IsComparable: true},
		},
	}

	var buf bytes.Buffer
	err := lenses.StructTmpl.Execute(&buf, s)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "type TestStructLenses struct")
	assert.Contains(t, out, "Name __lens.Lens[TestStruct, string]")
	assert.Contains(t, out, "NameO __lens_option.LensO[TestStruct, string]")
	assert.Contains(t, out, "Value __lens.Lens[TestStruct, *int]")
	assert.Contains(t, out, "ValueO __lens_option.LensO[TestStruct, *int]")
}

func TestStandaloneTemplate(t *testing.T) {
	s := lenses.StructInfo{
		Name:          "TestStruct",
		QualifiedName: "TestStruct",
		Fields: []lenses.FieldInfo{
			{Name: "Name", TypeName: "string", IsComparable: true},
			{Name: "Value", TypeName: "*int", IsOptional: true, IsComparable: true},
		},
	}

	var buf bytes.Buffer
	err := lenses.StandaloneTmpl.Execute(&buf, s)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "func MakeTestStructNameLens() __lens.Lens[TestStruct, string]")
	assert.Contains(t, out, "func MakeTestStructNameLensO() __lens_option.LensO[TestStruct, string]")
	assert.Contains(t, out, "func MakeTestStructValueLens() __lens.Lens[TestStruct, *int]")
	assert.Contains(t, out, "func MakeTestStructValueLensO() __lens_option.LensO[TestStruct, *int]")
	assert.Contains(t, out, "__iso_option.FromZero")
}

func TestConstructorTemplate(t *testing.T) {
	s := lenses.StructInfo{
		Name:          "TestStruct",
		QualifiedName: "TestStruct",
		Fields: []lenses.FieldInfo{
			{Name: "Name", TypeName: "string", IsComparable: true},
			{Name: "Value", TypeName: "*int", IsOptional: true, IsComparable: true},
		},
	}

	var buf bytes.Buffer
	err := lenses.ConstructorTmpl.Execute(&buf, s)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "func MakeTestStructLenses() TestStructLenses")
	assert.Contains(t, out, "return TestStructLenses{")
	assert.Contains(t, out, "Name: MakeTestStructNameLens(),")
	assert.Contains(t, out, "NameO: MakeTestStructNameLensO(),")
	assert.Contains(t, out, "Value: MakeTestStructValueLens(),")
	assert.Contains(t, out, "ValueO: MakeTestStructValueLensO(),")
}

// TestRefLensTemplateUsesCorrectConstructor verifies that the standalone template
// selects MakeLensStrictWithName for comparable fields and MakeLensRefWithName for
// non-comparable fields.
func TestRefLensTemplateUsesCorrectConstructor(t *testing.T) {
	s := lenses.StructInfo{
		Name:          "TestStruct",
		QualifiedName: "TestStruct",
		Fields: []lenses.FieldInfo{
			{Name: "Name", TypeName: "string", IsComparable: true},
			{Name: "Data", TypeName: "[]byte", IsComparable: false},
		},
	}

	var buf bytes.Buffer
	require.NoError(t, lenses.StandaloneTmpl.Execute(&buf, s))
	out := buf.String()

	assert.Contains(t, out, "return __lens.MakeLensStrictWithName(", "comparable field should use MakeLensStrictWithName")
	assert.Contains(t, out, "return __lens.MakeLensRefWithName(", "non-comparable field should use MakeLensRefWithName")
}

// TestTemplateOptionFieldNoLensO verifies that IsOption=true suppresses LensO generation.
func TestTemplateOptionFieldNoLensO(t *testing.T) {
	s := lenses.StructInfo{
		Name:          "Config",
		QualifiedName: "Config",
		Fields: []lenses.FieldInfo{
			{Name: "Name", TypeName: "string", IsComparable: true},
			{Name: "Value", TypeName: "option.Option[string]", IsComparable: true, IsOption: true},
			{Name: "Count", TypeName: "int", IsComparable: true},
		},
	}

	var structBuf bytes.Buffer
	require.NoError(t, lenses.StructTmpl.Execute(&structBuf, s))
	structOut := structBuf.String()

	// Name and Count should have LensO, Value must not.
	assert.Contains(t, structOut, "NameO __lens_option.LensO[Config, string]")
	assert.Contains(t, structOut, "CountO __lens_option.LensO[Config, int]")
	assert.NotContains(t, structOut, "ValueO __lens_option.LensO", "Option-typed field must not produce LensO")
}

// ---------------------------------------------------------------------------
// GenerateLensFile
// ---------------------------------------------------------------------------

func TestGenerateLensFileImportDisambiguation(t *testing.T) {
	newAliasA := lenses.AliasFromPath("k8s.io/api/apps/v1")
	newAliasB := lenses.AliasFromPath("k8s.io/apimachinery/pkg/apis/meta/v1")

	structs := []lenses.StructInfo{
		{
			Name:          "DeploymentWrapper",
			QualifiedName: "v1.DeploymentWrapper",
			Fields:        []lenses.FieldInfo{{Name: "Dep", TypeName: "v1.Deployment", BaseType: "v1.Deployment", IsComparable: false}},
			Imports:       map[string]string{"k8s.io/api/apps/v1": "v1"},
		},
		{
			Name:          "MetaWrapper",
			QualifiedName: "v1.MetaWrapper",
			Fields:        []lenses.FieldInfo{{Name: "Meta", TypeName: "v1.ObjectMeta", BaseType: "v1.ObjectMeta", IsComparable: false}},
			Imports:       map[string]string{"k8s.io/apimachinery/pkg/apis/meta/v1": "v1"},
		},
	}

	tmpDir := t.TempDir()
	err := lenses.GenerateLensFile(tmpDir, "gen_lens.go", "mypkg", structs, false)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, "gen_lens.go"))
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, newAliasA+` "k8s.io/api/apps/v1"`)
	assert.Contains(t, contentStr, newAliasB+` "k8s.io/apimachinery/pkg/apis/meta/v1"`)
	assert.NotContains(t, contentStr, "\tv1 \"k8s.io/api/apps/v1\"")
	assert.NotContains(t, contentStr, "\tv1 \"k8s.io/apimachinery/pkg/apis/meta/v1\"")
	assert.Contains(t, contentStr, newAliasA+".Deployment")
	assert.Contains(t, contentStr, newAliasB+".ObjectMeta")
}
