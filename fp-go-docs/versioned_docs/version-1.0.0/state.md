---
sidebar_position: 13
---

# State (v1)

The `State` type represents a stateful computation that can read and modify state.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [State v2](../v2/state).

**Key differences in v2:**
- Simplified API
- Better type inference
- Improved composition
:::

## Overview

`State` represents a computation that:
- Takes an initial state
- Returns a value and a new state
- Enables functional state management

```go
type State[S, A any] func(S) (A, S)
```

## Creating State Values

### Of (Pure Value)

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

func main() {
    // Create State that returns value without changing state
    state := S.Of[int, string]("hello")
    
    value, newState := state(42)
    
    fmt.Println("Value:", value)      // hello
    fmt.Println("State:", newState)   // 42 (unchanged)
}
```

### Get

Read the current state:

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

func main() {
    // Get returns the state as the value
    state := S.Get[int]()
    
    value, newState := state(42)
    
    fmt.Println("Value:", value)      // 42
    fmt.Println("State:", newState)   // 42
}
```

### Put

Replace the state:

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

func main() {
    // Put replaces the state
    state := S.Put[int](100)
    
    value, newState := state(42)
    
    fmt.Println("Value:", value)      // {} (unit)
    fmt.Println("State:", newState)   // 100
}
```

### Modify

Transform the state:

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

func main() {
    // Modify transforms the state
    state := S.Modify[int](func(s int) int {
        return s * 2
    })
    
    value, newState := state(21)
    
    fmt.Println("Value:", value)      // {} (unit)
    fmt.Println("State:", newState)   // 42
}
```

## Basic Operations

### Map

Transform the value:

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

func main() {
    state := S.Of[int, int](5)
    
    // Map transforms the value
    doubled := S.Map(func(n int) int {
        return n * 2
    })(state)
    
    value, newState := doubled(42)
    
    fmt.Println("Value:", value)      // 10
    fmt.Println("State:", newState)   // 42
}
```

### Chain

Chain stateful computations:

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

func increment() S.State[int, int] {
    return func(s int) (int, int) {
        newState := s + 1
        return newState, newState
    }
}

func double() S.State[int, int] {
    return func(s int) (int, int) {
        newState := s * 2
        return newState, newState
    }
}

func main() {
    // Chain operations
    state := S.Chain(double)(increment())
    
    value, finalState := state(5)
    
    fmt.Println("Value:", value)         // 12
    fmt.Println("Final State:", finalState) // 12
}
```

## Practical Examples

### Counter

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

type Counter struct {
    Count int
}

func increment() S.State[Counter, int] {
    return func(c Counter) (int, Counter) {
        newCount := c.Count + 1
        return newCount, Counter{Count: newCount}
    }
}

func decrement() S.State[Counter, int] {
    return func(c Counter) (int, Counter) {
        newCount := c.Count - 1
        return newCount, Counter{Count: newCount}
    }
}

func getCount() S.State[Counter, int] {
    return func(c Counter) (int, Counter) {
        return c.Count, c
    }
}

func main() {
    initial := Counter{Count: 0}
    
    // Increment
    value1, state1 := increment()(initial)
    fmt.Println("After increment:", value1) // 1
    
    // Increment again
    value2, state2 := increment()(state1)
    fmt.Println("After increment:", value2) // 2
    
    // Decrement
    value3, state3 := decrement()(state2)
    fmt.Println("After decrement:", value3) // 1
    
    // Get final count
    final, _ := getCount()(state3)
    fmt.Println("Final count:", final) // 1
}
```

### Stack Operations

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

type Stack struct {
    Items []int
}

func push(item int) S.State[Stack, struct{}] {
    return func(s Stack) (struct{}, Stack) {
        newItems := append([]int{item}, s.Items...)
        return struct{}{}, Stack{Items: newItems}
    }
}

func pop() S.State[Stack, *int] {
    return func(s Stack) (*int, Stack) {
        if len(s.Items) == 0 {
            return nil, s
        }
        item := s.Items[0]
        newItems := s.Items[1:]
        return &item, Stack{Items: newItems}
    }
}

func peek() S.State[Stack, *int] {
    return func(s Stack) (*int, Stack) {
        if len(s.Items) == 0 {
            return nil, s
        }
        return &s.Items[0], s
    }
}

func main() {
    initial := Stack{Items: []int{}}
    
    // Push items
    _, state1 := push(1)(initial)
    _, state2 := push(2)(state1)
    _, state3 := push(3)(state2)
    
    // Peek
    top, state4 := peek()(state3)
    if top != nil {
        fmt.Println("Top:", *top) // 3
    }
    
    // Pop
    popped, state5 := pop()(state4)
    if popped != nil {
        fmt.Println("Popped:", *popped) // 3
    }
    
    // Peek again
    newTop, _ := peek()(state5)
    if newTop != nil {
        fmt.Println("New top:", *newTop) // 2
    }
}
```

