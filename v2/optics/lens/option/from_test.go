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

package option

import (
	"testing"

	EQT "github.com/IBM/fp-go/v2/eq/testing"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	ISO "github.com/IBM/fp-go/v2/optics/iso"
	L "github.com/IBM/fp-go/v2/optics/lens"
	LT "github.com/IBM/fp-go/v2/optics/lens/testing"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Test types
type Config struct {
	timeout int
	retries int
}

type Settings struct {
	maxConnections int
	bufferSize     int
}

// TestFromIsoBasic tests basic functionality of FromIso
func TestFromIsoBasic(t *testing.T) {
	// Create an isomorphism that treats 0 as None
	zeroAsNone := ISO.MakeIso(
		func(t int) O.Option[int] {
			if t == 0 {
				return O.None[int]()
			}
			return O.Some(t)
		},
		func(opt O.Option[int]) int {
			return O.GetOrElse(F.Constant(0))(opt)
		},
	)

	// Create a lens to the timeout field
	timeoutLens := L.MakeLens(
		func(c Config) int { return c.timeout },
		func(c Config, t int) Config { c.timeout = t; return c },
	)

	// Convert to optional lens using FromIso
	optTimeoutLens := FromIso[Config](zeroAsNone)(timeoutLens)

	t.Run("GetNone", func(t *testing.T) {
		config := Config{timeout: 0, retries: 3}
		result := optTimeoutLens.Get(config)
		assert.True(t, O.IsNone(result))
	})

	t.Run("GetSome", func(t *testing.T) {
		config := Config{timeout: 30, retries: 3}
		result := optTimeoutLens.Get(config)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 30, O.GetOrElse(F.Constant(0))(result))
	})

	t.Run("SetNone", func(t *testing.T) {
		config := Config{timeout: 30, retries: 3}
		updated := optTimeoutLens.Set(O.None[int]())(config)
		assert.Equal(t, 0, updated.timeout)
		assert.Equal(t, 3, updated.retries) // Other fields unchanged
	})

	t.Run("SetSome", func(t *testing.T) {
		config := Config{timeout: 0, retries: 3}
		updated := optTimeoutLens.Set(O.Some(60))(config)
		assert.Equal(t, 60, updated.timeout)
		assert.Equal(t, 3, updated.retries) // Other fields unchanged
	})

	t.Run("SetPreservesOriginal", func(t *testing.T) {
		original := Config{timeout: 30, retries: 3}
		_ = optTimeoutLens.Set(O.Some(60))(original)
		// Original should be unchanged
		assert.Equal(t, 30, original.timeout)
		assert.Equal(t, 3, original.retries)
	})
}

// TestFromIsoWithNegativeSentinel tests using -1 as a sentinel value
func TestFromIsoWithNegativeSentinel(t *testing.T) {
	// Create an isomorphism that treats -1 as None
	negativeOneAsNone := ISO.MakeIso(
		func(n int) O.Option[int] {
			if n == -1 {
				return O.None[int]()
			}
			return O.Some(n)
		},
		func(opt O.Option[int]) int {
			return O.GetOrElse(F.Constant(-1))(opt)
		},
	)

	retriesLens := L.MakeLens(
		func(c Config) int { return c.retries },
		func(c Config, r int) Config { c.retries = r; return c },
	)

	optRetriesLens := FromIso[Config](negativeOneAsNone)(retriesLens)

	t.Run("GetNoneForNegativeOne", func(t *testing.T) {
		config := Config{timeout: 30, retries: -1}
		result := optRetriesLens.Get(config)
		assert.True(t, O.IsNone(result))
	})

	t.Run("GetSomeForZero", func(t *testing.T) {
		config := Config{timeout: 30, retries: 0}
		result := optRetriesLens.Get(config)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 0, O.GetOrElse(F.Constant(-1))(result))
	})

	t.Run("SetNoneToNegativeOne", func(t *testing.T) {
		config := Config{timeout: 30, retries: 5}
		updated := optRetriesLens.Set(O.None[int]())(config)
		assert.Equal(t, -1, updated.retries)
	})
}

