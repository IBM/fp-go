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
	"go/types"
	"log"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
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
		Usage:   "Enable verbose output",
	}

	flagIncludeTestFiles = &C.BoolFlag{
		Name:    keyIncludeTestFile,
		Aliases: []string{"t"},
		Usage:   "Include test files (*_test.go) when scanning for annotated types",
	}

	// flagTypeNames follows the stringer convention: a comma-separated list of
	// type names that bypasses annotation scanning and uses go/packages for full
	// type resolution (generics, external field types, struct tags).
	flagTypeNames = &C.StringFlag{
		Name:  keyTypeNames,
		Usage: "Comma-separated list of struct type names to generate lenses for (replaces annotation scanning)",
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
	Name     string
	TypeName string
	BaseType string // TypeName without leading * for pointer types
	// IsOptional is true if the field is a pointer or has a json omitempty tag
	IsOptional bool
	// IsComparable is true if the type satisfies the "comparable" constraint,
	// i.e. it can be used as a type argument for helpers such as
	// option.FromZero. Interface types qualify even though comparing two
	// interface values can panic at run time.
	IsComparable bool
	// IsStrictlyComparable is true if == on the type can never panic at run
	// time. Interface types (including error) are comparable but not strictly
	// comparable, so lenses for them must not rely on value equality.
	IsStrictlyComparable bool
	// IsEmbedded is true if this field comes from an embedded struct
	IsEmbedded bool
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
{{- if .IsStrictlyComparable}}
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

// getTypeName extracts the type name from a field type expression.
//
// The returned string is valid Go source for the type, so that it can be
// embedded verbatim into the generated lens declarations. Types this function
// does not special-case (functions, channels, inline structs, constraint
// unions, ...) are rendered with go/types.ExprString, which prints the
// expression as it was written.
func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeName(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			// Slice type
			return "[]" + getTypeName(t.Elt)
		}
		// Fixed-size array: the length must be preserved, otherwise the
		// generated code refers to a slice instead of an array.
		return "[" + types.ExprString(t.Len) + "]" + getTypeName(t.Elt)
	case *ast.MapType:
		return "map[" + getTypeName(t.Key) + "]" + getTypeName(t.Value)
	case *ast.SelectorExpr:
		return getTypeName(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		if t.Methods == nil || len(t.Methods.List) == 0 {
			return "any"
		}
		return types.ExprString(t)
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
		return types.ExprString(expr)
	}
}

// extractImports extracts package imports from a type expression
// Returns a map of package path -> package name
//
// The whole expression is walked, so qualified identifiers nested in function
// signatures, channel element types or inline structs are found as well.
func extractImports(expr ast.Expr, imports map[string]string) {
	ast.Inspect(expr, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		// This is a qualified identifier like "option.Option"
		if ident, ok := sel.X.(*ast.Ident); ok {
			// ident.Name is the package name (e.g., "option")
			// We need to track this for import resolution
			imports[ident.Name] = ident.Name
			return false
		}
		return true
	})
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

// typeDecls maps a package level type name to the type expression it is
// declared as, for example "Names" -> []string. It is used to decide whether a
// named type is comparable: a name alone carries no information, the underlying
// type does.
type typeDecls = map[string]ast.Expr

// comparability describes how a type behaves with respect to Go's == operator.
type comparability int

const (
	// notComparable means == does not compile for the type: slices, maps,
	// functions, and composites that contain one of those.
	notComparable comparability = iota
	// interfaceComparable means == compiles, but can panic at run time when the
	// dynamic value is not comparable. Interface types (including error) are in
	// this category. They do satisfy the "comparable" constraint since Go 1.20.
	interfaceComparable
	// strictlyComparable means == compiles and can never panic.
	strictlyComparable
)

// min returns the weaker of two comparability values. A composite type is only
// as comparable as its weakest component.
func (c comparability) min(other comparability) comparability {
	if other < c {
		return other
	}
	return c
}

// isComparableType reports whether a type expression satisfies Go's
// "comparable" constraint, which is what helpers such as option.FromZero and
// option.FromNonZero require.
//
// Comparable types in Go include:
// - Basic types (bool, numeric types, string)
// - Pointer types
// - Channel types
// - Interface types
// - Structs where all fields are comparable
// - Arrays where the element type is comparable
// - Named types whose underlying type is comparable
//
// Non-comparable types include:
// - Slices
// - Maps
// - Functions
// - Named types whose underlying type is one of those
//
// typeParams is a map of type parameter names to their constraints (e.g., "T" -> "any", "K" -> "comparable")
// decls is an optional map of named types in the package, used to resolve named type comparability
func isComparableType(expr ast.Expr, typeParams map[string]string, decls ...typeDecls) bool {
	return comparabilityOf(expr, typeParams, optionalDecls(decls), nil) != notComparable
}

