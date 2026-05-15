---
sidebar_position: 15
title: Higher-Kinded Types
hide_title: true
description: Understand how fp-go simulates higher-kinded types in Go, and why this matters for generic functional programming.
---

<PageHeader
  eyebrow="Concepts · 05 / 06"
  title="Higher-Kinded"
  titleAccent="Types."
  lede="Understand how fp-go simulates higher-kinded types in Go. Learn the patterns, limitations, and practical approaches for generic functional programming."
  meta={[
    {label: '// Difficulty', value: 'Advanced'},
    {label: '// Reading time', value: '12 min · 7 sections'},
    {label: '// Prereqs', value: 'Monads, composition'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Go's limitation" prose value={<>No type constructor <em>parameters</em>.</>} />
  <TLDRCard label="// fp-go's solution" prose value={<>Consistent API across <em>all types</em>.</>} />
  <TLDRCard label="// Trade-off" prose value={<>Some duplication for <em>simplicity</em>.</>} variant="up" />
</TLDR>

<Section id="what-are-hkts" number="01" title="What are" titleAccent="HKTs?">

**Higher-Kinded Types (HKTs)** are types that take other types as parameters.

### Regular Types (Kind *)

<CodeCard file="regular-types.go">
{`// Regular types - concrete
int           // A number
string        // Text
User          // A struct`}
</CodeCard>

### Generic Types (Kind * → *)

<CodeCard file="generic-types.go">
{`// Generic types - take one type parameter
Option[int]      // Option of int
Result[string]   // Result of string
Array[User]      // Array of User

// The type constructor:
Option[_]        // Takes a type, returns a type
Result[_]        // Takes a type, returns a type`}
</CodeCard>

### Higher-Kinded Types (Kind (* → *) → *)

<CodeCard file="hkts.go">
{`// HKTs - take a type constructor as parameter
// (This is what Go doesn't support)

// Hypothetical syntax:
Functor[F[_]]    // F is a type constructor
Monad[M[_]]      // M is a type constructor

// Would let us write:
func Map[F[_], A, B](f func(A) B) func(F[A]) F[B]

// Works for ANY F: Option, Result, Array, etc.`}
</CodeCard>

</Section>

<Section id="why-no-hkts" number="02" title="Why Go doesn't have" titleAccent="HKTs.">

### Go's Type System

Go 1.18+ has generics, but they're limited:

<Compare>
<CompareCol kind="good">

<CodeCard file="go-can-do.go">
{`// ✅ Can do: concrete type parameters
func Map[A, B any](f func(A) B, slice []A) []B`}
</CodeCard>

</CompareCol>
<CompareCol kind="bad">

<CodeCard file="go-cannot-do.go">
{`// ❌ Can't do: type constructor parameters
func Map[F[_], A, B any](f func(A) B, fa F[A]) F[B]
//         ^^^
//         Not allowed in Go`}
</CodeCard>

</CompareCol>
</Compare>

**Why?**
- Simpler type system
- Easier to implement
- Faster compilation
- Go philosophy: simplicity over power

</Section>

<Section id="workarounds" number="03" title="How fp-go works" titleAccent="around this.">

fp-go uses several techniques to simulate HKTs:

### Technique 1: Separate Functions per Type

Instead of one generic `Map`, we have:

<CodeCard file="separate-functions.go">
{`// Option.Map
func Map[B, A any](f func(A) B) func(Option[A]) Option[B]

// Result.Map
func Map[B, A any](f func(A) B) func(Result[A]) Result[B]

// Array.Map
func Map[B, A any](f func(A) B) func([]A) []B

// IO.Map
func Map[B, A any](f func(A) B) func(IO[A]) IO[B]`}
</CodeCard>

**Pros:**
- ✅ Works in Go
- ✅ Type-safe
- ✅ Clear which type you're using

**Cons:**
- ⚠️ Code duplication
- ⚠️ Can't write generic algorithms

### Technique 2: Consistent API

All types follow the same pattern:

<CodeCard file="consistent-api.go">
{`// Every monad has these functions:
Of[A](a A) M[A]                              // Put value in monad
Map[B, A](f func(A) B) func(M[A]) M[B]       // Transform value
Chain[B, A](f func(A) M[B]) func(M[A]) M[B]  // Chain operations
Fold[B, A](/* ... */) func(M[A]) B           // Extract value

// Example with Option:
option.Of(42)
option.Map(func(x int) int { return x * 2 })
option.Chain(func(x int) option.Option[int] { return option.Some(x) })
option.Fold(func() int { return 0 }, func(x int) int { return x })

// Example with Result:
result.Ok(42)
result.Map(func(x int) int { return x * 2 })
result.Chain(func(x int) result.Result[int] { return result.Ok(x) })
result.Fold(func(err error) int { return 0 }, func(x int) int { return x })`}
</CodeCard>

<Callout type="success">
**Benefit:** Learn once, use everywhere.
</Callout>

### Technique 3: Code Generation

fp-go uses code generation to create similar functions for each type:

<CodeCard file="codegen.go">
{`// Generated from template
//go:generate go run gen.go

// Generates:
// - option/map.go
// - result/map.go
// - array/map.go
// - etc.`}
</CodeCard>

**Benefit:** Consistency without manual duplication.

</Section>

<Section id="type-parameters" number="04" title="Understanding type" titleAccent="parameters.">

### Type Parameter Order in fp-go v2

fp-go v2 reordered type parameters for better inference:

<Compare>
<CompareCol kind="bad">

<CodeCard file="v1-ordering.go">
{`// v1: inferrable parameters first
func Map[A, B any](f func(A) B) func(Option[A]) Option[B]
//       ^  ^
//       |  Can't infer B from function signature
//       Can infer A from function argument`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="v2-ordering.go">
{`// v2: non-inferrable parameters first
func Map[B, A any](f func(A) B) func(Option[A]) Option[B]
//       ^  ^
//       |  Can infer A from function argument
//       Can't infer B, so comes first`}
</CodeCard>

</CompareCol>
</Compare>

**Why?** Go can infer trailing type parameters but not leading ones.

### Example

<CodeCard file="inference-example.go">
{`// With v2 ordering
opt := option.Some(5)

// Go can infer types
doubled := option.Map(func(x int) string {
    return strconv.Itoa(x * 2)
})(opt)
// Go infers: Map[string, int]
//                 ^      ^
//                 B      A (from function)

// No need to specify:
// option.Map[string, int](...)`}
</CodeCard>

</Section>

<Section id="practical-implications" number="05" title="Practical" titleAccent="implications.">

### 1. Learn the Pattern Once

All fp-go types follow the same pattern:

<CodeCard file="pattern.go">
{`// Pattern for any monad M:

// Create
M.Of(value)              // or M.Some, M.Ok, etc.

// Transform
M.Map(transform)         // Change the value

// Chain
M.Chain(operation)       // Sequence operations

// Extract
M.Fold(onError, onSuccess)  // Get the value out`}
</CodeCard>

### 2. Use Type-Specific Functions

<Compare>
<CompareCol kind="bad">

<CodeCard file="generic-impossible.go">
{`// Can't write generic code like:
func Process[M[_], A, B](m M[A], f func(A) B) M[B] {
    return M.Map(f)(m)  // Not possible in Go
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="specific-functions.go">
{`// Instead, write specific functions:
func ProcessOption[A, B](opt option.Option[A], f func(A) B) option.Option[B] {
    return option.Map(f)(opt)
}

func ProcessResult[A, B](res result.Result[A], f func(A) B) result.Result[B] {
    return result.Map(f)(res)
}`}
</CodeCard>

</CompareCol>
</Compare>

### 3. Embrace Go's Simplicity

<Callout type="info">
Instead of fighting Go's type system, work with it:
- ✅ Use specific types
- ✅ Accept some duplication
- ✅ Focus on clarity
</Callout>

</Section>

<Section id="comparison" number="06" title="Comparing with other" titleAccent="languages.">

### Haskell (Has HKTs)

<CodeCard file="haskell.hs" lang="haskell">
{`-- Generic map for any Functor
fmap :: Functor f => (a -> b) -> f a -> f b

-- Works for all:
fmap (+1) (Just 5)        -- Maybe Int
fmap (+1) [1,2,3]         -- List Int
fmap (+1) (Right 5)       -- Either e Int`}
</CodeCard>

### TypeScript (Simulates HKTs)

<CodeCard file="typescript.ts" lang="typescript">
{`// Type-level programming
interface Functor<F> {
  map<A, B>(f: (a: A) => B): (fa: F<A>) => F<B>
}

// Works for Option, Result, Array, etc.`}
</CodeCard>

### Go (No HKTs)

<CodeCard file="go-approach.go">
{`// Separate functions
option.Map[B, A](f func(A) B) func(option.Option[A]) option.Option[B]
result.Map[B, A](f func(A) B) func(result.Result[A]) result.Result[B]
array.Map[B, A](f func(A) B) func([]A) []B

// More verbose, but simpler`}
</CodeCard>

</Section>

<Section id="practical-advice" number="07" title="Practical" titleAccent="advice.">

### 1. Don't Fight the Type System

<Compare>
<CompareCol kind="bad">

<CodeCard file="too-generic.go">
{`// ❌ Don't try to be too generic
func Process[???](m ???) ??? {
    // Impossible in Go
}`}
</CodeCard>

</CompareCol>
<CompareCol kind="good">

<CodeCard file="specific-clear.go">
{`// ✅ Write specific, clear code
func ProcessUser(res result.Result[User]) result.Result[UserDTO] {
    return result.Map(toDTO)(res)
}`}
</CodeCard>

</CompareCol>
</Compare>

### 2. Embrace Duplication When Needed

<CodeCard file="duplication-ok.go">
{`// Some duplication is okay
func ProcessUsers(users []User) []UserDTO {
    return array.Map(toDTO)(users)
}

func ProcessUserResult(res result.Result[User]) result.Result[UserDTO] {
    return result.Map(toDTO)(res)
}

// Clear and type-safe`}
</CodeCard>

### 3. Use Consistent Patterns

<CodeCard file="consistent-patterns.go">
{`// Learn the pattern once
// Apply to all types

// Option
option.Map(f)(opt)
option.Chain(g)(opt)

// Result
result.Map(f)(res)
result.Chain(g)(res)

// Array
array.Map(f)(arr)
array.Chain(g)(arr)`}
</CodeCard>

### 4. Focus on Value

<Callout type="success">
HKTs are a means, not an end. Focus on:
- Clear code
- Type safety
- Composability
- Maintainability

Not on:
- Maximum abstraction
- Minimal duplication
- Theoretical purity
</Callout>

</Section>
