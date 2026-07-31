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
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"maps"
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
	keyTypeNames       = "type"
	keyPackageName     = "package"
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

	// flagTypeNames follows the stringer convention: a comma-separated list of
	// type names that bypasses annotation scanning and uses go/packages for full
	// type resolution (generics, external field types, struct tags).
	flagTypeNames = &C.StringFlag{
		Name:  keyTypeNames,
		Usage: "Comma-separated list of struct type names, or @filename to read from file (replaces annotation scanning)",
	}

	flagPackageName = &C.StringFlag{
		Name:  keyPackageName,
		Usage: "Package name for generated code (defaults to package name of existing files in target directory)",
	}
)

// structInfo holds information about a struct that needs lens generation
type structInfo struct {
	Name           string
	QualifiedName  string // e.g., "openai.ChatCompletionUserMessageParam" when package differs
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
	IsOption     bool   // true if the field type is already an option.Option — LensO must not be generated (would produce Option[Option[A]])
	IsEmbedded   bool   // true if this field comes from an embedded struct
	IsDeprecated bool   // true if the field is marked as deprecated
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
{{- range .Fields}}
	// {{.Name}} is a [__lens.Lens] for the {{.Name}} field of [{{$.QualifiedName}}]
{{- if .IsDeprecated}}
	// Deprecated: This field is deprecated
{{- end}}
	{{.Name}} __lens.Lens[{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
{{- range .Fields}}
{{- if and .IsComparable (not .IsOption)}}
	// {{.Name}}O is a [__lens_option.LensO] for the {{.Name}} field of [{{$.QualifiedName}}], treating the zero value as absent
{{- if .IsDeprecated}}
	// Deprecated: This field is deprecated
{{- end}}
	{{.Name}}O __lens_option.LensO[{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
{{- end}}
}

// {{.Name}}RefLenses provides [lenses] for accessing fields of [{{.Name}}] via a pointer to [{{.Name}}]
//
// [lenses]: __lens.Lens
type {{.Name}}RefLenses{{.TypeParams}} struct {
{{- range .Fields}}
	// {{.Name}} is a [__lens.Lens] for the {{.Name}} field of [{{$.QualifiedName}}] via a pointer receiver
{{- if .IsDeprecated}}
	// Deprecated: This field is deprecated
{{- end}}
	{{.Name}} __lens.Lens[*{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
{{- range .Fields}}
{{- if and .IsComparable (not .IsOption)}}
	// {{.Name}}O is a [__lens_option.LensO] for the {{.Name}} field of [{{$.QualifiedName}}] via a pointer receiver, treating the zero value as absent
{{- if .IsDeprecated}}
	// Deprecated: This field is deprecated
{{- end}}
	{{.Name}}O __lens_option.LensO[*{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
{{- end}}
}

// {{.Name}}Prisms provides [prisms] for accessing fields of [{{.Name}}]
//
// [prisms]: __prism.Prism
type {{.Name}}Prisms{{.TypeParams}} struct {
{{- range .Fields}}
	// {{.Name}} is a [__prism.Prism] for the {{.Name}} field of [{{$.QualifiedName}}]
{{- if .IsDeprecated}}
	// Deprecated: This field is deprecated
{{- end}}
	{{.Name}} __prism.Prism[{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
}

// {{.Name}}RefPrisms provides [prisms] for accessing fields of [{{.Name}}] via a pointer to [{{.Name}}]
//
// [prisms]: __prism.Prism
type {{.Name}}RefPrisms{{.TypeParams}} struct {
{{- range .Fields}}
	// {{.Name}} is a [__prism.Prism] for the {{.Name}} field of [{{$.QualifiedName}}] via a pointer receiver
{{- if .IsDeprecated}}
	// Deprecated: This field is deprecated
{{- end}}
	{{.Name}} __prism.Prism[*{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}]
{{- end}}
}
`

// lensStandaloneTemplate generates one standalone Make<TYPE><FIELD>Lens / Make<TYPE><FIELD>LensO /
// Make<TYPE><FIELD>RefLens / Make<TYPE><FIELD>RefLensO / Make<TYPE><FIELD>Prism /
// Make<TYPE><FIELD>RefPrism function per field so that callers can import only
// the lenses/prisms they actually need, resulting in smaller binaries.
const lensStandaloneTemplate = `
{{- range .Fields}}
// Make{{$.Name}}{{.Name}}Lens returns a [__lens.Lens] for the {{.Name}} field of [{{$.QualifiedName}}]
{{- if .IsDeprecated}}
//
// Deprecated: This field is deprecated
{{- end}}
func Make{{$.Name}}{{.Name}}Lens{{$.TypeParams}}() __lens.Lens[{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}] {
	return __lens.MakeLensWithName(
		func(s {{$.QualifiedName}}{{$.TypeParamNames}}) {{.TypeName}} { return s.{{.Name}} },
		func(s {{$.QualifiedName}}{{$.TypeParamNames}}, v {{.TypeName}}) {{$.QualifiedName}}{{$.TypeParamNames}} { s.{{.Name}} = v; return s },
		"{{$.QualifiedName}}{{$.TypeParamNames}}.{{.Name}}",
	)
}
{{- if and .IsComparable (not .IsOption)}}

// Make{{$.Name}}{{.Name}}LensO returns a [__lens_option.LensO] for the {{.Name}} field of [{{$.QualifiedName}}]
{{- if .IsDeprecated}}
//
// Deprecated: This field is deprecated
{{- end}}
func Make{{$.Name}}{{.Name}}LensO{{$.TypeParams}}() __lens_option.LensO[{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}] {
	return __lens_option.FromIso[{{$.QualifiedName}}{{$.TypeParamNames}}](__iso_option.FromZero[{{.TypeName}}]())(Make{{$.Name}}{{.Name}}Lens{{$.TypeParamNames}}())
}
{{- end}}

// Make{{$.Name}}{{.Name}}RefLens returns a [__lens.Lens] for the {{.Name}} field of [{{$.QualifiedName}}] via a pointer receiver
{{- if .IsDeprecated}}
//
// Deprecated: This field is deprecated
{{- end}}
{{- if .IsComparable}}
func Make{{$.Name}}{{.Name}}RefLens{{$.TypeParams}}() __lens.Lens[*{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}] {
	return __lens.MakeLensStrictWithName(
		func(s *{{$.QualifiedName}}{{$.TypeParamNames}}) {{.TypeName}} { return s.{{.Name}} },
		func(s *{{$.QualifiedName}}{{$.TypeParamNames}}, v {{.TypeName}}) *{{$.QualifiedName}}{{$.TypeParamNames}} { s.{{.Name}} = v; return s },
		"(*{{$.QualifiedName}}{{$.TypeParamNames}}).{{.Name}}",
	)
}
{{- else}}
func Make{{$.Name}}{{.Name}}RefLens{{$.TypeParams}}() __lens.Lens[*{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}] {
	return __lens.MakeLensRefWithName(
		func(s *{{$.QualifiedName}}{{$.TypeParamNames}}) {{.TypeName}} { return s.{{.Name}} },
		func(s *{{$.QualifiedName}}{{$.TypeParamNames}}, v {{.TypeName}}) *{{$.QualifiedName}}{{$.TypeParamNames}} { s.{{.Name}} = v; return s },
		"(*{{$.QualifiedName}}{{$.TypeParamNames}}).{{.Name}}",
	)
}
{{- end}}
{{- if and .IsComparable (not .IsOption)}}

// Make{{$.Name}}{{.Name}}RefLensO returns a [__lens_option.LensO] for the {{.Name}} field of [{{$.QualifiedName}}] via a pointer receiver
{{- if .IsDeprecated}}
//
// Deprecated: This field is deprecated
{{- end}}
func Make{{$.Name}}{{.Name}}RefLensO{{$.TypeParams}}() __lens_option.LensO[*{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}] {
	return __lens_option.FromIso[*{{$.QualifiedName}}{{$.TypeParamNames}}](__iso_option.FromZero[{{.TypeName}}]())(Make{{$.Name}}{{.Name}}RefLens{{$.TypeParamNames}}())
}
{{- end}}

// Make{{$.Name}}{{.Name}}Prism returns a [__prism.Prism] for the {{.Name}} field of [{{$.QualifiedName}}]
{{- if .IsDeprecated}}
//
// Deprecated: This field is deprecated
{{- end}}
func Make{{$.Name}}{{.Name}}Prism{{$.TypeParams}}() __prism.Prism[{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}] {
{{- if .IsComparable}}
	_fromNonZero := __option.FromNonZero[{{.TypeName}}]()
	return __prism.MakePrismWithName(
		func(s {{$.QualifiedName}}{{$.TypeParamNames}}) __option.Option[{{.TypeName}}] { return _fromNonZero(s.{{.Name}}) },
		func(v {{.TypeName}}) {{$.QualifiedName}}{{$.TypeParamNames}} {
			{{- if .IsEmbedded}}
			var result {{$.QualifiedName}}{{$.TypeParamNames}}
			result.{{.Name}} = v
			return result
			{{- else}}
			return {{$.QualifiedName}}{{$.TypeParamNames}}{ {{.Name}}: v }
			{{- end}}
		},
		"{{$.QualifiedName}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- else}}
	return __prism.MakePrismWithName(
		func(s {{$.QualifiedName}}{{$.TypeParamNames}}) __option.Option[{{.TypeName}}] { return __option.Some(s.{{.Name}}) },
		func(v {{.TypeName}}) {{$.QualifiedName}}{{$.TypeParamNames}} {
			{{- if .IsEmbedded}}
			var result {{$.QualifiedName}}{{$.TypeParamNames}}
			result.{{.Name}} = v
			return result
			{{- else}}
			return {{$.QualifiedName}}{{$.TypeParamNames}}{ {{.Name}}: v }
			{{- end}}
		},
		"{{$.QualifiedName}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- end}}
}

// Make{{$.Name}}{{.Name}}RefPrism returns a [__prism.Prism] for the {{.Name}} field of [{{$.QualifiedName}}] via a pointer receiver
{{- if .IsDeprecated}}
//
// Deprecated: This field is deprecated
{{- end}}
func Make{{$.Name}}{{.Name}}RefPrism{{$.TypeParams}}() __prism.Prism[*{{$.QualifiedName}}{{$.TypeParamNames}}, {{.TypeName}}] {
{{- if .IsComparable}}
	_fromNonZero := __option.FromNonZero[{{.TypeName}}]()
	return __prism.MakePrismWithName(
		func(s *{{$.QualifiedName}}{{$.TypeParamNames}}) __option.Option[{{.TypeName}}] { return _fromNonZero(s.{{.Name}}) },
		func(v {{.TypeName}}) *{{$.QualifiedName}}{{$.TypeParamNames}} {
			{{- if .IsEmbedded}}
			var result {{$.QualifiedName}}{{$.TypeParamNames}}
			result.{{.Name}} = v
			return &result
			{{- else}}
			return &{{$.QualifiedName}}{{$.TypeParamNames}}{ {{.Name}}: v }
			{{- end}}
		},
		"{{$.QualifiedName}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- else}}
	return __prism.MakePrismWithName(
		func(s *{{$.QualifiedName}}{{$.TypeParamNames}}) __option.Option[{{.TypeName}}] { return __option.Some(s.{{.Name}}) },
		func(v {{.TypeName}}) *{{$.QualifiedName}}{{$.TypeParamNames}} {
			{{- if .IsEmbedded}}
			var result {{$.QualifiedName}}{{$.TypeParamNames}}
			result.{{.Name}} = v
			return &result
			{{- else}}
			return &{{$.QualifiedName}}{{$.TypeParamNames}}{ {{.Name}}: v }
			{{- end}}
		},
		"{{$.QualifiedName}}{{$.TypeParamNames}}.{{.Name}}",
	)
{{- end}}
}
{{- end}}
`

const lensConstructorTemplate = `
// Make{{.Name}}Lenses creates a new [{{.Name}}Lenses] with [lenses] for all fields
//
// [lenses]: __lens.Lens
func Make{{.Name}}Lenses{{.TypeParams}}() {{.Name}}Lenses{{.TypeParamNames}} {
	return {{.Name}}Lenses{{.TypeParamNames}}{
		// mandatory lenses
{{- range .Fields}}
		{{.Name}}: Make{{$.Name}}{{.Name}}Lens{{$.TypeParamNames}}(),
{{- end}}
		// optional lenses
{{- range .Fields}}
{{- if and .IsComparable (not .IsOption)}}
		{{.Name}}O: Make{{$.Name}}{{.Name}}LensO{{$.TypeParamNames}}(),
{{- end}}
{{- end}}
	}
}

// Make{{.Name}}RefLenses creates a new [{{.Name}}RefLenses] with [lenses] for all fields via a pointer to [{{.Name}}]
//
// [lenses]: __lens.Lens
func Make{{.Name}}RefLenses{{.TypeParams}}() {{.Name}}RefLenses{{.TypeParamNames}} {
	return {{.Name}}RefLenses{{.TypeParamNames}}{
		// mandatory lenses
{{- range .Fields}}
		{{.Name}}: Make{{$.Name}}{{.Name}}RefLens{{$.TypeParamNames}}(),
{{- end}}
		// optional lenses
{{- range .Fields}}
{{- if and .IsComparable (not .IsOption)}}
		{{.Name}}O: Make{{$.Name}}{{.Name}}RefLensO{{$.TypeParamNames}}(),
{{- end}}
{{- end}}
	}
}

// Make{{.Name}}Prisms creates a new [{{.Name}}Prisms] with [prisms] for all fields
//
// [prisms]: __prism.Prism
func Make{{.Name}}Prisms{{.TypeParams}}() {{.Name}}Prisms{{.TypeParamNames}} {
	return {{.Name}}Prisms{{.TypeParamNames}} {
{{- range .Fields}}
		{{.Name}}: Make{{$.Name}}{{.Name}}Prism{{$.TypeParamNames}}(),
{{- end}}
	}
}

// Make{{.Name}}RefPrisms creates a new [{{.Name}}RefPrisms] with [prisms] for all fields via a pointer to [{{.Name}}]
//
// [prisms]: __prism.Prism
func Make{{.Name}}RefPrisms{{.TypeParams}}() {{.Name}}RefPrisms{{.TypeParamNames}} {
	return {{.Name}}RefPrisms{{.TypeParamNames}} {
{{- range .Fields}}
		{{.Name}}: Make{{$.Name}}{{.Name}}RefPrism{{$.TypeParamNames}}(),
{{- end}}
	}
}
`

var (
	structTmpl      *template.Template
	standaloneTmpl  *template.Template
	constructorTmpl *template.Template
)

func init() {
	var err error
	structTmpl, err = template.New("struct").Parse(lensStructTemplate)
	if err != nil {
		panic(err)
	}
	standaloneTmpl, err = template.New("standalone").Parse(lensStandaloneTemplate)
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
		return "any"
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
	parts := strings.SplitSeq(jsonTag, ",")
	for part := range parts {
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

// fpgoValueStructTypes is the set of fp-go generic types that are plain value
// structs (no func or slice fields) and are therefore comparable whenever all
// their type arguments are comparable.
//
// The key is the package short name as it typically appears in source imports
// (e.g. "option", "either"); the value is the set of exported type names from
// that package that follow this rule.
//
// This table must only list types whose underlying struct contains no
// inherently non-comparable fields (no slices, maps, or funcs).  When in
// doubt, omit the type — the generator will conservatively treat it as
// non-comparable.
var fpgoValueStructTypes = map[string]map[string]bool{
	"option": {"Option": true},
	"either": {"Either": true},
	"pair":   {"Pair": true},
	"tuple": {
		"Tuple1": true, "Tuple2": true, "Tuple3": true, "Tuple4": true,
		"Tuple5": true, "Tuple6": true, "Tuple7": true, "Tuple8": true,
		"Tuple9": true, "Tuple10": true, "Tuple11": true, "Tuple12": true,
		"Tuple13": true, "Tuple14": true, "Tuple15": true,
	},
	"result": {"Result": true},
}

// isFpGoValueStructPkg reports whether pkgName.typeName names an fp-go
// value-struct generic type whose comparability depends solely on its type
// arguments (see fpgoValueStructTypes).
func isFpGoValueStructPkg(pkgName, typeName string) bool {
	types, ok := fpgoValueStructTypes[pkgName]
	return ok && types[typeName]
}

// isOptionType reports whether an AST type expression is an Option instantiation,
// i.e. `pkg.Option[A]` for any package alias `pkg` and any type argument `A`.
// A field with this type must not get a LensO generated because that would
// produce an `Option[Option[A]]` return type.
func isOptionType(expr ast.Expr) bool {
	idx, ok := expr.(*ast.IndexExpr)
	if !ok {
		return false
	}
	sel, ok := idx.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	return sel.Sel.Name == "Option"
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
		// Without full type information we cannot determine whether the named
		// type's underlying definition is comparable.  The only safe default is
		// false — callers that need accurate results for external packages should
		// use the go/packages path (--type flag) which uses types.Comparable.
		//
		// Exception: types that are definitively interfaces (and therefore always
		// comparable at the language level) can be white-listed here.
		if ident, ok := t.X.(*ast.Ident); ok {
			pkgName := ident.Name
			typeName := t.Sel.Name
			if pkgName == "context" && typeName == "Context" {
				// context.Context is an interface, which is comparable
				return true
			}
		}
		return false
	case *ast.IndexExpr, *ast.IndexListExpr:
		// Generic instantiation, e.g. option.Option[string] or tuple.Tuple2[int, bool].
		//
		// fp-go value-struct types (Option, Either, Pair, TupleN, …) are plain structs
		// whose comparability depends entirely on their type arguments: if all type
		// arguments are comparable then the instantiation is comparable.
		//
		// We white-list the fp-go package names that follow this rule. For any other
		// generic type we conservatively return false.
		var baseExpr ast.Expr
		var typeArgs []ast.Expr
		if idx, ok := t.(*ast.IndexExpr); ok {
			baseExpr = idx.X
			typeArgs = []ast.Expr{idx.Index}
		} else if idxList, ok := t.(*ast.IndexListExpr); ok {
			baseExpr = idxList.X
			typeArgs = idxList.Indices
		}

		if sel, ok := baseExpr.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				if isFpGoValueStructPkg(ident.Name, sel.Sel.Name) {
					// Comparable iff every type argument is comparable.
					for _, arg := range typeArgs {
						if !isComparableType(arg, typeParams) {
							return false
						}
					}
					return true
				}
			}
		}
		// Unknown generic type — conservatively not comparable.
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
			fieldTypeName := getTypeName(field.Type)
			isOptional := false
			baseType := fieldTypeName

			if isPointerType(field.Type) {
				isOptional = true
				baseType = strings.TrimPrefix(fieldTypeName, "*")
			} else if hasOmitEmpty(field.Tag) {
				isOptional = true
			}

			results = append(results, embeddedFieldResult{
				fieldInfo: fieldInfo{
					Name:         name.Name,
					TypeName:     fieldTypeName,
					BaseType:     baseType,
					IsOptional:   isOptional,
					IsComparable: isComparableType(field.Type, typeParamsMap),
					IsOption:     isOptionType(field.Type),
					IsEmbedded:   true,
				},
				fieldType: field.Type,
			})
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
				typeName := getTypeName(field.Type)
				isOptional := false
				baseType := typeName

				if isPointerType(field.Type) {
					isOptional = true
					baseType = strings.TrimPrefix(typeName, "*")
				} else if hasOmitEmpty(field.Tag) {
					isOptional = true
				}

				fieldImports := make(map[string]string)
				extractImports(field.Type, fieldImports)
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
					IsComparable: isComparableType(field.Type, typeParamsMap),
					IsOption:     isOptionType(field.Type),
				})
			}
		}

		if len(fields) > 0 {
			typeParams, typeParamNames := extractTypeParams(typeSpec)
			structs = append(structs, structInfo{
				Name:           typeSpec.Name.Name,
				QualifiedName:  typeSpec.Name.Name, // Same package, no prefix needed
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
// aliasFromPath converts an import path into a valid Go identifier that is unique
// per path. Each path component separator (/ and .) is replaced with an underscore.
// Example: "k8s.io/apimachinery/pkg/apis/meta/v1" → "k8s_io_apimachinery_pkg_apis_meta_v1"
func aliasFromPath(importPath string) string {
	r := strings.NewReplacer("/", "_", ".", "_", "-", "_")
	return r.Replace(importPath)
}

// resolveImportAliases detects cases where two or more different import paths share
// the same short alias (e.g. both "k8s.io/api/apps/v1" and
// "k8s.io/apimachinery/pkg/apis/meta/v1" declare `package v1`). For every such
// conflicting alias it returns a remapping table that maps each affected import path
// to a disambiguated alias derived from the full import path. The structs slice is
// updated in-place: TypeName / BaseType / QualifiedName strings are rewritten to use
// the new alias, and the Imports map of each struct is updated accordingly.
func resolveImportAliases(structs []structInfo) map[string]string {
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
	pathToNewAlias := make(map[string]string) // only conflicting paths
	for alias, paths := range aliasToPaths {
		if len(paths) <= 1 {
			continue
		}
		for path := range paths {
			newAlias := aliasFromPath(path)
			_ = alias // oldAlias recorded via struct.Imports below
			pathToNewAlias[path] = newAlias
		}
	}

	if len(pathToNewAlias) == 0 {
		return nil
	}

	// Step 3: rewrite each struct in-place
	for i := range structs {
		s := &structs[i]
		for path, newAlias := range pathToNewAlias {
			oldAlias, ok := s.Imports[path]
			if !ok || oldAlias == newAlias {
				continue
			}
			// Rewrite TypeName / BaseType in all fields
			old := oldAlias + "."
			new := newAlias + "."
			for j := range s.Fields {
				s.Fields[j].TypeName = strings.ReplaceAll(s.Fields[j].TypeName, old, new)
				s.Fields[j].BaseType = strings.ReplaceAll(s.Fields[j].BaseType, old, new)
			}
			// Rewrite QualifiedName (may carry a package prefix)
			s.QualifiedName = strings.ReplaceAll(s.QualifiedName, old, new)
			// Update the import alias
			s.Imports[path] = newAlias
		}
	}

	return pathToNewAlias
}

func generateLensFile(absDir, filename, packageName string, structs []structInfo, verbose bool) error {
	// Resolve any import-alias conflicts before collecting imports.
	resolveImportAliases(structs)

	// Collect all unique imports from all structs (conflicts already resolved above)
	allImports := make(map[string]string) // import path -> alias
	for _, s := range structs {
		maps.Copy(allImports, s.Imports)
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

		// Generate standalone per-field helpers
		if err := standaloneTmpl.Execute(&buf, s); err != nil {
			return err
		}

		// Generate bulk constructors (delegate to the standalone helpers)
		if err := constructorTmpl.Execute(&buf, s); err != nil {
			return err
		}

		// Write to file
		if _, err := f.Write(buf.Bytes()); err != nil {
			return err
		}
	}

	// Close the file before formatting
	f.Close()

	// Format the generated file using gofmt
	content, err := os.ReadFile(outPath)
	if err != nil {
		return err
	}

	formatted, err := format.Source(content)
	if err != nil {
		log.Printf("Warning: failed to format %s: %v", outPath, err)
		// Don't fail if formatting fails, the file is still valid Go code
		return nil
	}

	// Write the formatted content back
	if err := os.WriteFile(outPath, formatted, 0644); err != nil {
		return err
	}

	return nil
}

// LensCommand creates the CLI command for lens generation.
//
// Two modes are supported:
//
//  1. Annotation mode (default): scans Go files in --dir for structs annotated
//     with "fp-go:Lens" and generates lenses for them.
//
//  2. Type-name mode (--type flag, following the stringer convention): accepts a
//     comma-separated list of struct names and optional package patterns as
//     positional arguments (default "."). Uses go/packages for full type
//     resolution — generics, external field types, and struct tags are all handled
//     correctly without requiring source annotations.
//
// readTypeNamesFromFile reads type names from a file, one per line.
// Empty lines and lines starting with # are ignored.
func readTypeNamesFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var typeNames []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		typeNames = append(typeNames, line)
	}

	return typeNames, nil
}

func LensCommand() *C.Command {
	return &C.Command{
		Name:        "lens",
		Usage:       "Generate lens code for annotated structs or named types",
		Description: "Scans Go files for structs annotated with 'fp-go:Lens', or — when --type is given — uses go/packages to load named types directly (stringer-style). Pointer fields and json-omitempty fields produce LensO (optional lens).",
		Flags: []C.Flag{
			flagLensDir,
			flagFilename,
			flagVerbose,
			flagIncludeTestFiles,
			flagTypeNames,
			flagPackageName,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			if typesStr := cmd.String(keyTypeNames); typesStr != "" {
				// Type-name mode: parse type names from file or comma-separated list
				var typeNames []string
				var err error

				if strings.HasPrefix(typesStr, "@") {
					// Read type names from file (one per line)
					typeNames, err = readTypeNamesFromFile(strings.TrimPrefix(typesStr, "@"))
					if err != nil {
						return err
					}
				} else {
					// Split comma-separated names
					typeNames = strings.Split(typesStr, ",")
					for i, n := range typeNames {
						typeNames[i] = strings.TrimSpace(n)
					}
				}

				patterns := cmd.Args().Slice()
				if len(patterns) == 0 {
					patterns = []string{"."}
				}
				return generateLensHelpersByType(
					cmd.String(keyLensDir),
					cmd.String(keyFilename),
					patterns,
					typeNames,
					cmd.String(keyPackageName),
					cmd.Bool(keyVerbose),
				)
			}
			// Annotation mode: scan directory for fp-go:Lens annotations.
			return generateLensHelpers(
				cmd.String(keyLensDir),
				cmd.String(keyFilename),
				cmd.Bool(keyVerbose),
				cmd.Bool(keyIncludeTestFile),
			)
		},
	}
}
