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
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	C "github.com/urfave/cli/v3"
)

// Deprecated:
func generateNestedCallbacks(i, total int) string {
	var buf strings.Builder
	for j := i; j < total; j++ {
		if j > i {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("func(T%d)", j+1))
	}
	if i > 0 {
		buf.WriteString(" ")
	}
	buf.WriteString(tupleType("T")(total))
	return buf.String()
}

func generateNestedCallbacksPlain(i, total int) string {
	fs := A.MakeBy(total-i, func(j int) string {
		return fmt.Sprintf("func(T%d)", j+i+1)
	})
	ts := A.Of(tupleTypePlain("T")(total))
	return joinAll(" ")(fs, ts)
}

func generateContextReaderIOEitherEitherize(f, fg *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// Eitherize%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [ReaderIOEither[R]]\n// The inverse function is [Uneitherize%d]\n", i, i, i, i)
	fmt.Fprintf(f, "func Eitherize%d[F ~func(context.Context", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ") (R, error)")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ", R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, ") ReaderIOEither[R] {\n")
	fmt.Fprintf(f, "  return G.Eitherize%d[ReaderIOEither[R]](f)\n", i)
	fmt.Fprintln(f, "}")

	// generic version
	fmt.Fprintf(fg, "\n// Eitherize%d converts a function with %d parameters returning a tuple into a function with %d parameters returning a [GRA]\n// The inverse function is [Uneitherize%d]\n", i, i, i, i)
	fmt.Fprintf(fg, "func Eitherize%d[GRA ~func(context.Context) GIOA, F ~func(context.Context", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ") (R, error), GIOA ~func() E.Either[error, R]")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ", R any](f F) func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j)
	}
	fmt.Fprintf(fg, ") GRA {\n")
	fmt.Fprintf(fg, "  return RE.Eitherize%d[GRA](f)\n", i)
	fmt.Fprintln(fg, "}")
}

func generateContextReaderIOEitherUneitherize(f, fg *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// Uneitherize%d converts a function with %d parameters returning a [ReaderIOEither[R]] into a function with %d parameters returning a tuple.\n// The first parameter is considered to be the [context.Context].\n", i, i+1, i)
	fmt.Fprintf(f, "func Uneitherize%d[F ~func(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, ") ReaderIOEither[R]")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ", R any](f F) func(context.Context")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ") (R, error) {\n")
	fmt.Fprintf(f, "  return G.Uneitherize%d[ReaderIOEither[R]", i)

	fmt.Fprintf(f, ", func(context.Context")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ")(R, error)](f)\n")
	fmt.Fprintln(f, "}")

	// generic version
	fmt.Fprintf(fg, "\n// Uneitherize%d converts a function with %d parameters returning a [GRA] into a function with %d parameters returning a tuple.\n// The first parameter is considered to be the [context.Context].\n", i, i, i)
	fmt.Fprintf(fg, "func Uneitherize%d[GRA ~func(context.Context) GIOA, F ~func(context.Context", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ") (R, error), GIOA ~func() E.Either[error, R]")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ", R any](f func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "T%d", j)
	}
	fmt.Fprintf(fg, ") GRA) F {\n")

	fmt.Fprintf(fg, "  return func(c context.Context")
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

func nonGenericContextReaderIOEither(param string) string {
	return fmt.Sprintf("ReaderIOEither[%s]", param)
}

var extrasContextReaderIOEither = A.Empty[string]()

func generateContextReaderIOEitherSequenceT(f *os.File, i int) {
	generateGenericSequenceT("", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
	generateGenericSequenceT("Seq", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
	generateGenericSequenceT("Par", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
}

func generateContextReaderIOEitherSequenceTuple(f *os.File, i int) {
	generateGenericSequenceTuple("", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
	generateGenericSequenceTuple("Seq", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
	generateGenericSequenceTuple("Par", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
}

func generateContextReaderIOEitherTraverseTuple(f *os.File, i int) {
	generateGenericTraverseTuple("", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
	generateGenericTraverseTuple("Seq", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
	generateGenericTraverseTuple("Par", nonGenericContextReaderIOEither, extrasContextReaderIOEither)(f, i)
}

func generateContextReaderIOEitherHelpers(filename string, count int) error {
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

	writePackage(f, pkg)

	fmt.Fprintf(f, `
import (
	"context"

	G "github.com/IBM/fp-go/v2/context/%s/generic"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/tuple"
)
`, pkg)

	writePackage(fg, "generic")

	fmt.Fprintf(fg, `
import (
	"context"

	E "github.com/IBM/fp-go/v2/either"
	RE "github.com/IBM/fp-go/v2/readerioeither/generic"
)
`)

	generateContextReaderIOEitherEitherize(f, fg, 0)
	generateContextReaderIOEitherUneitherize(f, fg, 0)

	for i := 1; i <= count; i++ {
		// eitherize
		generateContextReaderIOEitherEitherize(f, fg, i)
		generateContextReaderIOEitherUneitherize(f, fg, i)
		// sequenceT
		generateContextReaderIOEitherSequenceT(f, i)
		// sequenceTuple
		generateContextReaderIOEitherSequenceTuple(f, i)
		// traverseTuple
		generateContextReaderIOEitherTraverseTuple(f, i)
	}

	return nil
}

func ContextReaderIOEitherCommand() *C.Command {
	return &C.Command{
		Name:  "contextreaderioeither",
		Usage: "generate code for ContextReaderIOEither",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generateContextReaderIOEitherHelpers(
				cmd.String(keyFilename),
				cmd.Int(keyCount),
			)
		},
	}
}
