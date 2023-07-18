package exec

import (
	"context"

	RIOE "github.com/ibm/fp-go/context/readerioeither"
	"github.com/ibm/fp-go/exec"
	F "github.com/ibm/fp-go/function"
	GE "github.com/ibm/fp-go/internal/exec"
	IOE "github.com/ibm/fp-go/ioeither"
)

var (
	// Command executes a cancelable command
	Command = F.Curry3(command)
)

func command(name string, args []string, in []byte) RIOE.ReaderIOEither[exec.CommandOutput] {
	return func(ctx context.Context) IOE.IOEither[error, exec.CommandOutput] {
		return IOE.TryCatchError(func() (exec.CommandOutput, error) {
			return GE.Exec(ctx, name, args, in)
		})
	}
}
