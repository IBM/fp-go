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
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	S "github.com/IBM/fp-go/v2/string"
	C "github.com/urfave/cli/v3"
)

const (
	keyLensDir         = "dir"
	keyVerbose         = "verbose"
	keyIncludeTestFile = "include-test-files"
	lensAnnotation     = "fp-go:Lens"
)

var (
	flagLensDir = &C.StringFlag{
		Name:  keyLensDir,
		Value: ".",
		Usage: "Directory to scan for Go files",
	}

	flagVerbose = &C.BoolFlag{
		Name:    keyVerbose,
		Aliases: []string{"v"},
		Value:   false,
		Usage:   "Enable verbose output",
	}

	flagIncludeTestFiles = &C.BoolFlag{
		Name:    keyIncludeTestFile,
		Aliases: []string{"t"},
		Value:   false,
		Usage:   "Include test files (*_test.go) when scanning for annotated types",
	}
)

// structInfo holds information about a struct that needs lens generation
type structInfo struct {
	Name           string
	TypeParams     string // e.g., "[T any]" or "[K comparable, V any]" - for type declarations
	TypeParamNames string // e.g., "[T]" or "[K, V]" - for type usage in function signatures
	Fields         []fieldInfo
	Imports        map[string]string // package path -> alias
}

// fieldInfo holds information about a struct field
type fieldInfo struct {
	Name         string
	TypeName     string
	BaseType     string // TypeName without leading * for pointer types
	IsOptional   bool   // true if field is a pointer or has json omitempty tag
	IsComparable bool   // true if the type is comparable (can use ==)
	IsEmbedded   bool   // true if this field comes from an embedded struct
}

// templateData holds data for template rendering
type templateData struct {
	PackageName string
	Structs     []structInfo
}

const lensStructTemplate = `
// {{.Name}}Lenses provides [lenses] for accessing fields of [{{.Name}}]
//
// [lenses]: __lens.Lens
type {{.Name}}Lenses{{.TypeParams}} struct {
	// mandatory fields
{{- range .Fields}}
	{{.Name}} __lens.Lens[{{$.Name}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
	// optional fields
{{- range .Fields}}
{{- if .IsComparable}}
	{{.Name}}O __lens_option.LensO[{{$.Name}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
{{- end}}
}

// {{.Name}}RefLenses provides [lenses] for accessing fields of [{{.Name}}] via a reference to [{{.Name}}]
//
//
// [lenses]: __lens.Lens
type {{.Name}}RefLenses{{.TypeParams}} struct {
	// mandatory fields
{{- range .Fields}}
	{{.Name}} __lens.Lens[*{{$.Name}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
	// optional fields
{{- range .Fields}}
{{- if .IsComparable}}
	{{.Name}}O __lens_option.LensO[*{{$.Name}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
{{- end}}
}

// {{.Name}}Prisms provides [prisms] for accessing fields of [{{.Name}}]
//
// [prisms]: __prism.Prism
type {{.Name}}Prisms{{.TypeParams}} struct {
{{- range .Fields}}
	{{.Name}} __prism.Prism[{{$.Name}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
}

// {{.Name}}RefPrisms provides [prisms] for accessing fields of [{{.Name}}] via a reference to [{{.Name}}]
//
// [prisms]: __prism.Prism
type {{.Name}}RefPrisms{{.TypeParams}} struct {
{{- range .Fields}}
	{{.Name}} __prism.Prism[*{{$.Name}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
}
`

