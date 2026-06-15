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
	"database/sql"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	C "github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

const (
	keyExamplesSrc     = "src"
	keyExamplesDB      = "db"
	keyExamplesVerbose = "verbose"
	defaultDBPath      = "examples.db"
)

var (
	flagExamplesSrc = &C.StringFlag{
		Name:     keyExamplesSrc,
		Usage:    "Source directory containing Go test files",
		Required: true,
	}

	flagExamplesDB = &C.StringFlag{
		Name:  keyExamplesDB,
		Usage: "Path to SQLite database",
		Value: defaultDBPath,
	}

	flagExamplesVerbose = &C.BoolFlag{
		Name:    keyExamplesVerbose,
		Aliases: []string{"v"},
		Usage:   "Enable verbose output",
	}
)

// Example represents a Go example function
type Example struct {
	ID         string // "package::ExampleName"
	Package    string
	Symbol     string // "Type.Method" parsed from func name
	Name       string // "ExampleType_Method"
	Code       string
	DocComment string
	Output     string // contents of // Output: block
	Imports    string // import statements needed by the example
	File       string
}

// ExamplesCommand returns the CLI command for managing Go examples
func ExamplesCommand() *C.Command {
	return &C.Command{
		Name:  "examples",
		Usage: "Manage Go examples with SQLite",
		Description: `Recursively scan a directory for Go test files, extract Example functions,
and store them in a SQLite database with FTS5 full-text search support.`,
		Flags: []C.Flag{
			flagExamplesDB,
			flagExamplesVerbose,
		},
		Commands: []*C.Command{
			{
				Name:      "ingest",
				Usage:     "Ingest examples from a Go module",
				ArgsUsage: "",
				Flags: []C.Flag{
					flagExamplesSrc,
				},
				Action: func(ctx context.Context, cmd *C.Command) error {
					srcDir := cmd.String(keyExamplesSrc)
					dbPath := cmd.String(keyExamplesDB)
					verbose := cmd.Bool(keyExamplesVerbose)

					if srcDir == "" {
						return fmt.Errorf("source directory is required")
					}
					if dbPath == "" {
						dbPath = defaultDBPath
					}

					return ingestExamples(ctx, srcDir, dbPath, verbose)
				},
			},
			{
				Name:      "search",
				Usage:     "Search for examples using full-text search",
				ArgsUsage: "<query>",
				Flags: []C.Flag{
					&C.StringFlag{
						Name:  "package",
						Usage: "Filter by package name",
					},
				},
				Action: func(ctx context.Context, cmd *C.Command) error {
					query := cmd.Args().First()
					if query == "" {
						return fmt.Errorf("search query is required")
					}

					dbPath := cmd.String(keyExamplesDB)
					if dbPath == "" {
						dbPath = defaultDBPath
					}

					verbose := cmd.Bool(keyExamplesVerbose)
					packageFilter := cmd.String("package")
					return searchExamples(ctx, dbPath, query, packageFilter, verbose)
				},
			},
			{
				Name:      "get",
				Usage:     "Get a specific example by symbol name",
				ArgsUsage: "<symbol>",
				Action: func(ctx context.Context, cmd *C.Command) error {
					symbol := cmd.Args().First()
					if symbol == "" {
						return fmt.Errorf("symbol name is required")
					}

					dbPath := cmd.String(keyExamplesDB)
					if dbPath == "" {
						dbPath = defaultDBPath
					}

					verbose := cmd.Bool(keyExamplesVerbose)
					return getExample(ctx, dbPath, symbol, verbose)
				},
			},
		},
	}
}

