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

	E "github.com/IBM/fp-go/v2/effect"
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	R "github.com/IBM/fp-go/v2/result"
	C "github.com/urfave/cli/v3"
)

// CommandEffect represents a CLI command action as an Effect.
// The Effect takes a *C.Command as context and produces a result.
type CommandEffect = E.Effect[*C.Command, F.Void]

// ToAction converts a CommandEffect into a standard urfave/cli Action function.
// This allows Effect-based command handlers to be used with the cli/v3 framework.
//
// The conversion process:
//  1. Takes the Effect which expects a *C.Command context
//  2. Executes it with the provided command
//  3. Runs the resulting IO operation
//  4. Converts the Result to either nil (success) or error (failure)
//
// # Parameters
//
//   - effect: The CommandEffect to convert
//
// # Returns
//
//   - A function compatible with C.Command.Action signature
//
// # Example Usage
//
//	effect := func(cmd *C.Command) E.Thunk[F.Void] {
//	    return func(ctx context.Context) E.IOResult[F.Void] {
//	        return func() R.Result[F.Void] {
//	            // Command logic here
//	            return R.Of(F.Void{})
//	        }
//	    }
//	}
//	action := ToAction(effect)
//	command := &C.Command{
//	    Name: "example",
//	    Action: action,
//	}
func ToAction(effect CommandEffect) func(context.Context, *C.Command) error {
	return func(ctx context.Context, cmd *C.Command) error {
		// Execute the effect: cmd -> ctx -> IO -> Result
		return F.Pipe3(
			ctx,
			effect(cmd),
			io.Run,
			// Convert Result[Void] to error
			ET.Fold(F.Identity[error], F.Constant1[F.Void, error](nil)),
		)
	}
}

// FromAction converts a standard urfave/cli Action function into a CommandEffect.
// This allows existing cli/v3 action handlers to be lifted into the Effect type.
//
// The conversion process:
//  1. Takes a standard action function (context.Context, *C.Command) -> error
//  2. Wraps it in the Effect structure
//  3. Converts the error result to a Result type
//
// # Parameters
//
//   - action: The standard cli/v3 action function to convert
//
// # Returns
//
//   - A CommandEffect that wraps the original action
//
// # Example Usage
//
//	standardAction := func(ctx context.Context, cmd *C.Command) error {
//	    // Existing command logic
//	    return nil
//	}
//	effect := FromAction(standardAction)
//	// Now can be composed with other Effects
func FromAction(action func(context.Context, *C.Command) error) CommandEffect {
	return func(cmd *C.Command) E.Thunk[F.Void] {
		return func(ctx context.Context) E.IOResult[F.Void] {
			return func() R.Result[F.Void] {
				err := action(ctx, cmd)
				if err != nil {
					return R.Left[F.Void](err)
				}
				return R.Of(F.Void{})
			}
		}
	}
}

// MakeCommand creates a new Command with an Effect-based action.
// This is a convenience function that combines command creation with Effect conversion.
//
// # Parameters
//
//   - name: The command name
//   - usage: The command usage description
//   - flags: The command flags
//   - effect: The CommandEffect to use as the action
//
// # Returns
//
//   - A *C.Command configured with the Effect-based action
//
// # Example Usage
//
//	cmd := MakeCommand(
//	    "process",
//	    "Process data files",
//	    []C.Flag{
//	        &C.StringFlag{Name: "input", Usage: "Input file"},
//	    },
//	    func(cmd *C.Command) E.Thunk[F.Void] {
//	        return func(ctx context.Context) E.IOResult[F.Void] {
//	            return func() R.Result[F.Void] {
//	                input := cmd.String("input")
//	                // Process input...
//	                return R.Of(F.Void{})
//	            }
//	        }
//	    },
//	)
func MakeCommand(
	name string,
	usage string,
	flags []C.Flag,
	effect CommandEffect,
) *C.Command {
	return &C.Command{
		Name:   name,
		Usage:  usage,
		Flags:  flags,
		Action: ToAction(effect),
	}
}

// MakeCommandWithSubcommands creates a new Command with subcommands and an Effect-based action.
//
// # Parameters
//
//   - name: The command name
//   - usage: The command usage description
//   - flags: The command flags
//   - commands: The subcommands
//   - effect: The CommandEffect to use as the action
//
// # Returns
//
//   - A *C.Command configured with subcommands and the Effect-based action
//
// # Example Usage
//
//	cmd := MakeCommandWithSubcommands(
//	    "app",
//	    "Application commands",
//	    []C.Flag{},
//	    []*C.Command{subCmd1, subCmd2},
//	    defaultEffect,
//	)
func MakeCommandWithSubcommands(
	name string,
	usage string,
	flags []C.Flag,
	commands []*C.Command,
	effect CommandEffect,
) *C.Command {
	return &C.Command{
		Name:     name,
		Usage:    usage,
		Flags:    flags,
		Commands: commands,
		Action:   ToAction(effect),
	}
}