// isStrictlyComparableType reports whether == on the type can never panic at
// run time. It is stricter than isComparableType: interface typed values
// satisfy the "comparable" constraint but panic when two values of the same
// non-comparable dynamic type are compared, so lenses that rely on value
// equality (lens.MakeLensStrict) must not be generated for them.
func isStrictlyComparableType(expr ast.Expr, typeParams map[string]string, decls ...typeDecls) bool {
	return comparabilityOf(expr, typeParams, optionalDecls(decls), nil) == strictlyComparable
}

func optionalDecls(decls []typeDecls) typeDecls {
	if len(decls) > 0 {
		return decls[0]
	}
	return nil
}

// structComparability returns the weakest comparability over all fields of a
// struct type.
func structComparability(st *ast.StructType, typeParams map[string]string, decls typeDecls, visited map[string]bool) comparability {
	result := strictlyComparable
	for _, field := range st.Fields.List {
		result = result.min(comparabilityOf(field.Type, typeParams, decls, visited))
		if result == notComparable {
			return notComparable
		}
	}
	return result
}

// comparabilityOf classifies a type expression.
//
// visited guards against cycles while resolving named types. It is created
// lazily, since the overwhelming majority of types resolve without recursion.
func comparabilityOf(expr ast.Expr, typeParams map[string]string, decls typeDecls, visited map[string]bool) comparability {
	switch t := expr.(type) {
	case *ast.Ident:
		// Check if this is a type parameter
		if constraint, isTypeParam := typeParams[t.Name]; isTypeParam {
			// Type parameter - check its constraint
			if constraint == "comparable" {
				return strictlyComparable
			}
			return notComparable
		}

		// any and error are interfaces: comparable, but == can panic.
		if t.Name == "any" || t.Name == "error" {
			return interfaceComparable
		}

		// If the identifier resolves to a type declared in this package, the
		// underlying type decides. A named slice, map or function type, or a
		// struct that contains one of those, is not comparable.
		if underlying, ok := decls[t.Name]; ok && !visited[t.Name] {
			if visited == nil {
				visited = make(map[string]bool)
			}
			visited[t.Name] = true
			result := comparabilityOf(underlying, typeParams, decls, visited)
			delete(visited, t.Name)
			return result
		}

		// Basic types and named types from outside this package.
		// We can't determine if such a type is comparable without full type
		// checking, so we assume it is.
		return strictlyComparable
	case *ast.StarExpr:
		// Pointer types are always comparable
		return strictlyComparable
	case *ast.ArrayType:
		// Arrays are comparable if their element type is comparable
		if t.Len == nil {
			// This is a slice (no length), slices are not comparable
			return notComparable
		}
		// Fixed-size array, check element type
		return comparabilityOf(t.Elt, typeParams, decls, visited)
	case *ast.MapType:
		// Maps are not comparable
		return notComparable
	case *ast.FuncType:
		// Functions are not comparable
		return notComparable
	case *ast.InterfaceType:
		// Interface types satisfy "comparable", but == can panic
		return interfaceComparable
	case *ast.StructType:
		// Inline struct literal: check all fields
		return structComparability(t, typeParams, decls, visited)
	case *ast.SelectorExpr:
		// Qualified identifier (e.g., pkg.Type) from an external package.
		// Without full type resolution we cannot inspect the type's fields, so we
		// conservatively return notComparable — a struct whose fields include a
		// slice, map, or function is not comparable even if its name looks
		// innocent.
		//
		// Exceptions: types we know are comparable by definition.
		if ident, ok := t.X.(*ast.Ident); ok {
			pkgName := ident.Name
			typeName := t.Sel.Name
			// context.Context is an interface — comparable, but == can panic.
			if pkgName == "context" && typeName == "Context" {
				return interfaceComparable
			}
			// time.Time is a struct with only comparable fields.
			if pkgName == "time" && typeName == "Time" {
				return strictlyComparable
			}
		}
		// Unknown cross-package type: conservatively not comparable.
		return notComparable
	case *ast.IndexExpr, *ast.IndexListExpr:
		// Generic instantiation: Base[T] or Base[T1, T2, ...]
		// Extract the base type name and type arguments.
		var baseExpr ast.Expr
		var typeArgs []ast.Expr
		if idx, ok := t.(*ast.IndexExpr); ok {
			baseExpr = idx.X
			typeArgs = []ast.Expr{idx.Index}
		} else if idxList, ok := t.(*ast.IndexListExpr); ok {
			baseExpr = idxList.X
			typeArgs = idxList.Indices
		}

		// Resolve the package name and type name from either a selector (pkg.Type)
		// or a bare identifier (Type used in the same package / dot-imported).
		var pkgName, typeName string
		if sel, ok := baseExpr.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				pkgName = ident.Name
				typeName = sel.Sel.Name
			}
		} else if ident, ok := baseExpr.(*ast.Ident); ok {
			typeName = ident.Name
		}

		// option.Option[A] / Option[A]: comparable iff A is comparable.
		if typeName == "Option" && (pkgName == "option" || pkgName == "") {
			if len(typeArgs) == 1 {
				return comparabilityOf(typeArgs[0], typeParams, decls, visited)
			}
			return notComparable
		}
		// either.Either[E,A] / Either[E,A]: comparable iff both E and A are comparable.
		// (Default implementation stores E and A as plain fields, not pointers.)
		if typeName == "Either" && (pkgName == "either" || pkgName == "result" || pkgName == "") {
			if len(typeArgs) == 2 {
				return comparabilityOf(typeArgs[0], typeParams, decls, visited).
					min(comparabilityOf(typeArgs[1], typeParams, decls, visited))
			}
			return notComparable
		}
		// pair.Pair[L,R] / Pair[L,R]: comparable iff both L and R are comparable.
		if typeName == "Pair" && (pkgName == "pair" || pkgName == "") {
			if len(typeArgs) == 2 {
				return comparabilityOf(typeArgs[0], typeParams, decls, visited).
					min(comparabilityOf(typeArgs[1], typeParams, decls, visited))
			}
			return notComparable
		}
		// For other generic types, conservatively assume not comparable
		return notComparable
	case *ast.ChanType:
		// Channel types are comparable
		return strictlyComparable
	case *ast.ParenExpr:
		return comparabilityOf(t.X, typeParams, decls, visited)
	default:
		// Unknown type, conservatively assume not comparable
		return notComparable
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
// allTypeDecls is the map of all named types in the package (for comparability checks)
func extractEmbeddedFields(embedType ast.Expr, fileImports map[string]string, file *ast.File, typeParamsMap map[string]string, allTypeDecls typeDecls) []embeddedFieldResult {
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

				// Check how the type behaves under ==
				fieldComparability := comparabilityOf(field.Type, typeParamsMap, allTypeDecls, nil)

				results = append(results, embeddedFieldResult{
					fieldInfo: fieldInfo{
						Name:                 name.Name,
						TypeName:             fieldTypeName,
						BaseType:             baseType,
						IsOptional:           isOptional,
						IsComparable:         fieldComparability != notComparable,
						IsStrictlyComparable: fieldComparability == strictlyComparable,
						IsEmbedded:           true,
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

// collectTypeDecls parses a Go file and returns a map of all named types it
// declares, mapping the name to the type it is declared as. This is used to
// build the package-wide type map that is passed to isComparableType for
// cross-file comparability checks. Named slice, map and function types matter
// as much as structs here: all of them are non-comparable.
func collectTypeDecls(filename string) (typeDecls, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}
	result := make(typeDecls)
	collectTypeDeclsInto(node, result)
	return result, nil
}

// collectTypeDeclsInto adds every type declaration found in node to decls.
func collectTypeDeclsInto(node ast.Node, decls typeDecls) {
	ast.Inspect(node, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			decls[ts.Name.Name] = ts.Type
		}
		return true
	})
}

// parseFile parses a Go file and extracts structs with lens annotations.
// pkgTypeDecls is a package-wide map of named types (collected from all files
// in the package) used to resolve cross-file comparability.
func parseFile(filename string, pkgTypeDecls typeDecls) ([]structInfo, string, error) {
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

	// Build the named type map: start from the package-wide map so that types
	// declared in other files of the same package are also visible.
	// Then overlay types from this file (in case of name shadowing, the
	// local definition wins — though that can't happen for package-level types).
	allTypeDecls := make(typeDecls)
	// range over nil map is a no-op
	maps.Copy(allTypeDecls, pkgTypeDecls)
	collectTypeDeclsInto(node, allTypeDecls)

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
				embeddedResults := extractEmbeddedFields(field.Type, fileImports, node, typeParamsMap, allTypeDecls)
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

					// Check how the type behaves under ==: optional lenses need a
					// type that satisfies "comparable", the strict reference
					// lenses additionally need == never to panic.
					fieldComparability := comparabilityOf(field.Type, typeParamsMap, allTypeDecls, nil)

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
						Name:                 name.Name,
						TypeName:             typeName,
						BaseType:             baseType,
						IsOptional:           isOptional,
						IsComparable:         fieldComparability != notComparable,
						IsStrictlyComparable: fieldComparability == strictlyComparable,
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

	// Pre-pass: collect all named types from every non-generated file in the
	// directory so that cross-file references are resolved when checking
	// comparability.
	pkgTypeDecls := make(typeDecls)
	for _, file := range files {
		baseName := filepath.Base(file)
		if strings.HasPrefix(baseName, "gen_lens") && strings.HasSuffix(baseName, ".go") {
			continue
		}
		isTestFile := strings.HasSuffix(file, "_test.go")
		if isTestFile && !includeTestFiles {
			continue
		}
		fileTypeDecls, err := collectTypeDecls(file)
		if err != nil {
			log.Printf("Warning: failed to collect type declarations from %s: %v", file, err)
			continue
		}
		maps.Copy(pkgTypeDecls, fileTypeDecls)
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

		structs, pkg, err := parseFile(file, pkgTypeDecls)
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

// hasComparableField reports whether any of the structs has at least one
// comparable field, i.e. whether any optional lens will be generated.
func hasComparableField(structs []structInfo) bool {
	for _, s := range structs {
		for _, f := range s.Fields {
			if f.IsComparable {
				return true
			}
		}
	}
	return false
}

// generateLensFile generates a lens file for the given structs
func generateLensFile(absDir, filename, packageName string, structs []structInfo, verbose bool) error {
	// Collect all unique imports from all structs
	allImports := make(map[string]string) // import path -> alias
	for _, s := range structs {
		maps.Copy(allImports, s.Imports)
	}

	outPath := filepath.Join(absDir, filename)

	log.Printf("Generating lens code in [%s] for package [%s] with [%d] structs ...", outPath, packageName, len(structs))

	var buf bytes.Buffer

	// Write header
	writePackage(&buf, packageName)

	// Write imports
	buf.WriteString("import (\n")
	// Standard fp-go imports always needed
	buf.WriteString("\t__lens \"github.com/IBM/fp-go/v2/optics/lens\"\n")
	buf.WriteString("\t__option \"github.com/IBM/fp-go/v2/option\"\n")
	buf.WriteString("\t__prism \"github.com/IBM/fp-go/v2/optics/prism\"\n")
	// The optional lenses only exist for comparable fields. A struct whose
	// fields are all non-comparable (slices, maps, functions) produces no
	// optional lens at all, and importing these packages would then leave the
	// generated file with unused imports.
	if hasComparableField(structs) {
		buf.WriteString("\t__lens_option \"github.com/IBM/fp-go/v2/optics/lens/option\"\n")
		buf.WriteString("\t__iso_option \"github.com/IBM/fp-go/v2/optics/iso/option\"\n")
	}

	// Add additional imports collected from field types, in a stable order so
	// that regenerating an unchanged package produces an unchanged file.
	for _, importPath := range slices.Sorted(maps.Keys(allImports)) {
		buf.WriteString("\t" + allImports[importPath] + " \"" + importPath + "\"\n")
	}

	buf.WriteString(")\n")

	// Generate lens code for each struct using templates
	for _, s := range structs {
		// Generate struct type
		if err := structTmpl.Execute(&buf, s); err != nil {
			return err
		}

		// Generate constructor
		if err := constructorTmpl.Execute(&buf, s); err != nil {
			return err
		}
	}

	// Format the result. If the generated code does not parse, write it out
	// unformatted anyway so that the compiler can point at the offending line.
	content := buf.Bytes()
	formatted, err := format.Source(content)
	if err != nil {
		log.Printf("Warning: generated code for [%s] is not valid Go: %v", outPath, err)
	} else {
		content = formatted
	}

	return os.WriteFile(filepath.Clean(outPath), content, 0o644)
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
// deprecated
func LensCommand() *C.Command {
	return &C.Command{
		Name:        "lens",
		Usage:       "generate lens code for annotated structs or named types",
		Description: "Scans Go files for structs annotated with 'fp-go:Lens'.",
		Flags: []C.Flag{
			flagLensDir,
			flagFilename,
			flagVerbose,
			flagIncludeTestFiles,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			// print warning
			log.Println("Deprecated")
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
