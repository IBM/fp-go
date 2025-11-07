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
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	C "github.com/urfave/cli/v2"
)

const (
	keyLensDir     = "dir"
	keyVerbose     = "verbose"
	lensAnnotation = "fp-go:Lens"
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
)

// structInfo holds information about a struct that needs lens generation
type structInfo struct {
	Name    string
	Fields  []fieldInfo
	Imports map[string]string // package path -> alias
}

// fieldInfo holds information about a struct field
type fieldInfo struct {
	Name       string
	TypeName   string
	BaseType   string // TypeName without leading * for pointer types
	IsOptional bool   // true if json tag has omitempty or field is a pointer
}

// templateData holds data for template rendering
type templateData struct {
	PackageName string
	Structs     []structInfo
}

const lensStructTemplate = `
// {{.Name}}Lenses provides lenses for accessing fields of {{.Name}}
type {{.Name}}Lenses struct {
{{- range .Fields}}
	{{.Name}} {{if .IsOptional}}LO.LensO[{{$.Name}}, {{.TypeName}}]{{else}}L.Lens[{{$.Name}}, {{.TypeName}}]{{end}}
{{- end}}
}

// {{.Name}}RefLenses provides lenses for accessing fields of {{.Name}} via a reference to {{.Name}}
type {{.Name}}RefLenses struct {
{{- range .Fields}}
	{{.Name}} {{if .IsOptional}}LO.LensO[*{{$.Name}}, {{.TypeName}}]{{else}}L.Lens[*{{$.Name}}, {{.TypeName}}]{{end}}
{{- end}}
}
`

const lensConstructorTemplate = `
// Make{{.Name}}Lenses creates a new {{.Name}}Lenses with lenses for all fields
func Make{{.Name}}Lenses() {{.Name}}Lenses {
{{- range .Fields}}
{{- if .IsOptional}}
	getOrElse{{.Name}} := O.GetOrElse(F.ConstNil[{{.BaseType}}])
{{- end}}
{{- end}}
	return {{.Name}}Lenses{
{{- range .Fields}}
{{- if .IsOptional}}
		{{.Name}}: L.MakeLens(
			func(s {{$.Name}}) O.Option[{{.TypeName}}] { return O.FromNillable(s.{{.Name}}) },
			func(s {{$.Name}}, v O.Option[{{.TypeName}}]) {{$.Name}} { s.{{.Name}} = getOrElse{{.Name}}(v); return s },
		),
{{- else}}
		{{.Name}}: L.MakeLens(
			func(s {{$.Name}}) {{.TypeName}} { return s.{{.Name}} },
			func(s {{$.Name}}, v {{.TypeName}}) {{$.Name}} { s.{{.Name}} = v; return s },
		),
{{- end}}
{{- end}}
	}
}

// Make{{.Name}}RefLenses creates a new {{.Name}}RefLenses with lenses for all fields
func Make{{.Name}}RefLenses() {{.Name}}RefLenses {
{{- range .Fields}}
{{- if .IsOptional}}
	getOrElse{{.Name}} := O.GetOrElse(F.ConstNil[{{.BaseType}}])
{{- end}}
{{- end}}
	return {{.Name}}RefLenses{
{{- range .Fields}}
{{- if .IsOptional}}
		{{.Name}}: L.MakeLensRef(
			func(s *{{$.Name}}) O.Option[{{.TypeName}}] { return O.FromNillable(s.{{.Name}}) },
			func(s *{{$.Name}}, v O.Option[{{.TypeName}}]) *{{$.Name}} { s.{{.Name}} = getOrElse{{.Name}}(v); return s },
		),
{{- else}}
		{{.Name}}: L.MakeLensRef(
			func(s *{{$.Name}}) {{.TypeName}} { return s.{{.Name}} },
			func(s *{{$.Name}}, v {{.TypeName}}) *{{$.Name}} { s.{{.Name}} = v; return s },
		),
{{- end}}
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

		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				// Embedded field, skip for now
				continue
			}
			for _, name := range field.Names {
				// Only export lenses for exported fields
				if name.IsExported() {
					typeName := getTypeName(field.Type)
					isOptional := false
					baseType := typeName

					// Only pointer types can be optional
					if isPointerType(field.Type) {
						isOptional = true
						// Strip leading * for base type
						baseType = strings.TrimPrefix(typeName, "*")
					}

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
						Name:       name.Name,
						TypeName:   typeName,
						BaseType:   baseType,
						IsOptional: isOptional,
					})
				}
			}
		}

		if len(fields) > 0 {
			structs = append(structs, structInfo{
				Name:    typeSpec.Name.Name,
				Fields:  fields,
				Imports: structImports,
			})
		}

		return true
	})

	return structs, packageName, nil
}

// generateLensHelpers scans a directory for Go files and generates lens code
func generateLensHelpers(dir, filename string, verbose bool) error {
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

	// Parse all files and collect structs
	var allStructs []structInfo
	var packageName string

	for _, file := range files {
		// Skip generated files and test files
		if strings.HasSuffix(file, "_test.go") || strings.Contains(file, "gen.go") {
			if verbose {
				log.Printf("Skipping file: %s", filepath.Base(file))
			}
			continue
		}

		if verbose {
			log.Printf("Parsing file: %s", filepath.Base(file))
		}

		structs, pkg, err := parseFile(file)
		if err != nil {
			log.Printf("Warning: failed to parse %s: %v", file, err)
			continue
		}

		if verbose && len(structs) > 0 {
			log.Printf("Found %d annotated struct(s) in %s", len(structs), filepath.Base(file))
			for _, s := range structs {
				log.Printf("  - %s (%d fields)", s.Name, len(s.Fields))
			}
		}

		if packageName == "" {
			packageName = pkg
		}

		allStructs = append(allStructs, structs...)
	}

	if len(allStructs) == 0 {
		log.Printf("No structs with %s annotation found in %s", lensAnnotation, absDir)
		return nil
	}

	// Collect all unique imports from all structs
	allImports := make(map[string]string) // import path -> alias
	for _, s := range allStructs {
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

	log.Printf("Generating lens code in [%s] for package [%s] with [%d] structs ...", outPath, packageName, len(allStructs))

	// Write header
	writePackage(f, packageName)

	// Write imports
	f.WriteString("import (\n")
	// Standard fp-go imports always needed
	f.WriteString("\tF \"github.com/IBM/fp-go/v2/function\"\n")
	f.WriteString("\tL \"github.com/IBM/fp-go/v2/optics/lens\"\n")
	f.WriteString("\tLO \"github.com/IBM/fp-go/v2/optics/lens/option\"\n")
	f.WriteString("\tO \"github.com/IBM/fp-go/v2/option\"\n")

	// Add additional imports collected from field types
	for importPath, alias := range allImports {
		f.WriteString("\t" + alias + " \"" + importPath + "\"\n")
	}

	f.WriteString(")\n")

	// Generate lens code for each struct using templates
	for _, s := range allStructs {
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
		Description: "Scans Go files for structs annotated with 'fp-go:Lens' and generates lens types. Fields with json omitempty tag or pointer types generate LensO (optional lens).",
		Flags: []C.Flag{
			flagLensDir,
			flagFilename,
			flagVerbose,
		},
		Action: func(ctx *C.Context) error {
			return generateLensHelpers(
				ctx.String(keyLensDir),
				ctx.String(keyFilename),
				ctx.Bool(keyVerbose),
			)
		},
	}
}

// Made with Bob
