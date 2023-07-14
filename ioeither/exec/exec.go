package exec

import (
	"github.com/ibm/fp-go/exec"
	F "github.com/ibm/fp-go/function"
	IOE "github.com/ibm/fp-go/ioeither"
	G "github.com/ibm/fp-go/ioeither/generic"
)

var (
	// Command executes a command
	Command = F.Curry3(G.Command[IOE.IOEither[error, exec.CommandOutput]])
)
