package lenses

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/common"
	"github.com/IBM/fp-go/v2/option"
)

type (
	Lens[S, A any]      = common.Lens[S, A]
	Option[A any]       = option.Option[A]
	Endomorphism[A any] = endomorphism.Endomorphism[A]
)
