package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	C "github.com/urfave/cli/v2"
)

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

func generateContextReaderIOEitherTraverseTuple(suffix string) func(f, fg *os.File, i int) {
	return func(f, fg *os.File, i int) {
		// tupleT type
		tupleT := tupleType("T")(i)
		tupleA := tupleType("A")(i)

		// non-generic version
		// generic version
		fmt.Fprintf(f, "\n// Traverse%sTuple%d converts a [T.Tuple%d] of [A] via transformer functions transforming [A] to a [ReaderIOEither] into a [ReaderIOEither] of a [T.Tuple%d].\n", suffix, i, i, i)
		fmt.Fprintf(f, "func Traverse%sTuple%d[", suffix, i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "F%d ~func(A%d) ReaderIOEither[T%d]", j+1, j+1, j+1)
		}
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
		fmt.Fprintf(f, ") func(%s) ReaderIOEither[%s] {\n", tupleA, tupleT)
		fmt.Fprintf(f, "  return G.Traverse%sTuple%d[ReaderIOEither[%s]](", suffix, i, tupleT)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "f%d", j+1)
		}
		fmt.Fprintf(f, ")\n")
		fmt.Fprintf(f, "}\n")

		// generic version
		fmt.Fprintf(fg, "\n// Traverse%sTuple%d converts a [T.Tuple%d] of readers into a reader of a [T.Tuple%d].\n", suffix, i, i, i)
		fmt.Fprintf(fg, "func Traverse%sTuple%d[\n", suffix, i)
		fmt.Fprintf(fg, "  GR_TUPLE%d ~func(context.Context) GIO_TUPLE%d,\n", i, i)
		// the transformation functions
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  F%d ~func(A%d) GR_T%d,\n", j+1, j+1, j+1)
		}
		// the readers
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  GR_T%d ~func(context.Context) GIO_T%d,\n", j+1, j+1)
		}
		// the tuples
		fmt.Fprintf(fg, "  GIO_TUPLE%d ~func() E.Either[error, %s],\n", i, tupleT)
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  GIO_T%d ~func() E.Either[error, T%d],\n", j+1, j+1)
		}
		// input and result parameters
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  A%d,\n", j+1)
			fmt.Fprintf(fg, "  T%d", j+1)
			if j < i-1 {
				fmt.Fprintf(fg, ",\n")
			}
		}
		fmt.Fprintf(fg, " any](")
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(fg, ", ")
			}
			fmt.Fprintf(fg, "f%d F%d", j+1, j+1)
		}
		fmt.Fprintf(fg, ") func(%s) GR_TUPLE%d {\n", tupleA, i)
		fmt.Fprintf(fg, "  return func(t %s) GR_TUPLE%d {\n", tupleA, i)
		fmt.Fprintf(fg, "    return A.TraverseTuple%d(\n", i)
		// map call
		var cr string
		if i > 1 {
			cb := generateNestedCallbacks(1, i)
			cio := fmt.Sprintf("func() E.Either[error, %s]", cb)
			cr = fmt.Sprintf("func(context.Context) %s", cio)
		} else {
			cr = fmt.Sprintf("GR_TUPLE%d", i)
		}
		fmt.Fprintf(fg, "      Map[GR_T%d, %s, GIO_T%d],\n", 1, cr, 1)
		// the apply calls
		for j := 1; j < i; j++ {
			if j < i-1 {
				cb := generateNestedCallbacks(j+1, i)
				cio := fmt.Sprintf("func() E.Either[error, %s]", cb)
				cr = fmt.Sprintf("func(context.Context) %s", cio)
			} else {
				cr = fmt.Sprintf("GR_TUPLE%d", i)
			}
			fmt.Fprintf(fg, "      Ap%s[%s, func(context.Context) func() E.Either[error, %s], GR_T%d],\n", suffix, cr, generateNestedCallbacks(j, i), j+1)
		}
		// function parameters
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "      f%d,\n", j+1)
		}

		// raw parameters
		fmt.Fprintf(fg, "      t,\n")

		fmt.Fprintf(fg, "    )\n")
		fmt.Fprintf(fg, "  }\n")
		fmt.Fprintf(fg, "}\n")
	}
}