// initDB initializes the SQLite database with schema
func initDB(dbPath string) (*sql.DB, error) {
	// Create directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if dbDir != "" && dbDir != "." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables
	schema := `
	CREATE TABLE IF NOT EXISTS examples (
		id          TEXT PRIMARY KEY,
		package     TEXT NOT NULL,
		symbol      TEXT,
		name        TEXT NOT NULL,
		code        TEXT NOT NULL,
		doc_comment TEXT,
		output      TEXT,
		imports     TEXT,
		file        TEXT
	);

	CREATE VIRTUAL TABLE IF NOT EXISTS examples_fts USING fts5(
		name, symbol, package, doc_comment, code,
		content='examples', content_rowid='rowid'
	);

	-- Triggers to keep FTS index in sync
	CREATE TRIGGER IF NOT EXISTS examples_ai AFTER INSERT ON examples BEGIN
		INSERT INTO examples_fts(rowid, name, symbol, package, doc_comment, code)
		VALUES (new.rowid, new.name, new.symbol, new.package, new.doc_comment, new.code);
	END;

	CREATE TRIGGER IF NOT EXISTS examples_ad AFTER DELETE ON examples BEGIN
		INSERT INTO examples_fts(examples_fts, rowid, name, symbol, package, doc_comment, code)
		VALUES('delete', old.rowid, old.name, old.symbol, old.package, old.doc_comment, old.code);
	END;

	CREATE TRIGGER IF NOT EXISTS examples_au AFTER UPDATE ON examples BEGIN
		INSERT INTO examples_fts(examples_fts, rowid, name, symbol, package, doc_comment, code)
		VALUES('delete', old.rowid, old.name, old.symbol, old.package, old.doc_comment, old.code);
		INSERT INTO examples_fts(rowid, name, symbol, package, doc_comment, code)
		VALUES (new.rowid, new.name, new.symbol, new.package, new.doc_comment, new.code);
	END;
	`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

// extractUsedImports extracts only the imports used by a specific function
func extractUsedImports(file *ast.File, funcDecl *ast.FuncDecl) string {
	if len(file.Imports) == 0 {
		return ""
	}

	// Build a map of import paths to their names/aliases
	importMap := make(map[string]string) // path -> name (or "" for unnamed)
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if imp.Name != nil {
			importMap[path] = imp.Name.Name
		} else {
			// Extract package name from path (last segment)
			parts := strings.Split(path, "/")
			importMap[path] = parts[len(parts)-1]
		}
	}

	// Find which imports are actually used in the function
	usedImports := make(map[string]bool)
	ast.Inspect(funcDecl, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			// Check if the selector's X is an identifier that matches an import
			if ident, ok := x.X.(*ast.Ident); ok {
				// Find the import path for this identifier
				for path, name := range importMap {
					if name == ident.Name {
						usedImports[path] = true
						break
					}
				}
			}
		}
		return true
	})

	if len(usedImports) == 0 {
		return ""
	}

	// Build the import statement with only used imports
	var imports []string
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if usedImports[path] {
			if imp.Name != nil {
				imports = append(imports, fmt.Sprintf("%s %s", imp.Name.Name, imp.Path.Value))
			} else {
				imports = append(imports, imp.Path.Value)
			}
		}
	}

	if len(imports) == 0 {
		return ""
	}

	// Format as import block
	if len(imports) == 1 {
		return fmt.Sprintf("import %s", imports[0])
	}

	return fmt.Sprintf("import (\n\t%s\n)", strings.Join(imports, "\n\t"))
}

// extractOutputFromCode extracts the // Output: block from function code
func extractOutputFromCode(code string) string {
	lines := strings.Split(code, "\n")
	var output []string
	inOutput := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "// Output:") {
			inOutput = true
			// Get the content after "// Output:"
			content := strings.TrimSpace(strings.TrimPrefix(trimmed, "// Output:"))
			if content != "" {
				output = append(output, content)
			}
			continue
		}

		if inOutput {
			if strings.HasPrefix(trimmed, "//") {
				// Continue collecting output lines
				content := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
				output = append(output, content)
			} else if trimmed == "}" || (!strings.HasPrefix(trimmed, "//") && trimmed != "") {
				// End of output block
				break
			}
		}
	}

	return strings.Join(output, "\n")
}

// extractDocComment extracts the doc comment for a function
func extractDocComment(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	var lines []string
	for _, c := range doc.List {
		line := strings.TrimPrefix(c.Text, "//")
		line = strings.TrimPrefix(line, "/*")
		line = strings.TrimSuffix(line, "*/")
		lines = append(lines, strings.TrimSpace(line))
	}
	return strings.Join(lines, "\n")
}

