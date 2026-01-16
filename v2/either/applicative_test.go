// Copyright (c) 2024 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package either

import (
	"errors"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	S "github.com/IBM/fp-go/v2/semigroup"
	"github.com/stretchr/testify/assert"
)

// TestApplicativeOf tests the Of operation of the Applicative type class
func TestApplicativeOf(t *testing.T) {
	app := Applicative[error, int, string]()

	t.Run("wraps a value in Right context", func(t *testing.T) {
		result := app.Of(42)
		assert.True(t, IsRight(result))
		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result))
	})

	t.Run("wraps string value", func(t *testing.T) {
		app := Applicative[error, string, int]()
		result := app.Of("hello")
		assert.True(t, IsRight(result))
		assert.Equal(t, "hello", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("wraps zero value", func(t *testing.T) {
		result := app.Of(0)
		assert.True(t, IsRight(result))
		assert.Equal(t, 0, GetOrElse(func(error) int { return -1 })(result))
	})

	t.Run("wraps nil pointer", func(t *testing.T) {
		app := Applicative[error, *int, *string]()
		var ptr *int = nil
		result := app.Of(ptr)
		assert.True(t, IsRight(result))
	})
}

// TestApplicativeMap tests the Map operation of the Applicative type class
func TestApplicativeMap(t *testing.T) {
	app := Applicative[error, int, int]()

	t.Run("maps a function over Right value", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		eitherValue := app.Of(21)
		result := app.Map(double)(eitherValue)
		assert.True(t, IsRight(result))
		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result))
	})

	t.Run("maps type conversion", func(t *testing.T) {
		app := Applicative[error, int, string]()
		eitherValue := app.Of(42)
		result := app.Map(strconv.Itoa)(eitherValue)
		assert.True(t, IsRight(result))
		assert.Equal(t, "42", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("maps identity function", func(t *testing.T) {
		identity := func(x int) int { return x }
		eitherValue := app.Of(42)
		result := app.Map(identity)(eitherValue)
		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result))
	})

	t.Run("preserves Left on map", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		eitherValue := Left[int](errors.New("error"))
		result := app.Map(double)(eitherValue)
		assert.True(t, IsLeft(result))
	})

	t.Run("maps with utils.Double", func(t *testing.T) {
		result := F.Pipe1(
			app.Of(21),
			app.Map(utils.Double),
		)
		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result))
	})
}

// TestApplicativeAp tests the Ap operation of the standard Applicative (fail-fast)
func TestApplicativeAp(t *testing.T) {
	app := Applicative[error, int, int]()

	t.Run("applies wrapped function to wrapped value", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}
		eitherFunc := Right[error](add(10))
		eitherValue := Right[error](32)
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsRight(result))
		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result))
	})

	t.Run("fails fast when function is Left", func(t *testing.T) {
		err1 := errors.New("function error")
		eitherFunc := Left[func(int) int](err1)
		eitherValue := Right[error](42)
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsLeft(result))
		assert.Equal(t, err1, ToError(result))
	})

	t.Run("fails fast when value is Left", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}
		err2 := errors.New("value error")
		eitherFunc := Right[error](add(10))
		eitherValue := Left[int](err2)
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsLeft(result))
		assert.Equal(t, err2, ToError(result))
	})

	t.Run("fails fast when both are Left - returns first error", func(t *testing.T) {
		err1 := errors.New("function error")
		err2 := errors.New("value error")
		eitherFunc := Left[func(int) int](err1)
		eitherValue := Left[int](err2)
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsLeft(result))
		// Should return the first error (function error)
		assert.Equal(t, err1, ToError(result))
	})

	t.Run("applies with type conversion", func(t *testing.T) {
		toStringAndAppend := func(suffix string) func(int) string {
			return func(n int) string {
				return strconv.Itoa(n) + suffix
			}
		}
		eitherFunc := Right[error](toStringAndAppend("!"))
		eitherValue := Right[error](42)
		result := Ap[string](eitherValue)(eitherFunc)
		assert.Equal(t, "42!", GetOrElse(func(error) string { return "" })(result))
	})
}

