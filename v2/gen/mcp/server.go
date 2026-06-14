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
	"os"
	"path/filepath"
	"strings"

	"github.com/IBM/fp-go/v2/gen/data"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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
func NewServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "fp-go-generator",
		Version: "1.0.0",
	}, nil)

	// Register the list_skills tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_skills",
		Description: "List all available fp-go skills. Each skill provides specialized knowledge and guidance for specific fp-go features and patterns. Returns a list of skills with their names, descriptions, and paths.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, handleListSkills)

	// Register the use_skill tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "use_skill",
		Description: "Retrieve the full content of a specific fp-go skill by name. Use this after list_skills to get detailed documentation, examples, and best practices for a particular fp-go feature or pattern.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, handleUseSkill)

	// Register the search_examples tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_examples",
		Description: "Search for Go example functions using full-text search. Searches across example names, symbols, packages, documentation comments, and code. Returns up to 10 ranked results with metadata.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, handleSearchExamples)

	// Register the get_example tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_example",
		Description: "Retrieve a specific Go example by symbol name. Performs exact lookup by symbol (e.g., 'Type.Method') or function name (e.g., 'ExampleType_Method'). Returns the complete example with code, documentation, and output.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, handleGetExample)

	return server
}

// Run starts the MCP server with stdio transport
func Run(ctx context.Context) error {
	server := NewServer()

	// Run the server on stdio transport
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

// parseSkillMetadata extracts name and description from YAML frontmatter
func parseSkillMetadata(content []byte) (name, description string) {
	lines := strings.Split(string(content), "\n")
	inFrontmatter := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for frontmatter delimiters
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				// End of frontmatter
				break
			}
		}

		if !inFrontmatter {
			continue
		}

		// Parse name field
		if after, ok := strings.CutPrefix(trimmed, "name:"); ok {
			name = strings.TrimSpace(after)
			continue
		}

		// Parse description field
		if strings.HasPrefix(trimmed, "description:") {
			description = strings.TrimSpace(strings.TrimPrefix(trimmed, "description:"))
			continue
		}
	}

	return name, description
}

// handleListSkills handles the list_skills tool call
func handleListSkills(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, ListSkillsOutput, error) {
	var skills []SkillInfo

	// Iterate through the Skills map from data package
	for path, content := range data.Skills {
		// Only process SKILL.md files
		if !strings.HasSuffix(path, "SKILL.md") {
			continue
		}

		// Extract directory name from path (e.g., "fp-go-logging/SKILL.md" -> "fp-go-logging")
		dirName := filepath.Dir(path)

		// Parse metadata from YAML frontmatter
		name, description := parseSkillMetadata(content)

		// Validate that the name matches the directory name
		if name != "" && name != dirName {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error: skill name '%s' does not match directory name '%s' in %s", name, dirName, path),
					},
				},
				IsError: true,
			}, ListSkillsOutput{}, fmt.Errorf("skill name mismatch: '%s' != '%s' in %s", name, dirName, path)
		}

		// Use the name from frontmatter if available, otherwise use directory name
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

	output := ListSkillsOutput{
		Skills: skills,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Found %d skill(s)", len(skills)),
			},
		},
	}, output, nil
}

// stripFrontmatter removes YAML frontmatter from content and returns the remaining content
func stripFrontmatter(content []byte) string {
	lines := strings.Split(string(content), "\n")
	inFrontmatter := false
	contentStart := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for frontmatter delimiters
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				// End of frontmatter - content starts after this line
				contentStart = i + 1
				break
			}
		}
	}

	// If we found frontmatter, return content after it
	if contentStart > 0 && contentStart < len(lines) {
		return strings.Join(lines[contentStart:], "\n")
	}

	// No frontmatter found, return original content
	return string(content)
}

// handleUseSkill handles the use_skill tool call
func handleUseSkill(ctx context.Context, req *mcp.CallToolRequest, args UseSkillArgs) (*mcp.CallToolResult, UseSkillOutput, error) {
	if args.Name == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: skill name is required",
				},
			},
			IsError: true,
		}, UseSkillOutput{}, fmt.Errorf("skill name is required")
	}

	// Look up the skill in the Skills map
	// Use forward slashes to match the embedded map keys (cross-platform)
	skillPath := args.Name + "/SKILL.md"
	content, found := data.Skills[skillPath]
	if !found {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: skill '%s' not found", args.Name),
				},
			},
			IsError: true,
		}, UseSkillOutput{}, fmt.Errorf("skill '%s' not found", args.Name)
	}

	// Parse metadata from YAML frontmatter
	_, description := parseSkillMetadata(content)

	// Strip frontmatter from content
	contentWithoutHeader := stripFrontmatter(content)

	output := UseSkillOutput{
		Name:        args.Name,
		Description: description,
		Content:     contentWithoutHeader,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Retrieved skill: %s", args.Name),
			},
		},
	}, output, nil
}