// TestFromIsoLaws verifies that FromIso satisfies lens laws
func TestFromIsoLaws(t *testing.T) {
	// Create an isomorphism
	zeroAsNone := ISO.MakeIso(
		func(t int) O.Option[int] {
			if t == 0 {
				return O.None[int]()
			}
			return O.Some(t)
		},
		func(opt O.Option[int]) int {
			return O.GetOrElse(F.Constant(0))(opt)
		},
	)

	timeoutLens := L.MakeLens(
		func(c Config) int { return c.timeout },
		func(c Config, t int) Config { c.timeout = t; return c },
	)

	optTimeoutLens := FromIso[Config](zeroAsNone)(timeoutLens)

	eqOptInt := O.Eq(EQT.Eq[int]())
	eqConfig := EQT.Eq[Config]()

	config := Config{timeout: 30, retries: 3}
	newValue := O.Some(60)

	// Law 1: GetSet - lens.Set(lens.Get(s))(s) == s
	t.Run("GetSetLaw", func(t *testing.T) {
		result := optTimeoutLens.Set(optTimeoutLens.Get(config))(config)
		assert.True(t, eqConfig.Equals(config, result))
	})

	// Law 2: SetGet - lens.Get(lens.Set(a)(s)) == a
	t.Run("SetGetLaw", func(t *testing.T) {
		result := optTimeoutLens.Get(optTimeoutLens.Set(newValue)(config))
		assert.True(t, eqOptInt.Equals(newValue, result))
	})

	// Law 3: SetSet - lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)
	t.Run("SetSetLaw", func(t *testing.T) {
		a1 := O.Some(60)
		a2 := O.None[int]()
		result1 := optTimeoutLens.Set(a2)(optTimeoutLens.Set(a1)(config))
		result2 := optTimeoutLens.Set(a2)(config)
		assert.True(t, eqConfig.Equals(result1, result2))
	})

	// Use the testing helper to verify all laws
	t.Run("AllLaws", func(t *testing.T) {
		laws := LT.AssertLaws(t, eqOptInt, eqConfig)(optTimeoutLens)
		assert.True(t, laws(config, O.Some(100)))
		assert.True(t, laws(Config{timeout: 0, retries: 5}, O.None[int]()))
	})
}

// TestFromIsoComposition tests composing FromIso with other lenses
func TestFromIsoComposition(t *testing.T) {
	type Application struct {
		config Config
	}

	// Isomorphism for zero as none
	zeroAsNone := ISO.MakeIso(
		func(t int) O.Option[int] {
			if t == 0 {
				return O.None[int]()
			}
			return O.Some(t)
		},
		func(opt O.Option[int]) int {
			return O.GetOrElse(F.Constant(0))(opt)
		},
	)

	// Lens to config field
	configLens := L.MakeLens(
		func(a Application) Config { return a.config },
		func(a Application, c Config) Application { a.config = c; return a },
	)

	// Lens to timeout field
	timeoutLens := L.MakeLens(
		func(c Config) int { return c.timeout },
		func(c Config, t int) Config { c.timeout = t; return c },
	)

	// Compose: Application -> Config -> timeout (as Option)
	optTimeoutFromConfig := FromIso[Config](zeroAsNone)(timeoutLens)
	optTimeoutFromApp := F.Pipe1(
		configLens,
		L.Compose[Application](optTimeoutFromConfig),
	)

	app := Application{config: Config{timeout: 0, retries: 3}}

	t.Run("ComposedGet", func(t *testing.T) {
		result := optTimeoutFromApp.Get(app)
		assert.True(t, O.IsNone(result))
	})

	t.Run("ComposedSet", func(t *testing.T) {
		updated := optTimeoutFromApp.Set(O.Some(45))(app)
		assert.Equal(t, 45, updated.config.timeout)
		assert.Equal(t, 3, updated.config.retries)
	})
}

