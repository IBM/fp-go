package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	EX "github.com/IBM/fp-go/exec"

	T "github.com/IBM/fp-go/tuple"
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
	return T.MakeTuple2(stdOut.Bytes(), stdErr.Bytes()), err
}