const lensConstructorTemplate = `
// Make{{.Name}}Lenses creates a new [{{.Name}}Lenses] with [lenses] for all fields
//
// [lenses]:__lens.Lens
func Make{{.Name}}Lenses{{.TypeParams}}() {{.Name}}Lenses{{.TypeParamNames}} {
	// mandatory lenses
{{- range .Fields}}
	lens{{.Name}} := __lens.MakeLensWithName(
		func(s {{$.Name}}{{$.TypeParamNames}}) {{.TypeName}} { return s.{{.Name}} },
		func(s {{$.Name}}{{$.TypeParamNames}}, v {{.TypeName}}) {{$.Name}}{{$.TypeParamNames}} { s.{{.Name}} = v; return s },
		"{{$.Name}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- end}}
	// optional lenses
{{- range .Fields}}
{{- if .IsComparable}}
	lens{{.Name}}O := __lens_option.FromIso[{{$.Name}}{{$.TypeParamNames}}](__iso_option.FromZero[{{.TypeName}}]())(lens{{.Name}})
{{- end}}
{{- end}}
	return {{.Name}}Lenses{{.TypeParamNames}}{
		// mandatory lenses
{{- range .Fields}}
		{{.Name}}: lens{{.Name}},
{{- end}}
		// optional lenses
{{- range .Fields}}
{{- if .IsComparable}}
		{{.Name}}O: lens{{.Name}}O,
{{- end}}
{{- end}}
	}
}

// Make{{.Name}}RefLenses creates a new [{{.Name}}RefLenses] with [lenses] for all fields
//
// [lenses]:__lens.Lens
func Make{{.Name}}RefLenses{{.TypeParams}}() {{.Name}}RefLenses{{.TypeParamNames}} {
	// mandatory lenses
{{- range .Fields}}
{{- if .IsComparable}}
	lens{{.Name}} := __lens.MakeLensStrictWithName(
		func(s *{{$.Name}}{{$.TypeParamNames}}) {{.TypeName}} { return s.{{.Name}} },
		func(s *{{$.Name}}{{$.TypeParamNames}}, v {{.TypeName}}) *{{$.Name}}{{$.TypeParamNames}} { s.{{.Name}} = v; return s },
		"(*{{$.Name}}{{$.TypeParamNames}}).{{.Name}}",
	)
{{- else}}
	lens{{.Name}} := __lens.MakeLensRefWithName(
		func(s *{{$.Name}}{{$.TypeParamNames}}) {{.TypeName}} { return s.{{.Name}} },
		func(s *{{$.Name}}{{$.TypeParamNames}}, v {{.TypeName}}) *{{$.Name}}{{$.TypeParamNames}} { s.{{.Name}} = v; return s },
		"(*{{$.Name}}{{$.TypeParamNames}}).{{.Name}}",
	)
{{- end}}
{{- end}}
	// optional lenses
{{- range .Fields}}
{{- if .IsComparable}}
	lens{{.Name}}O := __lens_option.FromIso[*{{$.Name}}{{$.TypeParamNames}}](__iso_option.FromZero[{{.TypeName}}]())(lens{{.Name}})
{{- end}}
{{- end}}
	return {{.Name}}RefLenses{{.TypeParamNames}}{
		// mandatory lenses
{{- range .Fields}}
		{{.Name}}: lens{{.Name}},
{{- end}}
		// optional lenses
{{- range .Fields}}
{{- if .IsComparable}}
		{{.Name}}O: lens{{.Name}}O,
{{- end}}
{{- end}}
	}
}

// Make{{.Name}}Prisms creates a new [{{.Name}}Prisms] with [prisms] for all fields
//
// [prisms]:__prism.Prism
func Make{{.Name}}Prisms{{.TypeParams}}() {{.Name}}Prisms{{.TypeParamNames}} {
{{- range .Fields}}
{{- if .IsComparable}}
	_fromNonZero{{.Name}} := __option.FromNonZero[{{.TypeName}}]()
	_prism{{.Name}} := __prism.MakePrismWithName(
		func(s {{$.Name}}{{$.TypeParamNames}}) __option.Option[{{.TypeName}}] { return _fromNonZero{{.Name}}(s.{{.Name}}) },
		func(v {{.TypeName}}) {{$.Name}}{{$.TypeParamNames}} {
			{{- if .IsEmbedded}}
			var result {{$.Name}}{{$.TypeParamNames}}
			result.{{.Name}} = v
			return result
			{{- else}}
			return {{$.Name}}{{$.TypeParamNames}}{ {{.Name}}: v }
			{{- end}}
		},
		"{{$.Name}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- else}}
	_prism{{.Name}} := __prism.MakePrismWithName(
		func(s {{$.Name}}{{$.TypeParamNames}}) __option.Option[{{.TypeName}}] { return __option.Some(s.{{.Name}}) },
		func(v {{.TypeName}}) {{$.Name}}{{$.TypeParamNames}} {
			{{- if .IsEmbedded}}
			var result {{$.Name}}{{$.TypeParamNames}}
			result.{{.Name}} = v
			return result
			{{- else}}
			return {{$.Name}}{{$.TypeParamNames}}{ {{.Name}}: v }
			{{- end}}
		},
		"{{$.Name}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- end}}
{{- end}}
	return {{.Name}}Prisms{{.TypeParamNames}} {
{{- range .Fields}}
		{{.Name}}: _prism{{.Name}},
{{- end}}
	}
}

// Make{{.Name}}RefPrisms creates a new [{{.Name}}RefPrisms] with [prisms] for all fields
//
// [prisms]:__prism.Prism
func Make{{.Name}}RefPrisms{{.TypeParams}}() {{.Name}}RefPrisms{{.TypeParamNames}} {
{{- range .Fields}}
{{- if .IsComparable}}
	_fromNonZero{{.Name}} := __option.FromNonZero[{{.TypeName}}]()
	_prism{{.Name}} := __prism.MakePrismWithName(
		func(s *{{$.Name}}{{$.TypeParamNames}}) __option.Option[{{.TypeName}}] { return _fromNonZero{{.Name}}(s.{{.Name}}) },
		func(v {{.TypeName}}) *{{$.Name}}{{$.TypeParamNames}} {
			{{- if .IsEmbedded}}
			var result {{$.Name}}{{$.TypeParamNames}}
			result.{{.Name}} = v
			return &result
			{{- else}}
			return &{{$.Name}}{{$.TypeParamNames}}{ {{.Name}}: v }
			{{- end}}
		},
		"{{$.Name}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- else}}
	_prism{{.Name}} := __prism.MakePrismWithName(
		func(s *{{$.Name}}{{$.TypeParamNames}}) __option.Option[{{.TypeName}}] { return __option.Some(s.{{.Name}}) },
		func(v {{.TypeName}}) *{{$.Name}}{{$.TypeParamNames}} {
			{{- if .IsEmbedded}}
			var result {{$.Name}}{{$.TypeParamNames}}
			result.{{.Name}} = v
			return &result
			{{- else}}
			return &{{$.Name}}{{$.TypeParamNames}}{ {{.Name}}: v }
			{{- end}}
		},
		"{{$.Name}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- end}}
{{- end}}
	return {{.Name}}RefPrisms{{.TypeParamNames}} {
{{- range .Fields}}
		{{.Name}}: _prism{{.Name}},
{{- end}}
	}
}
`

