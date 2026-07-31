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
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/IBM/fp-go/gen/v2/cli/lenses"
	S "github.com/IBM/fp-go/v2/string"
	C "github.com/urfave/cli/v3"
)

const (
	keyLensDir         = "dir"
	keyVerbose         = "verbose"
	keyIncludeTestFile = "include-test-files"
	keyTypeNames       = "type"
	keyPackageName     = "package"
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

// Type aliases so that the rest of the cli package (lens_types.go etc.)
// continues to use the short names without a package prefix.
type structInfo = lenses.StructInfo
type fieldInfo = lenses.FieldInfo

// generateLensHelpers scans a directory for Go files and generates lens code.
func generateLensHelpers(dir, filename string, verbose, includeTestFiles bool) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Scanning directory: %s", absDir)
	}

	files, err := filepath.Glob(filepath.Join(absDir, "*.go"))
	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Found %d Go files", len(files))
	}

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

		if isTestFile && !includeTestFiles {
			if verbose {
				log.Printf("Skipping test file: %s", baseName)
			}
			continue
		}

		if verbose {
			log.Printf("Parsing file: %s", baseName)
		}

		structs, pkg, err := lenses.ParseFile(file)
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

		if isTestFile {
			testStructs = append(testStructs, structs...)
		} else {
			regularStructs = append(regularStructs, structs...)
		}
	}

	if len(regularStructs) == 0 && len(testStructs) == 0 {
		log.Printf("No structs with %s annotation found in %s", lenses.LensAnnotation, absDir)
		return nil
	}

	if len(regularStructs) > 0 {
		if err := lenses.GenerateLensFile(absDir, filename, packageName, regularStructs, verbose); err != nil {
			return err
		}
	}

	if len(testStructs) > 0 {
		testFilename := strings.TrimSuffix(filename, ".go") + "_test.go"
		if err := lenses.GenerateLensFile(absDir, testFilename, packageName, testStructs, verbose); err != nil {
			return err
		}
	}

	return nil
}

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
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		typeNames = append(typeNames, line)
	}

	return typeNames, nil
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
				var typeNames []string
				var err error

				if strings.HasPrefix(typesStr, "@") {
					typeNames, err = readTypeNamesFromFile(strings.TrimPrefix(typesStr, "@"))
					if err != nil {
						return err
					}
				} else {
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
