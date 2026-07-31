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
	"bytes"
	"fmt"
	"go/format"
	"log"
	"maps"
	"os"
	"path/filepath"
	"text/template"
)

const LensStructTemplate = `
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

// LensStandaloneTemplate generates one standalone Make<TYPE><FIELD>Lens / Make<TYPE><FIELD>LensO /
// Make<TYPE><FIELD>RefLens / Make<TYPE><FIELD>RefLensO / Make<TYPE><FIELD>Prism /
// Make<TYPE><FIELD>RefPrism function per field so that callers can import only
// the lenses/prisms they actually need, resulting in smaller binaries.
const LensStandaloneTemplate = `
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

const LensConstructorTemplate = `
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
	// StructTmpl is the pre-compiled struct-declaration template.
	StructTmpl *template.Template
	// StandaloneTmpl is the pre-compiled per-field function template.
	StandaloneTmpl *template.Template
	// ConstructorTmpl is the pre-compiled bulk-constructor template.
	ConstructorTmpl *template.Template
)

func init() {
	funcMap := template.FuncMap{
		"not": func(b bool) bool { return !b },
	}
	var err error
	StructTmpl, err = template.New("struct").Funcs(funcMap).Parse(LensStructTemplate)
	if err != nil {
		panic(err)
	}
	StandaloneTmpl, err = template.New("standalone").Funcs(funcMap).Parse(LensStandaloneTemplate)
	if err != nil {
		panic(err)
	}
	ConstructorTmpl, err = template.New("constructor").Funcs(funcMap).Parse(LensConstructorTemplate)
	if err != nil {
		panic(err)
	}
}

// WritePackageHeader writes the standard "package + generated" preamble to f.
func WritePackageHeader(f *os.File, pkg string) {
	fmt.Fprintf(f, "package %s\n\n", pkg)
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, "// This file was generated by robots.")
	fmt.Fprintln(f)
}

// GenerateLensFile renders templates for all structs and writes the formatted
// Go source to absDir/filename.
func GenerateLensFile(absDir, filename, packageName string, structs []StructInfo, verbose bool) error {
	// Resolve any import-alias conflicts before collecting imports.
	ResolveImportAliases(structs)

	// Collect all unique imports from all structs (conflicts already resolved above)
	allImports := make(map[string]string) // import path -> alias
	for _, s := range structs {
		maps.Copy(allImports, s.Imports)
	}

	outPath := filepath.Join(absDir, filename)
	f, err := os.Create(filepath.Clean(outPath))
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Generating lens code in [%s] for package [%s] with [%d] structs ...", outPath, packageName, len(structs))

	WritePackageHeader(f, packageName)

	f.WriteString("import (\n")
	f.WriteString("\t__lens \"github.com/IBM/fp-go/v2/optics/lens\"\n")
	f.WriteString("\t__option \"github.com/IBM/fp-go/v2/option\"\n")
	f.WriteString("\t__prism \"github.com/IBM/fp-go/v2/optics/prism\"\n")
	f.WriteString("\t__lens_option \"github.com/IBM/fp-go/v2/optics/lens/option\"\n")
	f.WriteString("\t__iso_option \"github.com/IBM/fp-go/v2/optics/iso/option\"\n")
	for importPath, alias := range allImports {
		f.WriteString("\t" + alias + " \"" + importPath + "\"\n")
	}
	f.WriteString(")\n")

	for _, s := range structs {
		var buf bytes.Buffer
		if err := StructTmpl.Execute(&buf, s); err != nil {
			return err
		}
		if err := StandaloneTmpl.Execute(&buf, s); err != nil {
			return err
		}
		if err := ConstructorTmpl.Execute(&buf, s); err != nil {
			return err
		}
		if _, err := f.Write(buf.Bytes()); err != nil {
			return err
		}
	}

	// Close before formatting so the formatter reads a complete file.
	f.Close()

	content, err := os.ReadFile(outPath)
	if err != nil {
		return err
	}

	formatted, err := format.Source(content)
	if err != nil {
		log.Printf("Warning: failed to format %s: %v", outPath, err)
		return nil
	}

	return os.WriteFile(outPath, formatted, 0644)
}
