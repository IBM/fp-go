package string

import (
	"fmt"

	S "github.com/ibm/fp-go/semigroup"
)

func concat(left string, right string) string {
	return fmt.Sprintf("%s%s", left, right)
}

func Semigroup() S.Semigroup[string] {
	return S.MakeSemigroup(concat)
}
