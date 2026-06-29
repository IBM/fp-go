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

	"github.com/otiai10/copy"
	C "github.com/urfave/cli/v3"
)

const (
	keySrcDir = "src"
	keyDstDir = "dst"
)

var (
	flagSrcDir = &C.StringFlag{
		Name:     keySrcDir,
		Usage:    "Source directory to copy from",
		Required: true,
	}

	flagDstDir = &C.StringFlag{
		Name:     keyDstDir,
		Usage:    "Destination directory to copy to",
		Required: true,
	}
)

// CpCommand returns the CLI command for copying files from source to destination
func CpCommand() *C.Command {
	return &C.Command{
		Name:  "cp",
		Usage: "Copy files from source directory to destination directory",
		Description: `Copy all files and subdirectories from the source directory to the destination directory.
The destination directory will be created if it doesn't exist.`,
		Flags: []C.Flag{
			flagSrcDir,
			flagDstDir,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			srcDir := cmd.String(keySrcDir)
			dstDir := cmd.String(keyDstDir)

			return copy.Copy(srcDir, dstDir)
		},
	}
}