// TestApplicativeVOf tests the Of operation of ApplicativeV
func TestApplicativeVOf(t *testing.T) {
	sg := S.MakeSemigroup(func(a, b string) string { return a + "; " + b })
	app := ApplicativeV[string, int, string](sg)

	t.Run("wraps a value in Right context", func(t *testing.T) {
		result := app.Of(42)
		assert.True(t, IsRight(result))
		assert.Equal(t, 42, GetOrElse(func(string) int { return 0 })(result))
	})
}

// TestApplicativeVMap tests the Map operation of ApplicativeV
func TestApplicativeVMap(t *testing.T) {
	sg := S.MakeSemigroup(func(a, b string) string { return a + "; " + b })
	app := ApplicativeV[string, int, int](sg)

	t.Run("maps a function over Right value", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		eitherValue := app.Of(21)
		result := app.Map(double)(eitherValue)
		assert.True(t, IsRight(result))
		assert.Equal(t, 42, GetOrElse(func(string) int { return 0 })(result))
	})

	t.Run("preserves Left on map", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		eitherValue := Left[int]("error")
		result := app.Map(double)(eitherValue)
		assert.True(t, IsLeft(result))
	})
}

// TestApplicativeVAp tests the Ap operation of ApplicativeV (validation with error accumulation)
func TestApplicativeVAp(t *testing.T) {
	sg := S.MakeSemigroup(func(a, b string) string { return a + "; " + b })
	app := ApplicativeV[string, int, int](sg)

	t.Run("applies wrapped function to wrapped value", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}
		eitherFunc := Right[string](add(10))
		eitherValue := Right[string](32)
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsRight(result))
		assert.Equal(t, 42, GetOrElse(func(string) int { return 0 })(result))
	})

	t.Run("returns Left when function is Left", func(t *testing.T) {
		eitherFunc := Left[func(int) int]("function error")
		eitherValue := Right[string](42)
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsLeft(result))
		leftValue := Fold(F.Identity[string], F.Constant1[int](""))(result)
		assert.Equal(t, "function error", leftValue)
	})

	t.Run("returns Left when value is Left", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}
		eitherFunc := Right[string](add(10))
		eitherValue := Left[int]("value error")
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsLeft(result))
		leftValue := Fold(F.Identity[string], F.Constant1[int](""))(result)
		assert.Equal(t, "value error", leftValue)
	})

	t.Run("accumulates errors when both are Left", func(t *testing.T) {
		eitherFunc := Left[func(int) int]("function error")
		eitherValue := Left[int]("value error")
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsLeft(result))
		// Should combine both errors using the semigroup
		combined := Fold(F.Identity[string], F.Constant1[int](""))(result)
		assert.Equal(t, "function error; value error", combined)
	})

	t.Run("accumulates multiple validation errors", func(t *testing.T) {
		type ValidationErrors []string
		sg := S.MakeSemigroup(func(a, b ValidationErrors) ValidationErrors {
			return append(append(ValidationErrors{}, a...), b...)
		})
		app := ApplicativeV[ValidationErrors, int, int](sg)

		eitherFunc := Left[func(int) int](ValidationErrors{"error1", "error2"})
		eitherValue := Left[int](ValidationErrors{"error3", "error4"})
		result := app.Ap(eitherValue)(eitherFunc)
		assert.True(t, IsLeft(result))

		errors := Fold(F.Identity[ValidationErrors], F.Constant1[int](ValidationErrors{}))(result)
		assert.Equal(t, ValidationErrors{"error1", "error2", "error3", "error4"}, errors)
	})
}

