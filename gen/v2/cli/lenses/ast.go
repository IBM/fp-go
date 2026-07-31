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

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

const LensAnnotation = "fp-go:Lens"

// HasLensAnnotation checks if a comment group contains the lens annotation.
func HasLensAnnotation(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, comment := range doc.List {
		if strings.Contains(comment.Text, LensAnnotation) {
			return true
		}
	}
	return false
}

// GetTypeName extracts the type name from a field type expression.
func GetTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + GetTypeName(t.X)
	case *ast.ArrayType:
		return "[]" + GetTypeName(t.Elt)
	case *ast.MapType:
		return "map[" + GetTypeName(t.Key) + "]" + GetTypeName(t.Value)
	case *ast.SelectorExpr:
		return GetTypeName(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "any"
	case *ast.IndexExpr:
		// Generic type with single type parameter (Go 1.18+)
		// e.g., Option[string]
		return GetTypeName(t.X) + "[" + GetTypeName(t.Index) + "]"
	case *ast.IndexListExpr:
		// Generic type with multiple type parameters (Go 1.18+)
		// e.g., Map[string, int]
		var params []string
		for _, index := range t.Indices {
			params = append(params, GetTypeName(index))
		}
		return GetTypeName(t.X) + "[" + strings.Join(params, ", ") + "]"
	default:
		return "any"
	}
}

// ExtractImports extracts package imports from a type expression and records
// them in imports as packageName -> packageName (caller resolves the full path).
func ExtractImports(expr ast.Expr, imports map[string]string) {
	switch t := expr.(type) {
	case *ast.StarExpr:
		ExtractImports(t.X, imports)
	case *ast.ArrayType:
		ExtractImports(t.Elt, imports)
	case *ast.MapType:
		ExtractImports(t.Key, imports)
		ExtractImports(t.Value, imports)
	case *ast.SelectorExpr:
		// This is a qualified identifier like "option.Option"
		if ident, ok := t.X.(*ast.Ident); ok {
			// ident.Name is the package name (e.g., "option")
			imports[ident.Name] = ident.Name
		}
	case *ast.IndexExpr:
		ExtractImports(t.X, imports)
		ExtractImports(t.Index, imports)
	case *ast.IndexListExpr:
		ExtractImports(t.X, imports)
		for _, index := range t.Indices {
			ExtractImports(index, imports)
		}
	}
}

// HasOmitEmpty checks if a struct tag contains json omitempty.
func HasOmitEmpty(tag *ast.BasicLit) bool {
	if tag == nil {
		return false
	}
	tagValue := strings.Trim(tag.Value, "`")
	structTag := reflect.StructTag(tagValue)
	jsonTag := structTag.Get("json")

	for part := range strings.SplitSeq(jsonTag, ",") {
		if strings.TrimSpace(part) == "omitempty" {
			return true
		}
	}
	return false
}

// IsPointerType checks if a type expression is a pointer.
func IsPointerType(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}

// FpgoValueStructTypes is the set of fp-go generic types that are plain value
// structs (no func or slice fields) and are therefore comparable whenever all
// their type arguments are comparable.
//
// The key is the package short name as it typically appears in source imports
// (e.g. "option", "either"); the value is the set of exported type names from
// that package that follow this rule.
var FpgoValueStructTypes = map[string]map[string]bool{
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

// IsFpGoValueStructPkg reports whether pkgName.typeName names an fp-go
// value-struct generic type whose comparability depends solely on its type
// arguments.
func IsFpGoValueStructPkg(pkgName, typeName string) bool {
	types, ok := FpgoValueStructTypes[pkgName]
	return ok && types[typeName]
}

// IsOptionType reports whether an AST type expression is an Option instantiation,
// i.e. `pkg.Option[A]` for any package alias `pkg` and any type argument `A`.
func IsOptionType(expr ast.Expr) bool {
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

// IsComparableType checks if a type expression represents a comparable type.
//
// typeParams is a map of type parameter names to their constraints
// (e.g., "T" -> "any", "K" -> "comparable").
func IsComparableType(expr ast.Expr, typeParams map[string]string) bool {
	switch t := expr.(type) {
	case *ast.Ident:
		// Check if this is a type parameter
		if constraint, isTypeParam := typeParams[t.Name]; isTypeParam {
			return constraint == "comparable"
		}
		// error is an interface, which is comparable
		if t.Name == "error" {
			return true
		}
		// Most basic types and named types are comparable.
		return true
	case *ast.StarExpr:
		return true
	case *ast.ArrayType:
		if t.Len == nil {
			// slice — not comparable
			return false
		}
		return IsComparableType(t.Elt, typeParams)
	case *ast.MapType:
		return false
	case *ast.FuncType:
		return false
	case *ast.InterfaceType:
		return true
	case *ast.StructType:
		// Conservatively return false for anonymous struct literals.
		return false
	case *ast.SelectorExpr:
		// Qualified identifier (e.g., pkg.Type).
		// White-list definitively interface types.
		if ident, ok := t.X.(*ast.Ident); ok {
			if ident.Name == "context" && t.Sel.Name == "Context" {
				return true
			}
		}
		return false
	case *ast.IndexExpr, *ast.IndexListExpr:
		// Generic instantiation, e.g. option.Option[string] or tuple.Tuple2[int, bool].
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
				if IsFpGoValueStructPkg(ident.Name, sel.Sel.Name) {
					for _, arg := range typeArgs {
						if !IsComparableType(arg, typeParams) {
							return false
						}
					}
					return true
				}
			}
		}
		return false
	case *ast.ChanType:
		return true
	default:
		return false
	}
}

// embeddedFieldResult holds both the field info and its AST type for import extraction.
type embeddedFieldResult struct {
	FieldInfo FieldInfo
	FieldType ast.Expr
}

