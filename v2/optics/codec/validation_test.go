package codec

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestIsWithPrimitiveTypes tests the Is function with primitive types
func TestIsWithPrimitiveTypes(t *testing.T) {
	t.Run("string type succeeds with string value", func(t *testing.T) {
		isString := Is[string]()
		res := isString("hello")

		assert.Equal(t, R.Of("hello"), res)
	})

	t.Run("string type fails with int value", func(t *testing.T) {
		isString := Is[string]()
		res := isString(42)

		assert.True(t, either.IsLeft(res), "Expected Left for invalid type")
	})

	t.Run("int type succeeds with int value", func(t *testing.T) {
		isInt := Is[int]()
		res := isInt(42)

		assert.Equal(t, R.Of(42), res)
	})

	t.Run("int type fails with string value", func(t *testing.T) {
		isInt := Is[int]()
		res := isInt("42")

		assert.True(t, either.IsLeft(res))
	})

	t.Run("bool type succeeds with bool value", func(t *testing.T) {
		isBool := Is[bool]()
		res := isBool(true)

		assert.Equal(t, R.Of(true), res)
	})

	t.Run("bool type fails with int value", func(t *testing.T) {
		isBool := Is[bool]()
		res := isBool(1)

		assert.True(t, either.IsLeft(res))
	})

	t.Run("float64 type succeeds with float64 value", func(t *testing.T) {
		isFloat := Is[float64]()
		res := isFloat(3.14)

		assert.Equal(t, R.Of(3.14), res)
	})

	t.Run("float64 type fails with int value", func(t *testing.T) {
		isFloat := Is[float64]()
		res := isFloat(42)

		assert.True(t, either.IsLeft(res))
	})
}

// TestIsWithNumericTypes tests Is with different numeric types
func TestIsWithNumericTypes(t *testing.T) {
	t.Run("int8 type", func(t *testing.T) {
		isInt8 := Is[int8]()

		res := isInt8(int8(127))
		assert.Equal(t, R.Of(int8(127)), res)

		// Fails with regular int
		res = isInt8(127)
		assert.True(t, either.IsLeft(res))
	})

	t.Run("int16 type", func(t *testing.T) {
		isInt16 := Is[int16]()

		res := isInt16(int16(32767))
		assert.Equal(t, R.Of(int16(32767)), res)
	})

	t.Run("int32 type", func(t *testing.T) {
		isInt32 := Is[int32]()

		res := isInt32(int32(2147483647))
		assert.Equal(t, R.Of(int32(2147483647)), res)
	})

	t.Run("int64 type", func(t *testing.T) {
		isInt64 := Is[int64]()

		res := isInt64(int64(9223372036854775807))
		assert.Equal(t, R.Of(int64(9223372036854775807)), res)
	})

	t.Run("uint type", func(t *testing.T) {
		isUint := Is[uint]()

		res := isUint(uint(42))
		assert.Equal(t, R.Of(uint(42)), res)

		// Fails with int
		res = isUint(42)
		assert.True(t, either.IsLeft(res))
	})

	t.Run("float32 type", func(t *testing.T) {
		isFloat32 := Is[float32]()

		res := isFloat32(float32(3.14))
		assert.Equal(t, R.Of(float32(3.14)), res)

		// Fails with float64
		res = isFloat32(3.14)
		assert.True(t, either.IsLeft(res))
	})
}

// TestIsWithComplexTypes tests Is with complex and composite types
func TestIsWithComplexTypes(t *testing.T) {
	t.Run("slice type succeeds with slice", func(t *testing.T) {
		isSlice := Is[[]int]()
		res := isSlice([]int{1, 2, 3})

		assert.Equal(t, R.Of([]int{1, 2, 3}), res)
	})

	t.Run("slice type fails with array", func(t *testing.T) {
		isSlice := Is[[]int]()
		res := isSlice([3]int{1, 2, 3})

		assert.True(t, either.IsLeft(res))
	})

	t.Run("map type succeeds with map", func(t *testing.T) {
		isMap := Is[map[string]int]()
		testMap := map[string]int{"a": 1, "b": 2}
		res := isMap(testMap)

		assert.Equal(t, R.Of(testMap), res)
	})

	t.Run("map type fails with wrong key type", func(t *testing.T) {
		isMap := Is[map[string]int]()
		wrongMap := map[int]int{1: 1, 2: 2}
		res := isMap(wrongMap)

		assert.True(t, either.IsLeft(res))
	})

	t.Run("array type succeeds with array", func(t *testing.T) {
		isArray := Is[[3]int]()
		res := isArray([3]int{1, 2, 3})

		assert.Equal(t, R.Of([3]int{1, 2, 3}), res)
	})

	t.Run("array type fails with different size", func(t *testing.T) {
		isArray := Is[[3]int]()
		res := isArray([4]int{1, 2, 3, 4})

		assert.True(t, either.IsLeft(res))
	})
}

