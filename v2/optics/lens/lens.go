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

// Lens is an optic used to zoom inside a product.
package lens

import (
	"fmt"

	"github.com/IBM/fp-go/v2/endomorphism"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
)

// setCopy wraps a setter for a pointer into a setter that first creates a copy before
// modifying that copy
func setCopy[SET ~func(*S, A) *S, S, A any](setter SET) func(s *S, a A) *S {

	var empty S
	safeSet := func(s *S, a A) *S {
		// make sure we have a total implementation
		cpy := *s
		return setter(&cpy, a)
	}

	return func(s *S, a A) *S {
		// make sure we have a total implementation
		if s != nil {
			return safeSet(s, a)
		}
		// fallback to the empty object
		return safeSet(&empty, a)
	}
}

func setCopyWithEq[GET ~func(*S) A, SET ~func(*S, A) *S, S, A any](pred EQ.Eq[A], getter GET, setter SET) func(s *S, a A) *S {

	var empty S
	safeSet := func(s *S, a A) *S {
		if pred.Equals(getter(s), a) {
			return s
		}
		// we need to make a copy
		cpy := *s
		return setter(&cpy, a)
	}

	return func(s *S, a A) *S {
		// make sure we have a total implementation
		if s != nil {
			return safeSet(s, a)
		}
		// fallback to the empty object
		return safeSet(&empty, a)
	}
}

// setCopyCurried wraps a setter for a pointer into a setter that first creates a copy before
// modifying that copy
func setCopyCurried[SET ~func(A) Endomorphism[*S], S, A any](setter SET) func(A) Endomorphism[*S] {
	var empty S

	return func(a A) Endomorphism[*S] {
		seta := setter(a)

		safeSet := func(s *S) *S {
			// make sure we have a total implementation
			cpy := *s
			return seta(&cpy)
		}

		return func(s *S) *S {
			// make sure we have a total implementation
			if s != nil {
				return safeSet(s)
			}
			// fallback to the empty object
			return safeSet(&empty)
		}
	}
}

