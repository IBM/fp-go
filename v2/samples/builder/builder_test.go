package builder

import (
	"testing"

	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestBuilderPrism(t *testing.T) {
	b1 := MakePerson("Carsten", 55)

	// this should be a valid person
	p1, ok := option.Unwrap(PersonPrism.GetOption(b1))
	assert.True(t, ok)

	// convert back to a builder
	b2 := PersonPrism.ReverseGet(p1)

	// change the name
	b3 := endomorphism.Chain(WithName("Jan"))(b1)

	p2 := PersonPrism.GetOption(b2)
	p3 := PersonPrism.GetOption(b3)

	assert.Equal(t, p2, option.Of(p1))
	assert.NotEqual(t, p3, option.Of(p1))

}
