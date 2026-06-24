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

package mcp

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/IBM/fp-go/v2/gen/data"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

// SkillInfo represents information about a skill
type SkillInfo struct {
	Name        string `json:"name" jsonschema:"The name of the skill (directory name)"`
	Description string `json:"description,omitempty" jsonschema:"A brief description of the skill from its SKILL.md file"`
	Path        string `json:"path" jsonschema:"The relative path to the skill directory"`
}

// ListSkillsOutput represents the output of the list_skills tool
type ListSkillsOutput struct {
	Skills []SkillInfo `json:"skills" jsonschema:"List of available skills"`
}

// UseSkillArgs represents the arguments for the use_skill tool
type UseSkillArgs struct {
	Name string `json:"name" jsonschema:"The name of the skill to use (directory name from list_skills)"`
}

// UseSkillOutput represents the output of the use_skill tool
type UseSkillOutput struct {
	Name        string `json:"name" jsonschema:"The name of the skill"`
	Description string `json:"description,omitempty" jsonschema:"A brief description of the skill from its SKILL.md file"`
	Content     string `json:"content" jsonschema:"The content of the skill's SKILL.md file after the YAML frontmatter header"`
}

// SearchExamplesArgs represents the arguments for the search_examples tool
type SearchExamplesArgs struct {
	Query         string `json:"query" jsonschema:"The search query for full-text search across example names, symbols, packages, doc comments, and code"`
	PackageFilter string `json:"package_filter,omitempty" jsonschema:"Optional package name to filter results (e.g., 'option' or 'either')"`
}

// GoExample represents a Go example function
type GoExample struct {
	ID         string `json:"id" jsonschema:"Unique identifier in format 'package::ExampleName'"`
	Package    string `json:"package" jsonschema:"Package path"`
	Symbol     string `json:"symbol,omitempty" jsonschema:"Symbol name (e.g., 'Type.Method' parsed from function name)"`
	Name       string `json:"name" jsonschema:"Function name (e.g., 'ExampleType_Method')"`
	DocComment string `json:"doc_comment,omitempty" jsonschema:"Documentation comment for the example"`
	Code       string `json:"code,omitempty" jsonschema:"Full source code of the example function"`
	Output     string `json:"output,omitempty" jsonschema:"Expected output from // Output: block"`
	Imports    string `json:"imports,omitempty" jsonschema:"Import statements needed by the example"`
	File       string `json:"file,omitempty" jsonschema:"File path where the example is located"`
}

// SearchExamplesOutput represents the output of the search_examples tool
type SearchExamplesOutput struct {
	Examples []GoExample `json:"examples" jsonschema:"List of examples matching the search query"`
	Count    int         `json:"count" jsonschema:"Number of examples found"`
}

// GetExampleArgs represents the arguments for the get_example tool
type GetExampleArgs struct {
	Symbol string `json:"symbol" jsonschema:"The symbol name to look up (e.g., 'Type.Method' or 'ExampleType_Method')"`
}

// GetExampleOutput represents the output of the get_example tool
type GetExampleOutput struct {
	Examples []GoExample `json:"examples" jsonschema:"List of examples matching the symbol"`
	Count    int         `json:"count" jsonschema:"Number of examples found"`
}

