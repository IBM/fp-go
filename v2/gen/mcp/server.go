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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/IBM/fp-go/v2/gen/data"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FindExamplesArgs represents the arguments for the find_examples tool
type FindExamplesArgs struct {
	Query               string   `json:"query" jsonschema:"The primary, conceptual search query. This should capture the user's main goal or question (e.g., 'using Option monad' or 'how to compose functions'). The query will be processed by a full-text search engine. Key Syntax: AND (default with spaces), OR operator, NOT operator, Grouping with (), Phrase Search with quotes, Prefix Search with *."`
	Keywords            []string `json:"keywords,omitempty" jsonschema:"A list of specific, exact keywords to narrow the search. Use this for precise terms like function names or type names."`
	RequiredPackages    []string `json:"required_packages,omitempty" jsonschema:"A list of Go packages that an example must use. Use this when the user's request is specific to a feature within a certain package (e.g., if the user asks about Option, you should filter by github.com/IBM/fp-go/v2/option)."`
	RelatedConcepts     []string `json:"related_concepts,omitempty" jsonschema:"A list of high-level concepts to filter by. Use this to find examples related to broader functional programming ideas or patterns (e.g., monads, functors, composition, error handling)."`
	IncludeExperimental bool     `json:"include_experimental,omitempty" jsonschema:"By default, this tool returns only production-safe examples. Set this to true only if the user explicitly asks for a bleeding-edge feature or if a stable solution cannot be found. If set to true, you MUST warn the user that the example uses experimental APIs not suitable for production."`
}

// ExampleResult represents a single example result
type ExampleResult struct {
	Title            string   `json:"title" jsonschema:"The title of the example. Use this as a heading when presenting the example to the user."`
	Summary          string   `json:"summary" jsonschema:"A one-sentence summary of the example's purpose. Use this to help the user decide if the example is relevant to them."`
	Keywords         []string `json:"keywords,omitempty" jsonschema:"A list of keywords for the example. You can use these to explain why this example was a good match for the user's query."`
	RequiredPackages []string `json:"required_packages,omitempty" jsonschema:"A list of Go packages required for the example to work. Before presenting the code, you should inform the user if any of these packages need to be installed."`
	RelatedConcepts  []string `json:"related_concepts,omitempty" jsonschema:"A list of related concepts. You can suggest these to the user as topics for follow-up questions."`
	RelatedTools     []string `json:"related_tools,omitempty" jsonschema:"A list of related MCP tools. You can suggest these as potential next steps for the user."`
	Content          string   `json:"content" jsonschema:"A complete, self-contained Go code example in Markdown format. This should be presented to the user inside a markdown code block."`
	Snippet          string   `json:"snippet,omitempty" jsonschema:"A contextual snippet from the content showing the matched search term. This field is critical for efficiently evaluating a result's relevance."`
}

// FindExamplesOutput represents the output of the find_examples tool
type FindExamplesOutput struct {
	Examples []ExampleResult `json:"examples" jsonschema:"List of example results matching the search criteria"`
}

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
	Name    string `json:"name" jsonschema:"The name of the skill"`
	Content string `json:"content" jsonschema:"The full content of the skill's SKILL.md file"`
}

// NewServer creates a new MCP server with fp-go tools
func NewServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "fp-go-generator",
		Version: "1.0.0",
	}, nil)

	// Register the find_examples tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "find_examples",
		Description: "Find code examples matching a search query. Searches through test files, example directories, and documentation for relevant fp-go code snippets demonstrating functional programming patterns.",
	}, handleFindExamples)

	// Register the list_skills tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_skills",
		Description: "List all available fp-go skills. Each skill provides specialized knowledge and guidance for specific fp-go features and patterns. Returns a list of skills with their names, descriptions, and paths.",
	}, handleListSkills)

	// Register the use_skill tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "use_skill",
		Description: "Retrieve the full content of a specific fp-go skill by name. Use this after list_skills to get detailed documentation, examples, and best practices for a particular fp-go feature or pattern.",
	}, handleUseSkill)

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

// handleFindExamples handles the find_examples tool call
func handleFindExamples(ctx context.Context, req *mcp.CallToolRequest, args FindExamplesArgs) (*mcp.CallToolResult, FindExamplesOutput, error) {
	// TODO: Implement actual example search logic
	// This would search through:
	// - Test files (*_test.go) for example functions
	// - Example directories for code samples
	// - Documentation files for usage patterns
	// - Match against keywords, packages, and concepts

	// For now, return sample examples
	output := FindExamplesOutput{
		Examples: []ExampleResult{
			{
				Title:    "Using Option Monad for Null Safety",
				Summary:  "Demonstrates how to use the Option monad to handle nullable values safely without nil checks.",
				Keywords: []string{"Option", "monad", "null safety", "Some", "None"},
				RequiredPackages: []string{
					"github.com/IBM/fp-go/v2/option",
				},
				RelatedConcepts: []string{"monads", "functional error handling", "type safety"},
				RelatedTools:    []string{},
				Content:         "```go\npackage main\n\nimport (\n\t\"fmt\"\n\tO \"github.com/IBM/fp-go/v2/option\"\n)\n\nfunc main() {\n\t// Create Some and None values\n\tsome := O.Some(42)\n\tnone := O.None[int]()\n\n\t// Use Map to transform values\n\tresult := O.Map(func(x int) int { return x * 2 })(some)\n\tfmt.Println(O.GetOrElse(func() int { return 0 })(result)) // Output: 84\n}\n```",
				Snippet:         "O.Map(func(x int) int { return x * 2 })(some)",
			},
		},
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Found %d example(s) matching query: %s", len(output.Examples), args.Query),
			},
		},
	}, output, nil
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
		if strings.HasPrefix(trimmed, "name:") {
			name = strings.TrimSpace(strings.TrimPrefix(trimmed, "name:"))
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

	output := UseSkillOutput{
		Name:    args.Name,
		Content: string(content),
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Retrieved skill: %s", args.Name),
			},
		},
	}, output, nil
}
