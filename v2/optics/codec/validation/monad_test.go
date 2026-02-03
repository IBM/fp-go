package validation

import (
	"fmt"
	"strings"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	t.Run("creates successful validation", func(t *testing.T) {
		result := Of(42)

		assert.Equal(t, Of(42), result)
	})

	t.Run("works with different types", func(t *testing.T) {
		strResult := Of("hello")
		assert.Equal(t, Of("hello"), strResult)

		boolResult := Of(true)
		assert.Equal(t, Of(true), boolResult)

		type Custom struct{ Value int }
		customResult := Of(Custom{Value: 100})
		assert.Equal(t, Of(Custom{Value: 100}), customResult)
	})

	t.Run("is equivalent to Success", func(t *testing.T) {
		value := 42
		ofResult := Of(value)
		successResult := Success(value)

		assert.Equal(t, ofResult, successResult)
	})
}

func TestMap(t *testing.T) {
	t.Run("transforms successful validation", func(t *testing.T) {
		double := N.Mul(2)
		result := Map(double)(Of(21))

		assert.Equal(t, Of(42), result)
	})

	t.Run("preserves failure", func(t *testing.T) {
		errs := Errors{&ValidationError{Messsage: "error"}}
		failure := Failures[int](errs)

		double := N.Mul(2)
		result := Map(double)(failure)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "error", errors[0].Messsage)
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		add10 := func(x int) int { return x + 10 }
		double := N.Mul(2)
		toString := func(x int) string { return fmt.Sprintf("%d", x) }

		result := F.Pipe3(
			Of(5),
			Map(add10),
			Map(double),
			Map(toString),
		)

		assert.Equal(t, Of("30"), result)
	})

	t.Run("type transformation", func(t *testing.T) {
		length := func(s string) int { return len(s) }
		result := Map(length)(Of("hello"))

		assert.Equal(t, Of(5), result)
	})
}

func TestAp(t *testing.T) {
	t.Run("applies function to value when both succeed", func(t *testing.T) {
		double := N.Mul(2)
		funcValidation := Of(double)
		valueValidation := Of(21)

		result := Ap[int](valueValidation)(funcValidation)

		assert.Equal(t, Of(42), result)
	})

	t.Run("accumulates errors when value fails", func(t *testing.T) {
		double := N.Mul(2)
		funcValidation := Of(double)
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := Ap[int](valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "value error", errors[0].Messsage)
	})

	t.Run("accumulates errors when function fails", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := Of(21)

		result := Ap[int](valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "function error", errors[0].Messsage)
	})

	t.Run("accumulates all errors when both fail", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := Ap[int](valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2)
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "function error")
		assert.Contains(t, messages, "value error")
	})

	t.Run("applies with string transformation", func(t *testing.T) {
		toUpper := func(s string) string { return fmt.Sprintf("UPPER:%s", s) }
		funcValidation := Of(toUpper)
		valueValidation := Of("hello")

		result := Ap[string](valueValidation)(funcValidation)

		assert.Equal(t, Of("UPPER:hello"), result)
	})

	t.Run("accumulates multiple validation errors from different sources", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error 1"},
			&ValidationError{Messsage: "function error 2"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error 1"},
		})

		result := Ap[int](valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 3)
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "function error 1")
		assert.Contains(t, messages, "function error 2")
		assert.Contains(t, messages, "value error 1")
	})
}

func TestMonadLaws(t *testing.T) {
	t.Run("functor identity law", func(t *testing.T) {
		// Map(id) == id
		value := Of(42)
		mapped := Map(F.Identity[int])(value)

		assert.Equal(t, value, mapped)
	})

	t.Run("functor composition law", func(t *testing.T) {
		// Map(f . g) == Map(f) . Map(g)
		f := N.Mul(2)
		g := func(x int) int { return x + 10 }
		composed := func(x int) int { return f(g(x)) }

		value := Of(5)
		left := Map(composed)(value)
		right := F.Pipe2(value, Map(g), Map(f))

		assert.Equal(t, left, right)
	})

	t.Run("applicative identity law", func(t *testing.T) {
		// Ap(v)(Of(id)) == v
		v := Of(42)
		result := Ap[int](v)(Of(F.Identity[int]))

		assert.Equal(t, v, result)
	})

	t.Run("applicative homomorphism law", func(t *testing.T) {
		// Ap(Of(x))(Of(f)) == Of(f(x))
		f := N.Mul(2)
		x := 21

		left := Ap[int](Of(x))(Of(f))
		right := Of(f(x))

		assert.Equal(t, left, right)
	})
}

