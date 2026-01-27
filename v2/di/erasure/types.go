package erasure

import (
	"github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/record"
)

type (
	Option[T any]              = option.Option[T]
	IOResult[T any]            = ioresult.IOResult[T]
	IOOption[T any]            = iooption.IOOption[T]
	Entry[K comparable, V any] = record.Entry[K, V]
	ReaderIOResult[R, T any]   = readerioresult.ReaderIOResult[R, T]
)
