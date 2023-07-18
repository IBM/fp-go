package testing

import (
	EQ "github.com/IBM/fp-go/eq"
	"github.com/stretchr/testify/assert"
)

// Eq implements the equal operation based on `ObjectsAreEqualValues` from the assertion library
func Eq[A any]() EQ.Eq[A] {
	return EQ.FromEquals(func(l, r A) bool {
		return assert.ObjectsAreEqualValues(l, r)
	})
}
