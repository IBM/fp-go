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

// Package exec provides utilities for executing system commands with IOResult-based error handling.
//
// This package wraps system command execution in the IOResult monad, which represents
// IO operations that can fail with errors using Go's idiomatic (value, error) tuple pattern.
// Unlike the result/exec package, these operations explicitly acknowledge their side-effectful
// nature by returning IOResult instead of plain Result.
//
// # Overview
//
// The exec package is designed for executing system commands as IO operations that may fail.
// Since command execution is inherently side-effectful (it interacts with the operating system,
// may produce different results over time, and has observable effects), IOResult is the
// appropriate abstraction.
//
// # Basic Usage
//
// The primary function is Command, which executes a system command:
//
//	import (
//	    "github.com/IBM/fp-go/v2/bytes"
//	    "github.com/IBM/fp-go/v2/exec"
//	    "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/idiomatic/ioresult"
//	    ioexec "github.com/IBM/fp-go/v2/idiomatic/ioresult/exec"
//	)
//
//	// Execute a command and get the output
//	version := F.Pipe1(
//	    ioexec.Command("openssl")([]string{"version"})([]byte{}),
//	    ioresult.Map(F.Flow2(
//	        exec.StdOut,
//	        bytes.ToString,
//	    )),
//	)
//
//	// Run the IO operation
//	result, err := version()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result)
//
// # Command Output
//
// Commands return exec.CommandOutput, which contains both stdout and stderr as byte slices.
// Use exec.StdOut and exec.StdErr to extract the respective streams:
//
//	output := ioexec.Command("ls")([]string{"-la"})([]byte{})
//	result := ioresult.Map(func(out exec.CommandOutput) string {
//	    stdout := exec.StdOut(out)
//	    stderr := exec.StdErr(out)
//	    return bytes.ToString(stdout)
//	})(output)
//
// # Composing Commands
//
// Commands can be composed using IOResult combinators:
//
//	// Chain multiple commands together
//	pipeline := F.Pipe2(
//	    ioexec.Command("echo")([]string{"hello world"})([]byte{}),
//	    ioresult.Chain(func(out exec.CommandOutput) ioexec.IOResult[exec.CommandOutput] {
//	        input := exec.StdOut(out)
//	        return ioexec.Command("tr")([]string{"a-z", "A-Z"})(input)
//	    }),
//	    ioresult.Map(F.Flow2(exec.StdOut, bytes.ToString)),
//	)
//
// # Error Handling
//
// Commands return errors for various failure conditions:
//   - Command not found
//   - Non-zero exit status
//   - Permission errors
//   - System resource errors
//
// Handle errors using IOResult's error handling combinators:
//
//	safeCommand := F.Pipe1(
//	    ioexec.Command("risky-command")([]string{})([]byte{}),
//	    ioresult.Alt(func() (exec.CommandOutput, error) {
//	        // Fallback on error
//	        return exec.CommandOutput{}, nil
//	    }),
//	)
package exec

import (
	"context"

	"github.com/IBM/fp-go/v2/exec"
	"github.com/IBM/fp-go/v2/function"
	INTE "github.com/IBM/fp-go/v2/internal/exec"
)

var (
	// Command executes a system command with side effects and returns an IOResult.
	//
	// This function is curried to allow partial application. It takes three parameters:
	//   - name: The command name or path to execute
	//   - args: Command-line arguments as a slice of strings
	//   - in: Input bytes to send to the command's stdin
	//
	// Returns IOResult[exec.CommandOutput] which, when executed, will run the command
	// and return either the command output or an error.
	//
	// The command is executed using the system's default shell context. The output
	// contains both stdout and stderr as byte slices, accessible via exec.StdOut
	// and exec.StdErr respectively.
	//
	// Example:
	//
	//	// Simple command execution
	//	version := Command("node")([]string{"--version"})([]byte{})
	//	result, err := version()
	//
	//	// With input piped to stdin
	//	echo := Command("cat")([]string{})([]byte("hello world"))
	//
	//	// Partial application for reuse
	//	git := Command("git")
	//	status := git([]string{"status"})([]byte{})
	//	log := git([]string{"log", "--oneline"})([]byte{})
	//
	//	// Composed with IOResult combinators
	//	result := F.Pipe1(
	//	    Command("openssl")([]string{"version"})([]byte{}),
	//	    ioresult.Map(F.Flow2(exec.StdOut, bytes.ToString)),
	//	)
	Command = function.Curry3(command)
)

func command(name string, args []string, in []byte) IOResult[exec.CommandOutput] {
	return func() (exec.CommandOutput, error) {
		return INTE.Exec(context.Background(), name, args, in)
	}
}
