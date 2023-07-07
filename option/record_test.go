package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequenceRecord(t *testing.T) {
	assert.Equal(t, Of(map[string]string{
		"a": "A",
		"b": "B",
	}), SequenceRecord(map[string]Option[string]{
		"a": Of("A"),
		"b": Of("B"),
	}))
}
