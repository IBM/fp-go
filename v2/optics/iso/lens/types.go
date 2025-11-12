package lens

import (
	"github.com/IBM/fp-go/v2/optics/iso"
	L "github.com/IBM/fp-go/v2/optics/lens"
)

type (
	Lens[S, A any] = L.Lens[S, A]
	Iso[S, A any]  = iso.Iso[S, A]
)
