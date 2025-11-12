package di

import (
	"github.com/IBM/fp-go/v2/context/ioresult"
	"github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Option[T any]   = option.Option[T]
	Result[T any]   = result.Result[T]
	IOResult[T any] = ioresult.IOResult[T]
	IOOption[T any] = iooption.IOOption[T]
)
