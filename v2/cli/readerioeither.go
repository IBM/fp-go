// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	C "github.com/urfave/cli/v3"
)

func generateReaderIOEitherFrom(f, fg *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// From%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [ReaderIOEither[R]]\n// The first parameter is considered to be the context [C].\n", i, i+1, i)
	fmt.Fprintf(f, "func From%d[F ~func(C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ") func() (R, error)")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ", C, R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, ") ReaderIOEither[C, error, R] {\n")
	fmt.Fprintf(f, "  return G.From%d[ReaderIOEither[C, error, R]](f)\n", i)
	fmt.Fprintln(f, "}")

	// generic version
	fmt.Fprintf(fg, "\n// From%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [GRA]\n// The first parameter is considerd to be the context [C].\n", i, i+1, i)
	fmt.Fprintf(fg, "func From%d[GRA ~func(C) GIOA, F ~func(C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ") func() (R, error), GIOA ~func() E.Either[error, R]")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ", C, R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j)
	}
	fmt.Fprintf(fg, ") GRA {\n")

	fmt.Fprintf(fg, "  return RD.From%d[GRA](func(r C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", t%d T%d", j, j)
	}
	fmt.Fprintf(fg, ") GIOA {\n")
	fmt.Fprintf(fg, "    return E.Eitherize0(f(r")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", t%d", j)
	}
	fmt.Fprintf(fg, "))\n")
	fmt.Fprintf(fg, "  })\n")
	fmt.Fprintf(fg, "}\n")
}

func generateReaderIOEitherEitherize(f, fg *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// Eitherize%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [ReaderIOEither[C, error, R]]\n// The first parameter is considered to be the context [C].\n", i, i+1, i)
	fmt.Fprintf(f, "func Eitherize%d[F ~func(C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ") (R, error)")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ", C, R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, ") ReaderIOEither[C, error, R] {\n")
	fmt.Fprintf(f, "  return G.Eitherize%d[ReaderIOEither[C, error, R]](f)\n", i)
	fmt.Fprintln(f, "}")

	// generic version
	fmt.Fprintf(fg, "\n// Eitherize%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [GRA]\n// The first parameter is considered to be the context [C].\n", i, i, i)
	fmt.Fprintf(fg, "func Eitherize%d[GRA ~func(C) GIOA, F ~func(C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ") (R, error), GIOA ~func() E.Either[error, R]")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ", C, R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j)
	}
	fmt.Fprintf(fg, ") GRA {\n")

	fmt.Fprintf(fg, "  return From%d[GRA](func(r C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", t%d T%d", j, j)
	}
	fmt.Fprintf(fg, ") func() (R, error) {\n")
	fmt.Fprintf(fg, "    return func() (R, error) {\n")
	fmt.Fprintf(fg, "      return f(r")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", t%d", j)
	}
	fmt.Fprintf(fg, ")\n")
	fmt.Fprintf(fg, "    }})\n")
	fmt.Fprintf(fg, "}\n")
}

func generateReaderIOEitherUneitherize(f, fg *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// Uneitherize%d converts a function with %d parameters returning a [ReaderIOEither[C, error, R]] into a function with %d parameters returning a tuple.\n// The first parameter is considered to be the context [C].\n", i, i+1, i)
	fmt.Fprintf(f, "func Uneitherize%d[F ~func(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, ") ReaderIOEither[C, error, R]")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ", C, R any](f F) func(C")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ") (R, error) {\n")
	fmt.Fprintf(f, "  return G.Uneitherize%d[ReaderIOEither[C, error, R]", i)

	fmt.Fprintf(f, ", func(C")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ")(R, error)](f)\n")
	fmt.Fprintln(f, "}")

	// generic version
	fmt.Fprintf(fg, "\n// Uneitherize%d converts a function with %d parameters returning a [GRA] into a function with %d parameters returning a tuple.\n// The first parameter is considered to be the context [C].\n", i, i, i)
	fmt.Fprintf(fg, "func Uneitherize%d[GRA ~func(C) GIOA, F ~func(C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ") (R, error), GIOA ~func() E.Either[error, R]")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ", C, R any](f func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j)
	}
	fmt.Fprintf(fg, ") GRA) F {\n")

	fmt.Fprintf(fg, "  return func(c C")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", t%d T%d", j, j)
	}
	fmt.Fprintf(fg, ") (R, error) {\n")
	fmt.Fprintf(fg, "    return E.UnwrapError(f(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "t%d", j)
	}
	fmt.Fprintf(fg, ")(c)())\n")
	fmt.Fprintf(fg, "  }\n")
	fmt.Fprintf(fg, "}\n")
}

func generateReaderIOEitherHelpers(filename string, count int) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	pkg := filepath.Base(absDir)
	f, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer f.Close()
	// construct subdirectory
	genFilename := filepath.Join("generic", filename)
	err = os.MkdirAll("generic", os.ModePerm)
	if err != nil {
		return err
	}
	fg, err := os.Create(filepath.Clean(genFilename))
	if err != nil {
		return err
	}
	defer fg.Close()

	// log
	log.Printf("Generating code in [%s] for package [%s] with [%d] repetitions ...", filename, pkg, count)

	// some header
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, "// This file was generated by robots at")
	fmt.Fprintf(f, "// %s\n\n", time.Now())

	fmt.Fprintf(f, "package %s\n\n", pkg)

	fmt.Fprintf(f, `
import (
	G "github.com/IBM/fp-go/v2/%s/generic"	
)
`, pkg)

	// some header
	fmt.Fprintln(fg, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(fg, "// This file was generated by robots at")
	fmt.Fprintf(fg, "// %s\n", time.Now())

	fmt.Fprintf(fg, "package generic\n\n")

	fmt.Fprintf(fg, `
import (
	E "github.com/IBM/fp-go/v2/either"
	RD "github.com/IBM/fp-go/v2/reader/generic"
)
`)

	// from
	generateReaderIOEitherFrom(f, fg, 0)
	// eitherize
	generateReaderIOEitherEitherize(f, fg, 0)
	// uneitherize
	generateReaderIOEitherUneitherize(f, fg, 0)

	for i := 1; i <= count; i++ {
		// from
		generateReaderIOEitherFrom(f, fg, i)
		// eitherize
		generateReaderIOEitherEitherize(f, fg, i)
		// uneitherize
		generateReaderIOEitherUneitherize(f, fg, i)
	}

	return nil
}

func ReaderIOEitherCommand() *C.Command {
	return &C.Command{
		Name:  "readerioeither",
		Usage: "generate code for ReaderIOEither",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generateReaderIOEitherHelpers(
				cmd.String(keyFilename),
				cmd.Int(keyCount),
			)
		},
	}
}
