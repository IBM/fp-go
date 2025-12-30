package ioeither

import (
	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	Consumer[A any] = consumer.Consumer[A]

	Predicate[A any] = predicate.Predicate[A]
)
