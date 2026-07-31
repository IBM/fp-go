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
	"testing"

	"github.com/IBM/fp-go/gen/v2/cli/lenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// AliasFromPath
// ---------------------------------------------------------------------------

func TestAliasFromPath(t *testing.T) {
	tests := []struct {
		importPath string
		expected   string
	}{
		{"k8s.io/api/apps/v1", "k8s_io_api_apps_v1"},
		{"k8s.io/apimachinery/pkg/apis/meta/v1", "k8s_io_apimachinery_pkg_apis_meta_v1"},
		{"github.com/IBM/fp-go/v2/option", "github_com_IBM_fp_go_v2_option"},
		{"some/simple", "some_simple"},
		{"hyphen-pkg/v1beta1", "hyphen_pkg_v1beta1"},
	}

	for _, tt := range tests {
		t.Run(tt.importPath, func(t *testing.T) {
			assert.Equal(t, tt.expected, lenses.AliasFromPath(tt.importPath))
		})
	}
}

// ---------------------------------------------------------------------------
// ResolveImportAliases
// ---------------------------------------------------------------------------

func TestResolveImportAliasesNoConflict(t *testing.T) {
	structs := []lenses.StructInfo{
		{
			Name: "A", QualifiedName: "appsv1.A",
			Fields:  []lenses.FieldInfo{{Name: "F", TypeName: "appsv1.Deployment", BaseType: "appsv1.Deployment"}},
			Imports: map[string]string{"k8s.io/api/apps/v1": "appsv1"},
		},
		{
			Name: "B", QualifiedName: "metav1.B",
			Fields:  []lenses.FieldInfo{{Name: "G", TypeName: "metav1.ObjectMeta", BaseType: "metav1.ObjectMeta"}},
			Imports: map[string]string{"k8s.io/apimachinery/pkg/apis/meta/v1": "metav1"},
		},
	}

	result := lenses.ResolveImportAliases(structs)
	assert.Nil(t, result, "no conflict should yield nil remapping table")

	assert.Equal(t, "appsv1", structs[0].Imports["k8s.io/api/apps/v1"])
	assert.Equal(t, "appsv1.Deployment", structs[0].Fields[0].TypeName)
	assert.Equal(t, "metav1", structs[1].Imports["k8s.io/apimachinery/pkg/apis/meta/v1"])
	assert.Equal(t, "metav1.ObjectMeta", structs[1].Fields[0].TypeName)
}

func TestResolveImportAliasesConflict(t *testing.T) {
	structs := []lenses.StructInfo{
		{
			Name: "A", QualifiedName: "v1.A",
			Fields:  []lenses.FieldInfo{{Name: "Dep", TypeName: "v1.Deployment", BaseType: "v1.Deployment"}},
			Imports: map[string]string{"k8s.io/api/apps/v1": "v1"},
		},
		{
			Name: "B", QualifiedName: "v1.B",
			Fields: []lenses.FieldInfo{
				{Name: "Meta", TypeName: "v1.ObjectMeta", BaseType: "v1.ObjectMeta"},
				{Name: "PtrMeta", TypeName: "*v1.ObjectMeta", BaseType: "v1.ObjectMeta"},
			},
			Imports: map[string]string{"k8s.io/apimachinery/pkg/apis/meta/v1": "v1"},
		},
	}

	result := lenses.ResolveImportAliases(structs)
	require.NotNil(t, result)

	newAliasA := lenses.AliasFromPath("k8s.io/api/apps/v1")
	assert.Equal(t, newAliasA, structs[0].Imports["k8s.io/api/apps/v1"])
	assert.Equal(t, newAliasA+".A", structs[0].QualifiedName)
	assert.Equal(t, newAliasA+".Deployment", structs[0].Fields[0].TypeName)
	assert.Equal(t, newAliasA+".Deployment", structs[0].Fields[0].BaseType)

	newAliasB := lenses.AliasFromPath("k8s.io/apimachinery/pkg/apis/meta/v1")
	assert.Equal(t, newAliasB, structs[1].Imports["k8s.io/apimachinery/pkg/apis/meta/v1"])
	assert.Equal(t, newAliasB+".B", structs[1].QualifiedName)
	assert.Equal(t, newAliasB+".ObjectMeta", structs[1].Fields[0].TypeName)
	assert.Equal(t, "*"+newAliasB+".ObjectMeta", structs[1].Fields[1].TypeName)
}

