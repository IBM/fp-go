package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	C "github.com/urfave/cli/v2"
)

const (
	keyCount = "count"
)

func generateNullary(f *os.File, i int) {
	// Create the nullary version
	fmt.Fprintf(f, "\n// Nullary%d creates a parameter less function from a parameter less function and %d functions. When executed the first parameter less function gets executed and then the result is piped through the remaining functions\n", i, i-1)
	fmt.Fprintf(f, "func Nullary%d[", i)
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "T%d", j)
	}
	fmt.Fprintf(f, " any](f1 func() T1")
	for j := 2; j <= i; j++ {
		fmt.Fprintf(f, ", f%d func(T%d) T%d", j, j-1, j)
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
	fmt.Fprintf(f, "\n// Flow%d creates a function that takes an initial value t0 and sucessively applies %d functions where the input of a function is the return value of the previous function\n// The final return value is the result of the last function application\n", i, i)
	fmt.Fprintf(f, "func Flow%d[T0", i)
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, " any](")
	for j := 1; j <= i; j++ {
		if j > 1 {
			fmt.Fprintf(f, ", ")
		}
		fmt.Fprintf(f, "f%d func(T%d) T%d", j, j-1, j)
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
	fmt.Fprintf(f, "\n// Pipe%d takes an initial value t0 and sucessively applies %d functions where the input of a function is the return value of the previous function\n// The final return value is the result of the last function application\n", i, i)
	fmt.Fprintf(f, "func Pipe%d[T0", i)
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, " any](t0 T0")
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", f%d func(T%d) T%d", j, j-1, j)
	}
	fmt.Fprintf(f, ") T%d {\n", i)
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, "  t%d := f%d(t%d)\n", j, j, j-1)
	}
	fmt.Fprintf(f, "  return t%d\n", i)
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
	fmt.Fprintf(f, "func Curry%d[T0", i)
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, " any](f func(T0")
	for j := 2; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j-1)
	}
	fmt.Fprintf(f, ") T%d) func(T0)", i)
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
	fmt.Fprintf(f, "func Uncurry%d[T0", i)
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, ", T%d", j)
	}
	fmt.Fprintf(f, " any](f")
	for j := 1; j <= i; j++ {
		fmt.Fprintf(f, " func(T%d)", j-1)
	}
	fmt.Fprintf(f, " T%d) func(", i)
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

func generateHelpers(count int) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	pkg := filepath.Base(absDir)
	f, err := os.Create("gen.go")
	if err != nil {
		return err
	}
	defer f.Close()
	// log
	log.Printf("Generating code for package [%s] with [%d] repetitions ...", pkg, count)

	// some header
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, "// This file was generated by robots at")
	fmt.Fprintf(f, "// %s\n", time.Now())

	fmt.Fprintf(f, "package %s\n", pkg)

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
	}

	return nil
}

func PipeCommand() *C.Command {
	return &C.Command{
		Name: "pipe",
		Flags: []C.Flag{
			&C.IntFlag{
				Name:  keyCount,
				Value: 20,
				Usage: "Number of variations to create",
			},
		},
		Action: func(ctx *C.Context) error {
			return generateHelpers(ctx.Int(keyCount))
		},
	}
}
