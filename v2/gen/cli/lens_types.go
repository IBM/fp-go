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
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"reflect"
	"strings"

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
	if targetPackageName == "" {
		// Derive from existing files in target directory
		targetPackageName, err = derivePackageNameFromDirectory(absDir)
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

	// If target package differs from source package, add import for source package
	// and update QualifiedName to include package prefix
	if targetPackageName != sourcePackageName && sourcePackagePath != "" {
		for i := range structs {
			if structs[i].Imports == nil {
				structs[i].Imports = make(map[string]string)
			}
			structs[i].Imports[sourcePackagePath] = sourcePackageName
			// Update QualifiedName to include package prefix
			structs[i].QualifiedName = sourcePackageName + "." + structs[i].Name
		}
		if verbose {
			log.Printf("Added import for source package: %s (%s)", sourcePackageName, sourcePackagePath)
		}
	}

	return generateLensFile(absDir, filename, targetPackageName, structs, verbose)
}

// derivePackageNameFromDirectory scans existing Go files in the directory to determine
// the package name. Returns empty string if no Go files are found.
func derivePackageNameFromDirectory(dir string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return "", err
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
			return node.Name.Name, nil
		}
	}

	return "", nil
}

// parsePackageByTypeNames loads packages via go/packages and returns structInfo for
// each type name that resolves to a struct in those packages. Returns the structs,
// source package name, and source package path.
func parsePackageByTypeNames(dir string, patterns []string, typeNames []string, verbose bool) ([]structInfo, string, string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports,
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

			// importPkgs accumulates external packages referenced in field types.
			// The qualifier closure populates it as types are stringified.
			importPkgs := make(map[string]*types.Package)
			qualifier := func(p *types.Package) string {
				if p == pkg.Types {
					return "" // same package — no qualifier needed
				}
				importPkgs[p.Path()] = p
				return p.Name()
			}

			typeParams, typeParamNames := extractNamedTypeParams(named, qualifier)
			fields := extractStructFields(structType, qualifier)

			imports := make(map[string]string, len(importPkgs))
			for path, p := range importPkgs {
				imports[path] = p.Name()
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

// extractStructFields extracts fieldInfo for every field in a struct, promoting
// embedded struct fields (one level deep, same pattern as the annotation scanner).
func extractStructFields(structType *types.Struct, qualifier types.Qualifier) []fieldInfo {
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
					for _, embField := range extractStructFields(embStruct, qualifier) {
						embField.IsEmbedded = true
						fields = append(fields, embField)
					}
				}
			}
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

		fields = append(fields, fieldInfo{
			Name:         field.Name(),
			TypeName:     typeName,
			BaseType:     baseType,
			IsOptional:   isOptional,
			IsComparable: isComparable,
			IsEmbedded:   false,
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
