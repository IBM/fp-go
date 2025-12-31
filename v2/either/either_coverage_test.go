// Copyright (c) 2023 - 2025 IBM Corp.
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
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Test BiMap
func TestBiMap(t *testing.T) {
	errToStr := error.Error
	intToStr := strconv.Itoa

	// Test Right case
	result := BiMap(errToStr, intToStr)(Right[error](42))
	assert.Equal(t, Right[string]("42"), result)

	// Test Left case
	result = BiMap(errToStr, intToStr)(Left[int](errors.New("error")))
	assert.Equal(t, Left[string]("error"), result)
}

// Test MonadBiMap
func TestMonadBiMap(t *testing.T) {
	errToStr := error.Error
	intToStr := strconv.Itoa

	result := MonadBiMap(Right[error](42), errToStr, intToStr)
	assert.Equal(t, Right[string]("42"), result)

	result = MonadBiMap(Left[int](errors.New("error")), errToStr, intToStr)
	assert.Equal(t, Left[string]("error"), result)
}

// Test MapLeft
func TestMapLeft(t *testing.T) {
	errToStr := error.Error

	result := MapLeft[int](errToStr)(Left[int](errors.New("error")))
	assert.Equal(t, Left[int]("error"), result)

	result = MapLeft[int](errToStr)(Right[error](42))
	assert.Equal(t, Right[string](42), result)
}

// Test MonadMapLeft
func TestMonadMapLeft(t *testing.T) {
	errToStr := error.Error

	result := MonadMapLeft(Left[int](errors.New("error")), errToStr)
	assert.Equal(t, Left[int]("error"), result)

	result = MonadMapLeft(Right[error](42), errToStr)
	assert.Equal(t, Right[string](42), result)
}

// Test MonadMapTo
func TestMonadMapTo(t *testing.T) {
	result := MonadMapTo(Right[error](42), "success")
	assert.Equal(t, Right[error]("success"), result)

	result = MonadMapTo(Left[int](errors.New("error")), "success")
	assert.Equal(t, Left[string](errors.New("error")), result)
}

// Test MonadChainFirst
func TestMonadChainFirst(t *testing.T) {
	f := func(x int) Either[error, string] {
		return Right[error](strconv.Itoa(x))
	}

	result := MonadChainFirst(Right[error](42), f)
	assert.Equal(t, Right[error](42), result)

	result = MonadChainFirst(Left[int](errors.New("error")), f)
	assert.Equal(t, Left[int](errors.New("error")), result)
}

// Test MonadChainTo
func TestMonadChainTo(t *testing.T) {
	result := MonadChainTo(Right[error](42), Right[error]("hello"))
	assert.Equal(t, Right[error]("hello"), result)

	result = MonadChainTo(Left[int](errors.New("error")), Right[error]("hello"))
	assert.Equal(t, Right[error]("hello"), result)
}

// Test Flatten
func TestFlatten(t *testing.T) {
	nested := Right[error](Right[error](42))
	result := Flatten(nested)
	assert.Equal(t, Right[error](42), result)

	nestedLeft := Right[error](Left[int](errors.New("error")))
	result = Flatten(nestedLeft)
	assert.Equal(t, Left[int](errors.New("error")), result)

	outerLeft := Left[Either[error, int]](errors.New("outer error"))
	result = Flatten(outerLeft)
	assert.Equal(t, Left[int](errors.New("outer error")), result)
}

// Test ToOption
func TestToOption(t *testing.T) {
	result := ToOption(Right[error](42))
	assert.Equal(t, O.Some(42), result)

	result = ToOption(Left[int](errors.New("error")))
	assert.Equal(t, O.None[int](), result)
}

// Test FromError
func TestFromError(t *testing.T) {
	validate := func(x int) error {
		if x < 0 {
			return errors.New("negative")
		}
		return nil
	}

	toEither := FromError(validate)
	result := toEither(42)
	assert.Equal(t, Right[error](42), result)

	result = toEither(-1)
	assert.True(t, IsLeft(result))
}

