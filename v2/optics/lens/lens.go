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
	EQ "github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/common"
)

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
//	nameLens := common.MakeLens(
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
	return common.MakeLens(get, set)
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
//	nameLens := common.MakeLensWithName(
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
	return common.MakeLensWithName(get, set, name)
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
//	nameLens := common.MakeLensCurried(
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
	return common.MakeLensCurried(get, set)
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
//	nameLens := common.MakeLensCurriedWithName(
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
	return common.MakeLensCurriedWithName(get, set, name)
}

// MakeLensCurriedRefWithName creates a [Lens] for pointer-based structures with a curried setter
// and a custom name.
//
// This is the named variant of [MakeLensCurriedRefWithName] combined with automatic copy-on-write
// semantics for pointer receivers. The setter does not need to create a copy manually; the copy
// is applied by the wrapped curried setter produced by setCopyCurried. The name is used in
// String, Format, and LogValue for debugging and structured logging.
//
// Type Parameters:
//   - GET: Getter function type (*S → A)
//   - SET: Curried setter function type (A → *S → *S)
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
//go:inline
func MakeLensCurriedRefWithName[GET ~func(*S) A, SET ~func(A) Endomorphism[*S], S, A any](get GET, set SET, name string) Lens[*S, A] {
	return common.MakeLensCurriedRefWithName(get, set, name)
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
//	nameLens := common.MakeLensRef(
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
	return common.MakeLensRef(get, set)
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
//	nameLens := common.MakeLensRefWithName(
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
	return common.MakeLensRefWithName(get, set, name)
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
//	nameLens := common.MakeLensWithEq(
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
	return common.MakeLensWithEq(pred, get, set)
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
//	nameLens := common.MakeLensWithEqWithName(
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
	return common.MakeLensWithEqWithName(pred, get, set, name)
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
//	nameLens := common.MakeLensStrict(
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
	return common.MakeLensStrict(get, set)
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
//	nameLens := common.MakeLensStrictWithName(
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
	return common.MakeLensStrictWithName(get, set, name)
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
//	nameLens := common.MakeLensRefCurried(
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
	return common.MakeLensRefCurried(get, set)
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
//	nameLens := common.MakeLensRefCurriedWithName(
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
	return common.MakeLensRefCurriedWithName(get, set, name)
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
	return common.LensId[S]()
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
	return common.LensIdRef[S]()
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
//	addressLens := common.MakeLens(
//	    func(p Person) Address { return p.Address },
//	    func(p Person, a Address) Person { p.Address = a; return p },
//	)
//
//	streetLens := common.MakeLens(
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
	return common.LensComposeLens[S](ab)
}

// ComposeIso combines a lens with an isomorphism to focus on a value of a different type.
//
// Given a lens from S to A and an isomorphism between A and B, ComposeIso creates a lens
// from S to B. The isomorphism is converted to a lens internally, then composed with the
// outer lens using the same rules as Compose.
//
// This is useful when the intermediate type A and the desired focus type B are related by a
// total, invertible conversion — for example, a newtype wrapper, a unit conversion, or a
// canonical representation.
//
// The resulting lens satisfies all three lens laws whenever both the outer lens and the
// isomorphism individually satisfy them.
//
// Type Parameters:
//   - S: Outer structure type
//   - A: Intermediate type — the focus of the outer lens and the source of the iso
//   - B: Inner focus type — the target of the iso and the focus of the resulting lens
//
// Parameters:
//   - ab: Isomorphism from A to B
//
// Returns:
//   - A function that takes a Lens[S, A] and returns a Lens[S, B]
//
// Example:
//
//	type Celsius float64
//	type Fahrenheit float64
//
//	type Thermometer struct {
//	    Reading Celsius
//	}
//
//	readingLens := common.MakeLens(
//	    func(t Thermometer) Celsius { return t.Reading },
//	    func(t Thermometer, c Celsius) Thermometer { t.Reading = c; return t },
//	)
//
//	celsiusToFahrenheit := common.MakeIso(
//	    func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
//	    func(f Fahrenheit) Celsius { return Celsius((f-32)*5/9) },
//	)
//
//	// Compose to read and write the temperature in Fahrenheit
//	fahrenheitLens := F.Pipe1(readingLens, lens.ComposeIso[Thermometer](celsiusToFahrenheit))
//
//	t := Thermometer{Reading: 100}
//	f := fahrenheitLens.Get(t)  // Fahrenheit(212)
//	updated := fahrenheitLens.Set(Fahrenheit(32))(t)
//	// updated.Reading == Celsius(0)
//
// See Also:
//   - Compose: the variant that accepts a Lens[A, B] instead of an Iso[A, B]
//
//go:inline
func ComposeIso[S, A, B any](ab Iso[A, B]) Operator[S, A, B] {
	return common.LensComposeIso[S](ab)
}

