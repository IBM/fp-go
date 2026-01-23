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

	A "github.com/IBM/fp-go/v2/array"
	C "github.com/urfave/cli/v3"
)

// [GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GTAB ~func() ET.Either[E, T.Tuple2[A, B]], E, A, B any](a GA, b GB) GTAB {

func nonGenericIOEither(param string) string {
	return fmt.Sprintf("IOEither[E, %s]", param)
}

var extrasIOEither = A.From("E")

func generateIOEitherSequenceT(f, fg *os.File, i int) {
	generateGenericSequenceT("", nonGenericIOEither, extrasIOEither)(f, i)
	generateGenericSequenceT("Seq", nonGenericIOEither, extrasIOEither)(f, i)
	generateGenericSequenceT("Par", nonGenericIOEither, extrasIOEither)(f, i)
}

func generateIOEitherSequenceTuple(f, fg *os.File, i int) {
	generateGenericSequenceTuple("", nonGenericIOEither, extrasIOEither)(f, i)
	generateGenericSequenceTuple("Seq", nonGenericIOEither, extrasIOEither)(f, i)
	generateGenericSequenceTuple("Par", nonGenericIOEither, extrasIOEither)(f, i)
}

func generateIOEitherTraverseTuple(f, fg *os.File, i int) {
	generateGenericTraverseTuple("", nonGenericIOEither, extrasIOEither)(f, i)
	generateGenericTraverseTuple("Seq", nonGenericIOEither, extrasIOEither)(f, i)
	generateGenericTraverseTuple("Par", nonGenericIOEither, extrasIOEither)(f, i)
}

func generateIOEitherUneitherize(f, fg *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// Uneitherize%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [IOEither[error, R]]\n", i, i+1, i)
	fmt.Fprintf(f, "func Uneitherize%d[F ~func(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") IOEither[error, R]")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j+1)
	}
	fmt.Fprintf(f, ", R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") (R, error) {\n")
	fmt.Fprintf(f, "  return G.Uneitherize%d[IOEither[error, R]](f)\n", i)
	fmt.Fprintln(f, "}")

	// generic version
	fmt.Fprintf(fg, "\n// Uneitherize%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [GIOA]\n", i, i, i)
	fmt.Fprintf(fg, "func Uneitherize%d[GIOA ~func() ET.Either[error, R], GTA ~func(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j+1)
	}
	fmt.Fprintf(fg, ") GIOA")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j+1)
	}
	fmt.Fprintf(fg, ", R any](f GTA) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j+1)
	}
	fmt.Fprintf(fg, ") (R, error) {\n")
	fmt.Fprintf(fg, "  return func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "t%d T%d", j+1, j+1)
	}
	fmt.Fprintf(fg, ") (R, error) {\n")
	fmt.Fprintf(fg, "    return ET.Unwrap(f(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "t%d", j+1)
	}
	fmt.Fprintf(fg, ")())\n")
	fmt.Fprintf(fg, "  }\n")
	fmt.Fprintf(fg, "}\n")
}

func generateIOEitherEitherize(f, fg *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// Eitherize%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [IOEither[error, R]]\n", i, i+1, i)
	fmt.Fprintf(f, "func Eitherize%d[F ~func(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") (R, error)")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j+1)
	}
	fmt.Fprintf(f, ", R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") IOEither[error, R] {\n")
	fmt.Fprintf(f, "  return G.Eitherize%d[IOEither[error, R]](f)\n", i)
	fmt.Fprintln(f, "}")

	// generic version
	fmt.Fprintf(fg, "\n// Eitherize%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [GIOA]\n", i, i, i)
	fmt.Fprintf(fg, "func Eitherize%d[GIOA ~func() ET.Either[error, R], F ~func(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j+1)
	}
	fmt.Fprintf(fg, ") (R, error)")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j+1)
	}
	fmt.Fprintf(fg, ", R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j+1)
	}
	fmt.Fprintf(fg, ") GIOA {\n")
	fmt.Fprintf(fg, "  e := ET.Eitherize%d(f)\n", i)
	fmt.Fprintf(fg, "  return func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "t%d T%d", j+1, j+1)
	}
	fmt.Fprintf(fg, ") GIOA {\n")
	fmt.Fprintf(fg, "    return func() ET.Either[error, R] {\n")
	fmt.Fprintf(fg, "      return e(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "t%d", j+1)
	}
	fmt.Fprintf(fg, ")\n")
	fmt.Fprintf(fg, "    }}\n")
	fmt.Fprintf(fg, "}\n")
}

func generateIOEitherHelpers(filename string, count int) error {
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
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/tuple"
)
`, pkg)

	// some header
	fmt.Fprintln(fg, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(fg, "// This file was generated by robots at")
	fmt.Fprintf(fg, "// %s\n", time.Now())

	fmt.Fprintf(fg, "package generic\n\n")

	fmt.Fprintf(fg, `
import (
	ET "github.com/IBM/fp-go/v2/either"
)
`)

	// eitherize
	generateIOEitherEitherize(f, fg, 0)
	// uneitherize
	generateIOEitherUneitherize(f, fg, 0)

	for i := 1; i <= count; i++ {
		// eitherize
		generateIOEitherEitherize(f, fg, i)
		// uneitherize
		generateIOEitherUneitherize(f, fg, i)
		// sequenceT
		generateIOEitherSequenceT(f, fg, i)
		// sequenceTuple
		generateIOEitherSequenceTuple(f, fg, i)
		// traverseTuple
		generateIOEitherTraverseTuple(f, fg, i)
	}

	return nil
}

func IOEitherCommand() *C.Command {
	return &C.Command{
		Name:  "ioeither",
		Usage: "generate code for IOEither",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generateIOEitherHelpers(
				cmd.String(keyFilename),
				cmd.Int(keyCount),
			)
		},
	}
}