// TestFromIsoModify tests using Modify with FromIso-based lenses
func TestFromIsoModify(t *testing.T) {
	zeroAsNone := ISO.MakeIso(
		func(t int) O.Option[int] {
			if t == 0 {
				return O.None[int]()
			}
			return O.Some(t)
		},
		func(opt O.Option[int]) int {
			return O.GetOrElse(F.Constant(0))(opt)
		},
	)

	timeoutLens := L.MakeLens(
		func(c Config) int { return c.timeout },
		func(c Config, t int) Config { c.timeout = t; return c },
	)

	optTimeoutLens := FromIso[Config](zeroAsNone)(timeoutLens)

	t.Run("ModifyNoneToSome", func(t *testing.T) {
		config := Config{timeout: 0, retries: 3}
		// Map None to Some(10)
		modified := L.Modify[Config](O.Map(N.Add(10)))(optTimeoutLens)(config)
		// Since it was None, Map doesn't apply, stays None (0)
		assert.Equal(t, 0, modified.timeout)
	})

	t.Run("ModifySomeValue", func(t *testing.T) {
		config := Config{timeout: 30, retries: 3}
		// Double the timeout value
		modified := L.Modify[Config](O.Map(N.Mul(2)))(optTimeoutLens)(config)
		assert.Equal(t, 60, modified.timeout)
	})

	t.Run("ModifyWithAlt", func(t *testing.T) {
		config := Config{timeout: 0, retries: 3}
		// Use Alt to provide a default
		modified := L.Modify[Config](func(opt O.Option[int]) O.Option[int] {
			return O.Alt(F.Constant(O.Some(10)))(opt)
		})(optTimeoutLens)(config)
		assert.Equal(t, 10, modified.timeout)
	})
}

