package result

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertEq[A any](l A, lerr error) func(A, error) func(*testing.T) {
	return func(r A, rerr error) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, lerr, rerr)
			if (lerr != nil) && (rerr != nil) {
				assert.Equal(t, l, r)
			}
		}
	}
}
