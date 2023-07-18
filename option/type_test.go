package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeConversion(t *testing.T) {

	var src any = "Carsten"

	dst := ToType[string](src)
	assert.Equal(t, Some("Carsten"), dst)
}

func TestInvalidConversion(t *testing.T) {
	var src any = make(map[string]string)

	dst := ToType[int](src)
	assert.Equal(t, None[int](), dst)
}
