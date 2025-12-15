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
	"fmt"
	"os"
	"strings"

	S "github.com/IBM/fp-go/v2/string"
)

// Deprecated:
func tupleType(name string) func(i int) string {
	return func(i int) string {
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("T.Tuple%d[", i))
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

func tupleTypePlain(name string) func(i int) string {
	return func(i int) string {
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("tuple.Tuple%d[", i))
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

func monadGenerateSequenceTNonGeneric(
	hkt func(string) string,
	fmap func(string, string) string,
	fap func(string, string) string,
) func(f *os.File, i int) {
	return func(f *os.File, i int) {

		tuple := tupleType("T")(i)

		fmt.Fprintf(f, "SequenceT%d[", i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", j+1)
		}
		fmt.Fprintf(f, "](")
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d %s", j+1, hkt(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(f, ") %s {", hkt(tuple))
		// the actual apply callback
		fmt.Fprintf(f, "  return apply.SequenceT%d(\n", i)
		// map callback

		curried := func(count int) string {
			var buf strings.Builder
			for j := count; j < i; j++ {
				buf.WriteString(fmt.Sprintf("func(T%d)", j+1))
			}
			buf.WriteString(tuple)
			return buf.String()
		}

		fmt.Fprintf(f, "    %s,\n", fmap("T1", curried(1)))
		for j := 1; j < i; j++ {
			fmt.Fprintf(f, "    %s,\n", fap(curried(j+1), fmt.Sprintf("T%d", j)))
		}

		for j := 0; j < i; j++ {
			fmt.Fprintf(f, "    T%d,\n", j+1)
		}

		fmt.Fprintf(f, "  )\n")

		fmt.Fprintf(f, "}\n")

	}
}

func monadGenerateSequenceTGeneric(
	hkt func(string) string,
	fmap func(string, string) string,
	fap func(string, string) string,
) func(f *os.File, i int) {
	return func(f *os.File, i int) {

		tuple := tupleType("T")(i)

		fmt.Fprintf(f, "SequenceT%d[", i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", j+1)
		}
		fmt.Fprintf(f, "](")
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d %s", j+1, hkt(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(f, ") %s {", hkt(tuple))
		// the actual apply callback
		fmt.Fprintf(f, "  return apply.SequenceT%d(\n", i)
		// map callback

		curried := func(count int) string {
			var buf strings.Builder
			for j := count; j < i; j++ {
				buf.WriteString(fmt.Sprintf("func(T%d)", j+1))
			}
			buf.WriteString(tuple)
			return buf.String()
		}

		fmt.Fprintf(f, "    %s,\n", fmap("T1", curried(1)))
		for j := 1; j < i; j++ {
			fmt.Fprintf(f, "    %s,\n", fap(curried(j+1), fmt.Sprintf("T%d", j)))
		}

		for j := 0; j < i; j++ {
			fmt.Fprintf(f, "    T%d,\n", j+1)
		}

		fmt.Fprintf(f, "  )\n")

		fmt.Fprintf(f, "}\n")

	}
}

func generateTraverseTuple1(
	hkt func(string) string,
	infix string) func(f *os.File, i int) {

	return func(f *os.File, i int) {
		tuple := tupleType("T")(i)

		fmt.Fprintf(f, "\n// TraverseTuple%d converts a [Tuple%d] of [A] via transformation functions transforming [A] to [%s] into a [%s].\n", i, i, hkt("A"), hkt(fmt.Sprintf("Tuple%d", i)))
		fmt.Fprintf(f, "func TraverseTuple%d[", i)
		// functions
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "F%d ~func(A%d) %s", j+1, j+1, hkt(fmt.Sprintf("T%d", j+1)))
		}
		if S.IsNonEmpty(infix) {
			fmt.Fprintf(f, ", %s", infix)
		}
		// types
		for j := 0; j < i; j++ {
			fmt.Fprintf(f, ", A%d, T%d", j+1, j+1)
		}
		fmt.Fprintf(f, " any](")
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "f%d F%d", j+1, j+1)
		}
		fmt.Fprintf(f, ") func (T.Tuple%d[", i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "A%d", j+1)
		}
		fmt.Fprintf(f, "]) %s {\n", hkt(tuple))
		fmt.Fprintf(f, "  return func(t T.Tuple%d[", i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "A%d", j+1)
		}
		fmt.Fprintf(f, "]) %s {\n", hkt(tuple))
		fmt.Fprintf(f, "    return A.TraverseTuple%d(\n", i)
		// map
		fmt.Fprintf(f, "      Map[")
		if S.IsNonEmpty(infix) {
			fmt.Fprintf(f, "%s, T1,", infix)
		} else {
			fmt.Fprintf(f, "T1,")
		}
		for j := 1; j < i; j++ {
			fmt.Fprintf(f, " func(T%d)", j+1)
		}
		fmt.Fprintf(f, " %s],\n", tuple)
		// applicatives
		for j := 1; j < i; j++ {
			fmt.Fprintf(f, "      Ap[")
			for k := j + 1; k < i; k++ {
				if k > j+1 {
					fmt.Fprintf(f, " ")
				}
				fmt.Fprintf(f, "func(T%d)", k+1)
			}
			if j < i-1 {
				fmt.Fprintf(f, " ")
			}
			fmt.Fprintf(f, "%s", tuple)
			if S.IsNonEmpty(infix) {
				fmt.Fprintf(f, ", %s", infix)
			}
			fmt.Fprintf(f, ", T%d],\n", j+1)
		}
		for j := 0; j < i; j++ {
			fmt.Fprintf(f, "      f%d,\n", j+1)
		}
		fmt.Fprintf(f, "      t,\n")
		fmt.Fprintf(f, "    )\n")
		fmt.Fprintf(f, "  }\n")
		fmt.Fprintf(f, "}\n")
	}
}

