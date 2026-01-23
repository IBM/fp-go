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

func generateUnsliced(f *os.File, i int) {
	// Create the optionize version
	fmt.Fprintf(f, "\n// Unsliced%d converts a function taking a slice parameter into a function with %d parameters\n", i, i)
	fmt.Fprintf(f, "func Unsliced%d[F ~func([]T) R, T, R any](f F) func(", i)
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T")
	}
	fmt.Fprintf(f, ") R {\n")
	fmt.Fprintf(f, "  return func(")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d", j+1)
	}
	if i > 0 {
		fmt.Fprintf(f, " T")
	}
	fmt.Fprintf(f, ") R {\n")
	fmt.Fprintf(f, "    return f([]T{")
	for j := 0; j < i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d", j+1)
	}
	fmt.Fprintln(f, "})")
	fmt.Fprintln(f, "  }")
	fmt.Fprintln(f, "}")
}

func generateVariadic(f *os.File, i int) {
	// Create the nullary version
	fmt.Fprintf(f, "\n// Variadic%d converts a function taking %d parameters and a final slice into a function with %d parameters but a final variadic argument\n", i, i, i)
	fmt.Fprintf(f, "func Variadic%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "V, R any](f func(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "[]V) R) func(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "...V) R {\n")
	fmt.Fprintf(f, "  return func(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d T%d", j, j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "v ...V) R {\n")
	fmt.Fprintf(f, "    return f(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d", j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "v)\n")
	fmt.Fprintf(f, "  }\n")

	fmt.Fprintf(f, "}\n")
}

func generateUnvariadic(f *os.File, i int) {
	// Create the nullary version
	fmt.Fprintf(f, "\n// Unvariadic%d converts a function taking %d parameters and a final variadic argument into a function with %d parameters but a final slice argument\n", i, i, i)
	fmt.Fprintf(f, "func Unvariadic%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "V, R any](f func(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "...V) R) func(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "[]V) R {\n")
	fmt.Fprintf(f, "  return func(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d T%d", j, j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "v []V) R {\n")
	fmt.Fprintf(f, "    return f(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d", j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	fmt.Fprintf(f, "v...)\n")
	fmt.Fprintf(f, "  }\n")

	fmt.Fprintf(f, "}\n")
}

func generateNullary(f *os.File, i int) {
	// Create the nullary version
	fmt.Fprintf(f, "\n// Nullary%d creates a parameter less function from a parameter less function and %d functions. When executed the first parameter less function gets executed and then the result is piped through the remaining functions\n", i, i-1)
	fmt.Fprintf(f, "func Nullary%d[F1 ~func() T1", i)
	for j := 2; j <= i; j++ {
		fmt.Fprintf(f, ", F%d ~func(T%d) T%d", j, j-1, j)
	}
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, " any](f1 F1")
	for j := 2; j <= i; j++ {
		fmt.Fprintf(f, ", f%d F%d", j, j)
	}
	fmt.Fprintf(f, ") func() T%d {\n", i)
	fmt.Fprintf(f, "  return func() T%d {\n", i)
	fmt.Fprintf(f, "    return Pipe%d(f1()", i-1)
	for j := 2; j <= i; j++ {
		fmt.Fprintf(f, ", f%d", j)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintln(f, "  }")

	fmt.Fprintln(f, "}")
}

func generateFlow(f *os.File, i int) {
	// Create the flow version
	fmt.Fprintf(f, "\n// Flow%d creates a function that takes an initial value t0 and successively applies %d functions where the input of a function is the return value of the previous function\n// The final return value is the result of the last function application\n", i, i)
	fmt.Fprintf(f, "func Flow%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "F%d ~func(T%d) T%d", j, j-1, j)
	}
	for j := 0; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, " any](")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "f%d F%d", j, j)
	}
	fmt.Fprintf(f, ") func(T0) T%d {\n", i)
	fmt.Fprintf(f, "  return func(t0 T0) T%d {\n", i)
	fmt.Fprintf(f, "    return Pipe%d(t0", i)
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", f%d", j)
	}
	fmt.Fprintf(f, ")\n")
	fmt.Fprintln(f, "  }")

	fmt.Fprintln(f, "}")

}

