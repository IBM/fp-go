// Copyright (c) 2023 IBM Corp.
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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	C "github.com/urfave/cli/v2"
)

func generateReaderFrom(f, fg *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// From%d converts a function with %d parameters returning a [R] into a function with %d parameters returning a [Reader[C, R]]\n// The first parameter is considered to be the context [C] of the reader\n", i, i+1, i)
	fmt.Fprintf(f, "func From%d[F ~func(C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ") R")
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
	fmt.Fprintf(f, ") Reader[C, R] {\n")
	fmt.Fprintf(f, "  return G.From%d[Reader[C, R]](f)\n", i)
	fmt.Fprintln(f, "}")

	// generic version
	fmt.Fprintf(fg, "\n// From%d converts a function with %d parameters returning a [R] into a function with %d parameters returning a [GRA]\n// The first parameter is considered to be the context [C].\n", i, i+1, i)
	fmt.Fprintf(fg, "func From%d[GRA ~func(C) R, F ~func(C", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", T%d", j)
	}
	fmt.Fprintf(fg, ") R")
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

	fmt.Fprintf(fg, "  return func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(fg, ", ")
		}
		fmt.Fprintf(fg, "t%d T%d", j, j)
	}
	fmt.Fprintf(fg, ") GRA {\n")
	fmt.Fprintf(fg, "    return MakeReader[GRA](func(r C) R {\n")
	fmt.Fprintf(fg, "      return f(r")
	for j := 0; j < i; j++ {
		fmt.Fprintf(fg, ", t%d", j)
	}
	fmt.Fprintf(fg, ")\n")
	fmt.Fprintf(fg, "    })\n")
	fmt.Fprintf(fg, "  }\n")
	fmt.Fprintf(fg, "}\n")
}

func generateReaderHelpers(filename string, count int) error {
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
	G "github.com/IBM/fp-go/%s/generic"	
)
`, pkg)

	// some header
	fmt.Fprintln(fg, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(fg, "// This file was generated by robots at")
	fmt.Fprintf(fg, "// %s\n", time.Now())

	fmt.Fprintf(fg, "package generic\n\n")

	// from
	generateReaderFrom(f, fg, 0)

	for i := 1; i <= count; i++ {
		// from
		generateReaderFrom(f, fg, i)
	}

	return nil
}

func ReaderCommand() *C.Command {
	return &C.Command{
		Name:  "reader",
		Usage: "generate code for Reader",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx *C.Context) error {
			return generateReaderHelpers(
				ctx.String(keyFilename),
				ctx.Int(keyCount),
			)
		},
	}
}