// NewServer creates a new MCP server with fp-go tools
func NewServer(db *sql.DB, verbose bool) *mcp.Server {
	if verbose {
		log.Println("[MCP] Creating MCP server with fp-go tools")
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "fp-go-generator",
		Version: "1.0.0",
	}, nil)

	if verbose {
		log.Println("[MCP] Registering tool: list_skills")
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_skills",
		Description: "List all available fp-go skills. Each skill provides specialized knowledge and guidance for specific fp-go features and patterns. Returns a list of skills with their names, descriptions, and paths.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, handleListSkills(verbose))

	if verbose {
		log.Println("[MCP] Registering tool: use_skill")
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "use_skill",
		Description: "Retrieve the full content of a specific fp-go skill by name. Use this after list_skills to get detailed documentation, examples, and best practices for a particular fp-go feature or pattern.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, handleUseSkill(verbose))

	if verbose {
		log.Println("[MCP] Registering tool: search_examples")
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_examples",
		Description: "Search for Go example functions using full-text search. Searches across example names, symbols, packages, documentation comments, and code. Returns up to 10 ranked results with metadata.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, handleSearchExamples(db, verbose))

	if verbose {
		log.Println("[MCP] Registering tool: get_example")
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_example",
		Description: "Retrieve a specific Go example by symbol name. Performs exact lookup by symbol (e.g., 'Type.Method') or function name (e.g., 'ExampleType_Method'). Returns the complete example with code, documentation, and output.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, handleGetExample(db, verbose))

	if verbose {
		log.Println("[MCP] Server created with 4 tools registered")
	}

	return server
}

// Run starts the MCP server with stdio transport
func Run(ctx context.Context, verbose bool) error {
	if verbose {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		log.Println("[MCP] Verbose logging enabled")
	}

	tmpFile, err := os.CreateTemp("", "examples-*.db")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(data.EXAMPLES_DB); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write database: %w", err)
	}
	tmpFile.Close()

	if verbose {
		log.Println("[MCP] Opening examples DB...")
	}

	db, err := sql.Open("sqlite", tmpPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	server := NewServer(db, verbose)

	if verbose {
		log.Println("[MCP] Starting MCP server with stdio transport...")
	}

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		if verbose {
			log.Printf("[MCP] Server stopped with error: %v\n", err)
		}
		return fmt.Errorf("server failed: %w", err)
	}

	if verbose {
		log.Println("[MCP] Server stopped successfully")
	}

	return nil
}

// skillFrontmatter represents the YAML frontmatter structure in SKILL.md files
type skillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// parseFrontmatter extracts name, description, and body from YAML frontmatter in a single pass.
// If no frontmatter is found or parsing fails, name and description are empty and body is the full content.
func parseFrontmatter(content []byte) (name, description, body string) {
	lines := strings.Split(string(content), "\n")
	var fmLines []string
	inFrontmatter := false
	bodyStart := 0

	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			bodyStart = i + 1
			break
		}
		if inFrontmatter {
			fmLines = append(fmLines, line)
		}
	}

	if len(fmLines) == 0 {
		return "", "", string(content)
	}

	var fm skillFrontmatter
	if err := yaml.Unmarshal([]byte(strings.Join(fmLines, "\n")), &fm); err != nil {
		return "", "", string(content)
	}

	if bodyStart > 0 && bodyStart < len(lines) {
		body = strings.Join(lines[bodyStart:], "\n")
	}
	return fm.Name, fm.Description, body
}

// renderExamplesToMarkdown builds a markdown string for a list of examples.
// When full is true, also includes ID, file path, and imports (used by get_example).
func renderExamplesToMarkdown(examples []GoExample, full bool) string {
	var b strings.Builder
	for i, ex := range examples {
		fmt.Fprintf(&b, "## %d. %s\n\n", i+1, ex.Name)
		if full {
			fmt.Fprintf(&b, "**ID:** `%s`\n\n", ex.ID)
		}
		fmt.Fprintf(&b, "**Package:** `%s`\n\n", ex.Package)
		if ex.Symbol != "" {
			fmt.Fprintf(&b, "**Symbol:** `%s`\n\n", ex.Symbol)
		}
		if full && ex.File != "" {
			fmt.Fprintf(&b, "**File:** `%s`\n\n", ex.File)
		}
		if ex.DocComment != "" {
			fmt.Fprintf(&b, "**Documentation:**\n%s\n\n", ex.DocComment)
		}
		if full && ex.Imports != "" {
			b.WriteString("**Imports:**\n```go\n")
			b.WriteString(ex.Imports)
			b.WriteString("\n```\n\n")
		}
		if ex.Code != "" {
			b.WriteString("**Code:**\n```go\n")
			b.WriteString(ex.Code)
			b.WriteString("\n```\n\n")
		}
		if ex.Output != "" {
			b.WriteString("**Output:**\n```\n")
			b.WriteString(ex.Output)
			b.WriteString("\n```\n\n")
		}
		if i < len(examples)-1 {
			b.WriteString("---\n\n")
		}
	}
	return b.String()
}

func errorResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
		IsError: true,
	}
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