// TestIsWithStructTypes tests Is with struct types
func TestIsWithStructTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	type Employee struct {
		Name   string
		Salary float64
	}

	t.Run("struct type succeeds with matching struct", func(t *testing.T) {
		isPerson := Is[Person]()
		person := Person{Name: "Alice", Age: 30}
		res := isPerson(person)

		assert.Equal(t, R.Of(person), res)
	})

	t.Run("struct type fails with different struct", func(t *testing.T) {
		isPerson := Is[Person]()
		employee := Employee{Name: "Bob", Salary: 50000}
		res := isPerson(employee)

		assert.True(t, either.IsLeft(res))
	})

	t.Run("struct type fails with primitive", func(t *testing.T) {
		isPerson := Is[Person]()
		res := isPerson("not a person")

		assert.True(t, either.IsLeft(res))
	})
}

// TestIsWithPointerTypes tests Is with pointer types
func TestIsWithPointerTypes(t *testing.T) {
	t.Run("pointer type succeeds with pointer", func(t *testing.T) {
		isStringPtr := Is[*string]()
		str := "hello"
		res := isStringPtr(&str)

		assert.Equal(t, R.Of(&str), res)
	})

	t.Run("pointer type fails with non-pointer", func(t *testing.T) {
		isStringPtr := Is[*string]()
		res := isStringPtr("hello")

		assert.True(t, either.IsLeft(res))
	})

	t.Run("pointer type succeeds with nil pointer", func(t *testing.T) {
		isStringPtr := Is[*string]()
		var nilPtr *string = nil
		res := isStringPtr(nilPtr)

		assert.Equal(t, R.Of(nilPtr), res)
	})

	t.Run("non-pointer type fails with pointer", func(t *testing.T) {
		isString := Is[string]()
		str := "hello"
		res := isString(&str)

		assert.True(t, either.IsLeft(res))
	})
}

// TestIsWithEmptyValues tests Is with empty/zero values
func TestIsWithEmptyValues(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		isString := Is[string]()
		res := isString("")

		assert.Equal(t, R.Of(""), res)
	})

	t.Run("zero int", func(t *testing.T) {
		isInt := Is[int]()
		res := isInt(0)

		assert.Equal(t, R.Of(0), res)
	})

	t.Run("false bool", func(t *testing.T) {
		isBool := Is[bool]()
		res := isBool(false)

		assert.Equal(t, R.Of(false), res)
	})

	t.Run("nil slice", func(t *testing.T) {
		isSlice := Is[[]int]()
		var nilSlice []int = nil
		res := isSlice(nilSlice)

		assert.Equal(t, R.Of(nilSlice), res)
	})

	t.Run("empty slice", func(t *testing.T) {
		isSlice := Is[[]int]()
		emptySlice := []int{}
		res := isSlice(emptySlice)

		assert.Equal(t, R.Of(emptySlice), res)
	})

	t.Run("nil map", func(t *testing.T) {
		isMap := Is[map[string]int]()
		var nilMap map[string]int = nil
		res := isMap(nilMap)

		assert.Equal(t, R.Of(nilMap), res)
	})
}

// TestIsWithChannelTypes tests Is with channel types
func TestIsWithChannelTypes(t *testing.T) {
	t.Run("channel type succeeds with channel", func(t *testing.T) {
		isChan := Is[chan int]()
		ch := make(chan int)
		defer close(ch)

		res := isChan(ch)
		assert.Equal(t, R.Of(ch), res)
	})

	t.Run("channel type fails with wrong channel type", func(t *testing.T) {
		isChan := Is[chan int]()
		ch := make(chan string)
		defer close(ch)

		res := isChan(ch)
		assert.True(t, either.IsLeft(res))
	})

	t.Run("bidirectional vs unidirectional channels", func(t *testing.T) {
		isSendChan := Is[chan<- int]()
		ch := make(chan int)
		defer close(ch)

		// Bidirectional channel can be used as send-only
		sendCh := chan<- int(ch)
		res := isSendChan(sendCh)
		assert.Equal(t, R.Of(sendCh), res)
	})
}

// TestIsWithFunctionTypes tests Is with function types
func TestIsWithFunctionTypes(t *testing.T) {
	t.Run("function type succeeds with matching function", func(t *testing.T) {
		isFunc := Is[func(int) int]()
		fn := func(x int) int { return x * 2 }

		res := isFunc(fn)
		// Functions can't be compared for equality, so just check it's Right
		assert.True(t, either.IsRight(res))
	})

	t.Run("function type fails with different signature", func(t *testing.T) {
		isFunc := Is[func(int) int]()
		fn := func(x string) string { return x }

		res := isFunc(fn)
		assert.True(t, either.IsLeft(res))
	})

	t.Run("function type fails with non-function", func(t *testing.T) {
		isFunc := Is[func(int) int]()
		res := isFunc(42)

		assert.True(t, either.IsLeft(res))
	})
}

// TestIsErrorMessages tests that Is produces appropriate error messages
func TestIsErrorMessages(t *testing.T) {
	t.Run("error message for type mismatch", func(t *testing.T) {
		isString := Is[string]()
		res := isString(42)

		assert.True(t, either.IsLeft(res), "Expected Left for type mismatch")
	})

	t.Run("error for struct type mismatch", func(t *testing.T) {
		type CustomType struct {
			Field string
		}

		isCustom := Is[CustomType]()
		res := isCustom("not a custom type")

		assert.True(t, either.IsLeft(res), "Expected Left for struct type mismatch")
	})
}
