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
	"bytes"
	"context"
	"fmt"
	"os/exec"

	EX "github.com/IBM/fp-go/v2/exec"

	P "github.com/IBM/fp-go/v2/pair"
)

func Exec(ctx context.Context, name string, args []string, in []byte) (EX.CommandOutput, error) {
	// command input
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = bytes.NewReader(in)
	// command result
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	// execute the command
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("command execution of [%s][%s] failed, stdout [%s], stderr [%s], cause [%w]", name, args, stdOut.String(), stdErr.String(), err)
	}
	// return the outputs
	return P.MakePair(stdOut.Bytes(), stdErr.Bytes()), err
}
