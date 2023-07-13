package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	C "github.com/urfave/cli/v2"
)

func writeTupleType(f *os.File, i int) {
	fmt.Fprintf(f, "Tuple%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, "]")
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
	writeTupleType(f, i)
	fmt.Fprintf(f, "] {\n")

	fmt.Fprintf(f, "  return M.MakeMonoid(func(l, r ")
	writeTupleType(f, i)
	fmt.Fprintf(f, ") ")
	writeTupleType(f, i)
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
	writeTupleType(f, i)
	fmt.Fprintf(f, "] {\n")

	fmt.Fprintf(f, "  return O.MakeOrd(func(l, r ")
	writeTupleType(f, i)
	fmt.Fprintf(f, ") int {\n")

	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "    if c:= o%d.Compare(l.F%d, r.F%d); c != 0 {return c}\n", j, j, j)
	}
	fmt.Fprintf(f, "    return 0\n")
	fmt.Fprintf(f, "  }, func(l, r ")
	writeTupleType(f, i)
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
	writeTupleType(f, i)
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
	fmt.Fprintf(f, "\n// Tupled%d converts a function with %d parameters returning into a function taking a Tuple%d\n// The inverse function is [Untupled%d]\n", i, i, i, i)
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
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	// log
	log.Printf("Generating code in [%s] for package [%s] with [%d] repetitions ...", filename, pkg, count)

	// some header
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, "// This file was generated by robots at")
	fmt.Fprintf(f, "// %s\n", time.Now())

	fmt.Fprintf(f, "package %s\n\n", pkg)

	fmt.Fprintf(f, `
import (
	M "github.com/ibm/fp-go/monoid"
	O "github.com/ibm/fp-go/ord"
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
	}

	return nil
}

func TupleCommand() *C.Command {
	return &C.Command{
		Name:  "tuple",
		Usage: "generate code for Tuple",
		Flags: []C.Flag{
			flagCount,
			flagFilename,
		},
		Action: func(ctx *C.Context) error {
			return generateTupleHelpers(
				ctx.String(keyFilename),
				ctx.Int(keyCount),
			)
		},
	}
}