// handleSearchExamples handles the search_examples tool call
func handleSearchExamples(ctx context.Context, req *mcp.CallToolRequest, args SearchExamplesArgs) (*mcp.CallToolResult, SearchExamplesOutput, error) {
	if args.Query == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: search query is required",
				},
			},
			IsError: true,
		}, SearchExamplesOutput{}, fmt.Errorf("search query is required")
	}

	// Create a temporary file for the database
	tmpFile, err := os.CreateTemp("", "examples-*.db")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: failed to create temp file: %v", err),
				},
			},
			IsError: true,
		}, SearchExamplesOutput{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write embedded database to temp file
	if _, err := tmpFile.Write(data.EXAMPLES_DB); err != nil {
		tmpFile.Close()
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: failed to write database: %v", err),
				},
			},
			IsError: true,
		}, SearchExamplesOutput{}, fmt.Errorf("failed to write database: %w", err)
	}
	tmpFile.Close()

	// Open the database
	db, err := sql.Open("sqlite", tmpPath)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: failed to open database: %v", err),
				},
			},
			IsError: true,
		}, SearchExamplesOutput{}, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	sqlQuery := `
		SELECT e.id, e.package, e.symbol, e.name, e.code, e.doc_comment, e.output, e.imports, e.file
		FROM examples e
		JOIN examples_fts f ON e.rowid = f.rowid
		WHERE examples_fts MATCH ?
		ORDER BY rank
		LIMIT 10
	`

	var rows *sql.Rows
	if args.PackageFilter != "" {
		sqlQuery = `
			SELECT e.id, e.package, e.symbol, e.name, e.code, e.doc_comment, e.output, e.imports, e.file
			FROM examples e
			JOIN examples_fts f ON e.rowid = f.rowid
			WHERE examples_fts MATCH ? AND e.package = ?
			ORDER BY rank
			LIMIT 10
		`
		rows, err = db.Query(sqlQuery, args.Query, args.PackageFilter)
	} else {
		rows, err = db.Query(sqlQuery, args.Query)
	}
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: search failed: %v", err),
				},
			},
			IsError: true,
		}, SearchExamplesOutput{}, fmt.Errorf("search failed: %w", err)
	}
	defer rows.Close()

	var examples []GoExample
	for rows.Next() {
		var ex GoExample
		if err := rows.Scan(&ex.ID, &ex.Package, &ex.Symbol, &ex.Name, &ex.Code, &ex.DocComment, &ex.Output, &ex.Imports, &ex.File); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error: failed to scan row: %v", err),
					},
				},
				IsError: true,
			}, SearchExamplesOutput{}, fmt.Errorf("failed to scan row: %w", err)
		}
		examples = append(examples, ex)
	}

	output := SearchExamplesOutput{
		Examples: examples,
		Count:    len(examples),
	}

	resultText := fmt.Sprintf("Found %d example(s) matching query: %s", len(examples), args.Query)
	if args.PackageFilter != "" {
		resultText += fmt.Sprintf(" (filtered by package: %s)", args.PackageFilter)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: resultText,
			},
		},
	}, output, nil
}

// handleGetExample handles the get_example tool call
func handleGetExample(ctx context.Context, req *mcp.CallToolRequest, args GetExampleArgs) (*mcp.CallToolResult, GetExampleOutput, error) {
	if args.Symbol == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: symbol name is required",
				},
			},
			IsError: true,
		}, GetExampleOutput{}, fmt.Errorf("symbol name is required")
	}

	// Create a temporary file for the database
	tmpFile, err := os.CreateTemp("", "examples-*.db")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: failed to create temp file: %v", err),
				},
			},
			IsError: true,
		}, GetExampleOutput{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write embedded database to temp file
	if _, err := tmpFile.Write(data.EXAMPLES_DB); err != nil {
		tmpFile.Close()
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: failed to write database: %v", err),
				},
			},
			IsError: true,
		}, GetExampleOutput{}, fmt.Errorf("failed to write database: %w", err)
	}
	tmpFile.Close()

	// Open the database
	db, err := sql.Open("sqlite", tmpPath)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: failed to open database: %v", err),
				},
			},
			IsError: true,
		}, GetExampleOutput{}, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	query := `
		SELECT id, package, symbol, name, code, doc_comment, output, imports, file
		FROM examples
		WHERE symbol = ? OR name = ? OR symbol LIKE ? OR name LIKE ?
		ORDER BY name
	`

	// Add wildcards for pattern matching
	pattern := "%" + args.Symbol + "%"
	rows, err := db.Query(query, args.Symbol, args.Symbol, pattern, pattern)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: query failed: %v", err),
				},
			},
			IsError: true,
		}, GetExampleOutput{}, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var examples []GoExample
	for rows.Next() {
		var ex GoExample
		if err := rows.Scan(&ex.ID, &ex.Package, &ex.Symbol, &ex.Name,
			&ex.Code, &ex.DocComment, &ex.Output, &ex.Imports, &ex.File); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error: failed to scan row: %v", err),
					},
				},
				IsError: true,
			}, GetExampleOutput{}, fmt.Errorf("failed to scan row: %w", err)
		}
		examples = append(examples, ex)
	}

	if len(examples) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("No examples found for symbol: %s", args.Symbol),
				},
			},
		}, GetExampleOutput{Examples: []GoExample{}, Count: 0}, nil
	}

	output := GetExampleOutput{
		Examples: examples,
		Count:    len(examples),
	}

	resultText := fmt.Sprintf("Retrieved %d example(s) for symbol: %s", len(examples), args.Symbol)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: resultText,
			},
		},
	}, output, nil
}
