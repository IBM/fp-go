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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/IBM/fp-go/gen/v2/cli/lenses"
	"golang.org/x/tools/go/packages"
)

// generateLensHelpersByType generates lens code for explicitly named struct types,
// following the stringer pattern: type names are CLI parameters, package loading
// uses go/packages for full type resolution (generics, external field types, tags).
func generateLensHelpersByType(dir, filename string, patterns []string, typeNames []string, packageNameOverride string, verbose bool) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Loading package from %s for types: %s", absDir, strings.Join(typeNames, ", "))
	}

	structs, sourcePackageName, sourcePackagePath, err := parsePackageByTypeNames(absDir, patterns, typeNames, verbose)
	if err != nil {
		return err
	}

	if len(structs) == 0 {
		log.Printf("No matching struct types found for: %s", strings.Join(typeNames, ", "))
		return nil
	}

	// Determine the target package name for generated code
	targetPackageName := packageNameOverride
	targetPackagePath := ""
	if targetPackageName == "" {
		// Derive from existing files in target directory
		targetPackageName, targetPackagePath, err = derivePackageFromDirectory(absDir)
		if err != nil || targetPackageName == "" {
			// Fallback to source package name if no existing files
			targetPackageName = sourcePackageName
			if verbose {
				log.Printf("No existing files in target directory, using source package name: %s", targetPackageName)
			}
		} else if verbose {
			log.Printf("Derived target package name from existing files: %s", targetPackageName)
		}
	} else if verbose {
		log.Printf("Using explicitly provided package name: %s", targetPackageName)
	}

	// If target package path differs from source package path, the source types must
	// be imported and qualified. We compare paths (not names) to correctly handle
	// cases where both packages share the same short name (e.g. both are named "v1"
	// but come from different module paths such as k8s.io/api/apps/v1 vs a local v1).
	if sourcePackagePath != "" && sourcePackagePath != targetPackagePath {
		for i := range structs {
			if structs[i].Imports == nil {
				structs[i].Imports = make(map[string]string)
			}
			// Preserve any pre-assigned alias (e.g. full-path alias from conflict
			// detection) rather than unconditionally overwriting with the short name.
			srcAlias := structs[i].Imports[sourcePackagePath]
			if srcAlias == "" {
				srcAlias = sourcePackageName
			}
			structs[i].Imports[sourcePackagePath] = srcAlias
			// Update QualifiedName to include the (possibly disambiguated) package prefix
			structs[i].QualifiedName = srcAlias + "." + structs[i].Name
		}
		if verbose {
			log.Printf("Added import for source package: %s (%s)", sourcePackageName, sourcePackagePath)
		}
	}

	return lenses.GenerateLensFile(absDir, filename, targetPackageName, structs, verbose)
}

// derivePackageFromDirectory scans existing Go files in the directory to determine
// the package name and import path. Returns empty strings if no Go files are found.
func derivePackageFromDirectory(dir string) (string, string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return "", "", err
	}

	for _, file := range files {
		// Skip generated files and test files
		baseName := filepath.Base(file)
		if strings.HasPrefix(baseName, "gen_") || strings.HasSuffix(baseName, "_test.go") {
			continue
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.PackageClauseOnly)
		if err != nil {
			continue
		}

		if node.Name != nil {
			pkgName := node.Name.Name
			// Use go/packages to resolve the full import path of this directory.
			pkgPath := resolvePackagePath(dir)
			return pkgName, pkgPath, nil
		}
	}

	return "", "", nil
}

// resolvePackagePath returns the Go import path for the given directory by loading
// it with go/packages. Returns an empty string if it cannot be determined.
func resolvePackagePath(dir string) string {
	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  dir,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil || len(pkgs) == 0 {
		return ""
	}
	return pkgs[0].PkgPath
}

