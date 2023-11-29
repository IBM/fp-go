package either

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompactArray(t *testing.T) {
	ar := []Either[string, string]{
		Of[string]("ok"),
		Left[string]("err"),
		Of[string]("ok"),
	}

	res := CompactArray(ar)
	assert.Equal(t, 2, len(res))
}