// ComposeRef combines two lenses for pointer-based structures.
//
// Deprecated: ComposeRef is not needed. When the outer lens is already a
// pointer lens (Lens[*S, A], created with MakeLensRef or MakeLensRefCurried),
// its Set implementation already copies *S before writing. The copy that
// ComposeRef adds via MakeLensRefCurriedWithName is therefore redundant —
// the composed setter calls the outer lens's Set, which copies, so the
// original pointer is never mutated.
//
// Use [Compose][*S] instead:
//
//	// Before
//	personStreetLens := F.Pipe1(addressLens, lens.ComposeRef[Person](streetLens))
//
//	// After — identical behaviour, no extra copy
//	personStreetLens := F.Pipe1(addressLens, lens.Compose[*Person](streetLens))
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
func ComposeRef[S, A, B any](ab Lens[A, B]) Operator[*S, A, B] {
	return common.LensComposeLensRef[S](ab)
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
//	valueLens := common.MakeLens(
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
	return common.LensModify[S](f)
}

// Set returns a function that updates the focus of a lens to a constant value.
//
// This is a convenience helper for partially applying a value before supplying the lens,
// making it useful in composition pipelines with F.Pipe.
//
// Example:
//
//	type Counter struct {
//	    Value int
//	}
//
//	valueLens := common.MakeLens(
//	    func(c Counter) int { return c.Value },
//	    func(c Counter, value int) Counter {
//	        c.Value = value
//	        return c
//	    },
//	)
//
//	counter := Counter{Value: 5}
//	updated := F.Pipe2(
//	    10,
//	    lens.Set[Counter](10),
//	    F.Ap(valueLens),
//	)
//	// updated.Value == 10
func Set[S any, A any](a A) func(Lens[S, A]) Endomorphism[S] {
	return common.LensSet[S](a)
}

// ModifyF transforms a value through a lens using a function that returns a value in a functor context.
//
// This is the functorial version of Modify, allowing transformations that produce effects
// (like Option, Either, IO, etc.) while updating the focused value. The functor's map operation
// is used to apply the lens's setter to the transformed value, preserving the computational context.
//
// This function corresponds to modifyF from monocle-ts, enabling effectful updates through lenses.
//
// # Type Parameters
//
//   - S: Structure type
//   - A: Focus type (the value being transformed)
//   - HKTA: Higher-kinded type containing the transformed value (e.g., Option[A], Either[E, A])
//   - HKTS: Higher-kinded type containing the updated structure (e.g., Option[S], Either[E, S])
//
// # Parameters
//
//   - fmap: A functor map operation that transforms A to S within the functor context
//
// # Returns
//
//   - A curried function that takes:
//     1. A transformation function (A → HKTA)
//     2. A Lens[S, A]
//     3. A structure S
//     And returns the updated structure in the functor context (HKTS)
//
// # Example Usage
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	ageLens := common.MakeLens(
//	    func(p Person) int { return p.Age },
//	    func(p Person, age int) Person { p.Age = age; return p },
//	)
//
//	// Validate age is positive, returning Option
//	validateAge := func(age int) option.Option[int] {
//	    if age > 0 {
//	        return option.Some(age)
//	    }
//	    return option.None[int]()
//	}
//
//	// Create a modifier that validates while updating
//	modifyAge := lens.ModifyF[Person, int](option.Functor[int, Person]().Map)
//
//	person := Person{Name: "Alice", Age: 30}
//	result := modifyAge(validateAge)(ageLens)(person)
//	// result is Some(Person{Name: "Alice", Age: 30})
//
//	invalidResult := modifyAge(func(age int) option.Option[int] {
//	    return option.None[int]()
//	})(ageLens)(person)
//	// invalidResult is None[Person]()
//
// # See Also
//
//   - Modify: Non-functorial version for simple transformations
//   - functor.Functor: The functor interface used for mapping
func ModifyF[S, A, HKTA, HKTS any](
	fmap functor.MapType[A, S, HKTA, HKTS],
) func(func(A) HKTA) func(Lens[S, A]) func(S) HKTS {
	return common.LensModifyF(fmap)
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
//	tempCelsiusLens := common.MakeLens(
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
	return common.LensIMap[S](ab, ba)
}
