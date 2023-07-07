package string

import (
	"strings"

	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/ord"
)

var (
	// ToUpperCase converts the string to uppercase
	ToUpperCase = strings.ToUpper

	// ToLowerCase converts the string to lowercase
	ToLowerCase = strings.ToLower

	// Ord implements the default ordering for strings
	Ord = O.FromStrictCompare[string]()
)

func Eq(left string, right string) bool {
	return left == right
}

func ToBytes(s string) []byte {
	return []byte(s)
}

func IsEmpty(s string) bool {
	return len(s) == 0
}

func IsNonEmpty(s string) bool {
	return len(s) > 0
}

func Size(s string) int {
	return len(s)
}

// Includes returns a predicate that tests for the existence of the search string
func Includes(searchString string) func(string) bool {
	return F.Bind2nd(strings.Contains, searchString)
}

// Equals returns a predicate that tests if a string is equal
func Equals(other string) func(string) bool {
	return F.Bind2nd(Eq, other)
}