// insertExample inserts an example into the database
func insertExample(db *sql.DB, ex Example) error {
	query := `
		INSERT OR REPLACE INTO examples
		(id, package, symbol, name, code, doc_comment, output, imports, file)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query, ex.ID, ex.Package, ex.Symbol, ex.Name,
		ex.Code, ex.DocComment, ex.Output, ex.Imports, ex.File)
	return err
}

// ingestExamples scans the source directory for test files and ingests examples into SQLite
func ingestExamples(ctx context.Context, srcDir, dbPath string, verbose bool) error {
	// Verify source directory exists
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory: %s", srcDir)
	}

	fmt.Printf("Ingesting examples from: %s\n", srcDir)
	fmt.Printf("Using database: %s\n", dbPath)

	// Initialize database
	db, err := initDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	if verbose {
		fmt.Println("Database initialized successfully")
	}

	// Collect all example functions
	examples, err := collectExamples(srcDir)
	if err != nil {
		return fmt.Errorf("failed to collect examples: %w", err)
	}

	if len(examples) == 0 {
		fmt.Println("No example functions found in source directory")
		return nil
	}

	fmt.Printf("Found %d example functions\n", len(examples))

	// Insert examples into database
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for i, ex := range examples {
		if verbose {
			fmt.Printf("  [%d/%d] Inserting: %s\n", i+1, len(examples), ex.ID)
		}
		if err := insertExample(db, ex); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert example %s: %w", ex.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Successfully ingested %d examples into SQLite database\n", len(examples))
	return nil
}

// collectExamples walks the directory tree and collects all example functions
func collectExamples(srcDir string) ([]Example, error) {
	type packageFiles struct {
		pkgPath string
		files   []string
	}

	packages := make(map[string]*packageFiles)

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && (info.Name() == "vendor" || strings.HasPrefix(info.Name(), ".")) {
			return filepath.SkipDir
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		dir := filepath.Dir(path)
		relPath, err := filepath.Rel(srcDir, dir)
		if err != nil {
			return err
		}

		pkgPath := filepath.ToSlash(relPath)
		if pkgPath == "." {
			pkgPath = ""
		}

		group, ok := packages[dir]
		if !ok {
			group = &packageFiles{pkgPath: pkgPath}
			packages[dir] = group
		}
		group.files = append(group.files, path)

		return nil
	})
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0, len(packages))
	for dir := range packages {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	var examples []Example
	for _, dir := range dirs {
		group := packages[dir]
		groupExamples, err := collectExamplesFromPackage(srcDir, group.pkgPath, group.files)
		if err != nil {
			return nil, err
		}
		examples = append(examples, groupExamples...)
	}

	return examples, nil
}

func collectExamplesFromPackage(srcDir, pkgPath string, paths []string) ([]Example, error) {
	fset := token.NewFileSet()
	files := make([]*ast.File, 0, len(paths))
	fileContents := make(map[string][]byte, len(paths))
	fileByName := make(map[string]*ast.File, len(paths))
	funcDecls := make(map[string]*ast.FuncDecl)

	for _, path := range paths {
		fileContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", path, err)
			continue
		}

		file, err := parser.ParseFile(fset, path, fileContent, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", path, err)
			continue
		}

		files = append(files, file)
		fileContents[path] = fileContent
		fileByName[path] = file

		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Recv != nil {
				continue
			}
			funcDecls[funcDecl.Name.Name] = funcDecl
		}

		if pkgPath == "" {
			pkgPath = file.Name.Name
		}
	}

	if len(files) == 0 {
		return nil, nil
	}

	docPkg, err := doc.NewFromFiles(fset, files, pkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build documentation for package %s: %w", pkgPath, err)
	}

	// Collect all examples from different sources
	var allDocExamples []*doc.Example

	// Package-level examples
	allDocExamples = append(allDocExamples, docPkg.Examples...)

	// Function-level examples
	for _, fn := range docPkg.Funcs {
		allDocExamples = append(allDocExamples, fn.Examples...)
	}

	// Type-level examples (including methods)
	for _, typ := range docPkg.Types {
		// Type examples
		allDocExamples = append(allDocExamples, typ.Examples...)

		// Method examples
		for _, method := range typ.Methods {
			allDocExamples = append(allDocExamples, method.Examples...)
		}
	}

	var examples []Example
	for _, docExample := range allDocExamples {
		funcName := "Example" + docExample.Name

		funcDecl, ok := funcDecls[funcName]
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: failed to locate AST for example %s in package %s\n", funcName, pkgPath)
			continue
		}

		position := fset.Position(funcDecl.Pos())
		fileContent, ok := fileContents[position.Filename]
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: failed to locate source for example %s in %s\n", funcName, position.Filename)
			continue
		}

		start := fset.Position(funcDecl.Pos())
		end := fset.Position(funcDecl.End())
		code := string(fileContent[start.Offset:end.Offset])

		symbol := strings.TrimPrefix(funcName, "Example")
		if docExample.Suffix != "" {
			symbol = strings.TrimSuffix(symbol, "_"+docExample.Suffix)
		}

		file := fileByName[position.Filename]
		examples = append(examples, Example{
			ID:         fmt.Sprintf("%s::%s", pkgPath, funcName),
			Package:    pkgPath,
			Symbol:     symbol,
			Name:       funcName,
			Code:       code,
			DocComment: extractDocComment(funcDecl.Doc),
			Output:     extractOutputFromCode(code),
			Imports:    extractUsedImports(file, funcDecl),
			File:       position.Filename,
		})
	}

	return examples, nil
}

// searchExamples searches for examples using FTS5
func searchExamples(ctx context.Context, dbPath, query, packageFilter string, verbose bool) error {
	if verbose {
		fmt.Printf("Opening database: %s\n", dbPath)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	sqlQuery := `
		SELECT e.id, e.package, e.symbol, e.name, e.doc_comment
		FROM examples e
		JOIN examples_fts f ON e.rowid = f.rowid
		WHERE examples_fts MATCH ?
		ORDER BY rank
		LIMIT 10
	`

	var rows *sql.Rows
	if packageFilter != "" {
		if verbose {
			fmt.Printf("Filtering by package: %s\n", packageFilter)
		}
		sqlQuery = `
			SELECT e.id, e.package, e.symbol, e.name, e.doc_comment
			FROM examples e
			JOIN examples_fts f ON e.rowid = f.rowid
			WHERE examples_fts MATCH ? AND e.package = ?
			ORDER BY rank
			LIMIT 10
		`
		rows, err = db.QueryContext(ctx, sqlQuery, query, packageFilter)
	} else {
		rows, err = db.QueryContext(ctx, sqlQuery, query)
	}
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}
	defer rows.Close()

	fmt.Printf("Search results for: %s\n\n", query)
	count := 0
	for rows.Next() {
		var id, pkg, symbol, name, docComment string
		if err := rows.Scan(&id, &pkg, &symbol, &name, &docComment); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		count++
		fmt.Printf("%d. %s\n", count, id)
		fmt.Printf("   Package: %s\n", pkg)
		if symbol != "" {
			fmt.Printf("   Symbol: %s\n", symbol)
		}
		if docComment != "" {
			fmt.Printf("   Doc: %s\n", docComment)
		}
		fmt.Println()
	}

	if count == 0 {
		fmt.Println("No examples found")
	} else if verbose {
		fmt.Printf("\nTotal results: %d\n", count)
	}

	return nil
}

// getExample retrieves a specific example by symbol
func getExample(ctx context.Context, dbPath, symbol string, verbose bool) error {
	if verbose {
		fmt.Printf("Opening database: %s\n", dbPath)
		fmt.Printf("Looking up symbol: %s\n", symbol)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	query := `
		SELECT id, package, symbol, name, code, doc_comment, output, imports, file
		FROM examples
		WHERE symbol = ? OR name = ? OR symbol LIKE ? OR name LIKE ?
		ORDER BY name
	`

	// Add wildcards for pattern matching
	pattern := "%" + symbol + "%"
	rows, err := db.QueryContext(ctx, query, symbol, symbol, pattern, pattern)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var examples []Example
	for rows.Next() {
		var ex Example
		if err := rows.Scan(&ex.ID, &ex.Package, &ex.Symbol, &ex.Name,
			&ex.Code, &ex.DocComment, &ex.Output, &ex.Imports, &ex.File); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		examples = append(examples, ex)
	}

	if len(examples) == 0 {
		return fmt.Errorf("no examples found for symbol: %s", symbol)
	}

	if verbose {
		fmt.Printf("Found %d example(s)!\n\n", len(examples))
	}

	for i, ex := range examples {
		if i > 0 {
			fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
		}

		fmt.Printf("Example: %s\n", ex.ID)
		fmt.Printf("Package: %s\n", ex.Package)
		if ex.Symbol != "" {
			fmt.Printf("Symbol: %s\n", ex.Symbol)
		}
		fmt.Printf("File: %s\n\n", ex.File)

		if ex.DocComment != "" {
			fmt.Printf("Documentation:\n%s\n\n", ex.DocComment)
		}

		if ex.Imports != "" {
			fmt.Printf("Imports:\n%s\n\n", ex.Imports)
		}

		fmt.Printf("Code:\n%s\n", ex.Code)

		if ex.Output != "" {
			fmt.Printf("\nExpected Output:\n%s\n", ex.Output)
		}
	}

	return nil
}
