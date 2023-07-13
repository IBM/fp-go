package cli

import (
	C "github.com/urfave/cli/v2"
)

func Commands() []*C.Command {
	return []*C.Command{
		PipeCommand(),
		OptionCommand(),
		EitherCommand(),
	}
}
