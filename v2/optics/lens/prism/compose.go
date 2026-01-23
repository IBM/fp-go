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

package prism

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/option"
)

func compose[S, A, B any](
	creator func(get option.Kleisli[S, B], set func(B) Endomorphism[S], name string) Optional[S, B],
	p Prism[A, B]) func(Lens[S, A]) Optional[S, B] {

	return func(l Lens[S, A]) Optional[S, B] {
		// GetOption: Lens.Get followed by Prism.GetOption
		// This extracts A from S, then tries to extract B from A
		getOption := F.Flow2(l.Get, p.GetOption)

		// Set: Constructs a setter that respects the Optional laws
		setOption := func(b B) func(S) S {
			// Pre-compute the new A value by using Prism.ReverseGet
			// This constructs an A from the given B
			setl := l.Set(p.ReverseGet(b))

			return func(s S) S {
				// Check if the Prism matches the current value
				return F.Pipe1(
					getOption(s),
					option.Fold(
						// None case: Prism doesn't match, return s unchanged (no-op)
						// This satisfies the GetSet law for Optional
						lazy.Of(s),
						// Some case: Prism matches, update the value
						// This satisfies the SetGet law for Optional
						func(_ B) S {
							return setl(s)
						},
					),
				)
			}
		}

		return creator(
			getOption,
			setOption,
			fmt.Sprintf("Compose[%s -> %s]", l, p),
		)
	}

}

// Compose composes a Lens with a Prism to create an Optional.
//
// This composition allows you to focus on a part of a structure (using a Lens)
// and then optionally extract a variant from that part (using a Prism). The result
// is an Optional because the Prism may not match the focused value.
//
// The composition follows the Optional laws (a relaxed form of lens laws):
//
// SetGet Law (GetSet for Optional):
//   - If optional.GetOption(s) = Some(b), then optional.GetOption(optional.Set(b)(s)) = Some(b)
//   - This ensures that setting a value and then getting it returns the same value
//
// GetSet Law (for Optional):
//   - If optional.GetOption(s) = None, then optional.Set(b)(s) = s (no-op)
//   - This ensures that setting a value when the optional doesn't match leaves the structure unchanged
//
// These laws are documented in the official fp-ts documentation:
// https://gcanti.github.io/monocle-ts/modules/Optional.ts.html
//
// Type Parameters:
//   - S: The source/outer structure type
//   - A: The intermediate type (focused by the Lens)
//   - B: The target type (focused by the Prism within A)
//
// Parameters:
//   - p: A Prism[A, B] that optionally extracts B from A
//
// Returns:
//   - A function that takes a Lens[S, A] and returns an Optional[S, B]
//
// Behavior:
//   - GetOption: First uses the Lens to get A from S, then uses the Prism to try to extract B from A.
//     Returns Some(b) if both operations succeed, None otherwise.
//   - Set: When setting a value b:
//   - If GetOption(s) returns Some(_), it means the Prism matches, so we:
//     1. Use Prism.ReverseGet to construct an A from b
//     2. Use Lens.Set to update S with the new A
//   - If GetOption(s) returns None, the Prism doesn't match, so we return s unchanged (no-op)
//
// Example:
//
//	type Config struct {
//	    Database DatabaseConfig
//	}
//
//	type DatabaseConfig struct {
//	    Connection ConnectionType
//	}
//
//	type ConnectionType interface{ isConnection() }
//	type PostgreSQL struct{ Host string }
//	type MySQL struct{ Host string }
//
//	// Lens to focus on Database field
//	dbLens := lens.MakeLens(
//	    func(c Config) DatabaseConfig { return c.Database },
//	    func(c Config, db DatabaseConfig) Config { c.Database = db; return c },
//	)
//
//	// Prism to extract PostgreSQL from ConnectionType
//	pgPrism := prism.MakePrism(
//	    func(ct ConnectionType) option.Option[PostgreSQL] {
//	        if pg, ok := ct.(PostgreSQL); ok {
//	            return option.Some(pg)
//	        }
//	        return option.None[PostgreSQL]()
//	    },
//	    func(pg PostgreSQL) ConnectionType { return pg },
//	)
//
//	// Compose to create Optional[Config, PostgreSQL]
//	configPgOptional := Compose[Config, DatabaseConfig, PostgreSQL](pgPrism)(dbLens)
//
//	config := Config{Database: DatabaseConfig{Connection: PostgreSQL{Host: "localhost"}}}
//	host := configPgOptional.GetOption(config)  // Some(PostgreSQL{Host: "localhost"})
//
//	updated := configPgOptional.Set(PostgreSQL{Host: "remote"})(config)
//	// updated.Database.Connection = PostgreSQL{Host: "remote"}
//
//	configMySQL := Config{Database: DatabaseConfig{Connection: MySQL{Host: "localhost"}}}
//	none := configPgOptional.GetOption(configMySQL)  // None (Prism doesn't match)
//	unchanged := configPgOptional.Set(PostgreSQL{Host: "remote"})(configMySQL)
//	// unchanged == configMySQL (no-op because Prism doesn't match)
func Compose[S, A, B any](p Prism[A, B]) func(Lens[S, A]) Optional[S, B] {
	return compose(O.MakeOptionalCurriedWithName[S, B], p)
}

