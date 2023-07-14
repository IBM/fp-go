package cli

import (
	"fmt"
	"os"
	"strings"
)

func tupleType(i int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("T.Tuple%d[", i))
	for j := 0; j < i; j++ {
		if j > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("T%d", j+1))
	}
	buf.WriteString("]")

	return buf.String()
}

func monadGenerateSequenceTNonGeneric(
	hkt func(string) string,
	fmap func(string, string) string,
	fap func(string, string) string,
) func(f *os.File, i int) {
	return func(f *os.File, i int) {

		tuple := tupleType(i)

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
		fmt.Fprintf(f, "  return Apply.SequenceT%d(\n", i)
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

		tuple := tupleType(i)

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
		fmt.Fprintf(f, "  return Apply.SequenceT%d(\n", i)
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
