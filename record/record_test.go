package record

import (
	"sort"
	"testing"

	"github.com/ibm/fp-go/internal/utils"
	O "github.com/ibm/fp-go/option"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	data := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	keys := Keys(data)
	sort.Strings(keys)

	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestValues(t *testing.T) {
	data := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	keys := Values(data)
	sort.Strings(keys)

	assert.Equal(t, []string{"A", "B", "C"}, keys)
}

func TestMap(t *testing.T) {
	data := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	expected := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	assert.Equal(t, expected, Map[string](utils.Upper)(data))
}

func TestLookup(t *testing.T) {
	data := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	assert.Equal(t, O.Some("a"), Lookup[string, string]("a")(data))
	assert.Equal(t, O.None[string](), Lookup[string, string]("a1")(data))
}