func TestMapWithOperator(t *testing.T) {
	t.Run("Map returns an Operator", func(t *testing.T) {
		double := N.Mul(2)
		operator := Map(double)

		// Operator can be applied to different validations
		result1 := operator(Of(10))
		result2 := operator(Of(20))

		val1 := either.MonadFold(result1,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		val2 := either.MonadFold(result2,
			func(Errors) int { return 0 },
			F.Identity[int],
		)

		assert.Equal(t, 20, val1)
		assert.Equal(t, 40, val2)
	})
}

func TestApWithOperator(t *testing.T) {
	t.Run("Ap returns an Operator", func(t *testing.T) {
		valueValidation := Of(21)
		operator := Ap[int](valueValidation)

		// Operator can be applied to different function validations
		double := N.Mul(2)
		triple := func(x int) int { return x * 3 }

		result1 := operator(Of(double))
		result2 := operator(Of(triple))

		assert.Equal(t, Of(42), result1)
		assert.Equal(t, Of(63), result2)
	})
}

func TestApplicative(t *testing.T) {
	t.Run("returns non-nil instance", func(t *testing.T) {
		app := Applicative[int, string]()
		assert.NotNil(t, app)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		app1 := Applicative[int, string]()
		app2 := Applicative[int, string]()

		// Both should work independently
		result1 := app1.Of(42)
		result2 := app2.Of(43)

		assert.Equal(t, Of(42), result1)
		assert.Equal(t, Of(43), result2)
	})
}

func TestApplicativeOf(t *testing.T) {
	app := Applicative[int, string]()

	t.Run("wraps a value in Validation context", func(t *testing.T) {
		result := app.Of(42)
		assert.Equal(t, Of(42), result)
	})

	t.Run("wraps string value", func(t *testing.T) {
		app := Applicative[string, int]()
		result := app.Of("hello")
		assert.Equal(t, Of("hello"), result)
	})

	t.Run("wraps zero value", func(t *testing.T) {
		result := app.Of(0)
		assert.Equal(t, Of(0), result)
	})

	t.Run("wraps complex types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		app := Applicative[User, string]()
		user := User{Name: "Alice", Age: 30}
		result := app.Of(user)

		assert.Equal(t, Of(user), result)
	})
}

func TestApplicativeMap(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("maps a function over successful validation", func(t *testing.T) {
		double := N.Mul(2)
		result := app.Map(double)(app.Of(21))

		assert.Equal(t, Of(42), result)
	})

	t.Run("maps type conversion", func(t *testing.T) {
		app := Applicative[int, string]()
		toString := func(x int) string { return fmt.Sprintf("%d", x) }
		result := app.Map(toString)(app.Of(42))

		assert.Equal(t, Of("42"), result)
	})

	t.Run("maps identity function", func(t *testing.T) {
		result := app.Map(F.Identity[int])(app.Of(42))

		assert.Equal(t, Of(42), result)
		value := either.MonadFold(result,
			func(Errors) int { return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 42, value)
	})

	t.Run("preserves failure", func(t *testing.T) {
		errs := Errors{&ValidationError{Messsage: "error"}}
		failure := Failures[int](errs)

		double := N.Mul(2)
		result := app.Map(double)(failure)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "error", errors[0].Messsage)
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		app := Applicative[int, string]()
		add10 := func(x int) int { return x + 10 }
		double := N.Mul(2)
		toString := func(x int) string { return fmt.Sprintf("%d", x) }

		result := F.Pipe3(
			app.Of(5),
			Map(add10),
			Map(double),
			app.Map(toString),
		)

		assert.Equal(t, Of("30"), result)
	})
}

func TestApplicativeAp(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("applies wrapped function to wrapped value when both succeed", func(t *testing.T) {
		double := N.Mul(2)
		funcValidation := Of(double)
		valueValidation := app.Of(21)

		result := app.Ap(valueValidation)(funcValidation)

		assert.Equal(t, Of(42), result)
	})

	t.Run("accumulates errors when value fails", func(t *testing.T) {
		double := N.Mul(2)
		funcValidation := Of(double)
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "value error", errors[0].Messsage)
	})

	t.Run("accumulates errors when function fails", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := app.Of(21)

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "function error", errors[0].Messsage)
	})

	t.Run("accumulates all errors when both fail", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2)
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "function error")
		assert.Contains(t, messages, "value error")
	})

	t.Run("applies with type conversion", func(t *testing.T) {
		app := Applicative[int, string]()
		toString := func(x int) string { return fmt.Sprintf("value:%d", x) }
		funcValidation := Of(toString)
		valueValidation := app.Of(42)

		result := app.Ap(valueValidation)(funcValidation)

		assert.Equal(t, Of("value:42"), result)
	})

	t.Run("accumulates multiple errors from different sources", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error 1"},
			&ValidationError{Messsage: "function error 2"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error 1"},
			&ValidationError{Messsage: "value error 2"},
		})

		result := app.Ap(valueValidation)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 4)
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "function error 1")
		assert.Contains(t, messages, "function error 2")
		assert.Contains(t, messages, "value error 1")
		assert.Contains(t, messages, "value error 2")
	})
}

func TestApplicativeComposition(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("composes Map and Of", func(t *testing.T) {
		double := N.Mul(2)
		result := F.Pipe1(
			app.Of(21),
			app.Map(double),
		)

		assert.Equal(t, Of(42), result)
	})

	t.Run("composes multiple Map operations", func(t *testing.T) {
		app := Applicative[int, string]()
		double := N.Mul(2)
		toString := func(x int) string { return fmt.Sprintf("%d", x) }

		result := F.Pipe2(
			app.Of(21),
			Map(double),
			app.Map(toString),
		)

		assert.Equal(t, Of("42"), result)
	})

	t.Run("composes Map and Ap", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		ioFunc := F.Pipe1(
			app.Of(5),
			Map(add),
		)
		valueValidation := app.Of(16)

		result := app.Ap(valueValidation)(ioFunc)

		assert.Equal(t, Of(21), result)
	})
}

