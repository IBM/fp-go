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

package codec

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/record"
	"github.com/stretchr/testify/assert"
)

// atKeyStr is a minimal comparable+Stringer key type backed by a plain string.
type atKeyStr string

func (k atKeyStr) String() string { return string(k) }

// TestAtKey_Decode_Success verifies that AtKey decodes successfully when the key exists.
func TestAtKey_Decode_Success(t *testing.T) {
	t.Run("present key returns success", func(t *testing.T) {
		c := AtKey[int](atKeyStr("count"))
		r := record.Record[atKeyStr, int]{"count": 42}
		result := c.Decode(r)
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("zero value is a valid success", func(t *testing.T) {
		c := AtKey[int](atKeyStr("count"))
		r := record.Record[atKeyStr, int]{"count": 0}
		result := c.Decode(r)
		assert.Equal(t, validation.Of(0), result)
	})

	t.Run("string value type", func(t *testing.T) {
		c := AtKey[string](atKeyStr("name"))
		r := record.Record[atKeyStr, string]{"name": "alice"}
		result := c.Decode(r)
		assert.Equal(t, validation.Of("alice"), result)
	})

	t.Run("record with multiple keys decodes targeted key", func(t *testing.T) {
		c := AtKey[int](atKeyStr("b"))
		r := record.Record[atKeyStr, int]{"a": 1, "b": 2, "c": 3}
		result := c.Decode(r)
		assert.Equal(t, validation.Of(2), result)
	})
}

// TestAtKey_Decode_Failure verifies that AtKey fails when the key is absent.
func TestAtKey_Decode_Failure(t *testing.T) {
	c := AtKey[int](atKeyStr("count"))

	t.Run("absent key returns Left", func(t *testing.T) {
		r := record.Record[atKeyStr, int]{"other": 7}
		result := c.Decode(r)
		assert.True(t, either.IsLeft(result), "expected Left for missing key")
	})

	t.Run("empty record returns Left", func(t *testing.T) {
		r := record.Record[atKeyStr, int]{}
		result := c.Decode(r)
		assert.True(t, either.IsLeft(result), "expected Left for empty record")
	})

	t.Run("nil record returns Left", func(t *testing.T) {
		result := c.Decode(nil)
		assert.True(t, either.IsLeft(result), "expected Left for nil record")
	})
}

// TestAtKey_Encode verifies that AtKey encodes a value into a single-entry record.
func TestAtKey_Encode(t *testing.T) {
	t.Run("encodes to single-entry record", func(t *testing.T) {
		c := AtKey[int](atKeyStr("score"))
		encoded := c.Encode(99)
		assert.Equal(t, record.Record[atKeyStr, int]{"score": 99}, encoded)
	})

	t.Run("encodes zero value", func(t *testing.T) {
		c := AtKey[int](atKeyStr("score"))
		encoded := c.Encode(0)
		assert.Equal(t, record.Record[atKeyStr, int]{"score": 0}, encoded)
	})

	t.Run("encodes string value", func(t *testing.T) {
		c := AtKey[string](atKeyStr("city"))
		encoded := c.Encode("berlin")
		assert.Equal(t, record.Record[atKeyStr, string]{"city": "berlin"}, encoded)
	})
}

// TestAtKey_Name verifies the human-readable codec name includes the key.
func TestAtKey_Name(t *testing.T) {
	t.Run("name contains key string", func(t *testing.T) {
		c := AtKey[int](atKeyStr("age"))
		assert.Equal(t, "AtKey[age]", c.Name())
	})

	t.Run("name reflects different keys", func(t *testing.T) {
		c1 := AtKey[string](atKeyStr("first"))
		c2 := AtKey[string](atKeyStr("second"))
		assert.Equal(t, "AtKey[first]", c1.Name())
		assert.Equal(t, "AtKey[second]", c2.Name())
	})
}

// TestAtKey_RoundTrip verifies that encoding then decoding returns the original value.
func TestAtKey_RoundTrip(t *testing.T) {
	c := AtKey[int](atKeyStr("x"))

	for i := range 5 {
		encoded := c.Encode(i)
		decoded := c.Decode(encoded)
		assert.Equal(t, validation.Of(i), decoded, "round-trip failed for %d", i)
	}
}

// TestAtKey_ErrorMessage verifies that the failure message mentions the key name.
func TestAtKey_ErrorMessage(t *testing.T) {
	c := AtKey[int](atKeyStr("myKey"))
	result := c.Decode(record.Record[atKeyStr, int]{})

	assert.True(t, either.IsLeft(result))
	msg := either.MonadFold(result,
		func(errs validation.Errors) string { return fmt.Sprintf("%v", errs) },
		func(_ int) string { return "" },
	)
	assert.Contains(t, msg, "myKey", "error should reference the missing key name")
}

// TestAtKey_Integration verifies two-level nested record lookup using AtKey.
func TestAtKey_Integration(t *testing.T) {
	outerCodec := AtKey[record.Record[atKeyStr, string]](atKeyStr("inner"))
	innerCodec := AtKey[string](atKeyStr("value"))

	t.Run("two-level nested record lookup", func(t *testing.T) {
		r := record.Record[atKeyStr, record.Record[atKeyStr, string]]{
			"inner": {"value": "hello"},
		}
		outerResult := outerCodec.Decode(r)
		assert.True(t, either.IsRight(outerResult))

		innerRecord, _ := either.Unwrap(outerResult)
		innerResult := innerCodec.Decode(innerRecord)
		assert.Equal(t, validation.Of("hello"), innerResult)
	})

	t.Run("missing outer key propagates failure", func(t *testing.T) {
		r := record.Record[atKeyStr, record.Record[atKeyStr, string]]{}
		outerResult := outerCodec.Decode(r)
		assert.True(t, either.IsLeft(outerResult))
	})

	t.Run("missing inner key propagates failure", func(t *testing.T) {
		r := record.Record[atKeyStr, record.Record[atKeyStr, string]]{
			"inner": {},
		}
		outerResult := outerCodec.Decode(r)
		innerRecord, _ := either.Unwrap(outerResult)
		innerResult := innerCodec.Decode(innerRecord)
		assert.True(t, either.IsLeft(innerResult))
	})
}

// TestAtKey_Validate verifies the context-aware Validate function.
func TestAtKey_Validate(t *testing.T) {
	c := AtKey[int](atKeyStr("n"))

	t.Run("validate succeeds for present key", func(t *testing.T) {
		r := record.Record[atKeyStr, int]{"n": 7}
		result := c.Validate(r)(validation.Context{})
		assert.Equal(t, validation.Of(7), result)
	})

	t.Run("validate fails for absent key", func(t *testing.T) {
		r := record.Record[atKeyStr, int]{}
		result := c.Validate(r)(validation.Context{})
		assert.True(t, either.IsLeft(result))
	})
}

// TestAtKey_TypeCheck verifies the Is field: it checks the value type V, not the outer record.
// Is[V]() succeeds when the raw any argument is of type V (here int).
func TestAtKey_TypeCheck(t *testing.T) {
	c := AtKey[int](atKeyStr("k"))

	t.Run("int value passes Is check", func(t *testing.T) {
		res := c.Is(42)
		assert.True(t, either.IsRight(res))
	})

	t.Run("string fails Is check", func(t *testing.T) {
		res := c.Is("not an int")
		assert.True(t, either.IsLeft(res))
	})

	t.Run("nil fails Is check", func(t *testing.T) {
		res := c.Is(nil)
		assert.True(t, either.IsLeft(res), "nil should fail the Is check")
	})
}

// TestAtKey_BenchmarkDecode benchmarks hot-path decoding.
func BenchmarkAtKey_Decode(b *testing.B) {
	c := AtKey[int](atKeyStr("key"))
	r := record.Record[atKeyStr, int]{"key": 1}
	b.ResetTimer()
	for range b.N {
		_ = c.Decode(r)
	}
}

// BenchmarkAtKey_Encode benchmarks hot-path encoding.
func BenchmarkAtKey_Encode(b *testing.B) {
	c := AtKey[int](atKeyStr("key"))
	b.ResetTimer()
	for range b.N {
		_ = c.Encode(42)
	}
}

