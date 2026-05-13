---
title: State
hide_title: true
description: Functional state management - stateful computations that transform state and produce values.
sidebar_position: 23
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="State"
  lede="Functional state management. State[S, A] represents a stateful computation that transforms state of type S and produces a value of type A."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/state' },
    { label: 'Type', value: 'func(S) Pair[A, S]' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

State encapsulates computations that read and modify state in a functional way:
- **Pure**: No mutable state
- **Composable**: Build complex state logic from simple pieces
- **Explicit**: State transformations are explicit

<CodeCard file="type_definition.go">
{`package state

// State transforms state S and produces value A
type State[S, A any] = func(S) pair.Pair[A, S]
`}
</CodeCard>

### When to Use

<ApiTable>
| Use Case | Example |
|----------|---------|
| Stateful computations | Counters, accumulators |
| Parsers | Stateful parsing |
| Functional state management | Avoid mutable state |
| Composable transformations | Build complex state logic |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[S, A any](value A) State[S, A]` | Wrap value, keep state unchanged |
| `Get` | `func Get[S any]() State[S, S]` | Get current state as value |
| `Put` | `func Put[S any](s S) State[S, unit.Unit]` | Set new state |
| `Modify` | `func Modify[S any](f func(S) S) State[S, unit.Unit]` | Modify state with function |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[S, A, B any](f func(A) B) func(State[S, A]) State[S, B]` | Transform value, keep state |
| `Chain` | `func Chain[S, A, B any](f func(A) State[S, B]) func(State[S, A]) State[S, B]` | Sequence stateful operations |
| `Ap` | `func Ap[S, A, B any](fa State[S, A]) func(State[S, func(A) B]) State[S, B]` | Apply wrapped function |
</ApiTable>

### Execution

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Evaluate` | `func Evaluate[S, A any](s S) func(State[S, A]) A` | Run and get value |
| `Execute` | `func Execute[S, A any](s S) func(State[S, A]) S` | Run and get final state |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "fmt"
    S "github.com/IBM/fp-go/v2/state"
    "github.com/IBM/fp-go/v2/pair"
)

func main() {
    // Get current state
    get := S.Get[int]()
    result := get(42)  // Pair(42, 42)
    fmt.Println(pair.Fst(result), pair.Snd(result))  // 42 42
    
    // Set new state
    set := S.Put[int](100)
    result = set(42)  // Pair(unit.VOID, 100)
    fmt.Println(pair.Snd(result))  // 100
    
    // Modify state
    modify := S.Modify[int](func(s int) int {
        return s + 1
    })
    result = modify(42)  // Pair(unit.VOID, 43)
    fmt.Println(pair.Snd(result))  // 43
}
`}
</CodeCard>

### Counter Example

<CodeCard file="counter.go">
{`package main

import (
    "fmt"
    S "github.com/IBM/fp-go/v2/state"
    F "github.com/IBM/fp-go/v2/function"
)

func Increment() S.State[int, unit.Unit] {
    return S.Modify[int](func(s int) int {
        return s + 1
    })
}

func Decrement() S.State[int, unit.Unit] {
    return S.Modify[int](func(s int) int {
        return s - 1
    })
}

func GetCount() S.State[int, int] {
    return S.Get[int]()
}

func main() {
    // Compose operations
    counter := F.Pipe3(
        Increment(),
        S.Chain(func(_ unit.Unit) S.State[int, unit.Unit] {
            return Increment()
        }),
        S.Chain(func(_ unit.Unit) S.State[int, int] {
            return GetCount()
        }),
    )
    
    result := counter(0)
    value := pair.Fst(result)
    state := pair.Snd(result)
    fmt.Printf("Value: %d, State: %d\n", value, state)  // Value: 2, State: 2
}
`}
</CodeCard>

### Stack Example

<CodeCard file="stack.go">
{`package main

import (
    "fmt"
    S "github.com/IBM/fp-go/v2/state"
    O "github.com/IBM/fp-go/v2/option"
)

type Stack[A any] []A

func Push[A any](item A) S.State[Stack[A], unit.Unit] {
    return S.Modify[Stack[A]](func(stack Stack[A]) Stack[A] {
        return append(stack, item)
    })
}

func Pop[A any]() S.State[Stack[A], O.Option[A]] {
    return func(stack Stack[A]) pair.Pair[O.Option[A], Stack[A]] {
        if len(stack) == 0 {
            return pair.MakePair(O.None[A](), stack)
        }
        item := stack[len(stack)-1]
        newStack := stack[:len(stack)-1]
        return pair.MakePair(O.Some(item), newStack)
    }
}

func Peek[A any]() S.State[Stack[A], O.Option[A]] {
    return func(stack Stack[A]) pair.Pair[O.Option[A], Stack[A]] {
        if len(stack) == 0 {
            return pair.MakePair(O.None[A](), stack)
        }
        return pair.MakePair(O.Some(stack[len(stack)-1]), stack)
    }
}

func main() {
    // Build stack operations
    operations := F.Pipe4(
        Push(1),
        S.Chain(func(_ unit.Unit) S.State[Stack[int], unit.Unit] {
            return Push(2)
        }),
        S.Chain(func(_ unit.Unit) S.State[Stack[int], unit.Unit] {
            return Push(3)
        }),
        S.Chain(func(_ unit.Unit) S.State[Stack[int], O.Option[int]] {
            return Pop[int]()
        }),
    )
    
    result := operations(Stack[int]{})
    value := pair.Fst(result)
    state := pair.Snd(result)
    
    fmt.Printf("Popped: %v, Stack: %v\n", value, state)
    // Popped: Some(3), Stack: [1 2]
}
`}
</CodeCard>

