package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEq[A any](l A, lok bool) func(A, bool) func(*testing.T) {
	return func(r A, rok bool) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, lok, rok)
			if lok && rok {
				assert.Equal(t, l, r)
			}
		}
	}
}
