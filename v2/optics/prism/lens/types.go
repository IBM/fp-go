package lens

import (
	"github.com/IBM/fp-go/v2/optics/common"
)

type (
	Lens[S, A any]     = common.Lens[S, A]
	Prism[S, A any]    = common.Prism[S, A]
	Optional[S, A any] = common.Optional[S, A]
)