// Test ToError
func TestToError(t *testing.T) {
	err := ToError(Left[int](errors.New("error")))
	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())

	err = ToError(Right[error](42))
	assert.NoError(t, err)
}

// Test OrElse
func TestOrElse(t *testing.T) {
	// Test basic recovery from Left
	recover := OrElse(func(e error) Either[error, int] {
		return Right[error](0)
	})

	result := recover(Left[int](errors.New("error")))
	assert.Equal(t, Right[error](0), result)

	// Test Right value passes through unchanged
	result = recover(Right[error](42))
	assert.Equal(t, Right[error](42), result)

	// Test selective recovery - recover some errors, propagate others
	selectiveRecover := OrElse(func(err error) Either[error, int] {
		if err.Error() == "not found" {
			return Right[error](0) // default value for "not found"
		}
		return Left[int](err) // propagate other errors
	})
	assert.Equal(t, Right[error](0), selectiveRecover(Left[int](errors.New("not found"))))
	permissionErr := errors.New("permission denied")
	assert.Equal(t, Left[int](permissionErr), selectiveRecover(Left[int](permissionErr)))

	// Test chaining multiple OrElse operations
	firstRecover := OrElse(func(err error) Either[error, int] {
		if err.Error() == "error1" {
			return Right[error](1)
		}
		return Left[int](err)
	})
	secondRecover := OrElse(func(err error) Either[error, int] {
		if err.Error() == "error2" {
			return Right[error](2)
		}
		return Left[int](err)
	})
	assert.Equal(t, Right[error](1), F.Pipe1(Left[int](errors.New("error1")), firstRecover))
	assert.Equal(t, Right[error](2), F.Pipe1(Left[int](errors.New("error2")), F.Flow2(firstRecover, secondRecover)))
}

// Test OrElseW
func TestOrElseW(t *testing.T) {
	type ValidationError string
	type AppError int

	// Test with Right value - should return Right with widened error type
	rightValue := Right[ValidationError]("success")
	recoverValidation := OrElse(func(ve ValidationError) Either[AppError, string] {
		return Left[string](AppError(400))
	})
	result := recoverValidation(rightValue)
	assert.True(t, IsRight(result))
	assert.Equal(t, "success", F.Pipe1(result, GetOrElse(F.Constant1[AppError](""))))

	// Test with Left value - should apply recovery with new error type
	leftValue := Left[string](ValidationError("invalid input"))
	result = recoverValidation(leftValue)
	assert.True(t, IsLeft(result))
	_, leftVal := Unwrap(result)
	assert.Equal(t, AppError(400), leftVal)

	// Test error type conversion - ValidationError to AppError
	convertError := OrElse(func(ve ValidationError) Either[AppError, int] {
		return Left[int](AppError(len(ve)))
	})
	converted := convertError(Left[int](ValidationError("short")))
	assert.True(t, IsLeft(converted))
	_, leftConv := Unwrap(converted)
	assert.Equal(t, AppError(5), leftConv)

	// Test recovery to Right with widened error type
	recoverToRight := OrElse(func(ve ValidationError) Either[AppError, int] {
		if ve == "recoverable" {
			return Right[AppError](99)
		}
		return Left[int](AppError(500))
	})
	assert.Equal(t, Right[AppError](99), recoverToRight(Left[int](ValidationError("recoverable"))))
	assert.True(t, IsLeft(recoverToRight(Left[int](ValidationError("fatal")))))

	// Test that Right values are preserved with widened error type
	preservedRight := Right[ValidationError](42)
	preserveRecover := OrElse(func(ve ValidationError) Either[AppError, int] {
		return Left[int](AppError(999))
	})
	preserved := preserveRecover(preservedRight)
	assert.Equal(t, Right[AppError](42), preserved)
}

// Test ToType
func TestToType(t *testing.T) {
	convert := ToType[int](func(v any) error {
		return fmt.Errorf("cannot convert %v to int", v)
	})

	result := convert(42)
	assert.Equal(t, Right[error](42), result)

	result = convert("string")
	assert.True(t, IsLeft(result))
}

