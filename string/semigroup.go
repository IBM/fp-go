package string

import (
	"fmt"

	S "github.com/IBM/fp-go/semigroup"
)

func concat(left string, right string) string {
	return fmt.Sprintf("%s%s", left, right)
}

func Semigroup() S.Semigroup[string] {
	return S.MakeSemigroup(concat)
}