func TestApplicativeLawsWithInstance(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("identity law: Ap(Of(id))(v) == v", func(t *testing.T) {
		identity := func(x int) int { return x }
		v := app.Of(42)

		left := app.Ap(v)(Of(identity))
		right := v

		assert.Equal(t, right, left)
	})

	t.Run("homomorphism law: Ap(Of(x))(Of(f)) == Of(f(x))", func(t *testing.T) {
		f := N.Mul(2)
		x := 21

		left := app.Ap(app.Of(x))(Of(f))
		right := app.Of(f(x))

		assert.Equal(t, right, left)
	})

	t.Run("interchange law: Ap(Of(y))(u) == Ap(u)(Of($ y))", func(t *testing.T) {
		double := N.Mul(2)
		u := Of(double)
		y := 21

		left := app.Ap(app.Of(y))(u)

		applyY := func(f func(int) int) int { return f(y) }
		right := Ap[int](u)(Of(applyY))

		assert.Equal(t, right, left)
	})

	t.Run("identity law with failure", func(t *testing.T) {
		identity := func(x int) int { return x }
		v := Failures[int](Errors{&ValidationError{Messsage: "error"}})

		left := app.Ap(v)(Of(identity))
		right := v

		assert.Equal(t, right, left)
	})
}

func TestApplicativeMultipleArguments(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("applies curried two-argument function", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		// Create validation with curried function
		funcValidation := F.Pipe1(
			app.Of(10),
			Map(add),
		)

		// Apply to second argument
		result := app.Ap(app.Of(32))(funcValidation)

		assert.Equal(t, Of(42), result)
	})

	t.Run("applies curried three-argument function", func(t *testing.T) {
		add3 := func(a int) func(int) func(int) int {
			return func(b int) func(int) int {
				return func(c int) int {
					return a + b + c
				}
			}
		}

		// Build up the computation step by step
		funcValidation1 := F.Pipe1(
			app.Of(10),
			Map(add3),
		)

		funcValidation2 := Ap[func(int) int](app.Of(20))(funcValidation1)
		result := Ap[int](app.Of(12))(funcValidation2)

		assert.Equal(t, Of(42), result)
	})

	t.Run("accumulates errors from multiple arguments", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		// First argument fails
		arg1 := Failures[int](Errors{&ValidationError{Messsage: "arg1 error"}})
		// Second argument fails
		arg2 := Failures[int](Errors{&ValidationError{Messsage: "arg2 error"}})

		funcValidation := F.Pipe1(arg1, Map(add))
		result := app.Ap(arg2)(funcValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2)
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "arg1 error")
		assert.Contains(t, messages, "arg2 error")
	})
}

func TestApplicativeWithDifferentTypes(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		app := Applicative[int, string]()
		toString := func(x int) string { return fmt.Sprintf("%d", x) }
		result := app.Map(toString)(app.Of(42))

		assert.Equal(t, Of("42"), result)
	})

	t.Run("string to int", func(t *testing.T) {
		app := Applicative[string, int]()
		toLength := func(s string) int { return len(s) }
		result := app.Map(toLength)(app.Of("hello"))

		assert.Equal(t, Of(5), result)
	})

	t.Run("bool to string", func(t *testing.T) {
		app := Applicative[bool, string]()
		toString := func(b bool) string {
			if b {
				return "true"
			}
			return "false"
		}
		result := app.Map(toString)(app.Of(true))

		assert.Equal(t, Of("true"), result)
	})
}

func TestApplicativeRealWorldScenario(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	t.Run("validates user with all valid fields", func(t *testing.T) {
		validateName := func(name string) Validation[string] {
			if len(name) < 3 {
				return Failures[string](Errors{&ValidationError{Messsage: "Name must be at least 3 characters"}})
			}
			return Success(name)
		}

		validateAge := func(age int) Validation[int] {
			if age < 18 {
				return Failures[int](Errors{&ValidationError{Messsage: "Must be 18 or older"}})
			}
			return Success(age)
		}

		validateEmail := func(email string) Validation[string] {
			if len(email) == 0 {
				return Failures[string](Errors{&ValidationError{Messsage: "Email is required"}})
			}
			return Success(email)
		}

		makeUser := func(name string) func(int) func(string) User {
			return func(age int) func(string) User {
				return func(email string) User {
					return User{Name: name, Age: age, Email: email}
				}
			}
		}

		name := validateName("Alice")
		age := validateAge(25)
		email := validateEmail("alice@example.com")

		// Use the standalone Ap function with proper type parameters
		result := Ap[User](email)(Ap[func(string) User](age)(Ap[func(int) func(string) User](name)(Of(makeUser))))

		expectedUser := User{Name: "Alice", Age: 25, Email: "alice@example.com"}
		assert.Equal(t, Of(expectedUser), result)
	})

	t.Run("accumulates all validation errors", func(t *testing.T) {
		validateName := func(name string) Validation[string] {
			if len(name) < 3 {
				return Failures[string](Errors{&ValidationError{Messsage: "Name must be at least 3 characters"}})
			}
			return Success(name)
		}

		validateAge := func(age int) Validation[int] {
			if age < 18 {
				return Failures[int](Errors{&ValidationError{Messsage: "Must be 18 or older"}})
			}
			return Success(age)
		}

		validateEmail := func(email string) Validation[string] {
			if len(email) == 0 {
				return Failures[string](Errors{&ValidationError{Messsage: "Email is required"}})
			}
			return Success(email)
		}

		makeUser := func(name string) func(int) func(string) User {
			return func(age int) func(string) User {
				return func(email string) User {
					return User{Name: name, Age: age, Email: email}
				}
			}
		}

		// All validations fail
		name := validateName("ab")
		age := validateAge(16)
		email := validateEmail("")

		// Use the standalone Ap function with proper type parameters
		result := Ap[User](email)(Ap[func(string) User](age)(Ap[func(int) func(string) User](name)(Of(makeUser))))

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(User) Errors { return nil },
		)
		assert.Len(t, errors, 3)
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "Name must be at least 3 characters")
		assert.Contains(t, messages, "Must be 18 or older")
		assert.Contains(t, messages, "Email is required")
	})
}