// Test Memoize
func TestMemoize(t *testing.T) {
	val := Right[error](42)
	result := Memoize(val)
	assert.Equal(t, val, result)
}

// Test Swap
func TestSwap(t *testing.T) {
	result := Swap(Right[error](42))
	assert.Equal(t, Left[error](42), result)

	result = Swap(Left[int](errors.New("error")))
	assert.Equal(t, Right[int](errors.New("error")), result)
}

// Test MonadFlap and Flap
func TestFlap(t *testing.T) {
	fab := Right[error](strconv.Itoa)
	result := MonadFlap(fab, 42)
	assert.Equal(t, Right[error]("42"), result)

	result = Flap[error, string](42)(fab)
	assert.Equal(t, Right[error]("42"), result)

	fabLeft := Left[func(int) string](errors.New("error"))
	result = MonadFlap(fabLeft, 42)
	assert.Equal(t, Left[string](errors.New("error")), result)
}

// Test Sequence2 and MonadSequence2
func TestSequence2(t *testing.T) {
	f := func(a int, b int) Either[error, int] {
		return Right[error](a + b)
	}

	result := Sequence2(f)(Right[error](1), Right[error](2))
	assert.Equal(t, Right[error](3), result)

	result = Sequence2(f)(Left[int](errors.New("error")), Right[error](2))
	assert.Equal(t, Left[int](errors.New("error")), result)

	result = MonadSequence2(Right[error](1), Right[error](2), f)
	assert.Equal(t, Right[error](3), result)
}

// Test Sequence3 and MonadSequence3
func TestSequence3(t *testing.T) {
	f := func(a, b, c int) Either[error, int] {
		return Right[error](a + b + c)
	}

	result := Sequence3(f)(Right[error](1), Right[error](2), Right[error](3))
	assert.Equal(t, Right[error](6), result)

	result = Sequence3(f)(Left[int](errors.New("error")), Right[error](2), Right[error](3))
	assert.Equal(t, Left[int](errors.New("error")), result)

	result = MonadSequence3(Right[error](1), Right[error](2), Right[error](3), f)
	assert.Equal(t, Right[error](6), result)
}

// Test Let
func TestLet(t *testing.T) {
	type State struct{ value int }
	result := F.Pipe2(
		Right[error](State{value: 10}),
		Let[error](
			func(v int) func(State) State {
				return func(s State) State { return State{value: s.value + v} }
			},
			func(s State) int { return 32 },
		),
		Map[error](F.Identity[State]),
	)
	assert.Equal(t, Right[error](State{value: 42}), result)
}

// Test LetTo
func TestLetTo(t *testing.T) {
	type State struct{ name string }
	result := F.Pipe2(
		Right[error](State{}),
		LetTo[error](
			func(n string) func(State) State {
				return func(s State) State { return State{name: n} }
			},
			"Alice",
		),
		Map[error](F.Identity[State]),
	)
	assert.Equal(t, Right[error](State{name: "Alice"}), result)
}

// Test BindTo
func TestBindTo(t *testing.T) {
	type State struct{ value int }
	result := F.Pipe2(
		Right[error](42),
		BindTo[error](func(v int) State { return State{value: v} }),
		Map[error](F.Identity[State]),
	)
	assert.Equal(t, Right[error](State{value: 42}), result)
}

// Test TraverseArray
func TestTraverseArray(t *testing.T) {
	parse := func(s string) Either[error, int] {
		v, err := strconv.Atoi(s)
		return TryCatchError(v, err)
	}

	result := TraverseArray(parse)([]string{"1", "2", "3"})
	assert.Equal(t, Right[error]([]int{1, 2, 3}), result)

	result = TraverseArray(parse)([]string{"1", "bad", "3"})
	assert.True(t, IsLeft(result))
}

