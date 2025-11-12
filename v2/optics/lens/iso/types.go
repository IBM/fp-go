package iso

import (
	"github.com/IBM/fp-go/v2/optics/iso"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
)

type (
	Option[A any]         = option.Option[A]
	Iso[S, A any]         = iso.Iso[S, A]
	Lens[S, A any]        = lens.Lens[S, A]
	Operator[S, A, B any] = lens.Operator[S, A, B]
)
