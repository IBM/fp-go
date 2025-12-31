package ioeither

import (
	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	// Consumer represents a function that consumes a value of type A.
	Consumer[A any] = consumer.Consumer[A]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	Predicate[A any] = predicate.Predicate[A]
)