// TestResolveImportAliasesConflictSingleStruct is the regression test for the
// double-replacement bug: when a single struct's Imports map contains both
// conflicting paths (e.g. a struct from k8s.io/api/apps/v1 whose fields also
// reference types from k8s.io/apimachinery/pkg/apis/meta/v1), the second
// replacement must not corrupt the already-rewritten QualifiedName / TypeName
// strings by matching the "v1." suffix embedded in "k8s_io_api_apps_v1.".
//
// Note: when a single struct has two conflicting paths both aliased as "v1",
// and all TypeName strings also use "v1.", ResolveImportAliases cannot
// determine which "v1.X" belongs to which package.  The post-fix invariant is
// therefore: (a) both paths get unambiguous full-path aliases, (b) all
// TypeName strings are rewritten to use one of those aliases (no double-mangle
// and no leftover "v1." prefix), and (c) the chosen alias comes from the
// actual import-path set — it will not be a fabricated string.
func TestResolveImportAliasesConflictSingleStruct(t *testing.T) {
	newAliasApps := lenses.AliasFromPath("k8s.io/api/apps/v1")
	newAliasMeta := lenses.AliasFromPath("k8s.io/apimachinery/pkg/apis/meta/v1")

	structs := []lenses.StructInfo{
		{
			Name:          "ControllerRevision",
			QualifiedName: "v1.ControllerRevision",
			Fields: []lenses.FieldInfo{
				{Name: "TypeMeta", TypeName: "v1.TypeMeta", BaseType: "v1.TypeMeta"},
				{Name: "ObjectMeta", TypeName: "v1.ObjectMeta", BaseType: "v1.ObjectMeta"},
			},
			Imports: map[string]string{
				"k8s.io/api/apps/v1":                   "v1",
				"k8s.io/apimachinery/pkg/apis/meta/v1": "v1",
			},
		},
	}

	result := lenses.ResolveImportAliases(structs)
	require.NotNil(t, result)

	// Both import paths must be present with their new (full-path) aliases.
	assert.Equal(t, newAliasApps, structs[0].Imports["k8s.io/api/apps/v1"])
	assert.Equal(t, newAliasMeta, structs[0].Imports["k8s.io/apimachinery/pkg/apis/meta/v1"])

	// QualifiedName must not be double-mangled and must use one of the two valid
	// unambiguous aliases.
	qn := structs[0].QualifiedName
	assert.True(t,
		qn == newAliasApps+".ControllerRevision" || qn == newAliasMeta+".ControllerRevision",
		"QualifiedName %q must be one of the two unambiguous forms, not double-mangled", qn)

	// Every TypeName must have had its "v1." prefix replaced by exactly one of
	// the full-path aliases.  None may still carry the ambiguous "v1." prefix and
	// none may carry a fabricated / double-mangled string.
	for _, f := range structs[0].Fields {
		assert.True(t,
			f.TypeName == newAliasApps+"."+f.Name || f.TypeName == newAliasMeta+"."+f.Name,
			"field %q TypeName %q must use one of the two unambiguous aliases, not double-mangled",
			f.Name, f.TypeName,
		)
		assert.NotContains(t, f.TypeName, "v1.v1.", "TypeName must not be double-prefixed")
	}
}