// Test TraverseArrayWithIndex
func TestTraverseArrayWithIndex(t *testing.T) {
	validate := func(i int, s string) Either[error, string] {
		if S.IsNonEmpty(s) {
			return Right[error](fmt.Sprintf("%d:%s", i, s))
		}
		return Left[string](fmt.Errorf("empty at index %d", i))
	}

	result := TraverseArrayWithIndex(validate)([]string{"a", "b"})
	assert.Equal(t, Right[error]([]string{"0:a", "1:b"}), result)

	result = TraverseArrayWithIndex(validate)([]string{"a", ""})
	assert.True(t, IsLeft(result))
}

// Test TraverseRecord
func TestTraverseRecord(t *testing.T) {
	parse := func(s string) Either[error, int] {
		v, err := strconv.Atoi(s)
		return TryCatchError(v, err)
	}

	input := map[string]string{"a": "1", "b": "2"}
	result := TraverseRecord[string](parse)(input)
	expected := Right[error](map[string]int{"a": 1, "b": 2})
	assert.Equal(t, expected, result)
}

// Test TraverseRecordWithIndex
func TestTraverseRecordWithIndex(t *testing.T) {
	validate := func(k string, v string) Either[error, string] {
		if S.IsNonEmpty(v) {
			return Right[error](k + ":" + v)
		}
		return Left[string](fmt.Errorf("empty value for key %s", k))
	}

	input := map[string]string{"a": "1"}
	result := TraverseRecordWithIndex(validate)(input)
	expected := Right[error](map[string]string{"a": "a:1"})
	assert.Equal(t, expected, result)
}

// Test SequenceRecord
func TestSequenceRecord(t *testing.T) {
	eithers := map[string]Either[error, int]{
		"a": Right[error](1),
		"b": Right[error](2),
	}
	result := SequenceRecord(eithers)
	expected := Right[error](map[string]int{"a": 1, "b": 2})
	assert.Equal(t, expected, result)

	eithersWithError := map[string]Either[error, int]{
		"a": Right[error](1),
		"b": Left[int](errors.New("error")),
	}
	result = SequenceRecord(eithersWithError)
	assert.True(t, IsLeft(result))
}

// Test Curry functions
func TestCurry0(t *testing.T) {
	getConfig := func() (string, error) { return "config", nil }
	curried := Curry0(getConfig)
	result := curried()
	assert.Equal(t, Right[error]("config"), result)
}

func TestCurry1(t *testing.T) {
	parse := strconv.Atoi
	curried := Curry1(parse)
	result := curried("42")
	assert.Equal(t, Right[error](42), result)

	result = curried("bad")
	assert.True(t, IsLeft(result))
}

func TestCurry2(t *testing.T) {
	divide := func(a, b int) (int, error) {
		if b == 0 {
			return 0, errors.New("div by zero")
		}
		return a / b, nil
	}
	curried := Curry2(divide)
	result := curried(10)(2)
	assert.Equal(t, Right[error](5), result)

	result = curried(10)(0)
	assert.True(t, IsLeft(result))
}

func TestCurry3(t *testing.T) {
	sum3 := func(a, b, c int) (int, error) {
		return a + b + c, nil
	}
	curried := Curry3(sum3)
	result := curried(1)(2)(3)
	assert.Equal(t, Right[error](6), result)
}

func TestCurry4(t *testing.T) {
	sum4 := func(a, b, c, d int) (int, error) {
		return a + b + c + d, nil
	}
	curried := Curry4(sum4)
	result := curried(1)(2)(3)(4)
	assert.Equal(t, Right[error](10), result)
}

// Test Uncurry functions
func TestUncurry0(t *testing.T) {
	curried := func() Either[error, string] { return Right[error]("value") }
	uncurried := Uncurry0(curried)
	result, err := uncurried()
	assert.NoError(t, err)
	assert.Equal(t, "value", result)
}

func TestUncurry1(t *testing.T) {
	curried := func(x int) Either[error, string] { return Right[error](strconv.Itoa(x)) }
	uncurried := Uncurry1(curried)
	result, err := uncurried(42)
	assert.NoError(t, err)
	assert.Equal(t, "42", result)
}

