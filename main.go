package main

import (
	"log"
	"os"

	"github.com/ibm/fp-go/cli"

	C "github.com/urfave/cli/v2"
)

func main() {

	app := &C.App{
		Name:  "fp-go",
		Usage: "Code generation for fp-go",
		Commands: []*C.Command{
			cli.PipeCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
