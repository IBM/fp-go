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
	"time"

	C "github.com/urfave/cli/v3"
)

func writeTupleType(f *os.File, symbol string, i int) {
	fmt.Fprintf(f, "Tuple%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "%s%d", symbol, j)
	}
	fmt.Fprintf(f, "]")
}

func makeTupleType(name string) func(i int) string {
	return func(i int) string {
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("Tuple%d[", i))
		for j := 0; j < i; j++ {
			if j > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(fmt.Sprintf("%s%d", name, j+1))
		}
		buf.WriteString("]")

		return buf.String()
	}
}

func generatePush(f *os.File, i int) {
	tuple1 := makeTupleType("T")(i)
	tuple2 := makeTupleType("T")(i + 1)
	// Create the replicate version
	fmt.Fprintf(f, "\n// Push%d creates a [Tuple%d] from a [Tuple%d] by appending a constant value\n", i, i+1, i)
	fmt.Fprintf(f, "func Push%d[", i)
	// function prototypes
	for j := 0; j <= i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, " any](value T%d) func(%s) %s {\n", i+1, tuple1, tuple2)
	fmt.Fprintf(f, "  return func(t %s) %s {\n", tuple1, tuple2)
	fmt.Fprintf(f, "    return MakeTuple%d(", i+1)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t.F%d", j+1)
	}
	fmt.Fprintf(f, ", value)\n")
	fmt.Fprintf(f, "  }\n")
	fmt.Fprintf(f, "}\n")
}

func generateReplicate(f *os.File, i int) {
	// Create the replicate version
	fmt.Fprintf(f, "\n// Replicate%d creates a [Tuple%d] with all fields set to the input value `t`\n", i, i)
	fmt.Fprintf(f, "func Replicate%d[T any](t T) Tuple%d[", i, i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T")
	}
	fmt.Fprintf(f, "] {\n")
	// execute the mapping
	fmt.Fprintf(f, "  return MakeTuple%d(", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t")
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "}\n")
}

func generateMap(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// Map%d maps each value of a [Tuple%d] via a mapping function\n", i, i)
	fmt.Fprintf(f, "func Map%d[", i)
	// function prototypes
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "F%d ~func(T%d) R%d", j, j, j)
	}
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", T%d, R%d", j, j)
	}
	fmt.Fprintf(f, " any](")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "f%d F%d", j, j)
	}
	fmt.Fprintf(f, ") func(")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") ")
	writeTupleType(f, "R", i)
	fmt.Fprintf(f, " {\n")

	fmt.Fprintf(f, " return func(t ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") ")
	writeTupleType(f, "R", i)
	fmt.Fprintf(f, " {\n")

	// execute the mapping
	fmt.Fprintf(f, "    return MakeTuple%d(\n", i)
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "      f%d(t.F%d),\n", j, j)
	}
	fmt.Fprintf(f, "    )\n")

	fmt.Fprintf(f, " }\n")
	fmt.Fprintf(f, "}\n")
}

func generateMonoid(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// Monoid%d creates a [Monoid] for a [Tuple%d] based on %d monoids for the contained types\n", i, i, i)
	fmt.Fprintf(f, "func Monoid%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, " any](")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "m%d M.Monoid[T%d]", j, j)
	}
	fmt.Fprintf(f, ") M.Monoid[")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, "] {\n")

	fmt.Fprintf(f, "  return M.MakeMonoid(func(l, r ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, "{\n")

	fmt.Fprintf(f, "    return MakeTuple%d(", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "m%d.Concat(l.F%d, r.F%d)", j, j, j)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "  }, MakeTuple%d(", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "m%d.Empty()", j)
	}
	fmt.Fprintf(f, "))\n")

	fmt.Fprintf(f, "}\n")
}

func generateOrd(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// Ord%d creates n [Ord] for a [Tuple%d] based on %d [Ord]s for the contained types\n", i, i, i)
	fmt.Fprintf(f, "func Ord%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, " any](")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "o%d O.Ord[T%d]", j, j)
	}
	fmt.Fprintf(f, ") O.Ord[")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, "] {\n")

	fmt.Fprintf(f, "  return O.MakeOrd(func(l, r ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") int {\n")

	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "    if c:= o%d.Compare(l.F%d, r.F%d); c != 0 {return c}\n", j, j, j)
	}
	fmt.Fprintf(f, "    return 0\n")
	fmt.Fprintf(f, "  }, func(l, r ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") bool {\n")
	fmt.Fprintf(f, "    return ")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, " && ")
		}
		fmt.Fprintf(f, "o%d.Equals(l.F%d, r.F%d)", j, j, j)
	}
	fmt.Fprintf(f, "\n")
	fmt.Fprintf(f, "  })\n")

	fmt.Fprintf(f, "}\n")
}

