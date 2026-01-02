package di

import (
	"github.com/IBM/fp-go/v2/context/ioresult"
	"github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/record"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Option represents an optional value that may or may not be present.
	Option[T any] = option.Option[T]

	// Result represents a computation that may fail with an error.
	Result[T any] = result.Result[T]

	// IOResult represents a synchronous computation that may fail with an error.
	IOResult[T any] = ioresult.IOResult[T]

	// IOOption represents a synchronous computation that may not produce a value.
	IOOption[T any] = iooption.IOOption[T]

	// Entry represents a key-value pair in a record/map structure.
	Entry[K comparable, V any] = record.Entry[K, V]
)
