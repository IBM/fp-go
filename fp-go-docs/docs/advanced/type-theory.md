---
title: Type Theory
hide_title: true
description: Theoretical foundations of functional programming including category theory, functors, monads, monoids, and algebraic data types.
sidebar_position: 2
---

<PageHeader
  eyebrow="Advanced · 02 / 04"
  title="Type"
  titleAccent="Theory"
  lede="Explore the theoretical foundations of functional programming including category theory, functors, applicatives, monads, and algebraic data types."
  meta={[
    { label: 'Difficulty', value: 'Expert' },
    { label: 'Topics', value: '10' },
    { label: 'Prerequisites', value: 'Abstract algebra, Category theory basics' }
  ]}
/>

<TLDR>
  <TLDRCard title="Category Theory" icon="network">
    Categories provide the mathematical foundation—objects (types), morphisms (functions), composition, and identity.
  </TLDRCard>
  <TLDRCard title="Functor Laws" icon="shield">
    Functors preserve structure: identity (fmap id = id) and composition (fmap (f ∘ g) = fmap f ∘ fmap g).
  </TLDRCard>
  <TLDRCard title="Monad Laws" icon="link">
    Monads enable sequencing: left identity, right identity, and associativity ensure predictable composition.
  </TLDRCard>
</TLDR>

<Section id="category-theory" number="01" title="Category Theory" titleAccent="Basics">

Categories consist of objects (types), morphisms (functions), composition, and identity.

<CodeCard file="category_basics.go">
{`package main

import "fmt"

// Identity morphism
func identity[A any](a A) A {
    return a
}

// Composition of morphisms
func compose[A, B, C any](f func(A) B, g func(B) C) func(A) C {
    return func(a A) C {
        return g(f(a))
    }
}

func main() {
    double := func(n int) int { return n * 2 }
    addTen := func(n int) int { return n + 10 }
    
    // Compose: double then addTen
    composed := compose(double, addTen)
    
    result := composed(5) // (5 * 2) + 10 = 20
    fmt.Println(result)
}
`}
</CodeCard>

<Callout type="info">
**Category Laws**: Composition must be associative: `(f ∘ g) ∘ h = f ∘ (g ∘ h)`, and identity must be neutral: `f ∘ id = id ∘ f = f`.
</Callout>

</Section>

<Section id="functors" number="02" title="Functors" titleAccent="& Laws">

Functors map between categories while preserving structure.

<CodeCard file="functor_laws.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
)

func main() {
    opt := O.Some(5)
    
    // Identity law: fmap(id) = id
    mapped1 := O.Map(func(x int) int { return x })(opt)
    fmt.Println(O.IsSome(mapped1)) // true
    
    // Composition law
    f := func(x int) int { return x * 2 }
    g := func(x int) int { return x + 10 }
    
    // fmap(f ∘ g)
    composed := O.Map(func(x int) int { return g(f(x)) })(opt)
    
    // fmap(f) ∘ fmap(g)
    separate := O.Map(g)(O.Map(f)(opt))
    
    fmt.Println(O.GetOrElse(func() int { return 0 })(composed))  // 20
    fmt.Println(O.GetOrElse(func() int { return 0 })(separate))  // 20
}
`}
</CodeCard>

<Callout type="success">
**Functor Examples**: Option, Either, Array, IO—all preserve structure when mapping functions over their contents.
</Callout>

</Section>

<Section id="applicatives" number="03" title="Applicative" titleAccent="Functors">

Applicatives allow applying wrapped functions to wrapped values.

<CodeCard file="applicative_laws.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
)

func main() {
    // Applicative: apply wrapped function to wrapped value
    optFunc := O.Some(func(n int) int { return n * 2 })
    optValue := O.Some(5)
    
    result := O.Ap(optValue)(optFunc)
    fmt.Println(O.GetOrElse(func() int { return 0 })(result)) // 10
    
    // With None
    noneFunc := O.None[func(int) int]()
    result2 := O.Ap(optValue)(noneFunc)
    fmt.Println(O.IsNone(result2)) // true
}
`}
</CodeCard>

<Callout type="info">
**Applicative Laws**: Identity, composition, homomorphism, and interchange ensure predictable behavior when applying functions.
</Callout>

</Section>

<Section id="monads" number="04" title="Monads" titleAccent="& Laws">

Monads enable sequencing computations with context.

<CodeCard file="monad_laws.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
)