func generateTupleType(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// Tuple%d is a struct that carries %d independently typed values\n", i, i)
	fmt.Fprintf(f, "type Tuple%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, " any] struct {\n")
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "  F%d T%d\n", j, j)
	}
	fmt.Fprintf(f, "}\n")
}

func generateMakeTupleType(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// MakeTuple%d is a function that converts its %d parameters into a [Tuple%d]\n", i, i, i)
	fmt.Fprintf(f, "func MakeTuple%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, " any](")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d T%d", j, j)
	}
	fmt.Fprintf(f, ") ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, " {\n")
	fmt.Fprintf(f, "  return Tuple%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, "]{")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d", j)
	}
	fmt.Fprintf(f, "}\n")
	fmt.Fprintf(f, "}\n")
}

func generateUntupled(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// Untupled%d converts a function with a [Tuple%d] parameter into a function with %d parameters\n// The inverse function is [Tupled%d]\n", i, i, i, i)
	fmt.Fprintf(f, "func Untupled%d[F ~func(Tuple%d[", i, i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "]) R")
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
	fmt.Fprintf(f, ") R {\n")
	fmt.Fprintf(f, "  return func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d T%d", j+1, j+1)
	}
	fmt.Fprintf(f, ") R {\n")
	fmt.Fprintf(f, "    return f(MakeTuple%d(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d", j+1)
	}
	fmt.Fprintln(f, "))")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "}")
}

func generateTupled(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// Tupled%d converts a function with %d parameters into a function taking a Tuple%d\n// The inverse function is [Untupled%d]\n", i, i, i, i)
	fmt.Fprintf(f, "func Tupled%d[F ~func(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, ") R")
	for j := 0; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j+1)
	}
	fmt.Fprintf(f, ", R any](f F) func(Tuple%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "]) R {\n")
	fmt.Fprintf(f, "  return func(t Tuple%d[", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j+1)
	}
	fmt.Fprintf(f, "]) R {\n")
	fmt.Fprintf(f, "    return f(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t.F%d", j+1)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "  }\n")
	fmt.Fprintln(f, "}")
}

func generateTupleHelpers(filename string, count int) error {
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

	fmt.Fprintf(f, `
import (
	M "github.com/IBM/fp-go/v2/monoid"
	O "github.com/IBM/fp-go/v2/ord"	
)
`)

	for i := 1; i <= count; i++ {
		// tuple type
		generateTupleType(f, i)
	}

	for i := 1; i <= count; i++ {
		// tuple generator
		generateMakeTupleType(f, i)
		// tupled wrapper
		generateTupled(f, i)
		// untupled wrapper
		generateUntupled(f, i)
		// monoid
		generateMonoid(f, i)
		// generate order
		generateOrd(f, i)
		// generate map
		generateMap(f, i)
		// generate replicate
		generateReplicate(f, i)
		// generate tuple functions such as string and fmt
		generateTupleString(f, i)
		// generate json support
		generateTupleMarshal(f, i)
		// generate json support
		generateTupleUnmarshal(f, i)
		// generate toArray
		generateToArray(f, i)
		// generate fromArray
		generateFromArray(f, i)
		// generate push
		if i < count {
			generatePush(f, i)
		}
	}

	return nil
}

func generateTupleMarshal(f *os.File, i int) {
	// Create the stringify version
	fmt.Fprintf(f, "\n// MarshalJSON marshals the [Tuple%d] into a JSON array\n", i)
	fmt.Fprintf(f, "func (t ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") MarshalJSON() ([]byte, error) {\n")
	fmt.Fprintf(f, "  return tupleMarshalJSON(")
	// function prototypes
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t.F%d", j)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "}\n")
}