// parsePackageByTypeNames loads packages via go/packages and returns structInfo for
// each type name that resolves to a struct in those packages. Returns the structs,
// source package name, and source package path.
func parsePackageByTypeNames(dir string, patterns []string, typeNames []string, verbose bool) ([]structInfo, string, string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
		Dir:  dir,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, "", "", fmt.Errorf("loading packages %v: %w", patterns, err)
	}

	if n := packages.PrintErrors(pkgs); n > 0 {
		return nil, "", "", fmt.Errorf("%d error(s) loading packages", n)
	}

	if len(pkgs) == 0 {
		return nil, "", "", fmt.Errorf("no packages found matching %v", patterns)
	}

	// O(1) lookup set for requested type names
	typeSet := make(map[string]bool, len(typeNames))
	for _, name := range typeNames {
		if name = strings.TrimSpace(name); name != "" {
			typeSet[name] = true
		}
	}

	var structs []structInfo
	var packageName string
	var packagePath string

	for _, pkg := range pkgs {
		if pkg.Types == nil {
			continue
		}
		if packageName == "" {
			packageName = pkg.Name
			packagePath = pkg.PkgPath
		}

		scope := pkg.Types.Scope()
		for _, typName := range scope.Names() {
			if !typeSet[typName] {
				continue
			}

			obj := scope.Lookup(typName)
			typeNameObj, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}

			named, ok := typeNameObj.Type().(*types.Named)
			if !ok {
				continue
			}

			structType, ok := named.Underlying().(*types.Struct)
			if !ok {
				if verbose {
					log.Printf("Type %s is not a struct, skipping", typName)
				}
				continue
			}

			// Build a map of field names to their documentation for deprecation detection.
			// This AST scan only needs to run once regardless of qualifier passes.
			fieldDocs := make(map[string]string)
			for _, file := range pkg.Syntax {
				ast.Inspect(file, func(n ast.Node) bool {
					if ts, ok := n.(*ast.TypeSpec); ok && ts.Name.Name == typName {
						if st, ok := ts.Type.(*ast.StructType); ok {
							for _, field := range st.Fields.List {
								if field.Doc != nil {
									for _, name := range field.Names {
										fieldDocs[name.Name] = field.Doc.Text()
									}
								}
							}
						}
					}
					return true
				})
			}

			// Pass 1: collect all *types.Package references touched by this struct's
			// field types using a simple qualifier (short name). This populates
			// importPkgs with the full set of referenced packages, which we then
			// inspect to detect short-name conflicts (e.g. two different packages
			// both declaring "package v1"). NeedImports without NeedDeps only gives
			// placeholder entries in pkg.Imports, so we must discover conflicts
			// from the actual type-system traversal rather than from pkg.Imports.
			importPkgs := make(map[string]*types.Package)
			basicQualifier := func(p *types.Package) string {
				importPkgs[p.Path()] = p
				return p.Name()
			}
			extractNamedTypeParams(named, basicQualifier)
			extractStructFields(structType, basicQualifier, nil)

			// Detect conflicts: short names that map to more than one import path.
			nameToPath := make(map[string]string, len(importPkgs))
			disambiguate := make(map[string]bool)
			for path, p := range importPkgs {
				shortName := p.Name()
				if existing, seen := nameToPath[shortName]; seen && existing != path {
					disambiguate[shortName] = true
				} else if !seen {
					nameToPath[shortName] = path
				}
			}

			// Pass 2: re-run with a conflict-aware qualifier. Ambiguous short names
			// get the full import-path alias so that generated TypeName strings are
			// unambiguous from the start and require no post-hoc rewriting.
			importPkgs2 := make(map[string]*types.Package, len(importPkgs))
			qualifier := func(p *types.Package) string {
				importPkgs2[p.Path()] = p
				if disambiguate[p.Name()] {
					return lenses.AliasFromPath(p.Path())
				}
				return p.Name()
			}
			typeParams, typeParamNames := extractNamedTypeParams(named, qualifier)
			fields := extractStructFields(structType, qualifier, fieldDocs)

			// Build the imports map using the same alias logic as the qualifier.
			imports := make(map[string]string, len(importPkgs2))
			for path, p := range importPkgs2 {
				if disambiguate[p.Name()] {
					imports[path] = lenses.AliasFromPath(path)
				} else {
					imports[path] = p.Name()
				}
			}

			if len(fields) > 0 {
				structs = append(structs, structInfo{
					Name:           typName,
					QualifiedName:  typName, // Will be updated if target package differs
					TypeParams:     typeParams,
					TypeParamNames: typeParamNames,
					Fields:         fields,
					Imports:        imports,
				})
				if verbose {
					log.Printf("Found struct %s with %d fields", typName, len(fields))
				}
			}
		}
	}

	return structs, packageName, packagePath, nil
}