var (
	structTmpl      *template.Template
	constructorTmpl *template.Template
)

func init() {
	var err error
	structTmpl, err = template.New("struct").Parse(lensStructTemplate)
	if err != nil {
		panic(err)
	}
	constructorTmpl, err = template.New("constructor").Parse(lensConstructorTemplate)
	if err != nil {
		panic(err)
	}
}

// hasLensAnnotation checks if a comment group contains the lens annotation
func hasLensAnnotation(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, comment := range doc.List {
		if strings.Contains(comment.Text, lensAnnotation) {
			return true
		}
	}
	return false
}

// getTypeName extracts the type name from a field type expression
func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeName(t.X)
	case *ast.ArrayType:
		return "[]" + getTypeName(t.Elt)
	case *ast.MapType:
		return "map[" + getTypeName(t.Key) + "]" + getTypeName(t.Value)
	case *ast.SelectorExpr:
		return getTypeName(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.IndexExpr:
		// Generic type with single type parameter (Go 1.18+)
		// e.g., Option[string]
		return getTypeName(t.X) + "[" + getTypeName(t.Index) + "]"
	case *ast.IndexListExpr:
		// Generic type with multiple type parameters (Go 1.18+)
		// e.g., Map[string, int]
		var params []string
		for _, index := range t.Indices {
			params = append(params, getTypeName(index))
		}
		return getTypeName(t.X) + "[" + strings.Join(params, ", ") + "]"
	default:
		return "any"
	}
}

// extractImports extracts package imports from a type expression
// Returns a map of package path -> package name
func extractImports(expr ast.Expr, imports map[string]string) {
	switch t := expr.(type) {
	case *ast.StarExpr:
		extractImports(t.X, imports)
	case *ast.ArrayType:
		extractImports(t.Elt, imports)
	case *ast.MapType:
		extractImports(t.Key, imports)
		extractImports(t.Value, imports)
	case *ast.SelectorExpr:
		// This is a qualified identifier like "option.Option"
		if ident, ok := t.X.(*ast.Ident); ok {
			// ident.Name is the package name (e.g., "option")
			// We need to track this for import resolution
			imports[ident.Name] = ident.Name
		}
	case *ast.IndexExpr:
		// Generic type with single type parameter
		extractImports(t.X, imports)
		extractImports(t.Index, imports)
	case *ast.IndexListExpr:
		// Generic type with multiple type parameters
		extractImports(t.X, imports)
		for _, index := range t.Indices {
			extractImports(index, imports)
		}
	}
}

