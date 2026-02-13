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

package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/boolean"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/reader"
)

// ApplicativeMonoid returns a [monoid.Monoid] for combining test assertion [Reader]s.
//
// This monoid combines multiple test assertions using logical AND (conjunction) semantics,
// meaning all assertions must pass for the combined assertion to pass. It leverages the
// applicative structure of Reader to execute multiple assertions with the same testing.T
// context and combines their boolean results using boolean.MonoidAll (logical AND).
//
// The monoid provides:
//   - Concat: Combines two assertions such that both must pass (logical AND)
//   - Empty: Returns an assertion that always passes (identity element)
//
// This is particularly useful for:
//   - Composing multiple test assertions into a single assertion
//   - Building complex test conditions from simpler ones
//   - Creating reusable assertion combinators
//   - Implementing test assertion DSLs
//
// # Monoid Laws
//
// The returned monoid satisfies the standard monoid laws:
//
//  1. Associativity:
//     Concat(Concat(a1, a2), a3) ≡ Concat(a1, Concat(a2, a3))
//
//  2. Left Identity:
//     Concat(Empty(), a) ≡ a
//
//  3. Right Identity:
//     Concat(a, Empty()) ≡ a
//
// # Returns
//
//   - A [monoid.Monoid][Reader] that combines assertions using logical AND
//
// # Example - Basic Usage
//
//	func TestUserValidation(t *testing.T) {
//	    user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}
//	    m := assert.ApplicativeMonoid()
//
//	    // Combine multiple assertions
//	    assertion := m.Concat(
//	        assert.Equal("Alice")(user.Name),
//	        m.Concat(
//	            assert.Equal(30)(user.Age),
//	            assert.StringNotEmpty(user.Email),
//	        ),
//	    )
//
//	    // Execute combined assertion
//	    assertion(t) // All three assertions must pass
//	}
//
// # Example - Building Reusable Validators
//
//	func TestWithReusableValidators(t *testing.T) {
//	    m := assert.ApplicativeMonoid()
//
//	    // Create a reusable validator
//	    validateUser := func(u User) assert.Reader {
//	        return m.Concat(
//	            assert.StringNotEmpty(u.Name),
//	            m.Concat(
//	                assert.True(u.Age > 0),
//	                assert.StringContains("@")(u.Email),
//	            ),
//	        )
//	    }
//
//	    user := User{Name: "Bob", Age: 25, Email: "bob@test.com"}
//	    validateUser(user)(t)
//	}
//
// # Example - Using Empty for Identity
//
//	func TestEmptyIdentity(t *testing.T) {
//	    m := assert.ApplicativeMonoid()
//	    assertion := assert.Equal(42)(42)
//
//	    // Empty is the identity - these are equivalent
//	    result1 := m.Concat(m.Empty(), assertion)(t)
//	    result2 := m.Concat(assertion, m.Empty())(t)
//	    result3 := assertion(t)
//	    // All three produce the same result
//	}
//
// # Example - Combining with AllOf
//
//	func TestCombiningWithAllOf(t *testing.T) {
//	    // ApplicativeMonoid provides the underlying mechanism for AllOf
//	    arr := []int{1, 2, 3, 4, 5}
//
//	    // These are conceptually equivalent:
//	    m := assert.ApplicativeMonoid()
//	    manual := m.Concat(
//	        assert.ArrayNotEmpty(arr),
//	        m.Concat(
//	            assert.ArrayLength[int](5)(arr),
//	            assert.ArrayContains(3)(arr),
//	        ),
//	    )
//
//	    // AllOf uses ApplicativeMonoid internally
//	    convenient := assert.AllOf([]assert.Reader{
//	        assert.ArrayNotEmpty(arr),
//	        assert.ArrayLength[int](5)(arr),
//	        assert.ArrayContains(3)(arr),
//	    })
//
//	    manual(t)
//	    convenient(t)
//	}
//
// # Related Functions
//
//   - [AllOf]: Convenient wrapper for combining multiple assertions using this monoid
//   - [boolean.MonoidAll]: The underlying boolean monoid (logical AND with true as identity)
//   - [reader.ApplicativeMonoid]: Generic applicative monoid for Reader types
//
// # References
//
//   - Haskell Monoid: https://hackage.haskell.org/package/base/docs/Data-Monoid.html
//   - Applicative Functors: https://hackage.haskell.org/package/base/docs/Control-Applicative.html
//   - Boolean Monoid (All): https://hackage.haskell.org/package/base/docs/Data-Monoid.html#t:All
func ApplicativeMonoid() monoid.Monoid[Reader] {
	return reader.ApplicativeMonoid[*testing.T](boolean.MonoidAll)
}
