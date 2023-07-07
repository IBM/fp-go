package number

import (
	"strconv"

	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/option"
)

func atoi(value string) (int, bool) {
	data, err := strconv.Atoi(value)
	return data, err == nil
}

var (
	// Atoi converts a string to an integer
	Atoi = O.Optionize1(atoi)
	// Itoa converts an integer to a string
	Itoa = F.Flow2(strconv.Itoa, O.Of[string])
)