// ComposeRef composes a Lens operating on pointer types with a Prism to create an Optional.
//
// This is the pointer-safe variant of Compose, designed for working with pointer types (*S).
// It automatically handles nil pointer cases and creates copies before modification to ensure
// immutability and prevent unintended side effects.
//
// The composition follows the same Optional laws as Compose:
//
// SetGet Law (GetSet for Optional):
//   - If optional.GetOption(s) = Some(b), then optional.GetOption(optional.Set(b)(s)) = Some(b)
//   - This ensures that setting a value and then getting it returns the same value
//
// GetSet Law (for Optional):
//   - If optional.GetOption(s) = None, then optional.Set(b)(s) = s (no-op)
//   - This ensures that setting a value when the optional doesn't match leaves the structure unchanged
//
// Nil Pointer Handling:
//   - When s is nil and GetOption would return None, Set operations return nil (no-op)
//   - When s is nil and GetOption would return Some (after creating default), Set creates a new instance
//   - All Set operations create a shallow copy of *S before modification to preserve immutability
//
// These laws are documented in the official fp-ts documentation:
// https://gcanti.github.io/monocle-ts/modules/Optional.ts.html
//
// Type Parameters:
//   - S: The source/outer structure type (used as *S in the lens)
//   - A: The intermediate type (focused by the Lens)
//   - B: The target type (focused by the Prism within A)
//
// Parameters:
//   - p: A Prism[A, B] that optionally extracts B from A
//
// Returns:
//   - A function that takes a Lens[*S, A] and returns an Optional[*S, B]
//
// Behavior:
//   - GetOption: First uses the Lens to get A from *S, then uses the Prism to try to extract B from A.
//     Returns Some(b) if both operations succeed, None otherwise.
//   - Set: When setting a value b:
//   - Creates a shallow copy of *S before any modification (nil-safe)
//   - If GetOption(s) returns Some(_), it means the Prism matches, so we:
//     1. Use Prism.ReverseGet to construct an A from b
//     2. Use Lens.Set to update the copy of *S with the new A
//   - If GetOption(s) returns None, the Prism doesn't match, so we return s unchanged (no-op)
//
// Example:
//
//	type Config struct {
//	    Connection ConnectionType
//	    AppName    string
//	}
//
//	type ConnectionType interface{ isConnection() }
//	type PostgreSQL struct{ Host string }
//	type MySQL struct{ Host string }
//
//	// Lens to focus on Connection field (pointer-based)
//	connLens := lens.MakeLensRef(
//	    func(c *Config) ConnectionType { return c.Connection },
//	    func(c *Config, ct ConnectionType) *Config { c.Connection = ct; return c },
//	)
//
//	// Prism to extract PostgreSQL from ConnectionType
//	pgPrism := prism.MakePrism(
//	    func(ct ConnectionType) option.Option[PostgreSQL] {
//	        if pg, ok := ct.(PostgreSQL); ok {
//	            return option.Some(pg)
//	        }
//	        return option.None[PostgreSQL]()
//	    },
//	    func(pg PostgreSQL) ConnectionType { return pg },
//	)
//
//	// Compose to create Optional[*Config, PostgreSQL]
//	configPgOptional := ComposeRef[Config, ConnectionType, PostgreSQL](pgPrism)(connLens)
//
//	// Works with non-nil pointers
//	config := &Config{Connection: PostgreSQL{Host: "localhost"}}
//	host := configPgOptional.GetOption(config)  // Some(PostgreSQL{Host: "localhost"})
//	updated := configPgOptional.Set(PostgreSQL{Host: "remote"})(config)
//	// updated is a new *Config with Connection = PostgreSQL{Host: "remote"}
//	// original config is unchanged (immutability preserved)
//
//	// Handles nil pointers safely
//	var nilConfig *Config = nil
//	none := configPgOptional.GetOption(nilConfig)  // None (nil pointer)
//	unchanged := configPgOptional.Set(PostgreSQL{Host: "remote"})(nilConfig)
//	// unchanged == nil (no-op because source is nil)
//
//	// Works with mismatched prisms
//	configMySQL := &Config{Connection: MySQL{Host: "localhost"}}
//	none = configPgOptional.GetOption(configMySQL)  // None (Prism doesn't match)
//	unchanged = configPgOptional.Set(PostgreSQL{Host: "remote"})(configMySQL)
//	// unchanged == configMySQL (no-op because Prism doesn't match)
func ComposeRef[S, A, B any](p Prism[A, B]) func(Lens[*S, A]) Optional[*S, B] {
	return compose(O.MakeOptionalRefCurriedWithName[S, B], p)
}