func generateSequenceTuple1(
	hkt func(string) string,
	infix string) func(f *os.File, i int) {

	return func(f *os.File, i int) {

		tuple := tupleType("T")(i)

		fmt.Fprintf(f, "\n// SequenceTuple%d converts a [Tuple%d] of [%s] into an [%s].\n", i, i, hkt("T"), hkt(fmt.Sprintf("Tuple%d", i)))
		fmt.Fprintf(f, "func SequenceTuple%d[", i)
		if S.IsNonEmpty(infix) {
			fmt.Fprintf(f, "%s", infix)
		}
		for j := 0; j < i; j++ {
			if S.IsNonEmpty(infix) || j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", j+1)
		}
		fmt.Fprintf(f, " any](t T.Tuple%d[", i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "%s", hkt(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(f, "]) %s {\n", hkt(tuple))
		fmt.Fprintf(f, "  return A.SequenceTuple%d(\n", i)
		// map
		fmt.Fprintf(f, "    Map[")
		if S.IsNonEmpty(infix) {
			fmt.Fprintf(f, "%s, T1,", infix)
		} else {
			fmt.Fprintf(f, "T1,")
		}
		for j := 1; j < i; j++ {
			fmt.Fprintf(f, " func(T%d)", j+1)
		}
		fmt.Fprintf(f, " %s],\n", tuple)
		// applicatives
		for j := 1; j < i; j++ {
			fmt.Fprintf(f, "    Ap[")
			for k := j + 1; k < i; k++ {
				if k > j+1 {
					fmt.Fprintf(f, " ")
				}
				fmt.Fprintf(f, "func(T%d)", k+1)
			}
			if j < i-1 {
				fmt.Fprintf(f, " ")
			}
			fmt.Fprintf(f, "%s", tuple)
			if S.IsNonEmpty(infix) {
				fmt.Fprintf(f, ", %s", infix)
			}
			fmt.Fprintf(f, ", T%d],\n", j+1)
		}
		fmt.Fprintf(f, "    t,\n")
		fmt.Fprintf(f, "  )\n")
		fmt.Fprintf(f, "}\n")
	}
}

func generateSequenceT1(
	hkt func(string) string,
	infix string) func(f *os.File, i int) {

	return func(f *os.File, i int) {

		tuple := tupleType("T")(i)

		fmt.Fprintf(f, "\n// SequenceT%d converts %d parameters of [%s] into a [%s].\n", i, i, hkt("T"), hkt(fmt.Sprintf("Tuple%d", i)))
		fmt.Fprintf(f, "func SequenceT%d[", i)
		if S.IsNonEmpty(infix) {
			fmt.Fprintf(f, "%s", infix)
		}
		for j := 0; j < i; j++ {
			if S.IsNonEmpty(infix) || j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", j+1)
		}
		fmt.Fprintf(f, " any](")
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d %s", j+1, hkt(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(f, ") %s {\n", hkt(tuple))
		fmt.Fprintf(f, "  return A.SequenceT%d(\n", i)
		// map
		fmt.Fprintf(f, "    Map[")
		if S.IsNonEmpty(infix) {
			fmt.Fprintf(f, "%s, T1,", infix)
		} else {
			fmt.Fprintf(f, "T1,")
		}
		for j := 1; j < i; j++ {
			fmt.Fprintf(f, " func(T%d)", j+1)
		}
		fmt.Fprintf(f, " %s],\n", tuple)
		// applicatives
		for j := 1; j < i; j++ {
			fmt.Fprintf(f, "    Ap[")
			for k := j + 1; k < i; k++ {
				if k > j+1 {
					fmt.Fprintf(f, " ")
				}
				fmt.Fprintf(f, "func(T%d)", k+1)
			}
			if j < i-1 {
				fmt.Fprintf(f, " ")
			}
			fmt.Fprintf(f, "%s", tuple)
			if S.IsNonEmpty(infix) {
				fmt.Fprintf(f, ", %s", infix)
			}
			fmt.Fprintf(f, ", T%d],\n", j+1)
		}
		for j := 0; j < i; j++ {
			fmt.Fprintf(f, "    t%d,\n", j+1)
		}
		fmt.Fprintf(f, "  )\n")
		fmt.Fprintf(f, "}\n")
	}

}
