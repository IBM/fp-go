package exec

import (
	"context"

	E "github.com/ibm/fp-go/either"
	"github.com/ibm/fp-go/exec"
	F "github.com/ibm/fp-go/function"
	GE "github.com/ibm/fp-go/internal/exec"
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
