package testing

import (
	"testing"

	AR "github.com/ibm/fp-go/internal/array"
	M "github.com/ibm/fp-go/monoid"
	"github.com/stretchr/testify/assert"
)

func assertLaws[A any](t *testing.T, m M.Monoid[A]) func(a A) bool {
	e := m.Empty()
	return func(a A) bool {
		return assert.Equal(t, a, m.Concat(a, e), "Monoid right identity") &&
			assert.Equal(t, a, m.Concat(e, a), "Monoid left identity")
	}
}

// AssertLaws asserts the monoid laws for a dataset
func AssertLaws[A any](t *testing.T, m M.Monoid[A]) func(data []A) bool {
	law := assertLaws(t, m)

	return func(data []A) bool {
		return AR.Reduce(data, func(result bool, value A) bool {
			return result && law(value)
		}, true)
	}
}