// TestFromIsoWithStringEmpty tests using empty string as None
func TestFromIsoWithStringEmpty(t *testing.T) {
	type User struct {
		name  string
		email string
	}

	// Isomorphism that treats empty string as None
	emptyAsNone := ISO.MakeIso(
		func(s string) O.Option[string] {
			if S.IsEmpty(s) {
				return O.None[string]()
			}
			return O.Some(s)
		},
		func(opt O.Option[string]) string {
			return O.GetOrElse(F.Constant(""))(opt)
		},
	)

	emailLens := L.MakeLens(
		func(u User) string { return u.email },
		func(u User, e string) User { u.email = e; return u },
	)

	optEmailLens := FromIso[User](emptyAsNone)(emailLens)

	t.Run("EmptyStringAsNone", func(t *testing.T) {
		user := User{name: "Alice", email: ""}
		result := optEmailLens.Get(user)
		assert.True(t, O.IsNone(result))
	})

	t.Run("NonEmptyStringAsSome", func(t *testing.T) {
		user := User{name: "Alice", email: "alice@example.com"}
		result := optEmailLens.Get(user)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "alice@example.com", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("SetNoneToEmpty", func(t *testing.T) {
		user := User{name: "Alice", email: "alice@example.com"}
		updated := optEmailLens.Set(O.None[string]())(user)
		assert.Equal(t, "", updated.email)
	})
}

// TestFromIsoRoundTrip tests round-trip conversions
func TestFromIsoRoundTrip(t *testing.T) {
	zeroAsNone := ISO.MakeIso(
		func(t int) O.Option[int] {
			if t == 0 {
				return O.None[int]()
			}
			return O.Some(t)
		},
		func(opt O.Option[int]) int {
			return O.GetOrElse(F.Constant(0))(opt)
		},
	)

	maxConnectionsLens := L.MakeLens(
		func(s Settings) int { return s.maxConnections },
		func(s Settings, m int) Settings { s.maxConnections = m; return s },
	)

	optMaxConnectionsLens := FromIso[Settings](zeroAsNone)(maxConnectionsLens)

	t.Run("RoundTripThroughGet", func(t *testing.T) {
		settings := Settings{maxConnections: 100, bufferSize: 1024}
		// Get the value, then Set it back
		opt := optMaxConnectionsLens.Get(settings)
		restored := optMaxConnectionsLens.Set(opt)(settings)
		assert.Equal(t, settings, restored)
	})

	t.Run("RoundTripThroughSet", func(t *testing.T) {
		settings := Settings{maxConnections: 0, bufferSize: 1024}
		// Set a new value, then Get it
		newOpt := O.Some(200)
		updated := optMaxConnectionsLens.Set(newOpt)(settings)
		retrieved := optMaxConnectionsLens.Get(updated)
		assert.True(t, O.Eq(EQT.Eq[int]()).Equals(newOpt, retrieved))
	})

	t.Run("RoundTripWithNone", func(t *testing.T) {
		settings := Settings{maxConnections: 100, bufferSize: 1024}
		// Set None, then get it back
		updated := optMaxConnectionsLens.Set(O.None[int]())(settings)
		retrieved := optMaxConnectionsLens.Get(updated)
		assert.True(t, O.IsNone(retrieved))
	})
}

// TestFromIsoChaining tests chaining multiple FromIso transformations
func TestFromIsoChaining(t *testing.T) {
	// Create two different isomorphisms
	zeroAsNone := ISO.MakeIso(
		func(t int) O.Option[int] {
			if t == 0 {
				return O.None[int]()
			}
			return O.Some(t)
		},
		func(opt O.Option[int]) int {
			return O.GetOrElse(F.Constant(0))(opt)
		},
	)

	timeoutLens := L.MakeLens(
		func(c Config) int { return c.timeout },
		func(c Config, t int) Config { c.timeout = t; return c },
	)

	optTimeoutLens := FromIso[Config](zeroAsNone)(timeoutLens)

	config := Config{timeout: 30, retries: 3}

	t.Run("ChainedOperations", func(t *testing.T) {
		// Chain multiple operations
		result := F.Pipe2(
			config,
			optTimeoutLens.Set(O.Some(60)),
			optTimeoutLens.Set(O.None[int]()),
		)
		assert.Equal(t, 0, result.timeout)
	})
}

// TestFromIsoMultipleFields tests using FromIso on multiple fields
func TestFromIsoMultipleFields(t *testing.T) {
	zeroAsNone := ISO.MakeIso(
		func(t int) O.Option[int] {
			if t == 0 {
				return O.None[int]()
			}
			return O.Some(t)
		},
		func(opt O.Option[int]) int {
			return O.GetOrElse(F.Constant(0))(opt)
		},
	)

	timeoutLens := L.MakeLens(
		func(c Config) int { return c.timeout },
		func(c Config, t int) Config { c.timeout = t; return c },
	)

	retriesLens := L.MakeLens(
		func(c Config) int { return c.retries },
		func(c Config, r int) Config { c.retries = r; return c },
	)

	optTimeoutLens := FromIso[Config](zeroAsNone)(timeoutLens)
	optRetriesLens := FromIso[Config](zeroAsNone)(retriesLens)

	t.Run("IndependentFields", func(t *testing.T) {
		config := Config{timeout: 0, retries: 5}

		// Get both fields
		timeoutOpt := optTimeoutLens.Get(config)
		retriesOpt := optRetriesLens.Get(config)

		assert.True(t, O.IsNone(timeoutOpt))
		assert.True(t, O.IsSome(retriesOpt))
		assert.Equal(t, 5, O.GetOrElse(F.Constant(0))(retriesOpt))
	})

	t.Run("SetBothFields", func(t *testing.T) {
		config := Config{timeout: 0, retries: 0}

		// Set both fields
		updated := F.Pipe2(
			config,
			optTimeoutLens.Set(O.Some(30)),
			optRetriesLens.Set(O.Some(3)),
		)

		assert.Equal(t, 30, updated.timeout)
		assert.Equal(t, 3, updated.retries)
	})
}
