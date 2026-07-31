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

// TestParsePackageByTypeNamesFullPathAlias verifies that all imported packages
// always receive full-path aliases, regardless of whether short names conflict.
func TestParsePackageByTypeNamesFullPathAlias(t *testing.T) {
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

	// Both packages must use full-path aliases.
	const metaAlias = "example_com_meta_v1"
	assert.Equal(t, metaAlias, s.Imports["example.com/meta/v1"],
		"meta/v1 must have a full-path alias")

	// The Timestamp field must use the full-path alias.
	require.Len(t, s.Fields, 2)
	assert.Equal(t, "Name", s.Fields[0].Name)
	assert.Equal(t, "string", s.Fields[0].TypeName)

	assert.Equal(t, "CreatedAt", s.Fields[1].Name)
	assert.Equal(t, metaAlias+".Timestamp", s.Fields[1].TypeName,
		"CreatedAt TypeName must use the full-path alias for meta/v1")
	assert.Equal(t, metaAlias+".Timestamp", s.Fields[1].BaseType,
		"CreatedAt BaseType must use the full-path alias for meta/v1")
}

// TestParsePackageByTypeNamesAlwaysFullPath verifies that even when no short-name
// collision exists, all imports still receive full-path aliases.
func TestParsePackageByTypeNamesAlwaysFullPath(t *testing.T) {
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
	// crypto/tls is a field-type package of http.Server; it must use a full-path alias.
	if alias, ok := s.Imports["crypto/tls"]; ok {
		assert.Equal(t, "crypto_tls", alias,
			"crypto/tls must use the full-path alias")
	}
	// Verify that no field TypeName uses the bare "tls." package qualifier.
	// (The correct form is "crypto_tls." — the full-path alias.)
	for _, f := range s.Fields {
		assert.NotRegexp(t, `(^|[\s\*\[])tls\.`, f.TypeName,
			"field %q TypeName %q must not use the bare 'tls.' qualifier", f.Name, f.TypeName)
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

// TestGenerateLensHelpersByTypeSourcePkgAlias verifies that the generated file
// uses full-path aliases for all packages (both the source package and any
// field-type packages), regardless of whether short names conflict.
func TestGenerateLensHelpersByTypeSourcePkgAlias(t *testing.T) {
	_, root := makeTwoV1Module(t)

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

		// Both packages must appear with full-path aliases.
		assert.Contains(t, contentStr, `example_com_apps_v1 "example.com/apps/v1"`,
			"run %d: apps/v1 must use the full-path alias", i)
		assert.Contains(t, contentStr, `example_com_meta_v1 "example.com/meta/v1"`,
			"run %d: meta/v1 must use the full-path alias", i)
		assert.Contains(t, contentStr, "example_com_apps_v1.Widget",
			"run %d: Widget must use the apps/v1 full-path alias", i)
		assert.Contains(t, contentStr, "example_com_meta_v1.Timestamp",
			"run %d: Timestamp must use the meta/v1 full-path alias", i)
		assert.NotContains(t, contentStr, "\tv1.",
			"run %d: no bare 'v1.' qualifier must appear in generated code", i)
	}
}

