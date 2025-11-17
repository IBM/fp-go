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

// Package exec provides utilities for executing system commands with Either-based error handling.
package exec

import (
	"context"

	"github.com/IBM/fp-go/v2/exec"
	"github.com/IBM/fp-go/v2/idiomatic/result"
	GE "github.com/IBM/fp-go/v2/internal/exec"
)

var (
	// Command executes a system command and returns the result as an Either.
	// Use this version if the command does not produce any side effects,
	// i.e., if the output is uniquely determined by the input.
	// For commands with side effects, typically you'd use the IOEither version instead.
	//
	// Parameters (curried):
	//   - name: The command name/path
	//   - args: Command arguments
	//   - in: Input bytes to send to the command's stdin
	//
	// Returns Either[error, CommandOutput] containing the command's output or an error.
	//
	// Example:
	//
	//	result := exec.Command("echo")( []string{"hello"})([]byte{})
	//	// result is Right(CommandOutput{Stdout: "hello\n", ...})
	Command = result.Curry3(command)
)

func command(name string, args []string, in []byte) (exec.CommandOutput, error) {
	return GE.Exec(context.Background(), name, args, in)
}
