package cli

import (
	C "github.com/urfave/cli/v2"
)

const (
	keyFilename = "filename"
	keyCount    = "count"
)

var (
	flagFilename = &C.StringFlag{
		Name:  keyFilename,
		Value: "gen.go",
		Usage: "Name of the generated file",
	}

	flagCount = &C.IntFlag{
		Name:  keyCount,
		Value: 20,
		Usage: "Number of variations to create",
	}
)