// handleListSkills handles the list_skills tool call
func handleListSkills(verbose bool) func(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, ListSkillsOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, ListSkillsOutput, error) {
		if verbose {
			log.Println("[MCP] Executing tool: list_skills")
		}

		var skills []SkillInfo

		for path, content := range data.Skills {
			if !strings.HasSuffix(path, "SKILL.md") {
				continue
			}

			dirName := filepath.Dir(path)
			name, description, _ := parseFrontmatter(content)

			if name != "" && name != dirName {
				msg := fmt.Sprintf("Error: skill name '%s' does not match directory name '%s' in %s", name, dirName, path)
				return errorResult(msg), ListSkillsOutput{}, fmt.Errorf("skill name mismatch: '%s' != '%s' in %s", name, dirName, path)
			}

			skillName := name
			if skillName == "" {
				skillName = dirName
			}
			skills = append(skills, SkillInfo{
				Name:        skillName,
				Description: description,
				Path:        dirName,
			})
		}

		sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })

		if verbose {
			log.Printf("[MCP] list_skills completed: found %d skills\n", len(skills))
		}

		var markdown strings.Builder
		markdown.WriteString("# Available fp-go Skills\n\n")
		markdown.WriteString("| Skill | Description |\n")
		markdown.WriteString("|-------|-------------|\n")
		for _, skill := range skills {
			desc := skill.Description
			if desc == "" {
				desc = "_No description available_"
			}
			fmt.Fprintf(&markdown, "| **%s** | %s |\n", skill.Name, desc)
		}
		fmt.Fprintf(&markdown, "\n_Total: %d skill(s)_\n", len(skills))

		return textResult(markdown.String()), ListSkillsOutput{Skills: skills}, nil
	}
}

// handleUseSkill handles the use_skill tool call
func handleUseSkill(verbose bool) func(ctx context.Context, req *mcp.CallToolRequest, args UseSkillArgs) (*mcp.CallToolResult, UseSkillOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args UseSkillArgs) (*mcp.CallToolResult, UseSkillOutput, error) {
		if verbose {
			log.Printf("[MCP] Executing tool: use_skill (name=%s)\n", args.Name)
		}

		if args.Name == "" {
			if verbose {
				log.Println("[MCP] use_skill error: skill name is required")
			}
			return errorResult("Error: skill name is required"), UseSkillOutput{}, fmt.Errorf("skill name is required")
		}

		// Use forward slashes to match the embedded map keys (cross-platform)
		skillPath := args.Name + "/SKILL.md"
		content, found := data.Skills[skillPath]
		if !found {
			if verbose {
				log.Printf("[MCP] use_skill error: skill '%s' not found\n", args.Name)
			}
			return errorResult(fmt.Sprintf("Error: skill '%s' not found", args.Name)),
				UseSkillOutput{},
				fmt.Errorf("skill '%s' not found", args.Name)
		}

		_, description, body := parseFrontmatter(content)

		output := UseSkillOutput{
			Name:        args.Name,
			Description: description,
			Content:     body,
		}

		if verbose {
			log.Printf("[MCP] use_skill completed: retrieved skill '%s'\n", args.Name)
		}

		var markdown strings.Builder
		fmt.Fprintf(&markdown, "# Skill: %s\n\n", args.Name)
		if description != "" {
			fmt.Fprintf(&markdown, "**Description:** %s\n\n", description)
			markdown.WriteString("---\n\n")
		}
		markdown.WriteString(body)

		return textResult(markdown.String()), output, nil
	}
}

