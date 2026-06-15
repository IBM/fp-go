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
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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

// isValidExampleName checks if a function name follows Go's example naming conventions:
// 1. Must start with "Example"
// 2. Package-level examples: "Example" or "Example_suffix" (suffix all lowercase, can contain underscores)
// 3. Type examples: "ExampleType" or "ExampleType_suffix" (Type capitalized, suffix all lowercase)
// 4. Method examples: "ExampleType_Method" or "ExampleType_Method_suffix" (both Type and Method capitalized, suffix all lowercase)
// 5. Suffixes are used for disambiguation and must be all lowercase (letters, numbers, underscores allowed)
func isValidExampleName(name string) bool {
	if !strings.HasPrefix(name, "Example") {
		return false
	}

	rest := strings.TrimPrefix(name, "Example")

	// "Example" alone is valid (package-level example)
	if rest == "" {
		return true
	}

	// First character after "Example" must be uppercase or underscore
	// Underscore indicates package-level example with disambiguating suffix (e.g., Example_suffix)
	if len(rest) > 0 && !unicode.IsUpper(rune(rest[0])) && rest[0] != '_' {
		return false
	}

	// If starts with underscore, it's Example_suffix format
	// The entire suffix (after underscore) must be all lowercase
	if len(rest) > 0 && rest[0] == '_' {
		suffix := rest[1:] // Remove leading underscore
		if suffix == "" {
			return false // "Example_" alone is invalid
		}
		// Suffix can contain underscores, but all letters must be lowercase
		for _, r := range suffix {
			if unicode.IsUpper(r) {
				return false
			}
		}
		return true
	}

	// Split by underscores to check each part
	parts := strings.Split(rest, "_")

	// First part (Type or Function name) must start with uppercase - already checked above

	// If there are 2 parts: Type_Method (both must start uppercase) or Function_suffix
	if len(parts) == 2 {
		secondPart := parts[1]
		if secondPart == "" {
			return false
		}

		// Check if second part is all lowercase (disambiguating suffix)
		allLower := true
		hasUpperAfterFirst := false

		for i, r := range secondPart {
			if unicode.IsUpper(r) {
				if i == 0 {
					// Uppercase at start is OK for Type_Method
					allLower = false
				} else {
					// Uppercase after first character means it's neither
					// a valid suffix (must be all lowercase) nor
					// a valid method name (would need to be PascalCase throughout)
					hasUpperAfterFirst = true
					allLower = false
				}
			}
		}

		if allLower {
			// It's a disambiguating suffix (all lowercase - valid)
			return true
		}

		if hasUpperAfterFirst {
			// Has uppercase in the middle/end - invalid
			// (not a valid suffix, and not a valid method name pattern)
			return false
		}

		// Starts with uppercase and no uppercase after - assume Type_Method format (valid)
		if unicode.IsUpper(rune(secondPart[0])) {
			return true
		}

		// Doesn't start with uppercase - invalid
		return false
	}

	// If there are 3 parts: Type_Method_suffix
	if len(parts) == 3 {
		// Second part (Method) must start with uppercase
		if !unicode.IsUpper(rune(parts[1][0])) {
			return false
		}

		// Third part (suffix) must be all lowercase
		for _, r := range parts[2] {
			if unicode.IsUpper(r) {
				return false
			}
		}
		return true
	}

	// More than 3 parts is invalid
	if len(parts) > 3 {
		return false
	}

	return true
}

// parseExampleName extracts the symbol from an example function name
// Example_Foo -> "Foo"
// ExampleType_Method -> "Type.Method"
func parseExampleName(name string) string {
	if !strings.HasPrefix(name, "Example") {
		return ""
	}

	rest := strings.TrimPrefix(name, "Example")
	if rest == "" {
		return ""
	}

	// Replace first underscore with dot for method examples
	if idx := strings.Index(rest, "_"); idx > 0 {
		return rest[:idx] + "." + rest[idx+1:]
	}

	return rest
}

// extractUsedImports extracts only the imports used by a specific function
func extractUsedImports(file *ast.File, funcDecl *ast.FuncDecl) string {
	if file.Imports == nil || len(file.Imports) == 0 {
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
	var examples []Example
	fset := token.NewFileSet()

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor and hidden directories
		if info.IsDir() && (info.Name() == "vendor" || strings.HasPrefix(info.Name(), ".")) {
			return filepath.SkipDir
		}

		// Skip directories and non-test files
		if info.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Read file content for output extraction
		fileContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", path, err)
			return nil
		}

		// Parse the Go file
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", path, err)
			return nil // Continue with other files
		}

		// Determine package path relative to source directory
		relPath, err := filepath.Rel(srcDir, filepath.Dir(path))
		if err != nil {
			return err
		}
		pkgPath := filepath.ToSlash(relPath)
		if pkgPath == "." {
			pkgPath = file.Name.Name
		}

		// Extract example functions
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Check if function name starts with "Example" and follows naming conventions
			funcName := funcDecl.Name.Name
			if !strings.HasPrefix(funcName, "Example") {
				continue
			}

			// Validate example naming conventions
			if !isValidExampleName(funcName) {
				fmt.Fprintf(os.Stderr, "Warning: skipping invalid example name %s in %s (must follow Go example naming conventions)\n", funcName, path)
				continue
			}

			// Validate function signature (no parameters, no return values)
			if funcDecl.Type.Params.NumFields() > 0 {
				fmt.Fprintf(os.Stderr, "Warning: skipping %s in %s (example functions must have no parameters)\n", funcName, path)
				continue
			}
			if funcDecl.Type.Results != nil && funcDecl.Type.Results.NumFields() > 0 {
				fmt.Fprintf(os.Stderr, "Warning: skipping %s in %s (example functions must have no return values)\n", funcName, path)
				continue
			}

			// Extract function code
			start := fset.Position(funcDecl.Pos())
			end := fset.Position(funcDecl.End())
			code := string(fileContent[start.Offset:end.Offset])

			// Parse symbol from function name
			symbol := parseExampleName(funcName)

			// Extract doc comment
			docComment := extractDocComment(funcDecl.Doc)

			// Extract output - look for // Output: comment in the function
			output := extractOutputFromCode(code)

			// Extract only the imports used by this specific function
			imports := extractUsedImports(file, funcDecl)

			example := Example{
				ID:         fmt.Sprintf("%s::%s", pkgPath, funcName),
				Package:    pkgPath,
				Symbol:     symbol,
				Name:       funcName,
				Code:       code,
				DocComment: docComment,
				Output:     output,
				Imports:    imports,
				File:       path,
			}

			examples = append(examples, example)
		}

		return nil
	})

	return examples, err
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
		rows, err = db.Query(sqlQuery, query, packageFilter)
	} else {
		rows, err = db.Query(sqlQuery, query)
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
	rows, err := db.Query(query, symbol, symbol, pattern, pattern)
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
