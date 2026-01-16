package validation

import (
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Value:    "test",
		Messsage: "invalid value",
	}

	assert.Equal(t, "ValidationError", err.Error())
}

func TestValidationError_String(t *testing.T) {
	err := &ValidationError{
		Value:    "test",
		Messsage: "invalid value",
	}

	expected := "ValidationError: invalid value"
	assert.Equal(t, expected, err.String())
}

func TestValidationError_Unwrap(t *testing.T) {
	t.Run("with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := &ValidationError{
			Value:    "test",
			Messsage: "invalid value",
			Cause:    cause,
		}

		assert.Equal(t, cause, err.Unwrap())
	})

	t.Run("without cause", func(t *testing.T) {
		err := &ValidationError{
			Value:    "test",
			Messsage: "invalid value",
		}

		assert.Nil(t, err.Unwrap())
	})
}

func TestValidationError_Format(t *testing.T) {
	t.Run("simple format without context", func(t *testing.T) {
		err := &ValidationError{
			Value:    "test",
			Messsage: "invalid value",
		}

		result := fmt.Sprintf("%v", err)
		assert.Equal(t, "invalid value", result)
	})

	t.Run("with context path", func(t *testing.T) {
		err := &ValidationError{
			Value:    "test",
			Context:  []ContextEntry{{Key: "user"}, {Key: "name"}},
			Messsage: "must not be empty",
		}

		result := fmt.Sprintf("%v", err)
		assert.Equal(t, "at user.name: must not be empty", result)
	})

	t.Run("with context using type", func(t *testing.T) {
		err := &ValidationError{
			Value:    123,
			Context:  []ContextEntry{{Type: "User"}, {Key: "age"}},
			Messsage: "must be positive",
		}

		result := fmt.Sprintf("%v", err)
		assert.Equal(t, "at User.age: must be positive", result)
	})

	t.Run("with cause - simple format", func(t *testing.T) {
		cause := errors.New("parse error")
		err := &ValidationError{
			Value:    "abc",
			Messsage: "invalid number",
			Cause:    cause,
		}

		result := fmt.Sprintf("%v", err)
		assert.Equal(t, "invalid number (caused by: parse error)", result)
	})

	t.Run("with cause - verbose format", func(t *testing.T) {
		cause := errors.New("parse error")
		err := &ValidationError{
			Value:    "abc",
			Messsage: "invalid number",
			Cause:    cause,
		}

		result := fmt.Sprintf("%+v", err)
		assert.Contains(t, result, "invalid number")
		assert.Contains(t, result, "caused by: parse error")
		assert.Contains(t, result, `value: "abc"`)
	})

	t.Run("verbose format shows value", func(t *testing.T) {
		err := &ValidationError{
			Value:    42,
			Messsage: "out of range",
		}

		result := fmt.Sprintf("%+v", err)
		assert.Contains(t, result, "out of range")
		assert.Contains(t, result, "value: 42")
	})

	t.Run("complex context path", func(t *testing.T) {
		err := &ValidationError{
			Value: "invalid",
			Context: []ContextEntry{
				{Key: "user"},
				{Key: "address"},
				{Key: "zipCode"},
			},
			Messsage: "invalid format",
		}

		result := fmt.Sprintf("%v", err)
		assert.Equal(t, "at user.address.zipCode: invalid format", result)
	})
}

func TestFailures(t *testing.T) {
	t.Run("creates left either with errors", func(t *testing.T) {
		errs := Errors{
			&ValidationError{Value: "test", Messsage: "error 1"},
			&ValidationError{Value: "test", Messsage: "error 2"},
		}

		result := Failures[int](errs)

		assert.True(t, either.IsLeft(result))
		left := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, left, 2)
		assert.Equal(t, "error 1", left[0].Messsage)
		assert.Equal(t, "error 2", left[1].Messsage)
	})

	t.Run("preserves error details", func(t *testing.T) {
		errs := Errors{
			&ValidationError{
				Value:    "abc",
				Context:  []ContextEntry{{Key: "field"}},
				Messsage: "invalid",
				Cause:    errors.New("cause"),
			},
		}

		result := Failures[string](errs)

		left := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		require.Len(t, left, 1)
		assert.Equal(t, "abc", left[0].Value)
		assert.Equal(t, "invalid", left[0].Messsage)
		assert.NotNil(t, left[0].Cause)
		assert.Len(t, left[0].Context, 1)
	})
}

