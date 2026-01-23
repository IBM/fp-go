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

func generateTraverseTuple(f *os.File, i int) {
	fmt.Fprintf(f, "\n// TraverseTuple%d is a utility function used to implement the sequence operation for higher kinded types based only on map and ap.\n", i)
	fmt.Fprintf(f, "// The function takes a [Tuple%d] of base types and %d functions that transform these based types into higher higher kinded types. It returns a higher kinded type of a [Tuple%d] with the resolved values.\n", i, i, i)
	fmt.Fprintf(f, "func TraverseTuple%d[\n", i)
	// map as the starting point
	fmt.Fprintf(f, "  MAP ~func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, " ")
		}
		fmt.Fprintf(f, "func(T%d)", j+1)
	}
	fmt.Fprintf(f, " ")
	fmt.Fprintf(f, "T.")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") func(HKT_T1)")
	if i > 1 {
		fmt.Fprintf(f, " HKT_F")
		for k := 1; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
	} else {
		fmt.Fprintf(f, " HKT_TUPLE%d", i)
	}
	fmt.Fprintf(f, ",\n")
	// the applicatives
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  AP%d ~func(", j)
		fmt.Fprintf(f, "HKT_T%d) func(", j+1)
		fmt.Fprintf(f, "HKT_F")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
		fmt.Fprintf(f, ")")
		if j+1 < i {
			fmt.Fprintf(f, " HKT_F")
			for k := j + 1; k < i; k++ {
				fmt.Fprintf(f, "_T%d", k+1)
			}
		} else {
			fmt.Fprintf(f, " HKT_TUPLE%d", i)
		}
		fmt.Fprintf(f, ",\n")
	}
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  F%d ~func(A%d) HKT_T%d,\n", j+1, j+1, j+1)
	}
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  A%d, T%d,\n", j+1, j+1)
	}
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  HKT_T%d, // HKT[T%d]\n", j+1, j+1)
	}
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  HKT_F")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
		fmt.Fprintf(f, ", // HKT[")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "func(T%d) ", k+1)
		}
		fmt.Fprintf(f, "T.")
		writeTupleType(f, "T", i)
		fmt.Fprintf(f, "]\n")
	}
	fmt.Fprintf(f, "  HKT_TUPLE%d any, // HKT[", i)
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, "]\n")
	fmt.Fprintf(f, "](\n")

	// the callbacks
	fmt.Fprintf(f, "  fmap MAP,\n")
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  fap%d AP%d,\n", j, j)
	}
	// the transformer functions
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "  f%d F%d,\n", j, j)
	}
	// the parameters
	fmt.Fprintf(f, "  t T.Tuple%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "A%d", j+1)
	}
	fmt.Fprintf(f, "],\n")
	fmt.Fprintf(f, ") HKT_TUPLE%d {\n", i)

	fmt.Fprintf(f, "  return F.Pipe%d(\n", i)
	fmt.Fprintf(f, "    f1(t.F1),\n")
	fmt.Fprintf(f, "    fmap(tupleConstructor%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "]()),\n")
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "    fap%d(f%d(t.F%d)),\n", j, j+1, j+1)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "}\n")
}

func generateSequenceTuple(f *os.File, i int) {
	fmt.Fprintf(f, "\n// SequenceTuple%d is a utility function used to implement the sequence operation for higher kinded types based only on map and ap.\n", i)
	fmt.Fprintf(f, "// The function takes a [Tuple%d] of higher higher kinded types and returns a higher kinded type of a [Tuple%d] with the resolved values.\n", i, i)
	fmt.Fprintf(f, "func SequenceTuple%d[\n", i)
	// map as the starting point
	fmt.Fprintf(f, "  MAP ~func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, " ")
		}
		fmt.Fprintf(f, "func(T%d)", j+1)
	}
	fmt.Fprintf(f, " ")
	fmt.Fprintf(f, "T.")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") func(HKT_T1)")
	if i > 1 {
		fmt.Fprintf(f, " HKT_F")
		for k := 1; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
	} else {
		fmt.Fprintf(f, " HKT_TUPLE%d", i)
	}
	fmt.Fprintf(f, ",\n")
	// the applicatives
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  AP%d ~func(", j)
		fmt.Fprintf(f, "HKT_T%d) func(", j+1)
		fmt.Fprintf(f, "HKT_F")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
		fmt.Fprintf(f, ")")
		if j+1 < i {
			fmt.Fprintf(f, " HKT_F")
			for k := j + 1; k < i; k++ {
				fmt.Fprintf(f, "_T%d", k+1)
			}
		} else {
			fmt.Fprintf(f, " HKT_TUPLE%d", i)
		}
		fmt.Fprintf(f, ",\n")
	}

	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  T%d,\n", j+1)
	}
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  HKT_T%d, // HKT[T%d]\n", j+1, j+1)
	}
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  HKT_F")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
		fmt.Fprintf(f, ", // HKT[")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "func(T%d) ", k+1)
		}
		fmt.Fprintf(f, "T.")
		writeTupleType(f, "T", i)
		fmt.Fprintf(f, "]\n")
	}
	fmt.Fprintf(f, "  HKT_TUPLE%d any, // HKT[", i)
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, "]\n")
	fmt.Fprintf(f, "](\n")

	// the callbacks
	fmt.Fprintf(f, "  fmap MAP,\n")
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  fap%d AP%d,\n", j, j)
	}
	// the parameters
	fmt.Fprintf(f, "  t T.Tuple%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "HKT_T%d", j+1)
	}
	fmt.Fprintf(f, "],\n")
	fmt.Fprintf(f, ") HKT_TUPLE%d {\n", i)

	fmt.Fprintf(f, "  return F.Pipe%d(\n", i)
	fmt.Fprintf(f, "    t.F1,\n")
	fmt.Fprintf(f, "    fmap(tupleConstructor%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "]()),\n")
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "    fap%d(t.F%d),\n", j, j+1)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "}\n")
}

