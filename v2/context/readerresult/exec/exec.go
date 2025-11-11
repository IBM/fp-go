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

package exec

import (
	"context"

	"github.com/IBM/fp-go/v2/exec"
	"github.com/IBM/fp-go/v2/function"
	INTE "github.com/IBM/fp-go/v2/internal/exec"
	"github.com/IBM/fp-go/v2/result"
)

var (
	// Command executes a command
	// use this version if the command does not produce any side effect, i.e. if the output is uniquely determined by by the input
	// typically you'd rather use the ReaderIOEither version of the command
	Command = function.Curry3(command)
)

func command(name string, args []string, in []byte) ReaderResult[exec.CommandOutput] {
	return func(ctx context.Context) Result[exec.CommandOutput] {
		return result.TryCatchError(INTE.Exec(ctx, name, args, in))
	}
}
