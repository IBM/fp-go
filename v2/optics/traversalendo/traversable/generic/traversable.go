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

package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
)

// FromTraversableLens creates a traversal endomorphism from a lens that focuses on a traversable structure.
//
// This function converts a lens that accesses a traversable field (like an array, option, or any
// structure implementing the Traversable type class) into a traversal endomorphism that can traverse
// all elements within that structure. The resulting traversal works with Endomorphism[S] as its
// effect type, making it composable with other traversal endomorphisms using monoid operations.
//
// The function is a generalization of FromArrayLens that works with any traversable structure,
// not just arrays. This allows you to create traversals for fields containing options, results,
// trees, or any custom traversable data structure.
//
// The function works by:
//  1. Using the lens getter to extract the traversable structure from the source
//  2. Applying the traversable's traverse operation to all elements using the transformation function
//  3. Using the lens setter (mapped to work with endomorphisms) to update the source structure
//
// This is particularly useful when building complex traversals that need to access and modify
// traversable fields within a larger structure. The endomorphism-based approach allows you to
// combine multiple such traversals using Concat and Empty operations.
//
// Type Parameters:
//   - HKTA: Higher-kinded type for A in the effect context
//   - S: The source structure type containing the traversable field
//   - A: The element type within the traversable structure
//   - GA: The traversable structure type (e.g., []A, option.Option[A])
//   - HKTS: Higher-kinded type for Endomorphism[S]
//   - HKTRA: Higher-kinded type for the traversable structure in the effect context
//
// Parameters:
//   - fmap: Function to map the lens setter into the endomorphism context
//
// Returns:
//   - A function that takes a Traversable and returns a function that takes a Lens and returns a Traversal
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/option"
//	    thunk "github.com/IBM/fp-go/v2/context/readerioresult"
//	    TLG "github.com/IBM/fp-go/v2/optics/traversalendo/traversable/generic"
//	)
//
//	type Config struct {
//	    Name     string
//	    Database option.Option[string]
//	}
//
//	// Create a lens for the Database field
//	dbLens := lens.Lens[Config, option.Option[string]]{
//	    Get: func(c Config) option.Option[string] { return c.Database },
//	    Set: func(c Config, db option.Option[string]) Config {
//	        c.Database = db
//	        return c
//	    },
//	}
//
//	// Create a traversable for option.Option
//	optTraversable := option.Traverse[string, thunk.Thunk[string]]
//
//	// Convert to a traversal endomorphism
//	dbTrav := TLG.FromTraversableLens[
//	    thunk.Thunk[string],
//	    Config,
//	    string,
//	    option.Option[string],
//	    thunk.Thunk[endomorphism.Endomorphism[Config]],
//	    thunk.Thunk[option.Option[string]],
//	](
//	    thunk.Map[option.Option[string], endomorphism.Endomorphism[Config]],
//	)(optTraversable)(dbLens)
//
//	// Use to modify the database connection string if present
//	updateDB := func(db string) thunk.Thunk[string] {
//	    return thunk.Of("postgresql://" + db)
//	}
//	result := dbTrav(updateDB)(Config{
//	    Name:     "myapp",
//	    Database: option.Some("localhost"),
//	})
//
// See Also:
//   - FromArrayLens: Specialized version for array fields
//   - Traversable: The type class for traversable structures
//   - Lens: For focusing on a single field
func FromTraversableLens[A, HKTA, S, GA, HKTS, HKTRA any](
	fmap functor.MapType[GA, Endomorphism[S], HKTRA, HKTS],
) func(Traversable[A, HKTA, GA, HKTRA]) func(Lens[S, GA]) Traversal[S, A, HKTS, HKTA] {
	return func(trv Traversable[A, HKTA, GA, HKTRA]) func(Lens[S, GA]) Traversal[S, A, HKTS, HKTA] {
		return func(l Lens[S, GA]) Traversal[S, A, HKTS, HKTA] {
			lGet := l.Get
			lSet := fmap(l.Set)
			return func(f func(A) HKTA) func(S) HKTS {
				return F.Flow3(
					lGet,
					trv(f),
					lSet,
				)
			}
		}
	}
}