// TestMonadChainLeft tests the MonadChainLeft function with error aggregation
func TestMonadChainLeft(t *testing.T) {
	t.Run("Success value passes through unchanged", func(t *testing.T) {
		result := MonadChainLeft(
			Success(42),
			func(errs Errors) Validation[int] {
				return Failures[int](Errors{
					&ValidationError{Messsage: "should not be called"},
				})
			},
		)
		assert.Equal(t, Success(42), result)
	})

	t.Run("Failure is transformed to Success (recovery)", func(t *testing.T) {
		result := MonadChainLeft(
			Failures[int](Errors{
				&ValidationError{Messsage: "not found"},
			}),
			func(errs Errors) Validation[int] {
				if len(errs) > 0 && errs[0].Messsage == "not found" {
					return Success(0) // recover with default
				}
				return Failures[int](errs)
			},
		)
		assert.Equal(t, Success(0), result)
	})

	t.Run("Errors are aggregated when transformation fails", func(t *testing.T) {
		result := MonadChainLeft(
			Failures[int](Errors{
				&ValidationError{Messsage: "error 1"},
				&ValidationError{Messsage: "error 2"},
			}),
			func(errs Errors) Validation[int] {
				// Transformation also fails
				return Failures[int](Errors{
					&ValidationError{Messsage: "error 3"},
				})
			},
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 3, "Should aggregate all errors")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "error 1")
		assert.Contains(t, messages, "error 2")
		assert.Contains(t, messages, "error 3")
	})

	t.Run("Multiple errors aggregated from both sources", func(t *testing.T) {
		result := MonadChainLeft(
			Failures[string](Errors{
				&ValidationError{Messsage: "original error 1"},
				&ValidationError{Messsage: "original error 2"},
			}),
			func(errs Errors) Validation[string] {
				return Failures[string](Errors{
					&ValidationError{Messsage: "new error 1"},
					&ValidationError{Messsage: "new error 2"},
				})
			},
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 4, "Should aggregate all 4 errors")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "original error 1")
		assert.Contains(t, messages, "original error 2")
		assert.Contains(t, messages, "new error 1")
		assert.Contains(t, messages, "new error 2")
	})

	t.Run("Adding context to existing errors", func(t *testing.T) {
		result := MonadChainLeft(
			Failures[int](Errors{
				&ValidationError{
					Value:    "abc",
					Messsage: "invalid number",
				},
			}),
			func(errs Errors) Validation[int] {
				// Add contextual information
				return Failures[int](Errors{
					&ValidationError{
						Context: []ContextEntry{
							{Key: "user", Type: "User"},
							{Key: "age", Type: "int"},
						},
						Messsage: "failed to parse user age",
					},
				})
			},
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should have both original and context errors")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "invalid number")
		assert.Contains(t, messages, "failed to parse user age")
		// Find the error with context
		var contextError *ValidationError
		for _, err := range errors {
			if len(err.Context) > 0 {
				contextError = err
				break
			}
		}
		assert.NotNil(t, contextError, "Should have an error with context")
		assert.Len(t, contextError.Context, 2)
	})

	t.Run("Conditional recovery based on error content", func(t *testing.T) {
		handler := func(errs Errors) Validation[int] {
			for _, err := range errs {
				switch err.Messsage {
				case "not found":
					return Success(0)
				case "timeout":
					return Success(-1)
				}
			}
			// Add recovery attempt error
			return Failures[int](Errors{
				&ValidationError{Messsage: "recovery failed"},
			})
		}

		// Test recovery for "not found"
		result1 := MonadChainLeft(
			Failures[int](Errors{&ValidationError{Messsage: "not found"}}),
			handler,
		)
		assert.Equal(t, Success(0), result1)

		// Test recovery for "timeout"
		result2 := MonadChainLeft(
			Failures[int](Errors{&ValidationError{Messsage: "timeout"}}),
			handler,
		)
		assert.Equal(t, Success(-1), result2)

		// Test error aggregation for unknown error
		result3 := MonadChainLeft(
			Failures[int](Errors{&ValidationError{Messsage: "unknown error"}}),
			handler,
		)
		assert.True(t, either.IsLeft(result3))
		errors := either.MonadFold(result3,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should aggregate original and recovery errors")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "unknown error")
		assert.Contains(t, messages, "recovery failed")
	})

	t.Run("Chaining multiple MonadChainLeft operations", func(t *testing.T) {
		// First transformation
		step1 := MonadChainLeft(
			Failures[int](Errors{
				&ValidationError{Messsage: "step 1 error"},
			}),
			func(errs Errors) Validation[int] {
				return Failures[int](Errors{
					&ValidationError{Messsage: "step 2 error"},
				})
			},
		)

		// Second transformation
		step2 := MonadChainLeft(
			step1,
			func(errs Errors) Validation[int] {
				return Failures[int](Errors{
					&ValidationError{Messsage: "step 3 error"},
				})
			},
		)

		assert.True(t, either.IsLeft(step2))
		errors := either.MonadFold(step2,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 3, "Should aggregate errors from all steps")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "step 1 error")
		assert.Contains(t, messages, "step 2 error")
		assert.Contains(t, messages, "step 3 error")
	})
}