func TestUncurry2(t *testing.T) {
	curried := func(a int) func(int) Either[error, int] {
		return func(b int) Either[error, int] {
			return Right[error](a + b)
		}
	}
	uncurried := Uncurry2(curried)
	result, err := uncurried(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 3, result)
}

func TestUncurry3(t *testing.T) {
	curried := func(a int) func(int) func(int) Either[error, int] {
		return func(b int) func(int) Either[error, int] {
			return func(c int) Either[error, int] {
				return Right[error](a + b + c)
			}
		}
	}
	uncurried := Uncurry3(curried)
	result, err := uncurried(1, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, 6, result)
}

func TestUncurry4(t *testing.T) {
	curried := func(a int) func(int) func(int) func(int) Either[error, int] {
		return func(b int) func(int) func(int) Either[error, int] {
			return func(c int) func(int) Either[error, int] {
				return func(d int) Either[error, int] {
					return Right[error](a + b + c + d)
				}
			}
		}
	}
	uncurried := Uncurry4(curried)
	result, err := uncurried(1, 2, 3, 4)
	assert.NoError(t, err)
	assert.Equal(t, 10, result)
}

// Test Variadic functions
func TestVariadic0(t *testing.T) {
	sum := func(nums []int) (int, error) {
		total := 0
		for _, n := range nums {
			total += n
		}
		return total, nil
	}
	variadicSum := Variadic0(sum)
	result := variadicSum(1, 2, 3)
	assert.Equal(t, Right[error](6), result)
}

func TestVariadic1(t *testing.T) {
	multiply := func(factor int, nums []int) ([]int, error) {
		result := make([]int, len(nums))
		for i, n := range nums {
			result[i] = n * factor
		}
		return result, nil
	}
	variadicMultiply := Variadic1(multiply)
	result := variadicMultiply(2, 1, 2, 3)
	assert.Equal(t, Right[error]([]int{2, 4, 6}), result)
}

func TestVariadic2(t *testing.T) {
	combine := func(a, b int, nums []int) (int, error) {
		total := a + b
		for _, n := range nums {
			total += n
		}
		return total, nil
	}
	variadicCombine := Variadic2(combine)
	result := variadicCombine(1, 2, 3, 4)
	assert.Equal(t, Right[error](10), result)
}

func TestVariadic3(t *testing.T) {
	combine := func(a, b, c int, nums []int) (int, error) {
		total := a + b + c
		for _, n := range nums {
			total += n
		}
		return total, nil
	}
	variadicCombine := Variadic3(combine)
	result := variadicCombine(1, 2, 3, 4, 5)
	assert.Equal(t, Right[error](15), result)
}

func TestVariadic4(t *testing.T) {
	combine := func(a, b, c, d int, nums []int) (int, error) {
		total := a + b + c + d
		for _, n := range nums {
			total += n
		}
		return total, nil
	}
	variadicCombine := Variadic4(combine)
	result := variadicCombine(1, 2, 3, 4, 5, 6)
	assert.Equal(t, Right[error](21), result)
}

// Test Unvariadic functions
func TestUnvariadic0(t *testing.T) {
	variadic := func(nums ...int) (int, error) {
		total := 0
		for _, n := range nums {
			total += n
		}
		return total, nil
	}
	unvariadic := Unvariadic0(variadic)
	result := unvariadic([]int{1, 2, 3})
	assert.Equal(t, Right[error](6), result)
}

func TestUnvariadic1(t *testing.T) {
	variadic := func(factor int, nums ...int) ([]int, error) {
		result := make([]int, len(nums))
		for i, n := range nums {
			result[i] = n * factor
		}
		return result, nil
	}
	unvariadic := Unvariadic1(variadic)
	result := unvariadic(2, []int{1, 2, 3})
	assert.Equal(t, Right[error]([]int{2, 4, 6}), result)
}