// TestApplicativeLaws tests the applicative functor laws for standard Applicative
func TestApplicativeLaws(t *testing.T) {
	app := Applicative[error, int, int]()

	t.Run("identity law: Ap(Of(id))(v) = v", func(t *testing.T) {
		identity := func(x int) int { return x }
		v := app.Of(42)

		left := app.Ap(v)(Of[error](identity))
		right := v

		assert.Equal(t, GetOrElse(func(error) int { return 0 })(right),
			GetOrElse(func(error) int { return 0 })(left))
	})

	t.Run("homomorphism law: Ap(Of(x))(Of(f)) = Of(f(x))", func(t *testing.T) {
		f := func(x int) int { return x * 2 }
		x := 21

		left := app.Ap(app.Of(x))(Of[error](f))
		right := app.Of(f(x))

		assert.Equal(t, GetOrElse(func(error) int { return 0 })(right),
			GetOrElse(func(error) int { return 0 })(left))
	})

	t.Run("interchange law: Ap(Of(y))(u) = Ap(u)(Of(f => f(y)))", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		u := Of[error](double)
		y := 21

		left := app.Ap(app.Of(y))(u)

		// For interchange, we need to apply the value to the function
		// This test verifies the law holds for the applicative
		right := Map[error](func(f func(int) int) int { return f(y) })(u)

		assert.Equal(t, GetOrElse(func(error) int { return 0 })(right),
			GetOrElse(func(error) int { return 0 })(left))
	})

	t.Run("composition law", func(t *testing.T) {
		// For Either, we test a simpler version of composition
		f := func(x int) int { return x * 2 }
		g := func(x int) int { return x + 10 }
		x := 16

		// Apply g then f
		left := F.Pipe2(
			app.Of(x),
			app.Map(g),
			app.Map(f),
		)

		// Compose f and g, then apply
		composed := func(x int) int { return f(g(x)) }
		right := app.Map(composed)(app.Of(x))

		assert.Equal(t, GetOrElse(func(error) int { return 0 })(right),
			GetOrElse(func(error) int { return 0 })(left))
	})
}

// TestApplicativeVLaws tests the applicative functor laws for ApplicativeV
func TestApplicativeVLaws(t *testing.T) {
	sg := S.MakeSemigroup(func(a, b string) string { return a + "; " + b })
	app := ApplicativeV[string, int, int](sg)

	t.Run("identity law: Ap(Of(id))(v) = v", func(t *testing.T) {
		identity := func(x int) int { return x }
		v := app.Of(42)

		left := app.Ap(v)(Of[string](identity))
		right := v

		assert.Equal(t, GetOrElse(func(string) int { return 0 })(right),
			GetOrElse(func(string) int { return 0 })(left))
	})

	t.Run("homomorphism law: Ap(Of(x))(Of(f)) = Of(f(x))", func(t *testing.T) {
		f := func(x int) int { return x * 2 }
		x := 21

		left := app.Ap(app.Of(x))(Of[string](f))
		right := app.Of(f(x))

		assert.Equal(t, GetOrElse(func(string) int { return 0 })(right),
			GetOrElse(func(string) int { return 0 })(left))
	})

	t.Run("interchange law: Ap(Of(y))(u) = Ap(u)(Of(f => f(y)))", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		u := Of[string](double)
		y := 21

		left := app.Ap(app.Of(y))(u)

		// For interchange, we need to apply the value to the function
		right := Map[string](func(f func(int) int) int { return f(y) })(u)

		assert.Equal(t, GetOrElse(func(string) int { return 0 })(right),
			GetOrElse(func(string) int { return 0 })(left))
	})
}