func generatePipe(f *os.File, i int) {
	// Create the pipe version
	fmt.Fprintf(f, "\n// Pipe%d takes an initial value t0 and successively applies %d functions where the input of a function is the return value of the previous function\n// The final return value is the result of the last function application\n", i, i)
	fmt.Fprintf(f, "func Pipe%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "F%d ~func(T%d) T%d", j, j-1, j)
	}
	if i > 0 {
		fmt.Fprintf(f, ", ")
	}
	for j := 0; j <= i; j++ {
		if j > 0 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}

	fmt.Fprintf(f, " any](t0 T0")
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", f%d F%d", j, j)
	}
	fmt.Fprintf(f, ") T%d {\n", i)
	fmt.Fprintf(f, "  return ")
	for j := i; j >= 1; j-- {
		fmt.Fprintf(f, "f%d(", j)
	}
	fmt.Fprintf(f, "t0")
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ")")
	}
	fmt.Fprintf(f, "\n")
	fmt.Fprintln(f, "}")
}

func recurseCurry(f *os.File, indent string, total, count int) {
	if count == 1 {
		fmt.Fprintf(f, "%sreturn func(t%d T%d) T%d {\n", indent, total-1, total-1, total)
		fmt.Fprintf(f, "%s  return f(t0", indent)
		for i := 1; i < total; i++ {
			fmt.Fprintf(f, ", t%d", i)
		}
		fmt.Fprintf(f, ")\n")
		fmt.Fprintf(f, "%s}\n", indent)
	} else {
		fmt.Fprintf(f, "%sreturn", indent)
		for i := total - count + 1; i <= total; i++ {
			fmt.Fprintf(f, " func(t%d T%d)", i-1, i-1)
		}
		fmt.Fprintf(f, " T%d {\n", total)
		recurseCurry(f, fmt.Sprintf("  %s", indent), total, count-1)
		fmt.Fprintf(f, "%s}\n", indent)
	}
}

func generateCurry(f *os.File, i int) {
	// Create the curry version
	fmt.Fprintf(f, "\n// Curry%d takes a function with %d parameters and returns a cascade of functions each taking only one parameter.\n// The inverse function is [Uncurry%d]\n", i, i, i)
	fmt.Fprintf(f, "func Curry%d[FCT ~func(T0", i)
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, ") T%d", i)
	// type arguments
	for j := 0; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, " any](f FCT) func(T0)")
	for j := 2; j <= i; j++ {
		fmt.Fprintf(f, " func(T%d)", j-1)
	}
	fmt.Fprintf(f, " T%d {\n", i)
	recurseCurry(f, "  ", i, i)
	fmt.Fprintf(f, "}\n")
}

func generateUncurry(f *os.File, i int) {
	// Create the uncurry version
	fmt.Fprintf(f, "\n// Uncurry%d takes a cascade of %d functions each taking only one parameter and returns a function with %d parameters .\n// The inverse function is [Curry%d]\n", i, i, i, i)
	fmt.Fprintf(f, "func Uncurry%d[FCT ~func(T0)", i)
	for j := 1; j < i; j++ {
		fmt.Fprintf(f, " func(T%d)", j)
	}
	fmt.Fprintf(f, " T%d", i)
	// the type parameters
	for j := 0; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, " any](f FCT) func(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j-1)
	}
	fmt.Fprintf(f, ") T%d {\n", i)
	fmt.Fprintf(f, "  return func(")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "t%d T%d", j-1, j-1)
	}
	fmt.Fprintf(f, ") T%d {\n", i)
	fmt.Fprintf(f, "    return f")
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "(t%d)", j-1)
	}
	fmt.Fprintln(f)

	fmt.Fprintf(f, "  }\n")

	fmt.Fprintf(f, "}\n")
}

func generatePipeHelpers(filename string, count int) error {
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

	fmt.Fprintf(f, "package %s\n", pkg)

	// pipe
	generatePipe(f, 0)
	// variadic
	generateVariadic(f, 0)
	// unvariadic
	generateUnvariadic(f, 0)
	// unsliced
	generateUnsliced(f, 0)

	for i := 1; i <= count; i++ {

		// pipe
		generatePipe(f, i)
		// flow
		generateFlow(f, i)
		// nullary
		generateNullary(f, i)
		// curry
		generateCurry(f, i)
		// uncurry
		generateUncurry(f, i)
		// variadic
		generateVariadic(f, i)
		// unvariadic
		generateUnvariadic(f, i)
		// unsliced
		generateUnsliced(f, i)
	}

	return nil
}

func PipeCommand() *C.Command {
	return &C.Command{
		Name:        "pipe",
		Usage:       "generate code for pipe, flow, curry, etc",
		Description: "Code generation for pipe, flow, curry, etc",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generatePipeHelpers(
				cmd.String(keyFilename),
				cmd.Int(keyCount),
			)
		},
	}
}