func TestSuccess(t *testing.T) {
	t.Run("creates right either with value", func(t *testing.T) {
		result := Success(42)

		assert.True(t, either.IsRight(result))
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("works with different types", func(t *testing.T) {
		strResult := Success("hello")
		str := either.MonadFold(strResult,
			func(Errors) string { return "" },
			F.Identity[string],
		)
		assert.Equal(t, "hello", str)

		boolResult := Success(true)
		b := either.MonadFold(boolResult,
			func(Errors) bool { return false },
			F.Identity[bool],
		)
		assert.Equal(t, true, b)

		type Custom struct{ Name string }
		customResult := Success(Custom{Name: "test"})
		custom := either.MonadFold(customResult,
			func(Errors) Custom { return Custom{} },
			F.Identity[Custom],
		)
		assert.Equal(t, "test", custom.Name)
	})
}

func TestFailureWithMessage(t *testing.T) {
	t.Run("creates failure with context", func(t *testing.T) {
		fail := FailureWithMessage[int]("abc", "expected integer")
		context := []ContextEntry{{Key: "age", Type: "int"}}

		result := fail(context)

		assert.True(t, either.IsLeft(result))
		errs := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		require.Len(t, errs, 1)
		assert.Equal(t, "abc", errs[0].Value)
		assert.Equal(t, "expected integer", errs[0].Messsage)
		assert.Equal(t, context, errs[0].Context)
		assert.Nil(t, errs[0].Cause)
	})

	t.Run("works with empty context", func(t *testing.T) {
		fail := FailureWithMessage[string](123, "wrong type")
		result := fail(nil)

		errs := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		require.Len(t, errs, 1)
		assert.Equal(t, 123, errs[0].Value)
		assert.Nil(t, errs[0].Context)
	})

	t.Run("preserves complex context", func(t *testing.T) {
		fail := FailureWithMessage[bool]("not a bool", "type mismatch")
		context := []ContextEntry{
			{Key: "user"},
			{Key: "settings"},
			{Key: "enabled"},
		}

		result := fail(context)

		errs := either.MonadFold(result,
			F.Identity[Errors],
			func(bool) Errors { return nil },
		)
		require.Len(t, errs, 1)
		assert.Equal(t, context, errs[0].Context)
	})
}

func TestFailureWithError(t *testing.T) {
	t.Run("creates failure with cause and context", func(t *testing.T) {
		cause := errors.New("parse failed")
		fail := FailureWithError[int]("abc", "invalid number")
		context := []ContextEntry{{Key: "count"}}

		result := fail(cause)(context)

		assert.True(t, either.IsLeft(result))
		errs := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		require.Len(t, errs, 1)
		assert.Equal(t, "abc", errs[0].Value)
		assert.Equal(t, "invalid number", errs[0].Messsage)
		assert.Equal(t, context, errs[0].Context)
		assert.Equal(t, cause, errs[0].Cause)
	})

	t.Run("cause is unwrappable", func(t *testing.T) {
		cause := errors.New("underlying")
		fail := FailureWithError[string](nil, "wrapper")
		result := fail(cause)(nil)

		errs := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		require.Len(t, errs, 1)
		assert.True(t, errors.Is(errs[0], cause))
	})

	t.Run("works with complex error chain", func(t *testing.T) {
		root := errors.New("root cause")
		wrapped := fmt.Errorf("wrapped: %w", root)
		fail := FailureWithError[int](0, "validation failed")

		result := fail(wrapped)([]ContextEntry{{Key: "field"}})

		errs := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		require.Len(t, errs, 1)
		assert.True(t, errors.Is(errs[0], root))
		assert.True(t, errors.Is(errs[0], wrapped))
	})
}

func TestValidationIntegration(t *testing.T) {
	t.Run("success and failure can be combined", func(t *testing.T) {
		success := Success(42)
		failure := Failures[int](Errors{
			&ValidationError{Value: "bad", Messsage: "error"},
		})

		assert.True(t, either.IsRight(success))
		assert.True(t, either.IsLeft(failure))
	})

	t.Run("context provides meaningful error paths", func(t *testing.T) {
		fail := FailureWithMessage[string](nil, "required field")
		context := []ContextEntry{
			{Key: "request"},
			{Key: "body"},
			{Key: "user"},
			{Key: "email"},
		}

		result := fail(context)
		errs := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)

		formatted := fmt.Sprintf("%v", errs[0])
		assert.Contains(t, formatted, "request.body.user.email")
		assert.Contains(t, formatted, "required field")
	})

	t.Run("multiple errors can be collected", func(t *testing.T) {
		errs := Errors{
			&ValidationError{
				Context:  []ContextEntry{{Key: "name"}},
				Messsage: "too short",
			},
			&ValidationError{
				Context:  []ContextEntry{{Key: "age"}},
				Messsage: "must be positive",
			},
			&ValidationError{
				Context:  []ContextEntry{{Key: "email"}},
				Messsage: "invalid format",
			},
		}

		result := Failures[any](errs)
		collected := either.MonadFold(result,
			F.Identity[Errors],
			func(any) Errors { return nil },
		)

		assert.Len(t, collected, 3)
		messages := make([]string, len(collected))
		for i, err := range collected {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "too short")
		assert.Contains(t, messages, "must be positive")
		assert.Contains(t, messages, "invalid format")
	})
}

func TestValidationError_FormatEdgeCases(t *testing.T) {
	t.Run("empty message", func(t *testing.T) {
		err := &ValidationError{
			Value:    "test",
			Messsage: "",
		}

		result := fmt.Sprintf("%v", err)
		assert.Equal(t, "", result)
	})

	t.Run("context with empty keys", func(t *testing.T) {
		err := &ValidationError{
			Value:    "test",
			Context:  []ContextEntry{{Key: ""}, {Type: "Type"}, {Key: ""}},
			Messsage: "error",
		}

		result := fmt.Sprintf("%v", err)
		// Should handle empty keys gracefully
		assert.Contains(t, result, "error")
	})

	t.Run("nil value", func(t *testing.T) {
		err := &ValidationError{
			Value:    nil,
			Messsage: "nil not allowed",
		}

		result := fmt.Sprintf("%+v", err)
		assert.Contains(t, result, "nil not allowed")
		assert.Contains(t, result, "value: <nil>")
	})
}
