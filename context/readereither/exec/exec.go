package exec

import (
	"context"

	RE "github.com/IBM/fp-go/context/readereither"
	E "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/exec"
	F "github.com/IBM/fp-go/function"
	GE "github.com/IBM/fp-go/internal/exec"
)

var (
	// Command executes a command
	// use this version if the command does not produce any side effect, i.e. if the output is uniquely determined by by the input
	// typically you'd rather use the ReaderIOEither version of the command
	Command = F.Curry3(command)
)

func command(name string, args []string, in []byte) RE.ReaderEither[exec.CommandOutput] {
	return func(ctx context.Context) E.Either[error, exec.CommandOutput] {
		return E.TryCatchError(func() (exec.CommandOutput, error) {
			return GE.Exec(ctx, name, args, in)
		})
	}
}