// hasOmitEmpty checks if a struct tag contains json omitempty
func hasOmitEmpty(tag *ast.BasicLit) bool {
	if tag == nil {
		return false
	}
	// Parse the struct tag
	tagValue := strings.Trim(tag.Value, "`")
	structTag := reflect.StructTag(tagValue)
	jsonTag := structTag.Get("json")

	// Check if omitempty is present
	parts := strings.Split(jsonTag, ",")
	for _, part := range parts {
		if strings.TrimSpace(part) == "omitempty" {
			return true
		}
	}
	return false
}

// isPointerType checks if a type expression is a pointer
func isPointerType(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}

// isComparableType checks if a type expression represents a comparable type.
// Comparable types in Go include:
// - Basic types (bool, numeric types, string)
// - Pointer types
// - Channel types
// - Interface types
// - Structs where all fields are comparable
// - Arrays where the element type is comparable
//
// Non-comparable types include:
// - Slices
// - Maps
// - Functions
//
// typeParams is a map of type parameter names to their constraints (e.g., "T" -> "any", "K" -> "comparable")
func isComparableType(expr ast.Expr, typeParams map[string]string) bool {
	switch t := expr.(type) {
	case *ast.Ident:
		// Check if this is a type parameter
		if constraint, isTypeParam := typeParams[t.Name]; isTypeParam {
			// Type parameter - check its constraint
			return constraint == "comparable"
		}

		// Basic types and named types
		// We assume named types are comparable unless they're known non-comparable types
		name := t.Name
		// Known non-comparable built-in types
		if name == "error" {
			// error is an interface, which is comparable
			return true
		}
		// Most basic types and named types are comparable
		// We can't determine if a custom type is comparable without type checking,
		// so we assume it is (conservative approach)
		return true
	case *ast.StarExpr:
		// Pointer types are always comparable
		return true
	case *ast.ArrayType:
		// Arrays are comparable if their element type is comparable
		if t.Len == nil {
			// This is a slice (no length), slices are not comparable
			return false
		}
		// Fixed-size array, check element type
		return isComparableType(t.Elt, typeParams)
	case *ast.MapType:
		// Maps are not comparable
		return false
	case *ast.FuncType:
		// Functions are not comparable
		return false
	case *ast.InterfaceType:
		// Interface types are comparable
		return true
	case *ast.StructType:
		// Structs are comparable if all fields are comparable
		// We can't easily determine this without full type information,
		// so we conservatively return false for struct literals
		return false
	case *ast.SelectorExpr:
		// Qualified identifier (e.g., pkg.Type)
		// We can't determine comparability without type information
		// Check for known non-comparable types from standard library
		if ident, ok := t.X.(*ast.Ident); ok {
			pkgName := ident.Name
			typeName := t.Sel.Name
			// Check for known non-comparable types
			if pkgName == "context" && typeName == "Context" {
				// context.Context is an interface, which is comparable
				return true
			}
			// For other qualified types, we assume they're comparable
			// This is a conservative approach
		}
		return true
	case *ast.IndexExpr, *ast.IndexListExpr:
		// Generic types - we can't determine comparability without type information
		// For common generic types, we can make educated guesses
		var baseExpr ast.Expr
		if idx, ok := t.(*ast.IndexExpr); ok {
			baseExpr = idx.X
		} else if idxList, ok := t.(*ast.IndexListExpr); ok {
			baseExpr = idxList.X
		}

		if sel, ok := baseExpr.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				pkgName := ident.Name
				typeName := sel.Sel.Name
				// Check for known non-comparable generic types
				if pkgName == "option" && typeName == "Option" {
					// Option types are not comparable (they contain a slice internally)
					return false
				}
				if pkgName == "either" && typeName == "Either" {
					// Either types are not comparable
					return false
				}
			}
		}
		// For other generic types, conservatively assume not comparable
		log.Printf("Not comparable type: %v\n", t)
		return false
	case *ast.ChanType:
		// Channel types are comparable
		return true
	default:
		// Unknown type, conservatively assume not comparable
		return false
	}
}