// TestChainLeft tests the curried ChainLeft function with error aggregation
func TestChainLeft(t *testing.T) {
	t.Run("Curried function transforms failures", func(t *testing.T) {
		handler := ChainLeft(func(errs Errors) Validation[int] {
			if len(errs) > 0 && errs[0].Messsage == "not found" {
				return Success(0)
			}
			return Failures[int](errs)
		})

		result := handler(Failures[int](Errors{
			&ValidationError{Messsage: "not found"},
		}))
		assert.Equal(t, Success(0), result)
	})

	t.Run("Curried function with Success value", func(t *testing.T) {
		handler := ChainLeft(func(errs Errors) Validation[int] {
			return Failures[int](Errors{
				&ValidationError{Messsage: "should not be called"},
			})
		})

		result := handler(Success(42))
		assert.Equal(t, Success(42), result)
	})

	t.Run("Use in pipeline with error aggregation", func(t *testing.T) {
		addContext := ChainLeft(func(errs Errors) Validation[string] {
			return Failures[string](Errors{
				&ValidationError{Messsage: "context: validation failed"},
			})
		})

		result := F.Pipe1(
			Failures[string](Errors{
				&ValidationError{Messsage: "original error"},
			}),
			addContext,
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should aggregate both errors")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "original error")
		assert.Contains(t, messages, "context: validation failed")
	})

	t.Run("Compose multiple ChainLeft operations with aggregation", func(t *testing.T) {
		handler1 := ChainLeft(func(errs Errors) Validation[int] {
			return Failures[int](Errors{
				&ValidationError{Messsage: "handler 1 error"},
			})
		})

		handler2 := ChainLeft(func(errs Errors) Validation[int] {
			return Failures[int](Errors{
				&ValidationError{Messsage: "handler 2 error"},
			})
		})

		result := F.Pipe2(
			Failures[int](Errors{
				&ValidationError{Messsage: "original error"},
			}),
			handler1,
			handler2,
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 3, "Should aggregate all errors from pipeline")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "original error")
		assert.Contains(t, messages, "handler 1 error")
		assert.Contains(t, messages, "handler 2 error")
	})

	t.Run("Error recovery in pipeline", func(t *testing.T) {
		recoverFromTimeout := ChainLeft(func(errs Errors) Validation[int] {
			for _, err := range errs {
				if err.Messsage == "timeout" {
					return Success(0)
				}
			}
			return Failures[int](errs)
		})

		// Test with timeout error - should recover
		result1 := F.Pipe1(
			Failures[int](Errors{&ValidationError{Messsage: "timeout"}}),
			recoverFromTimeout,
		)
		assert.Equal(t, Success(0), result1)

		// Test with other error - should propagate
		result2 := F.Pipe1(
			Failures[int](Errors{&ValidationError{Messsage: "other error"}}),
			recoverFromTimeout,
		)
		assert.True(t, either.IsLeft(result2))
	})

	t.Run("ChainLeft with Map combination", func(t *testing.T) {
		errorHandler := ChainLeft(func(errs Errors) Validation[int] {
			return Failures[int](Errors{
				&ValidationError{Messsage: "handled error"},
			})
		})

		valueMapper := Map(func(n int) string {
			return fmt.Sprintf("Value: %d", n)
		})

		// Test with Failure - errors should aggregate
		result1 := F.Pipe2(
			Failures[int](Errors{
				&ValidationError{Messsage: "original error"},
			}),
			errorHandler,
			valueMapper,
		)
		assert.True(t, either.IsLeft(result1))
		errors := either.MonadFold(result1,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 2)

		// Test with Success
		result2 := F.Pipe2(
			Success(42),
			errorHandler,
			valueMapper,
		)
		assert.Equal(t, Success("Value: 42"), result2)
	})

	t.Run("Reusable error enrichment handlers", func(t *testing.T) {
		addFieldContext := func(field string) func(Errors) Validation[string] {
			return func(errs Errors) Validation[string] {
				return Failures[string](Errors{
					&ValidationError{
						Context:  []ContextEntry{{Key: field, Type: "string"}},
						Messsage: fmt.Sprintf("validation failed for field: %s", field),
					},
				})
			}
		}

		emailHandler := ChainLeft(addFieldContext("email"))
		nameHandler := ChainLeft(addFieldContext("name"))

		// Apply email context
		result1 := emailHandler(Failures[string](Errors{
			&ValidationError{Messsage: "invalid format"},
		}))
		errors1 := either.MonadFold(result1,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors1, 2)
		messages1 := make([]string, len(errors1))
		for i, err := range errors1 {
			messages1[i] = err.Messsage
		}
		assert.Contains(t, messages1, "invalid format")
		// Check that one of the messages contains "email"
		hasEmail := false
		for _, msg := range messages1 {
			if strings.Contains(msg, "email") {
				hasEmail = true
				break
			}
		}
		assert.True(t, hasEmail, "Should have an error message containing 'email'")

		// Apply name context
		result2 := nameHandler(Failures[string](Errors{
			&ValidationError{Messsage: "too short"},
		}))
		errors2 := either.MonadFold(result2,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors2, 2)
		messages2 := make([]string, len(errors2))
		for i, err := range errors2 {
			messages2[i] = err.Messsage
		}
		assert.Contains(t, messages2, "too short")
		// Check that one of the messages contains "name"
		hasName := false
		for _, msg := range messages2 {
			if strings.Contains(msg, "name") {
				hasName = true
				break
			}
		}
		assert.True(t, hasName, "Should have an error message containing 'name'")
	})

	t.Run("Complex error aggregation scenario", func(t *testing.T) {
		// Simulate a validation pipeline with multiple error sources
		validateFormat := ChainLeft(func(errs Errors) Validation[string] {
			return Failures[string](Errors{
				&ValidationError{Messsage: "format validation failed"},
			})
		})

		validateLength := ChainLeft(func(errs Errors) Validation[string] {
			return Failures[string](Errors{
				&ValidationError{Messsage: "length validation failed"},
			})
		})

		validateContent := ChainLeft(func(errs Errors) Validation[string] {
			return Failures[string](Errors{
				&ValidationError{Messsage: "content validation failed"},
			})
		})

		result := F.Pipe3(
			Failures[string](Errors{
				&ValidationError{Messsage: "initial error"},
			}),
			validateFormat,
			validateLength,
			validateContent,
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 4, "Should aggregate all errors from pipeline")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "initial error")
		assert.Contains(t, messages, "format validation failed")
		assert.Contains(t, messages, "length validation failed")
		assert.Contains(t, messages, "content validation failed")
	})
}

