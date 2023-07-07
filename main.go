package main

import (
	"fmt"

	O "github.com/ibm/fp-go/option"
)

func isNonemptyString(val string) bool {
	return val != ""
}

// var O = OptionModule{of: O_of, some: O_of, none: none, mp: OMap}

func main() {

	opt_string := O.FromPredicate(isNonemptyString)

	stringO1 := opt_string("Carsten")
	stringO2 := opt_string("")

	fmt.Println(stringO1)
	fmt.Println(stringO2)

}