// TestApplicativeComposition tests composition of applicative operations
func TestApplicativeComposition(t *testing.T) {
	app := Applicative[error, int, int]()

	t.Run("composes Map and Of", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		result := F.Pipe1(
			app.Of(21),
			app.Map(double),
		)
		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result))
	})

	t.Run("composes multiple Map operations", func(t *testing.T) {
		app := Applicative[error, int, string]()
		double := func(x int) int { return x * 2 }
		toString := func(x int) string { return strconv.Itoa(x) }

		result := F.Pipe2(
			app.Of(21),
			Map[error](double),
			app.Map(toString),
		)
		assert.Equal(t, "42", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("composes Map and Ap", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		eitherFunc := F.Pipe1(
			app.Of(5),
			Map[error](add),
		)
		eitherValue := app.Of(16)

		result := app.Ap(eitherValue)(eitherFunc)
		assert.Equal(t, 21, GetOrElse(func(error) int { return 0 })(result))
	})
}

// TestApplicativeMultipleArguments tests applying functions with multiple arguments
func TestApplicativeMultipleArguments(t *testing.T) {
	app := Applicative[error, int, int]()

	t.Run("applies curried two-argument function", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		eitherFunc := F.Pipe1(
			app.Of(10),
			Map[error](add),
		)

		result := app.Ap(app.Of(32))(eitherFunc)
		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result))
	})

	t.Run("applies curried three-argument function", func(t *testing.T) {
		add3 := func(a int) func(int) func(int) int {
			return func(b int) func(int) int {
				return func(c int) int {
					return a + b + c
				}
			}
		}

		eitherFunc1 := F.Pipe1(
			app.Of(10),
			Map[error](add3),
		)

		eitherFunc2 := Ap[func(int) int](app.Of(20))(eitherFunc1)
		result := Ap[int](app.Of(12))(eitherFunc2)

		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result))
	})
}

// TestApplicativeInstance tests that Applicative returns a valid instance
func TestApplicativeInstance(t *testing.T) {
	t.Run("returns non-nil instance", func(t *testing.T) {
		app := Applicative[error, int, string]()
		assert.NotNil(t, app)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		app1 := Applicative[error, int, string]()
		app2 := Applicative[error, int, string]()

		result1 := app1.Of(42)
		result2 := app2.Of(43)

		assert.Equal(t, 42, GetOrElse(func(error) int { return 0 })(result1))
		assert.Equal(t, 43, GetOrElse(func(error) int { return 0 })(result2))
	})
}

// TestApplicativeVInstance tests that ApplicativeV returns a valid instance
func TestApplicativeVInstance(t *testing.T) {
	sg := S.MakeSemigroup(func(a, b string) string { return a + "; " + b })

	t.Run("returns non-nil instance", func(t *testing.T) {
		app := ApplicativeV[string, int, string](sg)
		assert.NotNil(t, app)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		app1 := ApplicativeV[string, int, string](sg)
		app2 := ApplicativeV[string, int, string](sg)

		result1 := app1.Of(42)
		result2 := app2.Of(43)

		assert.Equal(t, 42, GetOrElse(func(string) int { return 0 })(result1))
		assert.Equal(t, 43, GetOrElse(func(string) int { return 0 })(result2))
	})
}

