package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniq(t *testing.T) {
	data := From(1, 2, 3, 2, 4, 1)

	uniq := StrictUniq(data)
	assert.Equal(t, From(1, 2, 3, 4), uniq)
}