func generateContextReaderIOEitherSequenceTuple(suffix string) func(f, fg *os.File, i int) {
	return func(f, fg *os.File, i int) {
		// tuple type
		tuple := tupleType("T")(i)

		// non-generic version
		// generic version
		fmt.Fprintf(f, "\n// Sequence%sTuple%d converts a [T.Tuple%d] of [ReaderIOEither] into a [ReaderIOEither] of a [T.Tuple%d].\n", suffix, i, i, i)
		fmt.Fprintf(f, "func Sequence%sTuple%d[", suffix, i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", j+1)
		}
		fmt.Fprintf(f, " any](t T.Tuple%d[", i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "ReaderIOEither[T%d]", j+1)
		}
		fmt.Fprintf(f, "]) ReaderIOEither[%s] {\n", tuple)
		fmt.Fprintf(f, "  return G.Sequence%sTuple%d[ReaderIOEither[%s]](t)\n", suffix, i, tuple)
		fmt.Fprintf(f, "}\n")

		// generic version
		fmt.Fprintf(fg, "\n// Sequence%sTuple%d converts a [T.Tuple%d] of readers into a reader of a [T.Tuple%d].\n", suffix, i, i, i)
		fmt.Fprintf(fg, "func Sequence%sTuple%d[\n", suffix, i)

		fmt.Fprintf(fg, "  GR_TUPLE%d ~func(context.Context) GIO_TUPLE%d,\n", i, i)
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  GR_T%d ~func(context.Context) GIO_T%d,\n", j+1, j+1)
		}

		fmt.Fprintf(fg, "  GIO_TUPLE%d ~func() E.Either[error, %s],\n", i, tuple)
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  GIO_T%d ~func() E.Either[error, T%d],\n", j+1, j+1)
		}

		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  T%d", j+1)
			if j < i-1 {
				fmt.Fprintf(fg, ",\n")
			}
		}
		fmt.Fprintf(fg, " any](t T.Tuple%d[", i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(fg, ", ")
			}
			fmt.Fprintf(fg, "GR_T%d", j+1)
		}
		fmt.Fprintf(fg, "]) GR_TUPLE%d {\n", i)
		fmt.Fprintf(fg, "  return A.SequenceTuple%d(\n", i)
		// map call
		var cr string
		if i > 1 {
			cb := generateNestedCallbacks(1, i)
			cio := fmt.Sprintf("func() E.Either[error, %s]", cb)
			cr = fmt.Sprintf("func(context.Context) %s", cio)
		} else {
			cr = fmt.Sprintf("GR_TUPLE%d", i)
		}
		fmt.Fprintf(fg, "    Map[GR_T%d, %s, GIO_T%d],\n", 1, cr, 1)
		// the apply calls
		for j := 1; j < i; j++ {
			if j < i-1 {
				cb := generateNestedCallbacks(j+1, i)
				cio := fmt.Sprintf("func() E.Either[error, %s]", cb)
				cr = fmt.Sprintf("func(context.Context) %s", cio)
			} else {
				cr = fmt.Sprintf("GR_TUPLE%d", i)
			}
			fmt.Fprintf(fg, "    Ap%s[%s, func(context.Context) func() E.Either[error, %s], GR_T%d],\n", suffix, cr, generateNestedCallbacks(j, i), j+1)
		}
		// raw parameters
		fmt.Fprintf(fg, "    t,\n")

		fmt.Fprintf(fg, " )\n")
		fmt.Fprintf(fg, "}\n")
	}
}

