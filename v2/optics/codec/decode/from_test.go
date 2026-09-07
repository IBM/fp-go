package decode

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/record"
	"github.com/stretchr/testify/assert"
)

// TestFromOption tests the FromOption function
func TestFromOption(t *testing.T) {
	missing := lazy.Of(validation.FailureWithMessage[int]("key", "key is unavailable"))

	t.Run("lifts Some into a successful decoder", func(t *testing.T) {
		decoder := FromOption[validation.Context](missing)(option.Of(42))

		assert.Equal(t, validation.Of(42), decoder(nil))
	})

	t.Run("uses onNone for None", func(t *testing.T) {
		decoder := FromOption[validation.Context](missing)(option.None[int]())

		assert.Equal(t, validation.Failures[int](validation.Errors{&validation.ValidationError{
			Value:    "key",
			Messsage: "key is unavailable",
		}}), decoder(nil))
	})

	t.Run("failure carries the validation context", func(t *testing.T) {
		ctx := validation.Context{{Key: "user", Type: "User"}}
		decoder := FromOption[validation.Context](missing)(option.None[int]())

		assert.Equal(t, validation.Failures[int](validation.Errors{&validation.ValidationError{
			Value:    "key",
			Context:  ctx,
			Messsage: "key is unavailable",
		}}), decoder(ctx))
	})

	t.Run("onNone may recover with a default", func(t *testing.T) {
		decoder := FromOption[validation.Context](lazy.Of(Of[validation.Context](0)))(option.None[int]())

		assert.Equal(t, validation.Of(0), decoder(nil))
	})

	t.Run("onNone is not evaluated for Some", func(t *testing.T) {
		called := 0
		onNone := func() Decode[validation.Context, int] {
			called++
			return Left[validation.Context, int](validation.Errors{})
		}

		decoder := FromOption[validation.Context](onNone)(option.Of(42))

		assert.Equal(t, validation.Of(42), decoder(nil))
		assert.Equal(t, 0, called)
	})

	t.Run("composes with a record lookup", func(t *testing.T) {
		lookup := F.Flow2(
			record.Lookup[int]("age"),
			FromOption[validation.Context](missing),
		)

		assert.Equal(t, validation.Of(30), lookup(map[string]int{"age": 30})(nil))
		assert.True(t, either.IsLeft(lookup(map[string]int{})(nil)))
	})
}
