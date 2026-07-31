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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// hasOmitEmptyStringTag
// ---------------------------------------------------------------------------

func TestHasOmitEmptyStringTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected bool
	}{
		{name: "empty tag", tag: "", expected: false},
		{name: "no json tag", tag: `xml:"field"`, expected: false},
		{name: "json without omitempty", tag: `json:"field"`, expected: false},
		{name: "json with omitempty", tag: `json:"field,omitempty"`, expected: true},
		{name: "json omitempty before other option", tag: `json:"field,omitempty,string"`, expected: true},
		{name: "json omitempty after other option", tag: `json:"field,string,omitempty"`, expected: true},
		{name: "json omitempty only", tag: `json:",omitempty"`, expected: true},
		{name: "json dash", tag: `json:"-"`, expected: false},
		{name: "mixed tags omitempty in json", tag: `xml:"field" json:"f,omitempty"`, expected: true},
		{name: "mixed tags omitempty not in json", tag: `json:"f" xml:"field,omitempty"`, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, hasOmitEmptyStringTag(tt.tag))
		})
	}
}

// ---------------------------------------------------------------------------
// parsePackageByTypeNames — disambiguation: source-package short name vs
// field-type package short name
//
// This test directly exercises the fix in parsePackageByTypeNames: the source
// package is now seeded into importPkgs so that its short name is considered
// during the conflict-detection pass (Pass 1).  Without the fix, when the
// source package's short name (e.g. "v1") matches a field-type package's short
// name, both collide only after generateLensHelpersByType adds the source
// package to the imports map, leaving ResolveImportAliases with two paths
// sharing the same alias and a non-deterministic strings.NewReplacer outcome.
// ---------------------------------------------------------------------------

// makeTwoV1Module writes a self-contained Go module with:
//
//	example.com/apps/v1  — defines Widget{ Name string; CreatedAt meta/v1.Timestamp }
//	example.com/meta/v1  — defines Timestamp{ Seconds int64 }
//
// and returns the path to the apps/v1 directory (used as the Dir for
// packages.Config) and the module root.
func makeTwoV1Module(t *testing.T) (appsDir, root string) {
	t.Helper()
	root = t.TempDir()

	require.NoError(t, os.WriteFile(
		filepath.Join(root, "go.mod"),
		[]byte("module example.com\n\ngo 1.21\n"),
		0o644,
	))

	metaDir := filepath.Join(root, "meta", "v1")
	require.NoError(t, os.MkdirAll(metaDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(metaDir, "types.go"),
		[]byte("package v1\n\ntype Timestamp struct{ Seconds int64 }\n"), 0o644))

	appsDir = filepath.Join(root, "apps", "v1")
	require.NoError(t, os.MkdirAll(appsDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(appsDir, "types.go"),
		[]byte("package v1\n\nimport metav1 \"example.com/meta/v1\"\n\ntype Widget struct {\n\tName      string\n\tCreatedAt metav1.Timestamp\n}\n"), 0o644))

	return appsDir, root
}

// TestParsePackageByTypeNamesSourcePkgConflictDetection verifies that when the
// source package's short name conflicts with a field-type package's short name,
// the returned StructInfo uses full-path aliases for the field types so that
// the subsequent generateLensHelpersByType step never creates two import paths
// with the same alias.
func TestParsePackageByTypeNamesSourcePkgConflictDetection(t *testing.T) {
	appsDir, _ := makeTwoV1Module(t)

	structs, srcPkgName, srcPkgPath, err := parsePackageByTypeNames(
		appsDir,
		[]string{"example.com/apps/v1"},
		[]string{"Widget"},
		false,
	)
	require.NoError(t, err)
	require.Len(t, structs, 1, "should find the Widget struct")

	assert.Equal(t, "v1", srcPkgName, "source package short name")
	assert.Equal(t, "example.com/apps/v1", srcPkgPath, "source package path")

	s := structs[0]
	assert.Equal(t, "Widget", s.Name)

	// The meta/v1 field-type package must use a full-path alias because it
	// collides with the source package's short name "v1".
	const metaAlias = "example_com_meta_v1"
	assert.Equal(t, metaAlias, s.Imports["example.com/meta/v1"],
		"meta/v1 must have a full-path alias due to short-name conflict with source pkg")

	// The Timestamp field must use the full-path alias — not the ambiguous "v1.".
	require.Len(t, s.Fields, 2)
	nameField := s.Fields[0]
	createdAtField := s.Fields[1]
	assert.Equal(t, "Name", nameField.Name)
	assert.Equal(t, "string", nameField.TypeName)

	assert.Equal(t, "CreatedAt", createdAtField.Name)
	assert.Equal(t, metaAlias+".Timestamp", createdAtField.TypeName,
		"CreatedAt TypeName must use the full-path alias for meta/v1")
	assert.Equal(t, metaAlias+".Timestamp", createdAtField.BaseType,
		"CreatedAt BaseType must use the full-path alias for meta/v1")
}