// embeddedFieldResult holds both the field info and its AST type for import extraction
type embeddedFieldResult struct {
	fieldInfo fieldInfo
	fieldType ast.Expr
}

// extractEmbeddedFields extracts fields from an embedded struct type
// It returns a slice of embeddedFieldResult for all exported fields in the embedded struct
// typeParamsMap contains the type parameters of the parent struct (for checking comparability)
func extractEmbeddedFields(embedType ast.Expr, fileImports map[string]string, file *ast.File, typeParamsMap map[string]string) []embeddedFieldResult {
	var results []embeddedFieldResult

	// Get the type name of the embedded field
	var typeName string
	var typeIdent *ast.Ident

	switch t := embedType.(type) {
	case *ast.Ident:
		// Direct embedded type: type MyStruct struct { EmbeddedType }
		typeName = t.Name
		typeIdent = t
	case *ast.StarExpr:
		// Pointer embedded type: type MyStruct struct { *EmbeddedType }
		if ident, ok := t.X.(*ast.Ident); ok {
			typeName = ident.Name
			typeIdent = ident
		}
	case *ast.SelectorExpr:
		// Qualified embedded type: type MyStruct struct { pkg.EmbeddedType }
		// We can't easily resolve this without full type information
		// For now, skip these
		return results
	}

	if S.IsEmpty(typeName) || typeIdent == nil {
		return results
	}

	// Find the struct definition in the same file
	var embeddedStructType *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if ts.Name.Name == typeName {
				if st, ok := ts.Type.(*ast.StructType); ok {
					embeddedStructType = st
					return false
				}
			}
		}
		return true
	})

	if embeddedStructType == nil {
		// Struct not found in this file, might be from another package
		return results
	}

	// Extract fields from the embedded struct
	for _, field := range embeddedStructType.Fields.List {
		// Skip embedded fields within embedded structs (for now, to avoid infinite recursion)
		if len(field.Names) == 0 {
			continue
		}

		for _, name := range field.Names {
			// Generate lenses for both exported and unexported fields
			fieldTypeName := getTypeName(field.Type)
			if true { // Keep the block structure for minimal changes
				isOptional := false
				baseType := fieldTypeName

				// Check if field is optional
				if isPointerType(field.Type) {
					isOptional = true
					baseType = strings.TrimPrefix(fieldTypeName, "*")
				} else if hasOmitEmpty(field.Tag) {
					isOptional = true
				}

				// Check if the type is comparable
				isComparable := isComparableType(field.Type, typeParamsMap)

				results = append(results, embeddedFieldResult{
					fieldInfo: fieldInfo{
						Name:         name.Name,
						TypeName:     fieldTypeName,
						BaseType:     baseType,
						IsOptional:   isOptional,
						IsComparable: isComparable,
						IsEmbedded:   true,
					},
					fieldType: field.Type,
				})
			}
		}
	}

	return results
}

// extractTypeParams extracts type parameters from a type spec
// Returns two strings: full params like "[T any]" and names only like "[T]"
func extractTypeParams(typeSpec *ast.TypeSpec) (string, string) {
	if typeSpec.TypeParams == nil || len(typeSpec.TypeParams.List) == 0 {
		return "", ""
	}

	var params []string
	var names []string
	for _, field := range typeSpec.TypeParams.List {
		for _, name := range field.Names {
			constraint := getTypeName(field.Type)
			params = append(params, name.Name+" "+constraint)
			names = append(names, name.Name)
		}
	}

	fullParams := "[" + strings.Join(params, ", ") + "]"
	nameParams := "[" + strings.Join(names, ", ") + "]"
	return fullParams, nameParams
}

// buildTypeParamsMap creates a map of type parameter names to their constraints
// e.g., for "type Box[T any, K comparable]", returns {"T": "any", "K": "comparable"}
func buildTypeParamsMap(typeSpec *ast.TypeSpec) map[string]string {
	typeParamsMap := make(map[string]string)
	if typeSpec.TypeParams == nil || len(typeSpec.TypeParams.List) == 0 {
		return typeParamsMap
	}

	for _, field := range typeSpec.TypeParams.List {
		constraint := getTypeName(field.Type)
		for _, name := range field.Names {
			typeParamsMap[name.Name] = constraint
		}
	}

	return typeParamsMap
}

