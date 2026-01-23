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

func createCombinations(n int, all, prev []int) [][]int {
	l := len(prev)
	if l == n {
		return [][]int{prev}
	}
	var res [][]int
	for idx, val := range all {
		cpy := make([]int, l+1)
		copy(cpy, prev)
		cpy[l] = val

		res = append(res, createCombinations(n, all[idx+1:], cpy)...)
	}
	return res
}

func remaining(comb []int, total int) []int {
	var res []int
	mp := make(map[int]int)
	for _, idx := range comb {
		mp[idx] = idx
	}
	for i := 1; i <= total; i++ {
		_, ok := mp[i]
		if !ok {
			res = append(res, i)
		}
	}
	return res
}

func generateCombSingleBind(f *os.File, comb [][]int, total int) {
	for _, c := range comb {
		// remaining indexes
		rem := remaining(c, total)

		// bind function
		fmt.Fprintf(f, "\n// Bind")
		for _, idx := range c {
			fmt.Fprintf(f, "%d", idx)
		}
		fmt.Fprintf(f, "of%d takes a function with %d parameters and returns a new function with %d parameters that will bind these parameters to the positions [", total, total, len(c))
		for i, idx := range c {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "%d", idx)
		}
		fmt.Fprintf(f, "] of the original function.\n// The return value of is a function with the remaining %d parameters at positions [", len(rem))
		for i, idx := range rem {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "%d", idx)
		}
		fmt.Fprintf(f, "] of the original function.\n")
		fmt.Fprintf(f, "func Bind")
		for _, idx := range c {
			fmt.Fprintf(f, "%d", idx)
		}
		fmt.Fprintf(f, "of%d[F ~func(", total)
		for i := 0; i < total; i++ {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", i+1)
		}
		fmt.Fprintf(f, ") R")
		for i := 0; i < total; i++ {
			fmt.Fprintf(f, ", T%d", i+1)
		}
		fmt.Fprintf(f, ", R any](f F) func(")
		for i, idx := range c {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", idx)
		}
		fmt.Fprintf(f, ") func(")
		for i, idx := range rem {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", idx)
		}
		fmt.Fprintf(f, ") R {\n")

		fmt.Fprintf(f, "  return func(")

		for i, idx := range c {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d T%d", idx, idx)
		}
		fmt.Fprintf(f, ") func(")
		for i, idx := range rem {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", idx)
		}
		fmt.Fprintf(f, ") R {\n")

		fmt.Fprintf(f, "    return func(")
		for i, idx := range rem {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d T%d", idx, idx)
		}
		fmt.Fprintf(f, ") R {\n")

		fmt.Fprintf(f, "      return f(")
		for i := 1; i <= total; i++ {
			if i > 1 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d", i)
		}
		fmt.Fprintf(f, ")\n")

		fmt.Fprintf(f, "    }\n")
		fmt.Fprintf(f, "  }\n")
		fmt.Fprintf(f, "}\n")

		// ignore function
		fmt.Fprintf(f, "\n// Ignore")
		for _, idx := range c {
			fmt.Fprintf(f, "%d", idx)
		}
		fmt.Fprintf(f, "of%d takes a function with %d parameters and returns a new function with %d parameters that will ignore the values at positions [", total, len(rem), total)
		for i, idx := range c {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "%d", idx)
		}
		fmt.Fprintf(f, "] and pass the remaining %d parameters to the original function\n", len(rem))
		fmt.Fprintf(f, "func Ignore")
		for _, idx := range c {
			fmt.Fprintf(f, "%d", idx)
		}
		fmt.Fprintf(f, "of%d[", total)
		// start with the undefined parameters
		for i, idx := range c {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", idx)
		}
		if len(c) > 0 {
			fmt.Fprintf(f, " any, ")
		}
		fmt.Fprintf(f, "F ~func(")
		for i, idx := range rem {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", idx)
		}
		fmt.Fprintf(f, ") R")
		for _, idx := range rem {
			fmt.Fprintf(f, ", T%d", idx)
		}
		fmt.Fprintf(f, ", R any](f F) func(")
		for i := 0; i < total; i++ {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", i+1)
		}
		fmt.Fprintf(f, ") R {\n")

		fmt.Fprintf(f, "  return func(")
		for i := 0; i < total; i++ {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d T%d", i+1, i+1)
		}
		fmt.Fprintf(f, ") R {\n")
		fmt.Fprintf(f, "      return f(")
		for i, idx := range rem {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d", idx)
		}
		fmt.Fprintf(f, ")\n")

		fmt.Fprintf(f, "  }\n")
		fmt.Fprintf(f, "}\n")

	}
}

func generateSingleBind(f *os.File, total int) {

	fmt.Fprintf(f, "// Combinations for a total of %d arguments\n", total)

	// construct the indexes
	all := make([]int, total)
	for i := 0; i < total; i++ {
		all[i] = i + 1
	}
	// for all permutations of a certain length
	for j := 0; j < total; j++ {
		// get combinations of that size
		comb := createCombinations(j+1, all, []int{})
		generateCombSingleBind(f, comb, total)
	}
}

func generateBind(f *os.File, i int) {
	for j := 1; j < i; j++ {
		generateSingleBind(f, j)
	}
}

func generateBindHelpers(filename string, count int) error {
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

	generateBind(f, count)

	return nil
}

func BindCommand() *C.Command {
	return &C.Command{
		Name:        "bind",
		Usage:       "generate code for binder functions etc",
		Description: "Code generation for bind, etc",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx context.Context, cmd *C.Command) error {
			return generateBindHelpers(
				cmd.String(keyFilename),
				cmd.Int(keyCount),
			)
		},
	}
}
