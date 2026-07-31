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

	// The critical regression check: QualifiedName must not be double-mangled.
	// It must be one of the two valid unambiguous aliases, not a nested form like
	// "k8s_io_apimachinery_pkg_apis_meta_k8s_io_api_apps_v1.ControllerRevision".
	qn := structs[0].QualifiedName
	valid := qn == newAliasApps+".ControllerRevision" || qn == newAliasMeta+".ControllerRevision"
	assert.True(t, valid, "QualifiedName %q must be one of the two unambiguous forms, not double-mangled", qn)
}