### Random Number Generator

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
)

type RNG struct {
    Seed int64
}

func nextInt() S.State[RNG, int] {
    return func(rng RNG) (int, RNG) {
        // Simple LCG algorithm
        newSeed := (rng.Seed*1103515245 + 12345) & 0x7fffffff
        value := int(newSeed % 100)
        return value, RNG{Seed: newSeed}
    }
}

func nextIntInRange(min, max int) S.State[RNG, int] {
    return S.Map(func(n int) int {
        return min + (n % (max - min + 1))
    })(nextInt())
}

func main() {
    rng := RNG{Seed: 42}
    
    // Generate random numbers
    n1, rng1 := nextInt()(rng)
    n2, rng2 := nextInt()(rng1)
    n3, _ := nextInt()(rng2)
    
    fmt.Println("Random numbers:", n1, n2, n3)
    
    // Generate in range
    dice, _ := nextIntInRange(1, 6)(rng)
    fmt.Println("Dice roll:", dice)
}
```

### Parser State

```go
package main

import (
    "fmt"
    "strings"
    S "github.com/IBM/fp-go/state"
)

type ParserState struct {
    Input string
    Pos   int
}

func char() S.State[ParserState, *rune] {
    return func(ps ParserState) (*rune, ParserState) {
        if ps.Pos >= len(ps.Input) {
            return nil, ps
        }
        c := rune(ps.Input[ps.Pos])
        return &c, ParserState{
            Input: ps.Input,
            Pos:   ps.Pos + 1,
        }
    }
}

func satisfy(pred func(rune) bool) S.State[ParserState, *rune] {
    return func(ps ParserState) (*rune, ParserState) {
        c, newState := char()(ps)
        if c != nil && pred(*c) {
            return c, newState
        }
        return nil, ps
    }
}

func digit() S.State[ParserState, *rune] {
    return satisfy(func(c rune) bool {
        return c >= '0' && c <= '9'
    })
}

func letter() S.State[ParserState, *rune] {
    return satisfy(func(c rune) bool {
        return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
    })
}

func main() {
    input := ParserState{Input: "a1b2c3", Pos: 0}
    
    // Parse letter
    c1, state1 := letter()(input)
    if c1 != nil {
        fmt.Printf("Parsed letter: %c\n", *c1) // a
    }
    
    // Parse digit
    c2, state2 := digit()(state1)
    if c2 != nil {
        fmt.Printf("Parsed digit: %c\n", *c2) // 1
    }
    
    // Parse letter
    c3, _ := letter()(state2)
    if c3 != nil {
        fmt.Printf("Parsed letter: %c\n", *c3) // b
    }
}
```

## Composition

### Sequential State Updates

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/state"
    F "github.com/IBM/fp-go/function"
)

type GameState struct {
    Score  int
    Lives  int
    Level  int
}

func addScore(points int) S.State[GameState, int] {
    return func(gs GameState) (int, GameState) {
        newScore := gs.Score + points
        return newScore, GameState{
            Score: newScore,
            Lives: gs.Lives,
            Level: gs.Level,
        }
    }
}

func loseLife() S.State[GameState, int] {
    return func(gs GameState) (int, GameState) {
        newLives := gs.Lives - 1
        return newLives, GameState{
            Score: gs.Score,
            Lives: newLives,
            Level: gs.Level,
        }
    }
}

func levelUp() S.State[GameState, int] {
    return func(gs GameState) (int, GameState) {
        newLevel := gs.Level + 1
        return newLevel, GameState{
            Score: gs.Score,
            Lives: gs.Lives,
            Level: newLevel,
        }
    }
}

func main() {
    initial := GameState{Score: 0, Lives: 3, Level: 1}
    
    // Chain operations
    _, state1 := addScore(100)(initial)
    _, state2 := addScore(50)(state1)
    _, state3 := levelUp()(state2)
    _, finalState := loseLife()(state3)
    
    fmt.Printf("Final state: Score=%d, Lives=%d, Level=%d\n",
        finalState.Score, finalState.Lives, finalState.Level)
    // Output: Final state: Score=150, Lives=2, Level=2
}
```

## Migration to v2

### Key Changes

```go
// v1 and v2 are very similar for State
// Main improvements are in type inference

// v1
func incrementV1() S.State[int, int] {
    return func(s int) (int, int) {
        newState := s + 1
        return newState, newState
    }
}

// v2 (same pattern)
func incrementV2() S.State[int, int] {
    return func(s int) (int, int) {
        newState := s + 1
        return newState, newState
    }
}
```

## See Also

- [StateReaderIOEither v1](./statereaderioeither) - State with Reader, IO, and Either
- [State v2](../v2/state) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2