// parseFile parses a Go file and extracts structs with lens annotations
func parseFile(filename string) ([]structInfo, string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, "", err
	}

	var structs []structInfo
	packageName := node.Name.Name

	// Build import map: package name -> import path
	fileImports := make(map[string]string)
	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		var name string
		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			// Extract package name from path (last component)
			parts := strings.Split(path, "/")
			name = parts[len(parts)-1]
		}
		fileImports[name] = path
	}

	// First pass: collect all GenDecls with their doc comments
	declMap := make(map[*ast.TypeSpec]*ast.CommentGroup)
	ast.Inspect(node, func(n ast.Node) bool {
		if gd, ok := n.(*ast.GenDecl); ok {
			for _, spec := range gd.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					declMap[ts] = gd.Doc
				}
			}
		}
		return true
	})

	// Second pass: process type specs
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for type declarations
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if it's a struct type
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Get the doc comment from our map
		doc := declMap[typeSpec]
		if !hasLensAnnotation(doc) {
			return true
		}

		// Extract field information and collect imports
		var fields []fieldInfo
		structImports := make(map[string]string)

		// Build type parameters map for this struct
		typeParamsMap := buildTypeParamsMap(typeSpec)

		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				// Embedded field - promote its fields
				embeddedResults := extractEmbeddedFields(field.Type, fileImports, node, typeParamsMap)
				for _, embResult := range embeddedResults {
					// Extract imports from embedded field's type
					fieldImports := make(map[string]string)
					extractImports(embResult.fieldType, fieldImports)

					// Resolve package names to full import paths
					for pkgName := range fieldImports {
						if importPath, ok := fileImports[pkgName]; ok {
							structImports[importPath] = pkgName
						}
					}

					fields = append(fields, embResult.fieldInfo)
				}
				continue
			}
			for _, name := range field.Names {
				// Generate lenses for both exported and unexported fields
				typeName := getTypeName(field.Type)
				if true { // Keep the block structure for minimal changes
					isOptional := false
					baseType := typeName
					isComparable := false

					// Check if field is optional:
					// 1. Pointer types are always optional
					// 2. Non-pointer types with json omitempty tag are optional
					if isPointerType(field.Type) {
						isOptional = true
						// Strip leading * for base type
						baseType = strings.TrimPrefix(typeName, "*")
					} else if hasOmitEmpty(field.Tag) {
						// Non-pointer type with omitempty is also optional
						isOptional = true
					}

					// Check if the type is comparable (for non-optional fields)
					// For optional fields, we don't need to check since they use LensO
					isComparable = isComparableType(field.Type, typeParamsMap)
					// log.Printf("field %s, type: %v, isComparable: %b\n", name, field.Type, isComparable)

					// Extract imports from this field's type
					fieldImports := make(map[string]string)
					extractImports(field.Type, fieldImports)

					// Resolve package names to full import paths
					for pkgName := range fieldImports {
						if importPath, ok := fileImports[pkgName]; ok {
							structImports[importPath] = pkgName
						}
					}

					fields = append(fields, fieldInfo{
						Name:         name.Name,
						TypeName:     typeName,
						BaseType:     baseType,
						IsOptional:   isOptional,
						IsComparable: isComparable,
					})
				}
			}
		}

		if len(fields) > 0 {
			typeParams, typeParamNames := extractTypeParams(typeSpec)
			structs = append(structs, structInfo{
				Name:           typeSpec.Name.Name,
				TypeParams:     typeParams,
				TypeParamNames: typeParamNames,
				Fields:         fields,
				Imports:        structImports,
			})
		}

		return true
	})

	return structs, packageName, nil
}