func generateTupleUnmarshal(f *os.File, i int) {
	// Create the stringify version
	fmt.Fprintf(f, "\n// UnmarshalJSON unmarshals a JSON array into a [Tuple%d]\n", i)
	fmt.Fprintf(f, "func (t *")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") UnmarshalJSON(data []byte) error {\n")
	fmt.Fprintf(f, "  return tupleUnmarshalJSON(data")
	// function prototypes
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", &t.F%d", j)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "}\n")
}

func generateToArray(f *os.File, i int) {
	// Create the stringify version
	fmt.Fprintf(f, "\n// ToArray converts the [Tuple%d] into an array of type [R] using %d transformation functions from [T] to [R]\n// The inverse function is [FromArray%d]\n", i, i, i)
	fmt.Fprintf(f, "func ToArray%d[", i)
	// function prototypes
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "F%d ~func(T%d) R", j, j)
	}
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ", R any](")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "f%d F%d", j, j)
	}
	fmt.Fprintf(f, ") func(t ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") []R {\n")
	fmt.Fprintf(f, "  return func(t ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") []R {\n")
	fmt.Fprintf(f, "    return []R{\n")
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "      f%d(t.F%d),\n", j, j)
	}
	fmt.Fprintf(f, "    }\n")
	fmt.Fprintf(f, "  }\n")
	fmt.Fprintf(f, "}\n")
}

func generateFromArray(f *os.File, i int) {
	// Create the stringify version
	fmt.Fprintf(f, "\n// FromArray converts an array of [R] into a [Tuple%d] using %d functions from [R] to [T]\n// The inverse function is [ToArray%d]\n", i, i, i)
	fmt.Fprintf(f, "func FromArray%d[", i)
	// function prototypes
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "F%d ~func(R) T%d", j, j)
	}
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ", R any](")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "f%d F%d", j, j)
	}
	fmt.Fprintf(f, ") func(r []R) ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, " {\n")
	fmt.Fprintf(f, "  return func(r []R) ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, " {\n")
	fmt.Fprintf(f, "    return MakeTuple%d(\n", i)
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "      f%d(r[%d]),\n", j, j-1)
	}
	fmt.Fprintf(f, "    )\n")
	fmt.Fprintf(f, "  }\n")
	fmt.Fprintf(f, "}\n")
}

func generateTupleString(f *os.File, i int) {
	// Create the stringify version
	fmt.Fprintf(f, "\n// String prints some debug info for the [Tuple%d]\n", i)
	fmt.Fprintf(f, "func (t ")
	writeTupleType(f, "T", i)
	fmt.Fprintf(f, ") String() string {\n")
	// convert to string
	fmt.Fprint(f, "  return tupleString(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t.F%d", j)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintf(f, "}\n")
}

// func generateTupleJson(f *os.File, i int) {
// 	// Create the stringify version
// 	fmt.Fprintf(f, "\n// MarshalJSON converts the [Tuple%d] into a JSON byte stream\n", i)
// 	fmt.Fprintf(f, "func (t ")
// 	writeTupleType(f, "T", i)
// 	fmt.Fprintf(f, ") MarshalJSON() ([]byte, error) {\n")
// 	// convert to string
// 	fmt.Fprintf(f, "  return fmt.Sprintf(\"Tuple%d[", i)
// 	for j := 1; j <= i; j++ {
// 		if j > 1 {
// 			fmt.Fprintf(f, ", ")
// 		}
// 		fmt.Fprintf(f, "%s", "%T")
// 	}
// 	fmt.Fprintf(f, "](")
// 	for j := 1; j <= i; j++ {
// 		if j > 1 {
// 			fmt.Fprintf(f, ", ")
// 		}
// 		fmt.Fprintf(f, "%s", "%v")
// 	}
// 	fmt.Fprintf(f, ")\", ")
// 	for j := 1; j <= i; j++ {
// 		if j > 1 {
// 			fmt.Fprintf(f, ", ")
// 		}
// 		fmt.Fprintf(f, "t.F%d", j)
// 	}
// 	for j := 1; j <= i; j++ {
// 		fmt.Fprintf(f, ", t.F%d", j)
// 	}
// 	fmt.Fprintf(f, ")\n")
// 	fmt.Fprintf(f, "}\n")
// }

func TupleCommand() *C.Command {
	return &C.Command{
		Name:  "tuple",
		Usage: "generate code for Tuple",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generateTupleHelpers(
				cmd.String(keyFilename),
				cmd.Int(keyCount),
			)
		},
	}
}