func TestUnvariadic2(t *testing.T) {
	variadic := func(a, b int, nums ...int) (int, error) {
		total := a + b
		for _, n := range nums {
			total += n
		}
		return total, nil
	}
	unvariadic := Unvariadic2(variadic)
	result := unvariadic(1, 2, []int{3, 4})
	assert.Equal(t, Right[error](10), result)
}

func TestUnvariadic3(t *testing.T) {
	variadic := func(a, b, c int, nums ...int) (int, error) {
		total := a + b + c
		for _, n := range nums {
			total += n
		}
		return total, nil
	}
	unvariadic := Unvariadic3(variadic)
	result := unvariadic(1, 2, 3, []int{4, 5})
	assert.Equal(t, Right[error](15), result)
}

func TestUnvariadic4(t *testing.T) {
	variadic := func(a, b, c, d int, nums ...int) (int, error) {
		total := a + b + c + d
		for _, n := range nums {
			total += n
		}
		return total, nil
	}
	unvariadic := Unvariadic4(variadic)
	result := unvariadic(1, 2, 3, 4, []int{5, 6})
	assert.Equal(t, Right[error](21), result)
}

// Test Monad
func TestMonad(t *testing.T) {
	m := Monad[error, int, string]()

	// Test Of
	result := m.Of(42)
	assert.Equal(t, Right[error](42), result)

	// Test Map
	mapFn := m.Map(strconv.Itoa)
	result2 := mapFn(Right[error](42))
	assert.Equal(t, Right[error]("42"), result2)

	// Test Chain
	chainFn := m.Chain(func(x int) Either[error, string] {
		return Right[error](strconv.Itoa(x))
	})
	result3 := chainFn(Right[error](42))
	assert.Equal(t, Right[error]("42"), result3)

	// Test Ap
	apFn := m.Ap(Right[error](42))
	result4 := apFn(Right[error](strconv.Itoa))
	assert.Equal(t, Right[error]("42"), result4)
}

// Test AltSemigroup
func TestAltSemigroup(t *testing.T) {
	sg := AltSemigroup[error, int]()

	result := sg.Concat(Left[int](errors.New("error")), Right[error](42))
	assert.Equal(t, Right[error](42), result)

	result = sg.Concat(Right[error](1), Right[error](2))
	assert.Equal(t, Right[error](1), result)
}

// Test AlternativeMonoid
func TestAlternativeMonoid(t *testing.T) {
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[error](intAdd)

	result := m.Concat(Right[error](1), Right[error](2))
	assert.Equal(t, Right[error](3), result)

	empty := m.Empty()
	assert.Equal(t, Right[error](0), empty)
}

// Test AltMonoid
func TestAltMonoid(t *testing.T) {
	zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
	m := AltMonoid(zero)

	result := m.Concat(Left[int](errors.New("err1")), Right[error](42))
	assert.Equal(t, Right[error](42), result)

	empty := m.Empty()
	assert.Equal(t, Left[int](errors.New("empty")), empty)
}

// Test core.go Format
func TestFormat(t *testing.T) {
	e := Right[error](42)
	formatted := fmt.Sprintf("%s", e)
	assert.Contains(t, formatted, "Right")
	assert.Contains(t, formatted, "42")

	e2 := Left[int](errors.New("error"))
	formatted2 := fmt.Sprintf("%v", e2)
	assert.Contains(t, formatted2, "Left")
}

// Test TryCatch with error
func TestTryCatchWithError(t *testing.T) {
	result := TryCatch(0, errors.New("error"), func(err error) string {
		return err.Error()
	})
	assert.Equal(t, Left[int]("error"), result)
}

// Test Unwrap with else branch
func TestUnwrapElseBranch(t *testing.T) {
	val, err := Unwrap(Right[error](42))
	assert.Equal(t, 42, val)
	assert.Equal(t, error(nil), err)
}

// Test eitherFormat with different rune
func TestEitherFormatDifferentRune(t *testing.T) {
	e := Right[error](42)
	formatted := fmt.Sprintf("%v", e)
	assert.Contains(t, formatted, "Right")
}