// TestMonadMap tests the MonadMap function
func TestMonadMap(t *testing.T) {
	t.Run("transforms successful validation", func(t *testing.T) {
		result := MonadMap(Of(21), N.Mul(2))
		assert.Equal(t, Of(42), result)
	})

	t.Run("preserves failure unchanged", func(t *testing.T) {
		failure := Failures[int](Errors{
			&ValidationError{Messsage: "error"},
		})
		result := MonadMap(failure, N.Mul(2))

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "error", errors[0].Messsage)
	})

	t.Run("type transformation", func(t *testing.T) {
		result := MonadMap(Of(42), func(x int) string {
			return fmt.Sprintf("Value: %d", x)
		})

		assert.Equal(t, Of("Value: 42"), result)
	})

	t.Run("computing derived values", func(t *testing.T) {
		type User struct {
			FirstName string
			LastName  string
		}

		result := MonadMap(
			Of(User{"John", "Doe"}),
			func(u User) string { return u.FirstName + " " + u.LastName },
		)

		assert.Equal(t, Of("John Doe"), result)
	})

	t.Run("chaining multiple MonadMap operations", func(t *testing.T) {
		step1 := MonadMap(Of(5), func(x int) int { return x + 10 })
		step2 := MonadMap(step1, N.Mul(2))
		step3 := MonadMap(step2, func(x int) string { return fmt.Sprintf("%d", x) })

		assert.Equal(t, Of("30"), step3)
	})

	t.Run("identity function", func(t *testing.T) {
		original := Of(42)
		result := MonadMap(original, F.Identity[int])
		assert.Equal(t, original, result)
	})

	t.Run("preserves multiple errors", func(t *testing.T) {
		failure := Failures[int](Errors{
			&ValidationError{Messsage: "error 1"},
			&ValidationError{Messsage: "error 2"},
		})
		result := MonadMap(failure, func(x int) string { return fmt.Sprintf("%d", x) })

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 2)
		assert.Equal(t, "error 1", errors[0].Messsage)
		assert.Equal(t, "error 2", errors[1].Messsage)
	})
}

// TestMonadAp tests the MonadAp function with error accumulation
func TestMonadAp(t *testing.T) {
	t.Run("applies function to value when both succeed", func(t *testing.T) {
		double := N.Mul(2)
		result := MonadAp(Of(double), Of(21))
		assert.Equal(t, Of(42), result)
	})

	t.Run("returns function errors when function fails", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := Of(21)

		result := MonadAp(funcValidation, valueValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "function error", errors[0].Messsage)
	})

	t.Run("returns value errors when value fails", func(t *testing.T) {
		double := N.Mul(2)
		funcValidation := Of(double)
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := MonadAp(funcValidation, valueValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "value error", errors[0].Messsage)
	})

	t.Run("accumulates all errors when both fail", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error"},
		})

		result := MonadAp(funcValidation, valueValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should accumulate errors from both sources")
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "function error")
		assert.Contains(t, messages, "value error")
	})

	t.Run("accumulates multiple errors from both sources", func(t *testing.T) {
		funcValidation := Failures[func(int) int](Errors{
			&ValidationError{Messsage: "function error 1"},
			&ValidationError{Messsage: "function error 2"},
		})
		valueValidation := Failures[int](Errors{
			&ValidationError{Messsage: "value error 1"},
			&ValidationError{Messsage: "value error 2"},
		})

		result := MonadAp(funcValidation, valueValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 4, "Should accumulate all 4 errors")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "function error 1")
		assert.Contains(t, messages, "function error 2")
		assert.Contains(t, messages, "value error 1")
		assert.Contains(t, messages, "value error 2")
	})

	t.Run("type transformation with success", func(t *testing.T) {
		toString := func(x int) string { return fmt.Sprintf("Value: %d", x) }
		result := MonadAp(Of(toString), Of(42))

		assert.Equal(t, Of("Value: 42"), result)
	})

	t.Run("curried function application", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		// First application
		step1 := MonadMap(Of(10), add)
		// Second application
		result := MonadAp(step1, Of(32))

		assert.Equal(t, Of(42), result)
	})

	t.Run("validating multiple fields accumulates all errors", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}

		makeUser := func(name string) func(int) User {
			return func(age int) User { return User{name, age} }
		}

		nameValidation := Failures[string](Errors{
			&ValidationError{Messsage: "name too short"},
		})
		ageValidation := Failures[int](Errors{
			&ValidationError{Messsage: "age too low"},
		})

		// Apply name first
		step1 := MonadAp(Of(makeUser), nameValidation)
		// Apply age second
		result := MonadAp(step1, ageValidation)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(User) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should accumulate errors from both fields")
		messages := []string{errors[0].Messsage, errors[1].Messsage}
		assert.Contains(t, messages, "name too short")
		assert.Contains(t, messages, "age too low")
	})
}