### Parser Example

<CodeCard file="parser.go">
{`package main

import (
    "fmt"
    "strings"
    S "github.com/IBM/fp-go/v2/state"
    O "github.com/IBM/fp-go/v2/option"
)

type ParserState struct {
    Input string
    Pos   int
}

func Char(c rune) S.State[ParserState, O.Option[rune]] {
    return func(state ParserState) pair.Pair[O.Option[rune], ParserState] {
        if state.Pos >= len(state.Input) {
            return pair.MakePair(O.None[rune](), state)
        }
        
        current := rune(state.Input[state.Pos])
        if current == c {
            newState := ParserState{
                Input: state.Input,
                Pos:   state.Pos + 1,
            }
            return pair.MakePair(O.Some(current), newState)
        }
        
        return pair.MakePair(O.None[rune](), state)
    }
}

func String(s string) S.State[ParserState, O.Option[string]] {
    return func(state ParserState) pair.Pair[O.Option[string], ParserState] {
        if state.Pos+len(s) > len(state.Input) {
            return pair.MakePair(O.None[string](), state)
        }
        
        substr := state.Input[state.Pos : state.Pos+len(s)]
        if substr == s {
            newState := ParserState{
                Input: state.Input,
                Pos:   state.Pos + len(s),
            }
            return pair.MakePair(O.Some(s), newState)
        }
        
        return pair.MakePair(O.None[string](), state)
    }
}

func main() {
    parser := String("hello")
    
    state := ParserState{Input: "hello world", Pos: 0}
    result := parser(state)
    
    value := pair.Fst(result)
    newState := pair.Snd(result)
    
    fmt.Printf("Parsed: %v, Remaining: %s\n", value, newState.Input[newState.Pos:])
    // Parsed: Some(hello), Remaining:  world
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Accumulator

<CodeCard file="accumulator.go">
{`package main

import (
    S "github.com/IBM/fp-go/v2/state"
    F "github.com/IBM/fp-go/v2/function"
)

func AddToSum(n int) S.State[int, unit.Unit] {
    return S.Modify[int](func(sum int) int {
        return sum + n
    })
}

func GetSum() S.State[int, int] {
    return S.Get[int]()
}

func main() {
    // Sum numbers
    sumNumbers := F.Pipe4(
        AddToSum(10),
        S.Chain(func(_ unit.Unit) S.State[int, unit.Unit] {
            return AddToSum(20)
        }),
        S.Chain(func(_ unit.Unit) S.State[int, unit.Unit] {
            return AddToSum(30)
        }),
        S.Chain(func(_ unit.Unit) S.State[int, int] {
            return GetSum()
        }),
    )
    
    result := sumNumbers(0)
    fmt.Println(pair.Fst(result))  // 60
}
`}
</CodeCard>

### Pattern 2: Random Number Generator

<CodeCard file="random.go">
{`package main

import (
    S "github.com/IBM/fp-go/v2/state"
)

type RNG struct {
    Seed int64
}

func NextInt() S.State[RNG, int] {
    return func(rng RNG) pair.Pair[int, RNG] {
        // Linear congruential generator
        newSeed := (rng.Seed*1103515245 + 12345) & 0x7fffffff
        value := int(newSeed % 100)
        return pair.MakePair(value, RNG{Seed: newSeed})
    }
}

func main() {
    // Generate 3 random numbers
    gen := F.Pipe3(
        NextInt(),
        S.Chain(func(n1 int) S.State[RNG, int] {
            return NextInt()
        }),
        S.Chain(func(n2 int) S.State[RNG, int] {
            return NextInt()
        }),
    )
    
    result := gen(RNG{Seed: 42})
    fmt.Println(pair.Fst(result))
}
`}
</CodeCard>

</Section>

<Callout type="info">

**Pure State Management**: State monad provides pure functional state management without mutable variables. All state transformations are explicit and composable.

</Callout>
