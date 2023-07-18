package testing

import (
	"fmt"
	"testing"

	EQ "github.com/IBM/fp-go/eq"
	"github.com/stretchr/testify/assert"
)

func TestMonadLaws(t *testing.T) {
	// some comparison
	eqa := EQ.FromStrictEquals[bool]()
	eqb := EQ.FromStrictEquals[int]()
	eqc := EQ.FromStrictEquals[string]()

	ab := func(a bool) int {
		if a {
			return 1
		}
		return 0
	}

	bc := func(b int) string {
		return fmt.Sprintf("value %d", b)
	}

	laws := AssertLaws(t, eqa, eqb, eqc, ab, bc)

	assert.True(t, laws(true))
	assert.True(t, laws(false))
}