// ExtractEmbeddedFields extracts fields from an embedded struct type. It returns
// a slice of (FieldInfo, ast.Expr) pairs for all fields in the embedded struct.
// typeParamsMap contains the type parameters of the parent struct.
func ExtractEmbeddedFields(embedType ast.Expr, fileImports map[string]string, file *ast.File, typeParamsMap map[string]string) []embeddedFieldResult {
	var results []embeddedFieldResult

	var typeName string
	var typeIdent *ast.Ident

	switch t := embedType.(type) {
	case *ast.Ident:
		typeName = t.Name
		typeIdent = t
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			typeName = ident.Name
			typeIdent = ident
		}
	case *ast.SelectorExpr:
		// Qualified embedded type — skip (cannot resolve without full type info)
		return results
	}

	if typeName == "" || typeIdent == nil {
		return results
	}

	// Find the struct definition in the same file.
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
		return results
	}

	for _, field := range embeddedStructType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		for _, name := range field.Names {
			fieldTypeName := GetTypeName(field.Type)
			isOptional := false
			baseType := fieldTypeName

			if IsPointerType(field.Type) {
				isOptional = true
				baseType = strings.TrimPrefix(fieldTypeName, "*")
			} else if HasOmitEmpty(field.Tag) {
				isOptional = true
			}

			results = append(results, embeddedFieldResult{
				FieldInfo: FieldInfo{
					Name:         name.Name,
					TypeName:     fieldTypeName,
					BaseType:     baseType,
					IsOptional:   isOptional,
					IsComparable: IsComparableType(field.Type, typeParamsMap),
					IsOption:     IsOptionType(field.Type),
					IsEmbedded:   true,
				},
				FieldType: field.Type,
			})
		}
	}

	return results
}

// ExtractTypeParams extracts type parameters from a type spec.
// Returns two strings: full params like "[T any]" and names only like "[T]".
func ExtractTypeParams(typeSpec *ast.TypeSpec) (string, string) {
	if typeSpec.TypeParams == nil || len(typeSpec.TypeParams.List) == 0 {
		return "", ""
	}

	var params []string
	var names []string
	for _, field := range typeSpec.TypeParams.List {
		for _, name := range field.Names {
			constraint := GetTypeName(field.Type)
			params = append(params, name.Name+" "+constraint)
			names = append(names, name.Name)
		}
	}

	return "[" + strings.Join(params, ", ") + "]", "[" + strings.Join(names, ", ") + "]"
}

// BuildTypeParamsMap creates a map of type parameter names to their constraints.
// e.g., for "type Box[T any, K comparable]", returns {"T": "any", "K": "comparable"}.
func BuildTypeParamsMap(typeSpec *ast.TypeSpec) map[string]string {
	typeParamsMap := make(map[string]string)
	if typeSpec.TypeParams == nil || len(typeSpec.TypeParams.List) == 0 {
		return typeParamsMap
	}
	for _, field := range typeSpec.TypeParams.List {
		constraint := GetTypeName(field.Type)
		for _, name := range field.Names {
			typeParamsMap[name.Name] = constraint
		}
	}
	return typeParamsMap
}

// ParseFile parses a Go file and extracts structs with lens annotations.
// Returns the list of structs, the package name, and any error.
func ParseFile(filename string) ([]StructInfo, string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, "", err
	}

	var structs []StructInfo
	packageName := node.Name.Name

	// Build import map: package name -> import path
	fileImports := make(map[string]string)
	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		var name string
		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			name = parts[len(parts)-1]
		}
		fileImports[name] = path
	}

	// First pass: collect all GenDecls with their doc comments.
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

	// Second pass: process type specs.
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		doc := declMap[typeSpec]
		if !HasLensAnnotation(doc) {
			return true
		}

		var fields []FieldInfo
		structImports := make(map[string]string)
		typeParamsMap := BuildTypeParamsMap(typeSpec)

		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				// Embedded field — promote its fields.
				embeddedResults := ExtractEmbeddedFields(field.Type, fileImports, node, typeParamsMap)
				for _, embResult := range embeddedResults {
					fieldImports := make(map[string]string)
					ExtractImports(embResult.FieldType, fieldImports)
					for pkgName := range fieldImports {
						if importPath, ok := fileImports[pkgName]; ok {
							structImports[importPath] = pkgName
						}
					}
					fields = append(fields, embResult.FieldInfo)
				}
				continue
			}
			for _, name := range field.Names {
				typeName := GetTypeName(field.Type)
				isOptional := false
				baseType := typeName

				if IsPointerType(field.Type) {
					isOptional = true
					baseType = strings.TrimPrefix(typeName, "*")
				} else if HasOmitEmpty(field.Tag) {
					isOptional = true
				}

				fieldImports := make(map[string]string)
				ExtractImports(field.Type, fieldImports)
				for pkgName := range fieldImports {
					if importPath, ok := fileImports[pkgName]; ok {
						structImports[importPath] = pkgName
					}
				}

				fields = append(fields, FieldInfo{
					Name:         name.Name,
					TypeName:     typeName,
					BaseType:     baseType,
					IsOptional:   isOptional,
					IsComparable: IsComparableType(field.Type, typeParamsMap),
					IsOption:     IsOptionType(field.Type),
				})
			}
		}

		if len(fields) > 0 {
			tp, tpNames := ExtractTypeParams(typeSpec)
			structs = append(structs, StructInfo{
				Name:           typeSpec.Name.Name,
				QualifiedName:  typeSpec.Name.Name,
				TypeParams:     tp,
				TypeParamNames: tpNames,
				Fields:         fields,
				Imports:        structImports,
			})
		}

		return true
	})

	return structs, packageName, nil
}
