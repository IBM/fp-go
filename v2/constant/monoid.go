package constant

import (
	"github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
)

// Monoid returns a [M.Monoid] that returns a constant value in all operations
func Monoid[A any](a A) M.Monoid[A] {
	return M.MakeMonoid(function.Constant2[A, A](a), a)
}
