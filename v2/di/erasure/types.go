package erasure

import (
	"github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/option"
)

type (
	Option[T any]   = option.Option[T]
	IOResult[T any] = ioresult.IOResult[T]
	IOOption[T any] = iooption.IOOption[T]
)