// generateLensHelpers scans a directory for Go files and generates lens code
func generateLensHelpers(dir, filename string, verbose, includeTestFiles bool) error {
	// Get absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Scanning directory: %s", absDir)
	}

	// Find all Go files in the directory
	files, err := filepath.Glob(filepath.Join(absDir, "*.go"))
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Found %d Go files", len(files))
	}

	// Parse all files and collect structs, separating test and non-test files
	var regularStructs []structInfo
	var testStructs []structInfo
	var packageName string

	for _, file := range files {
		baseName := filepath.Base(file)

		// Skip generated lens files (both regular and test)
		if strings.HasPrefix(baseName, "gen_lens") && strings.HasSuffix(baseName, ".go") {
			if verbose {
				log.Printf("Skipping generated lens file: %s", baseName)
			}
			continue
		}

		isTestFile := strings.HasSuffix(file, "_test.go")

		// Skip test files unless includeTestFiles is true
		if isTestFile && !includeTestFiles {
			if verbose {
				log.Printf("Skipping test file: %s", baseName)
			}
			continue
		}

		if verbose {
			log.Printf("Parsing file: %s", baseName)
		}

		structs, pkg, err := parseFile(file)
		if err != nil {
			log.Printf("Warning: failed to parse %s: %v", file, err)
			continue
		}

		if verbose && len(structs) > 0 {
			log.Printf("Found %d annotated struct(s) in %s", len(structs), baseName)
			for _, s := range structs {
				log.Printf("  - %s (%d fields)", s.Name, len(s.Fields))
			}
		}

		if S.IsEmpty(packageName) {
			packageName = pkg
		}

		// Separate structs based on source file type
		if isTestFile {
			testStructs = append(testStructs, structs...)
		} else {
			regularStructs = append(regularStructs, structs...)
		}
	}

	if len(regularStructs) == 0 && len(testStructs) == 0 {
		log.Printf("No structs with %s annotation found in %s", lensAnnotation, absDir)
		return nil
	}

	// Generate regular lens file if there are regular structs
	if len(regularStructs) > 0 {
		if err := generateLensFile(absDir, filename, packageName, regularStructs, verbose); err != nil {
			return err
		}
	}

	// Generate test lens file if there are test structs
	if len(testStructs) > 0 {
		testFilename := strings.TrimSuffix(filename, ".go") + "_test.go"
		if err := generateLensFile(absDir, testFilename, packageName, testStructs, verbose); err != nil {
			return err
		}
	}

	return nil
}

// generateLensFile generates a lens file for the given structs
func generateLensFile(absDir, filename, packageName string, structs []structInfo, verbose bool) error {
	// Collect all unique imports from all structs
	allImports := make(map[string]string) // import path -> alias
	for _, s := range structs {
		for importPath, alias := range s.Imports {
			allImports[importPath] = alias
		}
	}

	// Create output file
	outPath := filepath.Join(absDir, filename)
	f, err := os.Create(filepath.Clean(outPath))
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Generating lens code in [%s] for package [%s] with [%d] structs ...", outPath, packageName, len(structs))

	// Write header
	writePackage(f, packageName)

	// Write imports
	f.WriteString("import (\n")
	// Standard fp-go imports always needed
	f.WriteString("\t__lens \"github.com/IBM/fp-go/v2/optics/lens\"\n")
	f.WriteString("\t__option \"github.com/IBM/fp-go/v2/option\"\n")
	f.WriteString("\t__prism \"github.com/IBM/fp-go/v2/optics/prism\"\n")
	f.WriteString("\t__lens_option \"github.com/IBM/fp-go/v2/optics/lens/option\"\n")
	f.WriteString("\t__iso_option \"github.com/IBM/fp-go/v2/optics/iso/option\"\n")

	// Add additional imports collected from field types
	for importPath, alias := range allImports {
		f.WriteString("\t" + alias + " \"" + importPath + "\"\n")
	}

	f.WriteString(")\n")

	// Generate lens code for each struct using templates
	for _, s := range structs {
		var buf bytes.Buffer

		// Generate struct type
		if err := structTmpl.Execute(&buf, s); err != nil {
			return err
		}

		// Generate constructor
		if err := constructorTmpl.Execute(&buf, s); err != nil {
			return err
		}

		// Write to file
		if _, err := f.Write(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

// LensCommand creates the CLI command for lens generation
func LensCommand() *C.Command {
	return &C.Command{
		Name:        "lens",
		Usage:       "generate lens code for annotated structs",
		Description: "Scans Go files for structs annotated with 'fp-go:Lens' and generates lens types. Pointer types and non-pointer types with json omitempty tag generate LensO (optional lens).",
		Flags: []C.Flag{
			flagLensDir,
			flagFilename,
			flagVerbose,
			flagIncludeTestFiles,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generateLensHelpers(
				cmd.String(keyLensDir),
				cmd.String(keyFilename),
				cmd.Bool(keyVerbose),
				cmd.Bool(keyIncludeTestFile),
			)
		},
	}
}
