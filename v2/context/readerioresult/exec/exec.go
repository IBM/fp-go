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

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/exec"
	F "github.com/IBM/fp-go/v2/function"
	GE "github.com/IBM/fp-go/v2/internal/exec"
	IOE "github.com/IBM/fp-go/v2/ioeither"
)

var (
	// Command executes a cancelable command
	Command = F.Curry3(command)
)

func command(name string, args []string, in []byte) RIOE.ReaderIOResult[exec.CommandOutput] {
	return func(ctx context.Context) IOE.IOEither[error, exec.CommandOutput] {
		return IOE.TryCatchError(func() (exec.CommandOutput, error) {
			return GE.Exec(ctx, name, args, in)
		})
	}
}