// TestMonadChain tests the MonadChain function
func TestMonadChain(t *testing.T) {
	t.Run("chains successful validations", func(t *testing.T) {
		result := MonadChain(
			Of(42),
			func(x int) Validation[string] {
				return Of(fmt.Sprintf("Value: %d", x))
			},
		)

		assert.Equal(t, Of("Value: 42"), result)
	})

	t.Run("short-circuits on first failure", func(t *testing.T) {
		failure := Failures[int](Errors{
			&ValidationError{Messsage: "initial error"},
		})

		result := MonadChain(
			failure,
			func(x int) Validation[string] {
				return Of("should not be called")
			},
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "initial error", errors[0].Messsage)
	})

	t.Run("propagates failure from chained function", func(t *testing.T) {
		result := MonadChain(
			Of(42),
			func(x int) Validation[string] {
				return Failures[string](Errors{
					&ValidationError{Messsage: "chained error"},
				})
			},
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "chained error", errors[0].Messsage)
	})

	t.Run("chains multiple operations", func(t *testing.T) {
		step1 := MonadChain(
			Of(5),
			func(x int) Validation[int] {
				return Of(x + 10)
			},
		)
		step2 := MonadChain(
			step1,
			func(x int) Validation[int] {
				return Of(x * 2)
			},
		)
		step3 := MonadChain(
			step2,
			func(x int) Validation[string] {
				return Of(fmt.Sprintf("%d", x))
			},
		)

		assert.Equal(t, Of("30"), step3)
	})

	t.Run("conditional validation", func(t *testing.T) {
		validatePositive := func(x int) Validation[int] {
			if x > 0 {
				return Of(x)
			}
			return Failures[int](Errors{
				&ValidationError{Messsage: "must be positive"},
			})
		}

		result1 := MonadChain(Of(42), validatePositive)
		assert.Equal(t, Of(42), result1)

		result2 := MonadChain(Of(-5), validatePositive)
		assert.True(t, either.IsLeft(result2))
	})

	t.Run("dependent validation", func(t *testing.T) {
		validateRange := func(min int) func(int) Validation[int] {
			return func(max int) Validation[int] {
				if max > min {
					return Of(max)
				}
				return Failures[int](Errors{
					&ValidationError{
						Messsage: fmt.Sprintf("max (%d) must be greater than min (%d)", max, min),
					},
				})
			}
		}

		result1 := MonadChain(Of(10), validateRange(5))
		assert.Equal(t, Of(10), result1)

		result2 := MonadChain(Of(3), validateRange(5))
		assert.True(t, either.IsLeft(result2))
	})
}

