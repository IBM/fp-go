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

func generateMakeProvider(f *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// MakeProvider%d creates a [DIE.Provider] for an [InjectionToken] from a function with %d dependencies\n", i, i)
	fmt.Fprintf(f, "func MakeProvider%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, " any, R any](\n")
	fmt.Fprintf(f, "  token InjectionToken[R],\n")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  d%d Dependency[T%d],\n", j+1, j+1)
	}
	fmt.Fprintf(f, "  f func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") IOE.IOEither[error, R],\n")
	fmt.Fprintf(f, ") DIE.Provider {\n")
	fmt.Fprint(f, "  return DIE.MakeProvider(\n")
	fmt.Fprint(f, "    token,\n")
	fmt.Fprintf(f, "    MakeProviderFactory%d(\n", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "      d%d,\n", j+1)
	}
	fmt.Fprint(f, "      f,\n")
	fmt.Fprint(f, "  ))\n")
	fmt.Fprintf(f, "}\n")
}

func generateMakeTokenWithDefault(f *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// MakeTokenWithDefault%d creates an [InjectionToken] with a default implementation with %d dependencies\n", i, i)
	fmt.Fprintf(f, "func MakeTokenWithDefault%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, " any, R any](\n")
	fmt.Fprintf(f, "  name string,\n")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  d%d Dependency[T%d],\n", j+1, j+1)
	}
	fmt.Fprintf(f, "  f func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") IOE.IOEither[error, R],\n")
	fmt.Fprintf(f, ") InjectionToken[R] {\n")
	fmt.Fprintf(f, "  return MakeTokenWithDefault[R](name, MakeProviderFactory%d(\n", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "    d%d,\n", j+1)
	}
	fmt.Fprint(f, "    f,\n")
	fmt.Fprint(f, "  ))\n")
	fmt.Fprintf(f, "}\n")
}

func generateMakeProviderFactory(f *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// MakeProviderFactory%d creates a [DIE.ProviderFactory] from a function with %d arguments and %d dependencies\n", i, i, i)
	fmt.Fprintf(f, "func MakeProviderFactory%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, " any, R any](\n")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  d%d Dependency[T%d],\n", j+1, j+1)
	}
	fmt.Fprintf(f, "  f func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") IOE.IOEither[error, R],\n")
	fmt.Fprintf(f, ") DIE.ProviderFactory {\n")
	fmt.Fprint(f, "  return DIE.MakeProviderFactory(\n")
	fmt.Fprint(f, "    A.From[DIE.Dependency](\n")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "      d%d,\n", j+1)
	}
	fmt.Fprint(f, "    ),\n")
	fmt.Fprintf(f, "    eraseProviderFactory%d(\n", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "      d%d,\n", j+1)
	}
	fmt.Fprint(f, "      f,\n")
	fmt.Fprint(f, "    ),\n")
	fmt.Fprint(f, "  )\n")
	fmt.Fprintf(f, "}\n")
}

func generateEraseProviderFactory(f *os.File, i int) {
	// non generic version
	fmt.Fprintf(f, "\n// eraseProviderFactory%d creates a function that takes a variadic number of untyped arguments and from a function of %d strongly typed arguments and %d dependencies\n", i, i, i)
	fmt.Fprintf(f, "func eraseProviderFactory%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, " any, R any](\n")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  d%d Dependency[T%d],\n", j+1, j+1)
	}
	fmt.Fprintf(f, "  f func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") IOE.IOEither[error, R]) func(params ...any) IOE.IOEither[error, any] {\n")
	fmt.Fprintf(f, "  ft := eraseTuple(T.Tupled%d(f))\n", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  t%d := lookupAt[T%d](%d, d%d)\n", j+1, j+1, j, j+1)
	}
	fmt.Fprint(f, "  return func(params ...any) IOE.IOEither[error, any] {\n")
	fmt.Fprintf(f, "    return ft(E.SequenceT%d(\n", i)
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "      t%d(params),\n", j+1)
	}
	fmt.Fprint(f, "    ))\n")
	fmt.Fprint(f, "  }\n")
	fmt.Fprintf(f, "}\n")
}

func generateDIHelpers(filename string, count int) error {
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
	// log
	log.Printf("Generating code in [%s] for package [%s] with [%d] repetitions ...", filename, pkg, count)

	// some header
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, "// This file was generated by robots at")
	fmt.Fprintf(f, "// %s\n\n", time.Now())

	fmt.Fprintf(f, "package %s\n\n", pkg)

	fmt.Fprint(f, `
import (
	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	T "github.com/IBM/fp-go/v2/tuple"
	A "github.com/IBM/fp-go/v2/array"
	DIE "github.com/IBM/fp-go/v2/di/erasure"
)
`)

	for i := 1; i <= count; i++ {
		generateEraseProviderFactory(f, i)
		generateMakeProviderFactory(f, i)
		generateMakeTokenWithDefault(f, i)
		generateMakeProvider(f, i)
	}

	return nil
}

func DICommand() *C.Command {
	return &C.Command{
		Name:  "di",
		Usage: "generate code for the dependency injection package",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generateDIHelpers(
				cmd.String(keyFilename),
				cmd.Int(keyCount),
			)
		},
	}
}
