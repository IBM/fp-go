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

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GenerateLensArgs represents the arguments for the generate_lens tool
type GenerateLensArgs struct {
	StructName  string `json:"struct_name" jsonschema:"Name of the struct to generate lens for"`
	PackagePath string `json:"package_path" jsonschema:"Package path containing the struct"`
	OutputFile  string `json:"output_file,omitempty" jsonschema:"Output file name (default: gen.go)"`
}

// NewServer creates a new MCP server with lens generation tools
func NewServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "fp-go-generator",
		Version: "1.0.0",
	}, nil)

	// Register the generate_lens tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_lens",
		Description: "Generate lens code for a Go struct",
	}, handleGenerateLens)

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

// handleGenerateLens handles the generate_lens tool call
func handleGenerateLens(ctx context.Context, req *mcp.CallToolRequest, args GenerateLensArgs) (*mcp.CallToolResult, any, error) {
	// Set default output file if not provided
	outputFile := args.OutputFile
	if outputFile == "" {
		outputFile = "gen.go"
	}

	// TODO: Integrate with actual lens generation logic from cli/lens.go
	// For now, return a placeholder response
	message := fmt.Sprintf(
		"Would generate lens for struct '%s' from package '%s' to file '%s'",
		args.StructName,
		args.PackagePath,
		outputFile,
	)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: message},
		},
	}, nil, nil
}

// Made with Bob