// extractNamedTypeParams returns the full type-parameter list (e.g. "[T any, K comparable]")
// and the names-only list (e.g. "[T, K]") for a generic named type.
func extractNamedTypeParams(named *types.Named, qualifier types.Qualifier) (string, string) {
	tparams := named.TypeParams()
	if tparams == nil || tparams.Len() == 0 {
		return "", ""
	}

	params := make([]string, 0, tparams.Len())
	names := make([]string, 0, tparams.Len())

	for i := 0; i < tparams.Len(); i++ {
		tp := tparams.At(i)
		name := tp.Obj().Name()
		constraint := types.TypeString(tp.Constraint(), qualifier)
		// go/types renders the "any" alias as "interface{}" — normalize it back.
		if constraint == "interface{}" {
			constraint = "any"
		}
		params = append(params, name+" "+constraint)
		names = append(names, name)
	}

	return "[" + strings.Join(params, ", ") + "]", "[" + strings.Join(names, ", ") + "]"
}

// hasAnonymousStructWithUnexportedFields checks if a type is an anonymous struct
// with unexported fields, which cannot be used in function signatures outside
// the defining package.
func hasAnonymousStructWithUnexportedFields(t types.Type) bool {
	// Unwrap pointer types
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	// Check if it's a struct type (not a named type)
	structType, ok := t.(*types.Struct)
	if !ok {
		return false
	}

	// Check if any field is unexported
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if !field.Exported() {
			return true
		}
	}

	return false
}

// isTypesOptionType reports whether t is an instantiation of
// github.com/IBM/fp-go/v2/option.Option, i.e. a field whose type is already an
// Option.  Generating a LensO for such a field would produce Option[Option[A]].
func isTypesOptionType(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "Option" && obj.Pkg() != nil && obj.Pkg().Path() == "github.com/IBM/fp-go/v2/option"
}

// extractStructFields extracts fieldInfo for every field in a struct, promoting
// embedded struct fields (one level deep, same pattern as the annotation scanner).
// fieldDocs maps field names to their documentation comments for deprecation detection.
func extractStructFields(structType *types.Struct, qualifier types.Qualifier, fieldDocs map[string]string) []fieldInfo {
	var fields []fieldInfo

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)

		if field.Embedded() {
			// Promote fields from the embedded type.
			embType := field.Type()
			if ptr, ok := embType.(*types.Pointer); ok {
				embType = ptr.Elem()
			}
			if named, ok := embType.(*types.Named); ok {
				if embStruct, ok := named.Underlying().(*types.Struct); ok {
					// Embedded fields don't have their own docs in the parent struct
					for _, embField := range extractStructFields(embStruct, qualifier, nil) {
						embField.IsEmbedded = true
						fields = append(fields, embField)
					}
				}
			}
			continue
		}

		// Skip fields with anonymous struct types that have unexported fields
		// These cannot be used in function signatures outside their defining package
		if hasAnonymousStructWithUnexportedFields(field.Type()) {
			continue
		}

		typeName := types.TypeString(field.Type(), qualifier)
		isPointer := false
		baseType := typeName

		if _, ok := field.Type().(*types.Pointer); ok {
			isPointer = true
			baseType = strings.TrimPrefix(typeName, "*")
		}

		isOptional := isPointer || hasOmitEmptyStringTag(tag)
		isComparable := types.Comparable(field.Type())
		isOption := isTypesOptionType(field.Type())

		// Check if field is deprecated by looking for "Deprecated:" in its documentation
		isDeprecated := false
		if fieldDocs != nil {
			if doc, ok := fieldDocs[field.Name()]; ok {
				isDeprecated = strings.Contains(doc, "Deprecated:")
			}
		}

		fields = append(fields, fieldInfo{
			Name:         field.Name(),
			TypeName:     typeName,
			BaseType:     baseType,
			IsOptional:   isOptional,
			IsComparable: isComparable,
			IsOption:     isOption,
			IsEmbedded:   false,
			IsDeprecated: isDeprecated,
		})
	}

	return fields
}

// hasOmitEmptyStringTag reports whether a raw struct tag string contains json:"...,omitempty".
func hasOmitEmptyStringTag(tag string) bool {
	if tag == "" {
		return false
	}
	jsonTag := reflect.StructTag(tag).Get("json")
	for part := range strings.SplitSeq(jsonTag, ",") {
		if strings.TrimSpace(part) == "omitempty" {
			return true
		}
	}
	return false
}