// handleSearchExamples handles the search_examples tool call
func handleSearchExamples(db *sql.DB, verbose bool) func(ctx context.Context, req *mcp.CallToolRequest, args SearchExamplesArgs) (*mcp.CallToolResult, SearchExamplesOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args SearchExamplesArgs) (*mcp.CallToolResult, SearchExamplesOutput, error) {
		if verbose {
			log.Printf("[MCP] Executing tool: search_examples (query=%s, package_filter=%s)\n", args.Query, args.PackageFilter)
		}

		if args.Query == "" {
			if verbose {
				log.Println("[MCP] search_examples error: query is required")
			}
			return errorResult("Error: search query is required"), SearchExamplesOutput{}, fmt.Errorf("search query is required")
		}

		const baseQuery = `
			SELECT e.id, e.package, e.symbol, e.name, e.code, e.doc_comment, e.output, e.imports, e.file
			FROM examples e
			JOIN examples_fts f ON e.rowid = f.rowid
			WHERE examples_fts MATCH ?
			ORDER BY rank
			LIMIT 10
		`
		const filteredQuery = `
			SELECT e.id, e.package, e.symbol, e.name, e.code, e.doc_comment, e.output, e.imports, e.file
			FROM examples e
			JOIN examples_fts f ON e.rowid = f.rowid
			WHERE examples_fts MATCH ? AND e.package = ?
			ORDER BY rank
			LIMIT 10
		`

		var rows *sql.Rows
		var err error
		if args.PackageFilter != "" {
			rows, err = db.Query(filteredQuery, args.Query, args.PackageFilter)
		} else {
			rows, err = db.Query(baseQuery, args.Query)
		}
		if err != nil {
			return errorResult(fmt.Sprintf("Error: search failed: %v", err)),
				SearchExamplesOutput{},
				fmt.Errorf("search failed: %w", err)
		}
		defer rows.Close()

		examples := make([]GoExample, 0)
		for rows.Next() {
			var ex GoExample
			if err := rows.Scan(&ex.ID, &ex.Package, &ex.Symbol, &ex.Name, &ex.Code, &ex.DocComment, &ex.Output, &ex.Imports, &ex.File); err != nil {
				return errorResult(fmt.Sprintf("Error: failed to scan row: %v", err)),
					SearchExamplesOutput{},
					fmt.Errorf("failed to scan row: %w", err)
			}
			examples = append(examples, ex)
		}
		if err := rows.Err(); err != nil {
			return errorResult(fmt.Sprintf("Error: row iteration failed: %v", err)),
				SearchExamplesOutput{},
				fmt.Errorf("row iteration failed: %w", err)
		}

		if verbose {
			log.Printf("[MCP] search_examples completed: found %d examples\n", len(examples))
		}

		var markdown strings.Builder
		fmt.Fprintf(&markdown, "# Search Results: %s\n\n", args.Query)
		if args.PackageFilter != "" {
			fmt.Fprintf(&markdown, "**Package Filter:** `%s`\n\n", args.PackageFilter)
		}
		fmt.Fprintf(&markdown, "**Found:** %d example(s)\n\n", len(examples))
		if len(examples) > 0 {
			markdown.WriteString("---\n\n")
			markdown.WriteString(renderExamplesToMarkdown(examples, false))
		}

		return textResult(markdown.String()), SearchExamplesOutput{Examples: examples, Count: len(examples)}, nil
	}
}

// handleGetExample handles the get_example tool call
func handleGetExample(db *sql.DB, verbose bool) func(ctx context.Context, req *mcp.CallToolRequest, args GetExampleArgs) (*mcp.CallToolResult, GetExampleOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args GetExampleArgs) (*mcp.CallToolResult, GetExampleOutput, error) {
		if verbose {
			log.Printf("[MCP] Executing tool: get_example (symbol=%s)\n", args.Symbol)
		}

		if args.Symbol == "" {
			if verbose {
				log.Println("[MCP] get_example error: symbol is required")
			}
			return errorResult("Error: symbol name is required"), GetExampleOutput{}, fmt.Errorf("symbol name is required")
		}

		const query = `
			SELECT id, package, symbol, name, code, doc_comment, output, imports, file
			FROM examples
			WHERE symbol = ? OR name = ? OR symbol LIKE ? OR name LIKE ?
			ORDER BY name
		`
		pattern := "%" + args.Symbol + "%"
		rows, err := db.Query(query, args.Symbol, args.Symbol, pattern, pattern)
		if err != nil {
			return errorResult(fmt.Sprintf("Error: query failed: %v", err)),
				GetExampleOutput{},
				fmt.Errorf("query failed: %w", err)
		}
		defer rows.Close()

		examples := make([]GoExample, 0)
		for rows.Next() {
			var ex GoExample
			if err := rows.Scan(&ex.ID, &ex.Package, &ex.Symbol, &ex.Name,
				&ex.Code, &ex.DocComment, &ex.Output, &ex.Imports, &ex.File); err != nil {
				return errorResult(fmt.Sprintf("Error: failed to scan row: %v", err)),
					GetExampleOutput{},
					fmt.Errorf("failed to scan row: %w", err)
			}
			examples = append(examples, ex)
		}
		if err := rows.Err(); err != nil {
			return errorResult(fmt.Sprintf("Error: row iteration failed: %v", err)),
				GetExampleOutput{},
				fmt.Errorf("row iteration failed: %w", err)
		}

		if len(examples) == 0 {
			return textResult(fmt.Sprintf("No examples found for symbol: %s", args.Symbol)),
				GetExampleOutput{Examples: []GoExample{}, Count: 0}, nil
		}

		if verbose {
			log.Printf("[MCP] get_example completed: retrieved %d examples\n", len(examples))
		}

		var markdown strings.Builder
		fmt.Fprintf(&markdown, "# Example: %s\n\n", args.Symbol)
		fmt.Fprintf(&markdown, "**Found:** %d example(s)\n\n", len(examples))
		markdown.WriteString("---\n\n")
		markdown.WriteString(renderExamplesToMarkdown(examples, true))

		return textResult(markdown.String()), GetExampleOutput{Examples: examples, Count: len(examples)}, nil
	}
}