func half(n int) O.Option[int] {
    if n%2 == 0 {
        return O.Some(n / 2)
    }
    return O.None[int]()
}

func main() {
    // Left identity: return(a) >>= f = f(a)
    left1 := O.Chain(half)(O.Some(10))
    left2 := half(10)
    fmt.Println(O.GetOrElse(func() int { return 0 })(left1)) // 5
    fmt.Println(O.GetOrElse(func() int { return 0 })(left2)) // 5
    
    // Right identity: m >>= return = m
    m := O.Some(10)
    right1 := O.Chain(func(n int) O.Option[int] { return O.Some(n) })(m)
    fmt.Println(O.GetOrElse(func() int { return 0 })(right1)) // 10
    
    // Associativity
    double := func(n int) O.Option[int] { return O.Some(n * 2) }
    
    // (m >>= half) >>= double
    assoc1 := O.Chain(double)(O.Chain(half)(O.Some(20)))
    
    // m >>= (x => half(x) >>= double)
    assoc2 := O.Chain(func(x int) O.Option[int] {
        return O.Chain(double)(half(x))
    })(O.Some(20))
    
    fmt.Println(O.GetOrElse(func() int { return 0 })(assoc1)) // 20
    fmt.Println(O.GetOrElse(func() int { return 0 })(assoc2)) // 20
}
`}
</CodeCard>

</Section>

<Section id="monoids" number="05" title="Monoids" titleAccent="& Semigroups">

Monoids provide associative operations with identity elements.

<CodeCard file="monoid_laws.go">
{`package main

import (
    "fmt"
)

type Monoid[A any] struct {
    Empty  A
    Concat func(A, A) A
}

func main() {
    // Integer addition monoid
    intAdd := Monoid[int]{
        Empty: 0,
        Concat: func(a, b int) int {
            return a + b
        },
    }
    
    // Associativity: (1 + 2) + 3 = 1 + (2 + 3)
    left := intAdd.Concat(intAdd.Concat(1, 2), 3)
    right := intAdd.Concat(1, intAdd.Concat(2, 3))
    fmt.Println(left == right) // true (both are 6)
    
    // Identity: 5 + 0 = 0 + 5 = 5
    leftId := intAdd.Concat(5, intAdd.Empty)
    rightId := intAdd.Concat(intAdd.Empty, 5)
    fmt.Println(leftId == 5 && rightId == 5) // true
    
    // String concatenation monoid
    stringConcat := Monoid[string]{
        Empty: "",
        Concat: func(a, b string) string {
            return a + b
        },
    }
    
    result := stringConcat.Concat("Hello", stringConcat.Concat(" ", "World"))
    fmt.Println(result) // Hello World
}
`}
</CodeCard>

<CodeCard file="semigroup.go">
{`package main

import (
    "fmt"
)

type Semigroup[A any] struct {
    Concat func(A, A) A
}

func main() {
    // Max semigroup (no identity for all integers)
    maxSemigroup := Semigroup[int]{
        Concat: func(a, b int) int {
            if a > b {
                return a
            }
            return b
        },
    }
    
    // Associativity
    left := maxSemigroup.Concat(maxSemigroup.Concat(5, 10), 3)
    right := maxSemigroup.Concat(5, maxSemigroup.Concat(10, 3))
    fmt.Println(left == right) // true (both are 10)
}
`}
</CodeCard>

</Section>

<Section id="natural-transformations" number="06" title="Natural" titleAccent="Transformations">

Natural transformations map between functors while preserving structure.

<CodeCard file="natural_transformation.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
    E "github.com/IBM/fp-go/v2/either"
)

// Natural transformation: Option -> Either
func optionToEither[A any](opt O.Option[A]) E.Either[string, A] {
    if O.IsSome(opt) {
        return E.Right[string](O.GetOrElse(func() A { var zero A; return zero })(opt))
    }
    return E.Left[A]("none")
}

func main() {
    opt1 := O.Some(42)
    opt2 := O.None[int]()
    
    either1 := optionToEither(opt1)
    either2 := optionToEither(opt2)
    
    fmt.Println(E.IsRight(either1)) // true
    fmt.Println(E.IsLeft(either2))  // true
}
`}
</CodeCard>

</Section>

<Section id="algebraic-types" number="07" title="Algebraic Data" titleAccent="Types">

Sum types (OR) and product types (AND) form the basis of type algebra.