// TestChain tests the curried Chain function
func TestChain(t *testing.T) {
	t.Run("creates reusable validation operator", func(t *testing.T) {
		validatePositive := Chain(func(x int) Validation[int] {
			if x > 0 {
				return Of(x)
			}
			return Failures[int](Errors{
				&ValidationError{Messsage: "must be positive"},
			})
		})

		result1 := validatePositive(Of(42))
		assert.Equal(t, Of(42), result1)

		result2 := validatePositive(Of(-5))
		assert.True(t, either.IsLeft(result2))
	})

	t.Run("use in pipeline", func(t *testing.T) {
		validatePositive := Chain(func(x int) Validation[int] {
			if x > 0 {
				return Of(x)
			}
			return Failures[int](Errors{
				&ValidationError{Messsage: "must be positive"},
			})
		})

		validateEven := Chain(func(x int) Validation[int] {
			if x%2 == 0 {
				return Of(x)
			}
			return Failures[int](Errors{
				&ValidationError{Messsage: "must be even"},
			})
		})

		result := F.Pipe2(
			Of(42),
			validatePositive,
			validateEven,
		)
		assert.Equal(t, Of(42), result)
	})

	t.Run("short-circuits on first failure in pipeline", func(t *testing.T) {
		validatePositive := Chain(func(x int) Validation[int] {
			if x > 0 {
				return Of(x)
			}
			return Failures[int](Errors{
				&ValidationError{Messsage: "must be positive"},
			})
		})

		validateEven := Chain(func(x int) Validation[int] {
			if x%2 == 0 {
				return Of(x)
			}
			return Failures[int](Errors{
				&ValidationError{Messsage: "must be even"},
			})
		})

		result := F.Pipe2(
			Of(-5),
			validatePositive,
			validateEven,
		)

		assert.True(t, either.IsLeft(result))
		errors := either.MonadFold(result,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		assert.Len(t, errors, 1)
		assert.Equal(t, "must be positive", errors[0].Messsage)
	})

	t.Run("type transformation in chain", func(t *testing.T) {
		toString := Chain(func(x int) Validation[string] {
			return Of(fmt.Sprintf("Value: %d", x))
		})

		result := F.Pipe1(Of(42), toString)

		assert.Equal(t, Of("Value: 42"), result)
	})
}

func TestOrElse(t *testing.T) {
	t.Run("OrElse is equivalent to ChainLeft - Success case", func(t *testing.T) {
		handler := func(errs Errors) Validation[int] {
			return Failures[int](Errors{
				&ValidationError{Messsage: "should not be called"},
			})
		}

		// Test with OrElse
		resultOrElse := OrElse(handler)(Success(42))
		// Test with ChainLeft
		resultChainLeft := ChainLeft(handler)(Success(42))

		assert.Equal(t, resultChainLeft, resultOrElse, "OrElse and ChainLeft should produce identical results for Success")
		assert.Equal(t, Success(42), resultOrElse)
	})

	t.Run("OrElse is equivalent to ChainLeft - Failure recovery", func(t *testing.T) {
		handler := func(errs Errors) Validation[int] {
			if len(errs) > 0 && errs[0].Messsage == "not found" {
				return Success(0)
			}
			return Failures[int](errs)
		}

		input := Failures[int](Errors{
			&ValidationError{Messsage: "not found"},
		})

		// Test with OrElse
		resultOrElse := OrElse(handler)(input)
		// Test with ChainLeft
		resultChainLeft := ChainLeft(handler)(input)

		assert.Equal(t, resultChainLeft, resultOrElse, "OrElse and ChainLeft should produce identical results for recovery")
		assert.Equal(t, Success(0), resultOrElse)
	})

	t.Run("OrElse is equivalent to ChainLeft - Error aggregation", func(t *testing.T) {
		handler := func(errs Errors) Validation[string] {
			return Failures[string](Errors{
				&ValidationError{Messsage: "additional error"},
			})
		}

		input := Failures[string](Errors{
			&ValidationError{Messsage: "original error"},
		})

		// Test with OrElse
		resultOrElse := OrElse(handler)(input)
		// Test with ChainLeft
		resultChainLeft := ChainLeft(handler)(input)

		assert.Equal(t, resultChainLeft, resultOrElse, "OrElse and ChainLeft should produce identical results for error aggregation")

		// Verify both aggregate errors
		assert.True(t, either.IsLeft(resultOrElse))
		errors := either.MonadFold(resultOrElse,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)
		assert.Len(t, errors, 2, "Should aggregate both errors")
		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "original error")
		assert.Contains(t, messages, "additional error")
	})

	t.Run("OrElse in pipeline composition", func(t *testing.T) {
		addContext := OrElse(func(errs Errors) Validation[int] {
			return Failures[int](Errors{
				&ValidationError{Messsage: "context added"},
			})
		})

		recoverFromNotFound := OrElse(func(errs Errors) Validation[int] {
			for _, err := range errs {
				if err.Messsage == "not found" {
					return Success(0)
				}
			}
			return Failures[int](errs)
		})

		// Test error aggregation in pipeline
		// When chaining OrElse operations, errors accumulate at each step
		result1 := F.Pipe2(
			Failures[int](Errors{
				&ValidationError{Messsage: "database error"},
			}),
			addContext,
			recoverFromNotFound,
		)

		assert.True(t, either.IsLeft(result1))
		errors := either.MonadFold(result1,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)
		// First OrElse adds "context added" to "database error" = 2 errors
		// Second OrElse adds those 2 errors again (via recoverFromNotFound returning errs) = 4 total
		assert.Len(t, errors, 4, "Should aggregate errors from pipeline (errors accumulate at each step)")

		// Test recovery in pipeline
		result2 := F.Pipe2(
			Failures[int](Errors{
				&ValidationError{Messsage: "not found"},
			}),
			addContext,
			recoverFromNotFound,
		)

		assert.Equal(t, Success(0), result2, "Should recover from 'not found' error")
	})

	t.Run("OrElse semantic meaning - fallback validation", func(t *testing.T) {
		// OrElse provides a semantic name for fallback/alternative validation
		// It reads naturally: "try this validation, or else try this alternative"

		validatePositive := func(x int) Validation[int] {
			if x > 0 {
				return Success(x)
			}
			return Failures[int](Errors{
				&ValidationError{Messsage: "must be positive"},
			})
		}

		// Use OrElse to provide a fallback: if validation fails, use default value
		withDefault := OrElse(func(errs Errors) Validation[int] {
			return Success(1) // default to 1 if validation fails
		})

		result := F.Pipe1(
			validatePositive(-5),
			withDefault,
		)

		assert.Equal(t, Success(1), result, "OrElse provides fallback value")
	})

	t.Run("OrElse vs ChainLeft - identical behavior verification", func(t *testing.T) {
		// Create various test scenarios
		scenarios := []struct {
			name    string
			input   Validation[int]
			handler func(Errors) Validation[int]
		}{
			{
				name:  "Success value",
				input: Success(42),
				handler: func(errs Errors) Validation[int] {
					return Failures[int](Errors{&ValidationError{Messsage: "error"}})
				},
			},
			{
				name:  "Failure with recovery",
				input: Failures[int](Errors{&ValidationError{Messsage: "error"}}),
				handler: func(errs Errors) Validation[int] {
					return Success(0)
				},
			},
			{
				name:  "Failure with error transformation",
				input: Failures[int](Errors{&ValidationError{Messsage: "error1"}}),
				handler: func(errs Errors) Validation[int] {
					return Failures[int](Errors{&ValidationError{Messsage: "error2"}})
				},
			},
			{
				name: "Multiple errors aggregation",
				input: Failures[int](Errors{
					&ValidationError{Messsage: "error1"},
					&ValidationError{Messsage: "error2"},
				}),
				handler: func(errs Errors) Validation[int] {
					return Failures[int](Errors{
						&ValidationError{Messsage: "error3"},
						&ValidationError{Messsage: "error4"},
					})
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.name, func(t *testing.T) {
				resultOrElse := OrElse(scenario.handler)(scenario.input)
				resultChainLeft := ChainLeft(scenario.handler)(scenario.input)

				assert.Equal(t, resultChainLeft, resultOrElse,
					"OrElse and ChainLeft must produce identical results for: %s", scenario.name)
			})
		}
	})
}
