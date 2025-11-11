package builder

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/result"
)

// BuilderPrism createa a [Prism] that converts between a builder and its type
func BuilderPrism[T any, B Builder[T]](creator func(T) B) Prism[B, T] {
	return prism.MakePrism(F.Flow2(B.Build, result.ToOption[T]), creator)
}