// MakeLens creates a [Lens] based on a getter and a setter F.
//
// The setter must create a (shallow) copy of the data structure. This happens automatically
// when the data is passed by value. For pointer-based structures, use [MakeLensRef] instead.
// For other reference types (slices, maps), ensure the setter creates a copy.
//
// Type Parameters:
//   - GET: Getter function type (S → A)
//   - SET: Setter function type (S, A → S)
//   - S: Source structure type
//   - A: Focus/field type
//
// Parameters:
//   - get: Function to extract value A from structure S
//   - set: Function to update value A in structure S, returning a new S
//
// Returns:
//   - A Lens[S, A] that can get and set values immutably
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLens(
//	    func(p Person) string { return p.Name },
//	    func(p Person, name string) Person {
//	        p.Name = name
//	        return p
//	    },
//	)
//
//	person := Person{Name: "Alice", Age: 30}
//	name := nameLens.Get(person)           // "Alice"
//	updated := nameLens.Set("Bob")(person) // Person{Name: "Bob", Age: 30}
//
//go:inline
func MakeLens[GET ~func(S) A, SET ~func(S, A) S, S, A any](get GET, set SET) Lens[S, A] {
	return MakeLensCurried(get, F.Bind2of2(set))
}

// MakeLensWithName creates a [Lens] with a custom name for debugging and logging.
//
// This is identical to [MakeLens] but allows you to specify a name that will be used
// when the lens is printed or formatted. The name is useful for debugging complex lens
// compositions and understanding which lens is being used in error messages or logs.
//
// The setter must create a (shallow) copy of the data structure. This happens automatically
// when the data is passed by value. For pointer-based structures, use [MakeLensRef] instead.
//
// Type Parameters:
//   - GET: Getter function type (S → A)
//   - SET: Setter function type (S, A → S)
//   - S: Source structure type
//   - A: Focus/field type
//
// Parameters:
//   - get: Function to extract value A from structure S
//   - set: Function to update value A in structure S, returning a new S
//   - name: A descriptive name for the lens (used in String() and Format())
//
// Returns:
//   - A Lens[S, A] with the specified name
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLensWithName(
//	    func(p Person) string { return p.Name },
//	    func(p Person, name string) Person {
//	        p.Name = name
//	        return p
//	    },
//	    "Person.Name",
//	)
//
//	fmt.Printf("Using lens: %s\n", nameLens)  // Prints: "Using lens: Person.Name"
//
//go:inline
func MakeLensWithName[GET ~func(S) A, SET ~func(S, A) S, S, A any](get GET, set SET, name string) Lens[S, A] {
	return MakeLensCurriedWithName(get, F.Bind2of2(set), name)
}

// MakeLensCurried creates a [Lens] with a curried setter F.
//
// This is similar to [MakeLens] but accepts a curried setter (A → S → S) instead of
// an uncurried one (S, A → S). The curried form is more composable in functional pipelines.
//
// The setter must create a (shallow) copy of the data structure. This happens automatically
// when the data is passed by value. For pointer-based structures, use [MakeLensRefCurried].
//
// Type Parameters:
//   - GET: Getter function type (S → A)
//   - SET: Curried setter function type (A → S → S)
//   - S: Source structure type
//   - A: Focus/field type
//
// Parameters:
//   - get: Function to extract value A from structure S
//   - set: Curried function to update value A in structure S
//
// Returns:
//   - A Lens[S, A] that can get and set values immutably
//
// Example:
//
//	nameLens := lens.MakeLensCurried(
//	    func(p Person) string { return p.Name },
//	    func(name string) func(Person) Person {
//	        return func(p Person) Person {
//	            p.Name = name
//	            return p
//	        }
//	    },
//	)
//
//go:inline
func MakeLensCurried[GET ~func(S) A, SET ~func(A) Endomorphism[S], S, A any](get GET, set SET) Lens[S, A] {
	return MakeLensCurriedWithName(get, set, "Lens")
}

// MakeLensCurriedWithName creates a [Lens] with a curried setter and a custom name.
//
// This combines the benefits of [MakeLensCurried] (curried setter for better composition)
// with [MakeLensWithName] (custom name for debugging). The name is useful for debugging
// complex lens compositions and understanding which lens is being used in error messages or logs.
//
// The setter must create a (shallow) copy of the data structure. This happens automatically
// when the data is passed by value. For pointer-based structures, use [MakeLensRefCurried].
//
// Type Parameters:
//   - GET: Getter function type (S → A)
//   - SET: Curried setter function type (A → S → S)
//   - S: Source structure type
//   - A: Focus/field type
//
// Parameters:
//   - get: Function to extract value A from structure S
//   - set: Curried function to update value A in structure S
//   - name: A descriptive name for the lens (used in String() and Format())
//
// Returns:
//   - A Lens[S, A] with the specified name
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLensCurriedWithName(
//	    func(p Person) string { return p.Name },
//	    func(name string) func(Person) Person {
//	        return func(p Person) Person {
//	            p.Name = name
//	            return p
//	        }
//	    },
//	    "Person.Name",
//	)
//
//	fmt.Printf("Using lens: %s\n", nameLens)  // Prints: "Using lens: Person.Name"
//
//go:inline
func MakeLensCurriedWithName[GET ~func(S) A, SET ~func(A) Endomorphism[S], S, A any](get GET, set SET, name string) Lens[S, A] {
	return Lens[S, A]{Get: get, Set: set, name: name}
}

//go:inline
func MakeLensCurriedRefWithName[GET ~func(*S) A, SET ~func(A) Endomorphism[*S], S, A any](get GET, set SET, name string) Lens[*S, A] {
	return Lens[*S, A]{Get: get, Set: setCopyCurried(set), name: name}
}

// MakeLensRef creates a [Lens] for pointer-based structures.
//
// Unlike [MakeLens], the setter does not need to create a copy manually. This function
// automatically wraps the setter to create a shallow copy of the pointed-to value before
// modification, ensuring immutability.
//
// This lens assumes that property A always exists in structure S (i.e., it's not optional).
//
// Type Parameters:
//   - GET: Getter function type (*S → A)
//   - SET: Setter function type (*S, A → *S)
//   - S: Source structure type (will be used as *S)
//   - A: Focus/field type
//
// Parameters:
//   - get: Function to extract value A from pointer *S
//   - set: Function to update value A in pointer *S (copying handled automatically)
//
// Returns:
//   - A Lens[*S, A] that can get and set values immutably on pointers
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLensRef(
//	    func(p *Person) string { return p.Name },
//	    func(p *Person, name string) *Person {
//	        p.Name = name  // No manual copy needed
//	        return p
//	    },
//	)
//
//	person := &Person{Name: "Alice", Age: 30}
//	updated := nameLens.Set("Bob")(person)
//	// person.Name is still "Alice", updated is a new pointer with Name "Bob"
//
//go:inline
func MakeLensRef[GET ~func(*S) A, SET func(*S, A) *S, S, A any](get GET, set SET) Lens[*S, A] {
	return MakeLens(get, setCopy(set))
}

// MakeLensRefWithName creates a [Lens] for pointer-based structures with a custom name.
//
// This combines [MakeLensRef] (automatic copying for pointer structures) with
// [MakeLensWithName] (custom name for debugging). The setter does not need to create
// a copy manually; this function automatically wraps it to ensure immutability.
// The name is useful for debugging complex lens compositions and understanding which
// lens is being used in error messages or logs.
//
// This lens assumes that property A always exists in structure S (i.e., it's not optional).
//
// Type Parameters:
//   - GET: Getter function type (*S → A)
//   - SET: Setter function type (*S, A → *S)
//   - S: Source structure type (will be used as *S)
//   - A: Focus/field type
//
// Parameters:
//   - get: Function to extract value A from pointer *S
//   - set: Function to update value A in pointer *S (copying handled automatically)
//   - name: A descriptive name for the lens (used in String() and Format())
//
// Returns:
//   - A Lens[*S, A] with the specified name
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLensRefWithName(
//	    func(p *Person) string { return p.Name },
//	    func(p *Person, name string) *Person {
//	        p.Name = name  // No manual copy needed
//	        return p
//	    },
//	    "Person.Name",
//	)
//
//	person := &Person{Name: "Alice", Age: 30}
//	fmt.Printf("Using lens: %s\n", nameLens)  // Prints: "Using lens: Person.Name"
//	updated := nameLens.Set("Bob")(person)
//	// person.Name is still "Alice", updated is a new pointer with Name "Bob"
//
//go:inline
func MakeLensRefWithName[GET ~func(*S) A, SET func(*S, A) *S, S, A any](get GET, set SET, name string) Lens[*S, A] {
	return MakeLensWithName(get, setCopy(set), name)
}

// MakeLensWithEq creates a [Lens] for pointer-based structures with equality optimization.
//
// This is similar to [MakeLensRef] but includes an optimization: if the new value equals
// the current value (according to the provided Eq predicate), the original pointer is returned
// unchanged instead of creating a copy. This can improve performance and reduce allocations
// when setting values that don't actually change the structure.
//
// The setter does not need to create a copy manually; this function automatically wraps it
// to ensure immutability when changes are made.
//
// This lens assumes that property A always exists in structure S (i.e., it's not optional).
//
// Type Parameters:
//   - GET: Getter function type (*S → A)
//   - SET: Setter function type (*S, A → *S)
//   - S: Source structure type (will be used as *S)
//   - A: Focus/field type
//
// Parameters:
//   - pred: Equality predicate to compare values of type A
//   - get: Function to extract value A from pointer *S
//   - set: Function to update value A in pointer *S (copying handled automatically)
//
// Returns:
//   - A Lens[*S, A] that can get and set values immutably on pointers with equality optimization
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLensWithEq(
//	    eq.FromStrictEquals[string](),
//	    func(p *Person) string { return p.Name },
//	    func(p *Person, name string) *Person {
//	        p.Name = name  // No manual copy needed
//	        return p
//	    },
//	)
//
//	person := &Person{Name: "Alice", Age: 30}
//
//	// Setting the same value returns the original pointer (no copy)
//	same := nameLens.Set("Alice")(person)
//	// same == person (same pointer)
//
//	// Setting a different value creates a new copy
//	updated := nameLens.Set("Bob")(person)
//	// person.Name is still "Alice", updated is a new pointer with Name "Bob"
//
//go:inline
func MakeLensWithEq[GET ~func(*S) A, SET func(*S, A) *S, S, A any](pred EQ.Eq[A], get GET, set SET) Lens[*S, A] {
	return MakeLens(get, setCopyWithEq(pred, get, set))
}

// MakeLensWithEqWithName creates a [Lens] for pointer-based structures with equality optimization and a custom name.
//
// This combines [MakeLensWithEq] (equality optimization) with [MakeLensWithName] (custom name for debugging).
// If the new value equals the current value (according to the provided Eq predicate), the original pointer
// is returned unchanged instead of creating a copy. The name is useful for debugging complex lens compositions.
//
// The setter does not need to create a copy manually; this function automatically wraps it
// to ensure immutability when changes are made.
//
// This lens assumes that property A always exists in structure S (i.e., it's not optional).
//
// Type Parameters:
//   - GET: Getter function type (*S → A)
//   - SET: Setter function type (*S, A → *S)
//   - S: Source structure type (will be used as *S)
//   - A: Focus/field type
//
// Parameters:
//   - pred: Equality predicate to compare values of type A
//   - get: Function to extract value A from pointer *S
//   - set: Function to update value A in pointer *S (copying handled automatically)
//   - name: A descriptive name for the lens (used in String() and Format())
//
// Returns:
//   - A Lens[*S, A] with equality optimization and the specified name
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLensWithEqWithName(
//	    eq.FromStrictEquals[string](),
//	    func(p *Person) string { return p.Name },
//	    func(p *Person, name string) *Person {
//	        p.Name = name  // No manual copy needed
//	        return p
//	    },
//	    "Person.Name",
//	)
//
//	person := &Person{Name: "Alice", Age: 30}
//	fmt.Printf("Using lens: %s\n", nameLens)  // Prints: "Using lens: Person.Name"
//
//	// Setting the same value returns the original pointer (no copy)
//	same := nameLens.Set("Alice")(person)  // same == person
//
//	// Setting a different value creates a new copy
//	updated := nameLens.Set("Bob")(person)  // person.Name still "Alice"
//
//go:inline
func MakeLensWithEqWithName[GET ~func(*S) A, SET func(*S, A) *S, S, A any](pred EQ.Eq[A], get GET, set SET, name string) Lens[*S, A] {
	return MakeLensWithName(get, setCopyWithEq(pred, get, set), name)
}

// MakeLensStrict creates a [Lens] for pointer-based structures with strict equality optimization.
//
// This is a convenience function that combines [MakeLensWithEq] with strict equality comparison (==).
// It's suitable for comparable types (primitives, strings, pointers, etc.) and provides the same
// optimization as MakeLensWithEq: if the new value equals the current value, the original pointer
// is returned unchanged instead of creating a copy.
//
// The setter does not need to create a copy manually; this function automatically wraps it
// to ensure immutability when changes are made.
//
// This lens assumes that property A always exists in structure S (i.e., it's not optional).
//
// Type Parameters:
//   - GET: Getter function type (*S → A)
//   - SET: Setter function type (*S, A → *S)
//   - S: Source structure type (will be used as *S)
//   - A: Focus/field type (must be comparable)
//
// Parameters:
//   - get: Function to extract value A from pointer *S
//   - set: Function to update value A in pointer *S (copying handled automatically)
//
// Returns:
//   - A Lens[*S, A] that can get and set values immutably on pointers with strict equality optimization
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	// Using MakeLensStrict for a string field (comparable type)
//	nameLens := lens.MakeLensStrict(
//	    func(p *Person) string { return p.Name },
//	    func(p *Person, name string) *Person {
//	        p.Name = name  // No manual copy needed
//	        return p
//	    },
//	)
//
//	person := &Person{Name: "Alice", Age: 30}
//
//	// Setting the same value returns the original pointer (no copy)
//	same := nameLens.Set("Alice")(person)
//	// same == person (same pointer)
//
//	// Setting a different value creates a new copy
//	updated := nameLens.Set("Bob")(person)
//	// person.Name is still "Alice", updated is a new pointer with Name "Bob"
//
//go:inline
func MakeLensStrict[GET ~func(*S) A, SET func(*S, A) *S, S any, A comparable](get GET, set SET) Lens[*S, A] {
	return MakeLensWithEq(EQ.FromStrictEquals[A](), get, set)
}

// MakeLensStrictWithName creates a [Lens] for pointer-based structures with strict equality optimization and a custom name.
//
// This combines [MakeLensStrict] (strict equality optimization using ==) with [MakeLensWithName]
// (custom name for debugging). It's a convenience function suitable for comparable types
// (primitives, strings, pointers, etc.). If the new value equals the current value, the original
// pointer is returned unchanged instead of creating a copy. The name is useful for debugging.
//
// The setter does not need to create a copy manually; this function automatically wraps it
// to ensure immutability when changes are made.
//
// This lens assumes that property A always exists in structure S (i.e., it's not optional).
//
// Type Parameters:
//   - GET: Getter function type (*S → A)
//   - SET: Setter function type (*S, A → *S)
//   - S: Source structure type (will be used as *S)
//   - A: Focus/field type (must be comparable)
//
// Parameters:
//   - get: Function to extract value A from pointer *S
//   - set: Function to update value A in pointer *S (copying handled automatically)
//   - name: A descriptive name for the lens (used in String() and Format())
//
// Returns:
//   - A Lens[*S, A] with strict equality optimization and the specified name
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	// Using MakeLensStrictWithName for a string field (comparable type)
//	nameLens := lens.MakeLensStrictWithName(
//	    func(p *Person) string { return p.Name },
//	    func(p *Person, name string) *Person {
//	        p.Name = name  // No manual copy needed
//	        return p
//	    },
//	    "Person.Name",
//	)
//
//	person := &Person{Name: "Alice", Age: 30}
//	fmt.Printf("Using lens: %s\n", nameLens)  // Prints: "Using lens: Person.Name"
//
//	// Setting the same value returns the original pointer (no copy)
//	same := nameLens.Set("Alice")(person)  // same == person
//
//	// Setting a different value creates a new copy
//	updated := nameLens.Set("Bob")(person)  // person.Name still "Alice"
//
//go:inline
func MakeLensStrictWithName[GET ~func(*S) A, SET func(*S, A) *S, S any, A comparable](get GET, set SET, name string) Lens[*S, A] {
	return MakeLensWithEqWithName(EQ.FromStrictEquals[A](), get, set, name)
}

// MakeLensRefCurried creates a [Lens] for pointer-based structures with a curried setter.
//
// This combines the benefits of [MakeLensRef] (automatic copying) with [MakeLensCurried]
// (curried setter for better composition). The setter does not need to create a copy manually;
// this function automatically wraps it to ensure immutability.
//
// This lens assumes that property A always exists in structure S (i.e., it's not optional).
//
// Type Parameters:
//   - S: Source structure type (will be used as *S)
//   - A: Focus/field type
//
// Parameters:
//   - get: Function to extract value A from pointer *S
//   - set: Curried function to update value A in pointer *S (copying handled automatically)
//
// Returns:
//   - A Lens[*S, A] that can get and set values immutably on pointers
//
// Example:
//
//	nameLens := lens.MakeLensRefCurried(
//	    func(p *Person) string { return p.Name },
//	    func(name string) func(*Person) *Person {
//	        return func(p *Person) *Person {
//	            p.Name = name  // No manual copy needed
//	            return p
//	        }
//	    },
//	)
//
//go:inline
func MakeLensRefCurried[S, A any](get func(*S) A, set func(A) Endomorphism[*S]) Lens[*S, A] {
	return MakeLensCurried(get, setCopyCurried(set))
}

// MakeLensRefCurriedWithName creates a [Lens] for pointer-based structures with a curried setter and custom name.
//
// This combines the benefits of [MakeLensRefCurried] (automatic copying with curried setter)
// with [MakeLensWithName] (custom name for debugging). The setter does not need to create
// a copy manually; this function automatically wraps it to ensure immutability. The curried
// form is more composable in functional pipelines, and the name is useful for debugging.
//
// This lens assumes that property A always exists in structure S (i.e., it's not optional).
//
// Type Parameters:
//   - S: Source structure type (will be used as *S)
//   - A: Focus/field type
//
// Parameters:
//   - get: Function to extract value A from pointer *S
//   - set: Curried function to update value A in pointer *S (copying handled automatically)
//   - name: A descriptive name for the lens (used in String() and Format())
//
// Returns:
//   - A Lens[*S, A] with the specified name
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	nameLens := lens.MakeLensRefCurriedWithName(
//	    func(p *Person) string { return p.Name },
//	    func(name string) func(*Person) *Person {
//	        return func(p *Person) *Person {
//	            p.Name = name  // No manual copy needed
//	            return p
//	        }
//	    },
//	    "Person.Name",
//	)
//
//	person := &Person{Name: "Alice", Age: 30}
//	fmt.Printf("Using lens: %s\n", nameLens)  // Prints: "Using lens: Person.Name"
//	updated := nameLens.Set("Bob")(person)
//	// person.Name is still "Alice", updated is a new pointer with Name "Bob"
//
//go:inline
func MakeLensRefCurriedWithName[S, A any](get func(*S) A, set func(A) Endomorphism[*S], name string) Lens[*S, A] {
	return MakeLensCurriedWithName(get, setCopyCurried(set), name)
}

// id returns a [Lens] implementing the identity operation
func id[GET ~func(S) S, SET ~func(S, S) S, S any](creator func(get GET, set SET, name string) Lens[S, S]) Lens[S, S] {
	return creator(F.Identity[S], F.Second[S, S], "LensIdentity")
}

// Id returns an identity [Lens] that focuses on the entire structure.
//
// The identity lens is useful as a starting point for lens composition or when you need
// a lens that doesn't actually focus on a subpart. Get returns the structure unchanged,
// and Set replaces the entire structure.
//
// Type Parameters:
//   - S: The structure type
//
// Returns:
//   - A Lens[S, S] where both source and focus are the same type
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	idLens := lens.Id[Person]()
//	person := Person{Name: "Alice", Age: 30}
//
//	same := idLens.Get(person)  // Returns person unchanged
//	replaced := idLens.Set(Person{Name: "Bob", Age: 25})(person)
//	// replaced is Person{Name: "Bob", Age: 25}
func Id[S any]() Lens[S, S] {
	return id(MakeLensWithName[Endomorphism[S], func(S, S) S])
}

// IdRef returns an identity [Lens] for pointer-based structures.
//
// This is the pointer version of [Id]. It focuses on the entire pointer structure,
// with automatic copying to ensure immutability.
//
// Type Parameters:
//   - S: The structure type (will be used as *S)
//
// Returns:
//   - A Lens[*S, *S] where both source and focus are pointers to the same type
//
// Example:
//
//	idLens := lens.IdRef[Person]()
//	person := &Person{Name: "Alice", Age: 30}
//
//	same := idLens.Get(person)  // Returns person pointer
//	replaced := idLens.Set(&Person{Name: "Bob", Age: 25})(person)
//	// person.Name is still "Alice", replaced is a new pointer
func IdRef[S any]() Lens[*S, *S] {
	return id(MakeLensRefWithName[Endomorphism[*S], func(*S, *S) *S])
}

// Compose combines two lenses and allows to narrow down the focus to a sub-lens
func compose[GET ~func(S) B, SET ~func(B) func(S) S, S, A, B any](
	creator func(get GET, set SET, name string) Lens[S, B],
	ab Lens[A, B],
) Operator[S, A, B] {
	abget := ab.Get
	abset := ab.Set
	return func(sa Lens[S, A]) Lens[S, B] {
		saget := sa.Get
		saset := sa.Set
		return creator(
			F.Flow2(saget, abget),
			func(b B) func(S) S {
				return endomorphism.Join(F.Flow3(
					saget,
					abset(b),
					saset,
				))
			},
			fmt.Sprintf("LensCompose[%s -> %s]", sa, ab),
		)
	}
}

// Compose combines two lenses to focus on a deeply nested field.
//
// Given a lens from S to A and a lens from A to B, Compose creates a lens from S to B.
// This allows you to navigate through nested structures in a composable way.
//
// The composition follows the mathematical property: (sa ∘ ab).Get = ab.Get ∘ sa.Get
//
// Type Parameters:
//   - S: Outer structure type
//   - A: Intermediate structure type
//   - B: Inner focus type
//
// Parameters:
//   - ab: Lens from A to B (inner lens)
//
// Returns:
//   - A function that takes a Lens[S, A] and returns a Lens[S, B]
//
// Example:
//
//	type Address struct {
//	    Street string
//	    City   string
//	}
//
//	type Person struct {
//	    Name    string
//	    Address Address
//	}
//
//	addressLens := lens.MakeLens(
//	    func(p Person) Address { return p.Address },
//	    func(p Person, a Address) Person { p.Address = a; return p },
//	)
//
//	streetLens := lens.MakeLens(
//	    func(a Address) string { return a.Street },
//	    func(a Address, s string) Address { a.Street = s; return a },
//	)
//
//	// Compose to access street directly from person
//	personStreetLens := F.Pipe1(addressLens, lens.Compose[Person](streetLens))
//
//	person := Person{Name: "Alice", Address: Address{Street: "Main St"}}
//	street := personStreetLens.Get(person)  // "Main St"
//	updated := personStreetLens.Set("Oak Ave")(person)
//
//go:inline
func Compose[S, A, B any](ab Lens[A, B]) Operator[S, A, B] {
	return compose(MakeLensCurriedWithName[func(S) B, func(B) func(S) S], ab)
}

// ComposeRef combines two lenses for pointer-based structures.
//
// This is the pointer version of [Compose], automatically handling copying to ensure immutability.
// It allows you to navigate through nested pointer structures in a composable way.
//
// Type Parameters:
//   - S: Outer structure type (will be used as *S)
//   - A: Intermediate structure type
//   - B: Inner focus type
//
// Parameters:
//   - ab: Lens from A to B (inner lens)
//
// Returns:
//   - A function that takes a Lens[*S, A] and returns a Lens[*S, B]
//
// Example:
//
//	type Address struct {
//	    Street string
//	}
//
//	type Person struct {
//	    Name    string
//	    Address Address
//	}
//
//	addressLens := lens.MakeLensRef(
//	    func(p *Person) Address { return p.Address },
//	    func(p *Person, a Address) *Person { p.Address = a; return p },
//	)
//
//	streetLens := lens.MakeLens(
//	    func(a Address) string { return a.Street },
//	    func(a Address, s string) Address { a.Street = s; return a },
//	)
//
//	personStreetLens := F.Pipe1(addressLens, lens.ComposeRef[Person](streetLens))
func ComposeRef[S, A, B any](ab Lens[A, B]) Operator[*S, A, B] {
	return compose(MakeLensRefCurriedWithName[S, B], ab)
}

// Modify transforms a value through a lens using a transformation F.
//
// Instead of setting a specific value, Modify applies a function to the current value.
// This is useful for updates like incrementing a counter, appending to a string, etc.
// If the transformation doesn't change the value, the original structure is returned.
//
// Type Parameters:
//   - S: Structure type
//   - FCT: Transformation function type (A → A)
//   - A: Focus type
//
// Parameters:
//   - f: Transformation function to apply to the focused value
//
// Returns:
//   - A function that takes a Lens[S, A] and returns an Endomorphism[S]
//
// Example:
//
//	type Counter struct {
//	    Value int
//	}
//
//	valueLens := lens.MakeLens(
//	    func(c Counter) int { return c.Value },
//	    func(c Counter, v int) Counter { c.Value = v; return c },
//	)
//
//	counter := Counter{Value: 5}
//
//	// Increment the counter
//	incremented := F.Pipe2(
//	    valueLens,
//	    lens.Modify[Counter](func(v int) int { return v + 1 }),
//	    F.Ap(counter),
//	)
//	// incremented.Value == 6
//
//	// Double the counter
//	doubled := F.Pipe2(
//	    valueLens,
//	    lens.Modify[Counter](func(v int) int { return v * 2 }),
//	    F.Ap(counter),
//	)
//	// doubled.Value == 10
func Modify[S any, FCT ~func(A) A, A any](f FCT) func(Lens[S, A]) Endomorphism[S] {
	return func(la Lens[S, A]) Endomorphism[S] {
		return endomorphism.Join(F.Flow3(
			la.Get,
			f,
			la.Set,
		))
	}
}

// IMap transforms the focus type of a lens using an isomorphism.
//
// An isomorphism is a pair of functions (A → B, B → A) that are inverses of each other.
// IMap allows you to work with a lens in a different but equivalent type. This is useful
// for unit conversions, encoding/decoding, or any bidirectional transformation.
//
// Type Parameters:
//   - E: Structure type
//   - AB: Forward transformation function type (A → B)
//   - BA: Backward transformation function type (B → A)
//   - A: Original focus type
//   - B: Transformed focus type
//
// Parameters:
//   - ab: Forward transformation (A → B)
//   - ba: Backward transformation (B → A)
//
// Returns:
//   - A function that takes a Lens[E, A] and returns a Lens[E, B]
//
// Example:
//
//	type Celsius float64
//	type Fahrenheit float64
//
//	celsiusToFahrenheit := func(c Celsius) Fahrenheit {
//	    return Fahrenheit(c*9/5 + 32)
//	}
//
//	fahrenheitToCelsius := func(f Fahrenheit) Celsius {
//	    return Celsius((f - 32) * 5 / 9)
//	}
//
//	type Weather struct {
//	    Temperature Celsius
//	}
//
//	tempCelsiusLens := lens.MakeLens(
//	    func(w Weather) Celsius { return w.Temperature },
//	    func(w Weather, t Celsius) Weather { w.Temperature = t; return w },
//	)
//
//	// Create a lens that works with Fahrenheit
//	tempFahrenheitLens := F.Pipe1(
//	    tempCelsiusLens,
//	    lens.IMap[Weather](celsiusToFahrenheit, fahrenheitToCelsius),
//	)
//
//	weather := Weather{Temperature: 20} // 20°C
//	tempF := tempFahrenheitLens.Get(weather)  // 68°F
//	updated := tempFahrenheitLens.Set(86)(weather)  // Set to 86°F (30°C)
func IMap[S any, AB ~func(A) B, BA ~func(B) A, A, B any](ab AB, ba BA) Operator[S, A, B] {
	return func(ea Lens[S, A]) Lens[S, B] {
		return MakeLensCurriedWithName(F.Flow2(ea.Get, ab), F.Flow2(ba, ea.Set), fmt.Sprintf("IMap[%s]", ea))
	}
}

// String returns the name of the lens as a string.
//
// This implements the fmt.Stringer interface, allowing lenses to be printed
// in a human-readable format for debugging and logging purposes. The returned
// string is the name provided when creating the lens with [MakeLensWithName]
// or [MakeLensCurriedWithName], or "Lens" for lenses created with other constructors.
//
// Returns:
//   - The lens name as a string
//
// Example:
//
//	nameLens := lens.MakeLensWithName(
//	    func(p Person) string { return p.Name },
//	    func(p Person, name string) Person { p.Name = name; return p },
//	    "Person.Name",
//	)
//	fmt.Println(nameLens)  // Prints: "Person.Name"
