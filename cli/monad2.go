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
	"os"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	N "github.com/IBM/fp-go/number"
	S "github.com/IBM/fp-go/string"
)

var (
	concStrgs     = A.Monoid[string]().Concat
	intercalStrgs = A.Intercalate(S.Monoid)
	concAllStrgs  = A.ConcatAll(A.Monoid[string]())
)

func joinAll(middle string) func(all ...[]string) string {
	ic := intercalStrgs(middle)
	return func(all ...[]string) string {
		return ic(concAllStrgs(all))
	}
}

func generateGenericSequenceT(
	nonGenericType func(string) string,
	genericType func(string) string,
	extra []string,
) func(f, fg *os.File, i int) {
	return func(f, fg *os.File, i int) {
		// tuple
		tuple := tupleType("T")(i)
		// all types T
		typesT := A.MakeBy(i, F.Flow2(
			N.Inc[int],
			S.Format[int]("T%d"),
		))
		// non generic version
		fmt.Fprintf(f, "\n// SequenceT%d converts %d [%s] into a [%s]\n", i, i, nonGenericType("T"), nonGenericType(tuple))
		fmt.Fprintf(f, "func SequenceT%d[%s any](\n", i, joinAll(", ")(extra, typesT))
		for j := 0; j < i; j++ {
			fmt.Fprintf(f, "  t%d %s,\n", j+1, nonGenericType(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(f, ") %s {\n", nonGenericType(tuple))
		fmt.Fprintf(f, "  return G.SequenceT%d[\n", i)
		fmt.Fprintf(f, "    %s,\n", nonGenericType(tuple))
		for j := 0; j < i; j++ {
			fmt.Fprintf(f, "    %s,\n", nonGenericType(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(f, "  ](")
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d", j+1)
		}
		fmt.Fprintf(f, ")\n")
		fmt.Fprintf(f, "}\n")

		// generic version
		fmt.Fprintf(fg, "\n// SequenceT%d converts %d [%s] into a [%s]\n", i, i, genericType("T"), genericType(tuple))
		fmt.Fprintf(fg, "func SequenceT%d[\n", i)
		fmt.Fprintf(fg, "  G_TUPLE%d ~%s,\n", i, genericType(tuple))
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  G_T%d ~%s, \n", j+1, genericType(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(fg, "  %s any](\n", joinAll(", ")(extra, typesT))
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  t%d G_T%d,\n", j+1, j+1)
		}
		fmt.Fprintf(fg, ")  G_TUPLE%d {\n", i)
		fmt.Fprintf(fg, "  return A.SequenceT%d(\n", i)
		// map call
		var cio string
		cb := generateNestedCallbacks(1, i)
		if i > 1 {
			cio = genericType(cb)
		} else {
			cio = fmt.Sprintf("G_TUPLE%d", i)
		}
		fmt.Fprintf(fg, "    Map[%s],\n", joinAll(", ")(A.From("G_T1", cio), extra, A.From("T1", cb)))
		// the apply calls
		for j := 1; j < i; j++ {
			if j < i-1 {
				cb := generateNestedCallbacks(j+1, i)
				cio = genericType(cb)
			} else {
				cio = fmt.Sprintf("G_TUPLE%d", i)
			}
			fmt.Fprintf(fg, "    Ap[%s, %s, G_T%d],\n", cio, genericType(generateNestedCallbacks(j, i)), j+1)
		}
		// function parameters
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "    t%d,\n", j+1)
		}

		fmt.Fprintf(fg, "  )\n")
		fmt.Fprintf(fg, "}\n")
	}
}

func generateGenericSequenceTuple(
	nonGenericType func(string) string,
	genericType func(string) string,
	extra []string,
) func(f, fg *os.File, i int) {
	return func(f, fg *os.File, i int) {
		// tuple
		tuple := tupleType("T")(i)
		// all types T
		typesT := A.MakeBy(i, F.Flow2(
			N.Inc[int],
			S.Format[int]("T%d"),
		))
		// non generic version
		fmt.Fprintf(f, "\n// SequenceTuple%d converts a [T.Tuple%d[%s]] into a [%s]\n", i, i, nonGenericType("T"), nonGenericType(tuple))
		fmt.Fprintf(f, "func SequenceTuple%d[%s any](t T.Tuple%d[", i, joinAll(", ")(extra, typesT), i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "%s", nonGenericType(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(f, "]) %s {\n", nonGenericType(tuple))
		fmt.Fprintf(f, "  return G.SequenceTuple%d[\n", i)
		fmt.Fprintf(f, "    %s,\n", nonGenericType(tuple))
		for j := 0; j < i; j++ {
			fmt.Fprintf(f, "    %s,\n", nonGenericType(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(f, "  ](t)\n")
		fmt.Fprintf(f, "}\n")

		// generic version
		fmt.Fprintf(fg, "\n// SequenceTuple%d converts a [T.Tuple%d[%s]] into a [%s]\n", i, i, genericType("T"), genericType(tuple))
		fmt.Fprintf(fg, "func SequenceTuple%d[\n", i)
		fmt.Fprintf(fg, "  G_TUPLE%d ~%s,\n", i, genericType(tuple))
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  G_T%d ~%s, \n", j+1, genericType(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(fg, "  %s any](t T.Tuple%d[", joinAll(", ")(extra, typesT), i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(fg, ", ")
			}
			fmt.Fprintf(fg, "G_T%d", j+1)
		}
		fmt.Fprintf(fg, "])  G_TUPLE%d {\n", i)
		fmt.Fprintf(fg, "  return A.SequenceTuple%d(\n", i)
		// map call
		var cio string
		cb := generateNestedCallbacks(1, i)
		if i > 1 {
			cio = genericType(cb)
		} else {
			cio = fmt.Sprintf("G_TUPLE%d", i)
		}
		fmt.Fprintf(fg, "    Map[%s],\n", joinAll(", ")(A.From("G_T1", cio), extra, A.From("T1", cb)))
		// the apply calls
		for j := 1; j < i; j++ {
			if j < i-1 {
				cb := generateNestedCallbacks(j+1, i)
				cio = genericType(cb)
			} else {
				cio = fmt.Sprintf("G_TUPLE%d", i)
			}
			fmt.Fprintf(fg, "    Ap[%s, %s, G_T%d],\n", cio, genericType(generateNestedCallbacks(j, i)), j+1)
		}
		// function parameters
		fmt.Fprintf(fg, "    t)\n")
		fmt.Fprintf(fg, "}\n")
	}
}

func generateGenericTraverseTuple(
	nonGenericType func(string) string,
	genericType func(string) string,
	extra []string,
) func(f, fg *os.File, i int) {
	return func(f, fg *os.File, i int) {
		// tuple
		tupleT := tupleType("T")(i)
		tupleA := tupleType("A")(i)
		// all types T
		typesT := A.MakeBy(i, F.Flow2(
			N.Inc[int],
			S.Format[int]("T%d"),
		))
		// all types A
		typesA := A.MakeBy(i, F.Flow2(
			N.Inc[int],
			S.Format[int]("A%d"),
		))
		// all function types
		typesF := A.MakeBy(i, F.Flow2(
			N.Inc[int],
			func(j int) string {
				return fmt.Sprintf("F%d ~func(A%d) %s", j, j, nonGenericType(fmt.Sprintf("T%d", j)))
			},
		))
		// non generic version
		fmt.Fprintf(f, "\n// TraverseTuple%d converts a [T.Tuple%d[%s]] into a [%s]\n", i, i, nonGenericType("T"), nonGenericType(tupleT))
		fmt.Fprintf(f, "func TraverseTuple%d[%s any](", i, joinAll(", ")(typesF, extra, typesA, typesT))
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "f%d F%d", j+1, j+1)
		}
		fmt.Fprintf(f, ") func(%s) %s {\n", tupleA, nonGenericType(tupleT))
		fmt.Fprintf(f, "  return G.TraverseTuple%d[%s](", i, nonGenericType(tupleT))
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "f%d", j+1)
		}
		fmt.Fprintf(f, ")\n")
		fmt.Fprintf(f, "}\n")

		// generic version
		fmt.Fprintf(fg, "\n// TraverseTuple%d converts a [T.Tuple%d[%s]] into a [%s]\n", i, i, genericType("T"), genericType(tupleT))
		fmt.Fprintf(fg, "func TraverseTuple%d[\n", i)
		fmt.Fprintf(fg, "  G_TUPLE%d ~%s,\n", i, genericType(tupleT))
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  F%d ~func(A%d) G_T%d,\n", j+1, j+1, j+1)
		}
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  G_T%d ~%s, \n", j+1, genericType(fmt.Sprintf("T%d", j+1)))
		}
		fmt.Fprintf(fg, "  %s any](", joinAll(", ")(extra, typesA, typesT))
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(fg, ", ")
			}
			fmt.Fprintf(fg, "f%d F%d", j+1, j+1)
		}
		fmt.Fprintf(fg, ") func(%s) G_TUPLE%d {\n", tupleA, i)
		fmt.Fprintf(fg, "  return func(t %s) G_TUPLE%d {\n", tupleA, i)
		fmt.Fprintf(fg, "    return A.TraverseTuple%d(\n", i)
		// map call
		var cio string
		cb := generateNestedCallbacks(1, i)
		if i > 1 {
			cio = genericType(cb)
		} else {
			cio = fmt.Sprintf("G_TUPLE%d", i)
		}
		fmt.Fprintf(fg, "    Map[%s],\n", joinAll(", ")(A.From("G_T1", cio), extra, A.From("T1", cb)))
		// the apply calls
		for j := 1; j < i; j++ {
			if j < i-1 {
				cb := generateNestedCallbacks(j+1, i)
				cio = genericType(cb)
			} else {
				cio = fmt.Sprintf("G_TUPLE%d", i)
			}
			fmt.Fprintf(fg, "    Ap[%s, %s, G_T%d],\n", cio, genericType(generateNestedCallbacks(j, i)), j+1)
		}
		// function parameters
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "    f%d,\n", j+1)
		}
		// tuple parameter
		fmt.Fprintf(fg, "    t)\n")
		fmt.Fprintf(fg, "   }\n")
		fmt.Fprintf(fg, "}\n")
	}
}
