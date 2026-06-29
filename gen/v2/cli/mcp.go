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

	"github.com/IBM/fp-go/gen/v2/mcp"

	C "github.com/urfave/cli/v3"
)

// McpCommand returns the CLI command for starting the MCP server
func McpCommand() *C.Command {
	return &C.Command{
		Name:  "mcp",
		Usage: "Start the MCP server with stdio transport",
		Description: `Start the Model Context Protocol (MCP) server that provides tools for code generation.
The server communicates over stdin/stdout and can be used by MCP clients.`,
		Flags: []C.Flag{
			&C.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose logging to stderr",
			},
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			verbose := cmd.Bool("verbose")
			return mcp.Run(ctx, verbose)
		},
	}
}
