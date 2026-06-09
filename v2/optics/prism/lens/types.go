package lens

import (
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/optics/prism"
)

type (
	Lens[S, A any]     = lens.Lens[S, A]
	Prism[S, A any]    = prism.Prism[S, A]
	Optional[S, A any] = optional.Optional[S, A]
)
