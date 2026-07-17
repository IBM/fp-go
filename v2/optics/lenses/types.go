package lenses

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
)

type (
	Lens[S, A any]      = lens.Lens[S, A]
	Option[A any]       = option.Option[A]
	Endomorphism[A any] = endomorphism.Endomorphism[A]
)
