package lens

import (
	"github.com/IBM/fp-go/v2/optics/common"
	"github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/optics/prism"
)

type (
	Lens[S, A any]     = common.Lens[S, A]
	Prism[S, A any]    = prism.Prism[S, A]
	Optional[S, A any] = optional.Optional[S, A]
)