// TestApplicativeWithDifferentTypes tests applicative with various type combinations
func TestApplicativeWithDifferentTypes(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		app := Applicative[error, int, string]()
		result := app.Map(strconv.Itoa)(app.Of(42))
		assert.Equal(t, "42", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("string to int", func(t *testing.T) {
		app := Applicative[error, string, int]()
		toLength := func(s string) int { return len(s) }
		result := app.Map(toLength)(app.Of("hello"))
		assert.Equal(t, 5, GetOrElse(func(error) int { return 0 })(result))
	})

	t.Run("bool to string", func(t *testing.T) {
		app := Applicative[error, bool, string]()
		toString := func(b bool) string {
			if b {
				return "true"
			}
			return "false"
		}
		result := app.Map(toString)(app.Of(true))
		assert.Equal(t, "true", GetOrElse(func(error) string { return "" })(result))
	})
}

// TestApplicativeVFormValidationExample demonstrates a realistic form validation scenario
func TestApplicativeVFormValidationExample(t *testing.T) {
	type ValidationErrors []string

	sg := S.MakeSemigroup(func(a, b ValidationErrors) ValidationErrors {
		return append(append(ValidationErrors{}, a...), b...)
	})

	validateName := func(name string) Either[ValidationErrors, string] {
		if len(name) < 3 {
			return Left[string](ValidationErrors{"Name must be at least 3 characters"})
		}
		return Right[ValidationErrors](name)
	}

	validateAge := func(age int) Either[ValidationErrors, int] {
		if age < 18 {
			return Left[int](ValidationErrors{"Must be 18 or older"})
		}
		return Right[ValidationErrors](age)
	}

	validateEmail := func(email string) Either[ValidationErrors, string] {
		if len(email) == 0 {
			return Left[string](ValidationErrors{"Email is required"})
		}
		return Right[ValidationErrors](email)
	}

	t.Run("all validations pass", func(t *testing.T) {
		name := validateName("Alice")
		age := validateAge(25)
		email := validateEmail("alice@example.com")

		// Verify all individual validations passed
		assert.True(t, IsRight(name))
		assert.True(t, IsRight(age))
		assert.True(t, IsRight(email))

		// Combine validations - all pass
		result := F.Pipe2(
			name,
			Map[ValidationErrors](func(n string) string { return n }),
			Map[ValidationErrors](func(n string) string { return n + " validated" }),
		)

		assert.True(t, IsRight(result))
		value := GetOrElse(func(ValidationErrors) string { return "" })(result)
		assert.Equal(t, "Alice validated", value)
	})

	t.Run("all validations fail - accumulates all errors", func(t *testing.T) {
		name := validateName("ab")
		age := validateAge(16)
		email := validateEmail("")

		// Manually combine errors using the semigroup
		var allErrors ValidationErrors
		if IsLeft(name) {
			allErrors = Fold(F.Identity[ValidationErrors], F.Constant1[string](ValidationErrors{}))(name)
		}
		if IsLeft(age) {
			ageErrors := Fold(F.Identity[ValidationErrors], F.Constant1[int](ValidationErrors{}))(age)
			allErrors = sg.Concat(allErrors, ageErrors)
		}
		if IsLeft(email) {
			emailErrors := Fold(F.Identity[ValidationErrors], F.Constant1[string](ValidationErrors{}))(email)
			allErrors = sg.Concat(allErrors, emailErrors)
		}

		assert.Len(t, allErrors, 3)
		assert.Contains(t, allErrors, "Name must be at least 3 characters")
		assert.Contains(t, allErrors, "Must be 18 or older")
		assert.Contains(t, allErrors, "Email is required")
	})

	t.Run("partial validation failure", func(t *testing.T) {
		name := validateName("Alice")
		age := validateAge(16)
		email := validateEmail("")

		// Verify name passes
		assert.True(t, IsRight(name))

		// Manually combine errors using the semigroup
		var allErrors ValidationErrors
		if IsLeft(age) {
			allErrors = Fold(F.Identity[ValidationErrors], F.Constant1[int](ValidationErrors{}))(age)
		}
		if IsLeft(email) {
			emailErrors := Fold(F.Identity[ValidationErrors], F.Constant1[string](ValidationErrors{}))(email)
			if len(allErrors) > 0 {
				allErrors = sg.Concat(allErrors, emailErrors)
			} else {
				allErrors = emailErrors
			}
		}

		assert.Len(t, allErrors, 2)
		assert.Contains(t, allErrors, "Must be 18 or older")
		assert.Contains(t, allErrors, "Email is required")
	})
}

// TestApplicativeVsApplicativeV demonstrates the difference between fail-fast and validation
func TestApplicativeVsApplicativeV(t *testing.T) {
	t.Run("Applicative fails fast", func(t *testing.T) {
		app := Applicative[error, int, int]()

		err1 := errors.New("error1")
		err2 := errors.New("error2")

		eitherFunc := Left[func(int) int](err1)
		eitherValue := Left[int](err2)

		result := app.Ap(eitherValue)(eitherFunc)

		assert.True(t, IsLeft(result))
		// Only the first error is returned
		assert.Equal(t, err1, ToError(result))
	})

	t.Run("ApplicativeV accumulates errors", func(t *testing.T) {
		sg := S.MakeSemigroup(func(a, b string) string { return a + "; " + b })
		app := ApplicativeV[string, int, int](sg)

		eitherFunc := Left[func(int) int]("error1")
		eitherValue := Left[int]("error2")

		result := app.Ap(eitherValue)(eitherFunc)

		assert.True(t, IsLeft(result))
		// Both errors are accumulated
		combined := Fold(F.Identity[string], F.Constant1[int](""))(result)
		assert.Equal(t, "error1; error2", combined)
	})
}
