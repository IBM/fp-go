package generic

import (
	"context"

	ET "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/exec"
	GE "github.com/IBM/fp-go/internal/exec"
)

// Command executes a command
func Command[GA ~func() ET.Either[error, exec.CommandOutput]](name string, args []string, in []byte) GA {
	return TryCatchError[GA](func() (exec.CommandOutput, error) {
		return GE.Exec(context.Background(), name, args, in)
	})
}
