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

package lenses

import "strings"

// AliasFromPath converts an import path into a valid Go identifier that is unique
// per path. Each path component separator (/, . and -) is replaced with an underscore.
//
// Example: "k8s.io/apimachinery/pkg/apis/meta/v1" → "k8s_io_apimachinery_pkg_apis_meta_v1"
func AliasFromPath(importPath string) string {
	r := strings.NewReplacer("/", "_", ".", "_", "-", "_")
	return r.Replace(importPath)
}

// ResolveImportAliases detects cases where two or more different import paths share
// the same short alias (e.g. both "k8s.io/api/apps/v1" and
// "k8s.io/apimachinery/pkg/apis/meta/v1" declare `package v1`). For every such
// conflicting alias it returns a remapping table that maps each affected import path
// to a disambiguated alias derived from the full import path. The structs slice is
// updated in-place: TypeName / BaseType / QualifiedName strings are rewritten to use
// the new alias, and the Imports map of each struct is updated accordingly.
func ResolveImportAliases(structs []StructInfo) map[string]string {
	// Step 1: collect alias → set-of-paths across all structs
	aliasToPaths := make(map[string]map[string]struct{})
	for _, s := range structs {
		for path, alias := range s.Imports {
			if aliasToPaths[alias] == nil {
				aliasToPaths[alias] = make(map[string]struct{})
			}
			aliasToPaths[alias][path] = struct{}{}
		}
	}

	// Step 2: for every alias with >1 path, map each path to a unique alias
	pathToNewAlias := make(map[string]string)
	for _, paths := range aliasToPaths {
		if len(paths) <= 1 {
			continue
		}
		for path := range paths {
			pathToNewAlias[path] = AliasFromPath(path)
		}
	}

	if len(pathToNewAlias) == 0 {
		return nil
	}

	// Step 3: rewrite each struct in-place.
	// Build all replacements for a struct first, then apply them atomically via
	// strings.NewReplacer so that an already-rewritten prefix (e.g. the "v1." tail
	// of "k8s_io_api_apps_v1.") is never matched by a subsequent iteration.
	for i := range structs {
		s := &structs[i]
		var pairs []string // alternating old, new pairs for strings.NewReplacer
		for path, newAlias := range pathToNewAlias {
			oldAlias, ok := s.Imports[path]
			if !ok || oldAlias == newAlias {
				continue
			}
			pairs = append(pairs, oldAlias+".", newAlias+".")
			s.Imports[path] = newAlias
		}
		if len(pairs) == 0 {
			continue
		}
		r := strings.NewReplacer(pairs...)
		for j := range s.Fields {
			s.Fields[j].TypeName = r.Replace(s.Fields[j].TypeName)
			s.Fields[j].BaseType = r.Replace(s.Fields[j].BaseType)
		}
		s.QualifiedName = r.Replace(s.QualifiedName)
	}

	return pathToNewAlias
}
