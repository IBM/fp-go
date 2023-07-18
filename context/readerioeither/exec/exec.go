package exec

import (
	"context"

	RIOE "github.com/IBM/fp-go/context/readerioeither"
	"github.com/IBM/fp-go/exec"
	F "github.com/IBM/fp-go/function"
	GE "github.com/IBM/fp-go/internal/exec"
	IOE "github.com/IBM/fp-go/ioeither"
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