// TestParsePackageByTypeNamesNoConflict verifies that when no short-name
// collision exists the short name is preserved as-is (no unnecessary aliasing).
func TestParsePackageByTypeNamesNoConflict(t *testing.T) {
	// net/http.Server has fields from crypto/tls (short "tls"), context (short
	// "context"), time (short "time") — none of these share the short name
	// "http", so no disambiguation should fire.
	structs, srcPkgName, _, err := parsePackageByTypeNames(
		".", // current dir — stdlib packages are always resolvable
		[]string{"net/http"},
		[]string{"Server"},
		false,
	)
	require.NoError(t, err)
	require.Len(t, structs, 1)
	assert.Equal(t, "http", srcPkgName)

	s := structs[0]
	// tls.Config is a field of http.Server; its alias must remain the short name.
	if alias, ok := s.Imports["crypto/tls"]; ok {
		assert.Equal(t, "tls", alias,
			"crypto/tls must keep short alias 'tls' when no conflict exists")
	}
	// Verify that field TypeNames use the short name, not a full-path alias.
	for _, f := range s.Fields {
		assert.NotContains(t, f.TypeName, "crypto_tls",
			"TypeName should not contain a full-path alias when there is no conflict")
	}
}

// TestParsePackageByTypeNamesDisambiguationIsDeterministic runs
// parsePackageByTypeNames 20 times for the conflicting case and asserts that
// every run produces identical output, proving that the fix eliminated the
// non-determinism caused by map iteration over importPkgs.
func TestParsePackageByTypeNamesDisambiguationIsDeterministic(t *testing.T) {
	appsDir, _ := makeTwoV1Module(t)

	const metaAlias = "example_com_meta_v1"
	const runs = 20

	for i := 0; i < runs; i++ {
		structs, _, _, err := parsePackageByTypeNames(
			appsDir,
			[]string{"example.com/apps/v1"},
			[]string{"Widget"},
			false,
		)
		require.NoError(t, err, "run %d", i)
		require.Len(t, structs, 1, "run %d", i)

		s := structs[0]
		require.Len(t, s.Fields, 2, "run %d", i)
		createdAt := s.Fields[1]
		assert.Equal(t, metaAlias+".Timestamp", createdAt.TypeName,
			"run %d: CreatedAt.TypeName must be deterministic", i)
	}
}

// ---------------------------------------------------------------------------
// generateLensHelpersByType — source-package import alias after disambiguation
// ---------------------------------------------------------------------------

// TestGenerateLensHelpersByTypeSourcePkgAlias verifies the end-to-end alias
// assignment when the source package's short name collides with a field-type
// package's short name.  After the fix, the generated file must import meta/v1
// with a full-path alias and the Timestamp field must use that alias — never
// the apps/v1 alias.
func TestGenerateLensHelpersByTypeSourcePkgAlias(t *testing.T) {
	appsDir, root := makeTwoV1Module(t)

	// Output into a sibling "lenses" package so sourcePackagePath ≠ targetPackagePath,
	// forcing generateLensHelpersByType to add the source-package import.
	lensesDir := filepath.Join(root, "lenses")
	require.NoError(t, os.MkdirAll(lensesDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(lensesDir, "doc.go"),
		[]byte("// Package lenses holds generated lens code.\npackage lenses\n"), 0o644))

	// Run 20 times to catch any residual non-determinism.
	for i := 0; i < 20; i++ {
		err := generateLensHelpersByType(
			lensesDir,
			"gen_lens.go",
			[]string{"example.com/apps/v1"},
			[]string{"Widget"},
			"",
			false,
		)
		require.NoError(t, err, "run %d", i)

		content, err := os.ReadFile(filepath.Join(lensesDir, "gen_lens.go"))
		require.NoError(t, err, "run %d", i)
		contentStr := string(content)

		_ = appsDir // used above

		assert.Contains(t, contentStr, `"example.com/meta/v1"`,
			"run %d: meta/v1 import must be present", i)
		assert.Contains(t, contentStr, "example_com_meta_v1.Timestamp",
			"run %d: Timestamp must use the meta/v1 full-path alias", i)
		assert.NotContains(t, contentStr, "example_com_apps_v1.Timestamp",
			"run %d: Timestamp must NOT use the apps/v1 alias (swapped prefix regression)", i)
		assert.NotContains(t, contentStr, "\tv1.Timestamp",
			"run %d: Timestamp must not use the ambiguous short alias v1", i)
	}
}