<Compare>
<CompareCol kind="bad">
<CodeCard file="sum_type.go">
{`// Sum type: Success OR Failure
type Result[E, A any] interface {
    isResult()
}

type Success[E, A any] struct {
    Value A
}

func (Success[E, A]) isResult() {}

type Failure[E, A any] struct {
    Error E
}

func (Failure[E, A]) isResult() {}

func divide(a, b int) Result[string, int] {
    if b == 0 {
        return Failure[string, int]{Error: "division by zero"}
    }
    return Success[string, int]{Value: a / b}
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="product_type.go">
{`// Product type: A AND B
type Tuple[A, B any] struct {
    First  A
    Second B
}

func main() {
    // Product of string and int
    tuple := Tuple[string, int]{
        First:  "Alice",
        Second: 30,
    }
    
    fmt.Printf("Name: %s, Age: %d\\n", tuple.First, tuple.Second)
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>

<Section id="kleisli" number="08" title="Kleisli" titleAccent="Composition">

Compose monadic functions for elegant pipelines.

<CodeCard file="kleisli.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
)

// Kleisli composition: (a -> m b) -> (b -> m c) -> (a -> m c)
func kleisli[A, B, C any](
    f func(A) O.Option[B],
    g func(B) O.Option[C],
) func(A) O.Option[C] {
    return func(a A) O.Option[C] {
        return O.Chain(g)(f(a))
    }
}

func half(n int) O.Option[int] {
    if n%2 == 0 {
        return O.Some(n / 2)
    }
    return O.None[int]()
}

func double(n int) O.Option[int] {
    return O.Some(n * 2)
}

func main() {
    // Compose half and double
    composed := kleisli(half, double)
    
    result := composed(10)
    fmt.Println(O.GetOrElse(func() int { return 0 })(result)) // 10
}
`}
</CodeCard>

</Section>

<Section id="hkt" number="09" title="Higher-Kinded" titleAccent="Types">

Simulate higher-kinded types in Go for generic abstractions.

<CodeCard file="hkt_simulation.go">
{`package main

import "fmt"

// HKT interface
type HKT[F any, A any] interface {
    unwrap() F
}

// Option HKT
type OptionHKT[A any] struct {
    value *A
}

func (o OptionHKT[A]) unwrap() *A {
    return o.value
}

// Functor type class
type Functor[F any] interface {
    Map(f func(any) any, fa F) F
}

// Option functor instance
type OptionFunctor struct{}

func (OptionFunctor) Map(f func(any) any, fa any) any {
    opt := fa.(OptionHKT[any])
    if opt.value == nil {
        return opt
    }
    result := f(*opt.value)
    return OptionHKT[any]{value: &result}
}

func main() {
    opt := OptionHKT[int]{value: new(int)}
    *opt.value = 5
    
    functor := OptionFunctor{}
    mapped := functor.Map(func(x any) any {
        return x.(int) * 2
    }, opt)
    
    result := mapped.(OptionHKT[any])
    if result.value != nil {
        fmt.Println(*result.value) // 10
    }
}
`}
</CodeCard>

</Section>

<Section id="yoneda" number="10" title="Yoneda" titleAccent="Lemma">

The Yoneda lemma relates natural transformations to functor elements.

<CodeCard file="yoneda.go">
{`package main

import (
    "fmt"
)

// Yoneda encoding
type Yoneda[F any, A any] struct {
    Run func(func(A) any) F
}

// Lower: Yoneda F A -> F A
func Lower[F, A any](y Yoneda[F, A]) F {
    return y.Run(func(a A) any { return a })
}

// Lift: F A -> Yoneda F A
func Lift[F, A any](fa F, fmap func(func(A) any, F) F) Yoneda[F, A] {
    return Yoneda[F, A]{
        Run: func(f func(A) any) F {
            return fmap(f, fa)
        },
    }
}

func main() {
    // Example with slice as functor
    slice := []int{1, 2, 3}
    
    fmap := func(f func(int) any, s []int) []any {
        result := make([]any, len(s))
        for i, v := range s {
            result[i] = f(v)
        }
        return result
    }
    
    // Lift to Yoneda
    yoneda := Lift(slice, fmap)
    
    // Apply transformation
    doubled := yoneda.Run(func(n int) any { return n * 2 })
    
    fmt.Println(doubled) // [2 4 6]
}
`}
</CodeCard>

</Section>
