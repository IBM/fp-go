// Copyright (c) 2023 IBM Corp.
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

	E "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/exec"
	F "github.com/IBM/fp-go/function"
	GE "github.com/IBM/fp-go/internal/exec"
)

var (
	// Command executes a command
	// use this version if the command does not produce any side effect, i.e. if the output is uniquely determined by by the input
	// typically you'd rather use the IOEither version of the command
	Command = F.Curry3(command)
)

func command(name string, args []string, in []byte) E.Either[error, exec.CommandOutput] {
	return E.TryCatchError(func() (exec.CommandOutput, error) {
		return GE.Exec(context.Background(), name, args, in)
	})
}