func generateContextReaderIOEitherSequenceT(suffix string) func(f, fg *os.File, i int) {
	return func(f, fg *os.File, i int) {
		// tuple type
		tuple := tupleType("T")(i)

		// non-generic version
		// generic version
		fmt.Fprintf(f, "\n// Sequence%sT%d converts %d [ReaderIOEither] into a [ReaderIOEither] of a [T.Tuple%d].\n", suffix, i, i, i)
		fmt.Fprintf(f, "func Sequence%sT%d[", suffix, i)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "T%d", j+1)
		}
		fmt.Fprintf(f, " any](")
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d ReaderIOEither[T%d]", j+1, j+1)
		}
		fmt.Fprintf(f, ") ReaderIOEither[%s] {\n", tuple)
		fmt.Fprintf(f, "  return G.Sequence%sT%d[ReaderIOEither[%s]](", suffix, i, tuple)
		for j := 0; j < i; j++ {
			if j > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "t%d", j+1)
		}
		fmt.Fprintf(f, ")\n")
		fmt.Fprintf(f, "}\n")

		// generic version
		fmt.Fprintf(fg, "\n// Sequence%sT%d converts %d readers into a reader of a [T.Tuple%d].\n", suffix, i, i, i)
		fmt.Fprintf(fg, "func Sequence%sT%d[\n", suffix, i)

		fmt.Fprintf(fg, "  GR_TUPLE%d ~func(context.Context) GIO_TUPLE%d,\n", i, i)
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  GR_T%d ~func(context.Context) GIO_T%d,\n", j+1, j+1)
		}

		fmt.Fprintf(fg, "  GIO_TUPLE%d ~func() E.Either[error, %s],\n", i, tuple)
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  GIO_T%d ~func() E.Either[error, T%d],\n", j+1, j+1)
		}

		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  T%d", j+1)
			if j < i-1 {
				fmt.Fprintf(fg, ",\n")
			}
		}
		fmt.Fprintf(fg, " any](\n")
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "  t%d GR_T%d,\n", j+1, j+1)
		}
		fmt.Fprintf(fg, ") GR_TUPLE%d {\n", i)
		fmt.Fprintf(fg, "  return A.SequenceT%d(\n", i)
		// map call
		var cr string
		if i > 1 {
			cb := generateNestedCallbacks(1, i)
			cio := fmt.Sprintf("func() E.Either[error, %s]", cb)
			cr = fmt.Sprintf("func(context.Context) %s", cio)
		} else {
			cr = fmt.Sprintf("GR_TUPLE%d", i)
		}
		fmt.Fprintf(fg, "    Map[GR_T%d, %s, GIO_T%d],\n", 1, cr, 1)
		// the apply calls
		for j := 1; j < i; j++ {
			if j < i-1 {
				cb := generateNestedCallbacks(j+1, i)
				cio := fmt.Sprintf("func() E.Either[error, %s]", cb)
				cr = fmt.Sprintf("func(context.Context) %s", cio)
			} else {
				cr = fmt.Sprintf("GR_TUPLE%d", i)
			}
			fmt.Fprintf(fg, "    Ap%s[%s, func(context.Context) func() E.Either[error, %s], GR_T%d],\n", suffix, cr, generateNestedCallbacks(j, i), j+1)
		}
		// raw parameters
		for j := 0; j < i; j++ {
			fmt.Fprintf(fg, "    t%d,\n", j+1)
		}

		fmt.Fprintf(fg, " )\n")
		fmt.Fprintf(fg, "}\n")
	}
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

	G "github.com/IBM/fp-go/context/%s/generic"
	T "github.com/IBM/fp-go/tuple"
)
`, pkg)

	writePackage(fg, "generic")

	fmt.Fprintf(fg, `
import (
	"context"

	E "github.com/IBM/fp-go/either"
	RE "github.com/IBM/fp-go/readerioeither/generic"
	A "github.com/IBM/fp-go/internal/apply"
	T "github.com/IBM/fp-go/tuple"
)
`)

	generateContextReaderIOEitherEitherize(f, fg, 0)

	for i := 1; i <= count; i++ {
		// eitherize
		generateContextReaderIOEitherEitherize(f, fg, i)
		// sequenceT
		generateContextReaderIOEitherSequenceT("")(f, fg, i)
		generateContextReaderIOEitherSequenceT("Seq")(f, fg, i)
		generateContextReaderIOEitherSequenceT("Par")(f, fg, i)
		// sequenceTuple
		generateContextReaderIOEitherSequenceTuple("")(f, fg, i)
		generateContextReaderIOEitherSequenceTuple("Seq")(f, fg, i)
		generateContextReaderIOEitherSequenceTuple("Par")(f, fg, i)
		// traverseTuple
		generateContextReaderIOEitherTraverseTuple("")(f, fg, i)
		generateContextReaderIOEitherTraverseTuple("Seq")(f, fg, i)
		generateContextReaderIOEitherTraverseTuple("Par")(f, fg, i)
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
		Action: func(ctx *C.Context) error {
			return generateContextReaderIOEitherHelpers(
				ctx.String(keyFilename),
				ctx.Int(keyCount),
			)
		},
	}
}