func generateSequenceT(f *os.File, i int) {
	fmt.Fprintf(f, "\n// SequenceT%d is a utility function used to implement the sequence operation for higher kinded types based only on map and ap.\n", i)
	fmt.Fprintf(f, "// The function takes %d higher higher kinded types and returns a higher kinded type of a [Tuple%d] with the resolved values.\n", i, i)
	fmt.Fprintf(f, "func SequenceT%d[\n", i)
	// map as the starting point
	fmt.Fprintf(f, "  MAP ~func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, " ")
		}
		fmt.Fprintf(f, "func(T%d)", j+1)
	}
	fmt.Fprintf(f, " ")
	fmt.Fprintf(f, "T.")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") func(HKT_T1)")
	if i > 1 {
		fmt.Fprintf(f, " HKT_F")
		for k := 1; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
	} else {
		fmt.Fprintf(f, " HKT_TUPLE%d", i)
	}
	fmt.Fprintf(f, ",\n")
	// the applicatives
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  AP%d ~func(", j)
		fmt.Fprintf(f, "HKT_T%d) func(", j+1)
		fmt.Fprintf(f, "HKT_F")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
		fmt.Fprintf(f, ")")
		if j+1 < i {
			fmt.Fprintf(f, " HKT_F")
			for k := j + 1; k < i; k++ {
				fmt.Fprintf(f, "_T%d", k+1)
			}
		} else {
			fmt.Fprintf(f, " HKT_TUPLE%d", i)
		}
		fmt.Fprintf(f, ",\n")
	}

	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  T%d,\n", j+1)
	}
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  HKT_T%d, // HKT[T%d]\n", j+1, j+1)
	}
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  HKT_F")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "_T%d", k+1)
		}
		fmt.Fprintf(f, ", // HKT[")
		for k := j; k < i; k++ {
			fmt.Fprintf(f, "func(T%d) ", k+1)
		}
		fmt.Fprintf(f, "T.")
		writeTupleType(f, "T", i)
		fmt.Fprintf(f, "]\n")
	}
	fmt.Fprintf(f, "  HKT_TUPLE%d any, // HKT[", i)
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, "]\n")
	fmt.Fprintf(f, "](\n")

	// the callbacks
	fmt.Fprintf(f, "  fmap MAP,\n")
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "  fap%d AP%d,\n", j, j)
	}
	// the parameters
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, "  t%d HKT_T%d,\n", j+1, j+1)
	}
	fmt.Fprintf(f, ") HKT_TUPLE%d {\n", i)

	fmt.Fprintf(f, "  return F.Pipe%d(\n", i)
	fmt.Fprintf(f, "    t1,\n")
	fmt.Fprintf(f, "    fmap(tupleConstructor%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "]()),\n")
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, "    fap%d(t%d),\n", j, j+1)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "}\n")
}

func generateTupleConstructor(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// tupleConstructor%d returns a curried version of [T.MakeTuple%d]\n", i, i)
	fmt.Fprintf(f, "func tupleConstructor%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, " any]()")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, " func(T%d)", j+1)
	}
	fmt.Fprintf(f, " T.Tuple%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "] {\n")

	fmt.Fprintf(f, "  return F.Curry%d(T.MakeTuple%d[", i, i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "])\n")

	fmt.Fprintf(f, "}\n")
}

func generateApplyHelpers(filename string, count int) error {
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

	// print out some helpers
	fmt.Fprintf(f, `
import (
	F "github.com/IBM/fp-go/v2/function"
	T "github.com/IBM/fp-go/v2/tuple"
)
`)

	for i := 1; i <= count; i++ {
		// tuple constructor
		generateTupleConstructor(f, i)
		// sequenceT
		generateSequenceT(f, i)
		// sequenceTuple
		generateSequenceTuple(f, i)
		// traverseTuple
		generateTraverseTuple(f, i)
	}

	return nil
}

func ApplyCommand() *C.Command {
	return &C.Command{
		Name:  "apply",
		Usage: "generate code for the sequence operations of apply",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generateApplyHelpers(
				cmd.String(keyFilename),
				cmd.Int(keyCount),
			)
		},
	}
}