// TestResolveImportAliasesConflictCrossStruct tests the case where two structs
// each have exactly one conflicting "v1" import.  Because each struct owns a
// distinct import path, ResolveImportAliases can assign the exact right alias
// to each struct independently, and we can make precise per-struct assertions.
func TestResolveImportAliasesConflictCrossStruct(t *testing.T) {
	newAliasApps := lenses.AliasFromPath("k8s.io/api/apps/v1")
	newAliasMeta := lenses.AliasFromPath("k8s.io/apimachinery/pkg/apis/meta/v1")

	structs := []lenses.StructInfo{
		{
			Name:          "ControllerRevision",
			QualifiedName: "v1.ControllerRevision",
			Fields: []lenses.FieldInfo{
				{Name: "Revision", TypeName: "int64", BaseType: "int64"},
				{Name: "Data", TypeName: "v1.RawExtension", BaseType: "v1.RawExtension"},
			},
			Imports: map[string]string{"k8s.io/api/apps/v1": "v1"},
		},
		{
			Name:          "Time",
			QualifiedName: "v1.Time",
			Fields: []lenses.FieldInfo{
				{Name: "Seconds", TypeName: "int64", BaseType: "int64"},
				{Name: "Nanos", TypeName: "v1.Duration", BaseType: "v1.Duration"},
			},
			Imports: map[string]string{"k8s.io/apimachinery/pkg/apis/meta/v1": "v1"},
		},
	}

	result := lenses.ResolveImportAliases(structs)
	require.NotNil(t, result)

	// ControllerRevision belongs to apps/v1 — its QualifiedName and field types
	// referencing "v1." must use the apps alias, not the meta alias.
	assert.Equal(t, newAliasApps, structs[0].Imports["k8s.io/api/apps/v1"])
	assert.Equal(t, newAliasApps+".ControllerRevision", structs[0].QualifiedName)
	assert.Equal(t, "int64", structs[0].Fields[0].TypeName, "primitive field must be unchanged")
	assert.Equal(t, newAliasApps+".RawExtension", structs[0].Fields[1].TypeName)
	assert.Equal(t, newAliasApps+".RawExtension", structs[0].Fields[1].BaseType)

	// Time belongs to meta/v1 — its QualifiedName and field types must use the meta alias.
	assert.Equal(t, newAliasMeta, structs[1].Imports["k8s.io/apimachinery/pkg/apis/meta/v1"])
	assert.Equal(t, newAliasMeta+".Time", structs[1].QualifiedName)
	assert.Equal(t, "int64", structs[1].Fields[0].TypeName, "primitive field must be unchanged")
	assert.Equal(t, newAliasMeta+".Duration", structs[1].Fields[1].TypeName)
	assert.Equal(t, newAliasMeta+".Duration", structs[1].Fields[1].BaseType)
}

// TestResolveImportAliasesConflictPointerAndSliceTypes checks that pointer
// and slice TypeName strings are rewritten correctly under alias conflict.
func TestResolveImportAliasesConflictPointerAndSliceTypes(t *testing.T) {
	newAliasApps := lenses.AliasFromPath("k8s.io/api/apps/v1")

	structs := []lenses.StructInfo{
		{
			Name:          "Wrapper",
			QualifiedName: "v1.Wrapper",
			Fields: []lenses.FieldInfo{
				{Name: "PtrField", TypeName: "*v1.Deployment", BaseType: "v1.Deployment"},
				{Name: "SliceField", TypeName: "[]v1.Pod", BaseType: "v1.Pod"},
				{Name: "MapField", TypeName: "map[string]v1.Service", BaseType: "v1.Service"},
			},
			Imports: map[string]string{"k8s.io/api/apps/v1": "v1"},
		},
		{
			// Second struct creates the conflict: same alias, different path.
			Name:          "Other",
			QualifiedName: "v1.Other",
			Fields:        []lenses.FieldInfo{{Name: "X", TypeName: "v1.ObjectMeta", BaseType: "v1.ObjectMeta"}},
			Imports:       map[string]string{"k8s.io/apimachinery/pkg/apis/meta/v1": "v1"},
		},
	}

	result := lenses.ResolveImportAliases(structs)
	require.NotNil(t, result)

	// Pointer, slice, and map type prefixes must all be rewritten correctly.
	assert.Equal(t, "*"+newAliasApps+".Deployment", structs[0].Fields[0].TypeName)
	assert.Equal(t, newAliasApps+".Deployment", structs[0].Fields[0].BaseType)
	assert.Equal(t, "[]"+newAliasApps+".Pod", structs[0].Fields[1].TypeName)
	assert.Equal(t, "map[string]"+newAliasApps+".Service", structs[0].Fields[2].TypeName)
}
