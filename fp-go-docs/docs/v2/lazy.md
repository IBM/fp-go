---
title: Lazy
hide_title: true
description: Deferred computation with memoization - compute expensive values only when needed.
sidebar_position: 25
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Lazy"
  lede="Deferred computation that produces a value. Lazy[A] emphasizes memoization for expensive, pure computations."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/lazy' },
    { label: 'Type', value: 'Monad (func() A)' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

Lazy is identical to IO in structure but emphasizes pure, memoizable computations:

<CodeCard file="type_definition.go">
{`package lazy

// Lazy is a deferred computation
type Lazy[A any] = func() A
`}
</CodeCard>

### Lazy vs IO

<ApiTable>
| Lazy | IO |
|------|-----|
| Pure computations | Side effects |
| Emphasizes memoization | Emphasizes laziness |
| Expensive calculations | Time-based operations |
| Circular dependencies | Random values, logging |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[A any](f func() A) Lazy[A]` | Create lazy computation |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[A, B any](f func(A) B) func(Lazy[A]) Lazy[B]` | Transform result |
| `Chain` | `func Chain[A, B any](f func(A) Lazy[B]) func(Lazy[A]) Lazy[B]` | Sequence operations |
| `Flatten` | `func Flatten[A any](Lazy[Lazy[A]]) Lazy[A]` | Unwrap nested Lazy |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[A, B any](fa Lazy[A]) func(Lazy[func(A) B]) Lazy[B]` | Apply wrapped function |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Lazy Evaluation

<CodeCard file="basic.go">
{`package main

import (
    "fmt"
    L "github.com/IBM/fp-go/v2/lazy"
)

func main() {
    // Defer expensive computation
    expensive := L.Of(func() int {
        fmt.Println("Computing...")
        // Expensive operation
        return 42
    })
    
    fmt.Println("Lazy value created")
    
    // Only computed when called
    result := expensive()  // Prints "Computing..." then returns 42
}
`}
</CodeCard>

### Memoization Pattern

<CodeCard file="memoization.go">
{`package main

import (
    "sync"
    L "github.com/IBM/fp-go/v2/lazy"
)

var once sync.Once
var cached Data

func GetCachedData() L.Lazy[Data] {
    return func() Data {
        once.Do(func() {
            fmt.Println("Computing expensive data...")
            cached = computeExpensiveData()
        })
        return cached
    }
}

func main() {
    lazyData := GetCachedData()
    
    // First call computes
    data1 := lazyData()  // Prints "Computing expensive data..."
    
    // Subsequent calls use cache
    data2 := lazyData()  // No output, uses cached value
    data3 := lazyData()  // No output, uses cached value
}
`}
</CodeCard>

### Transformations

<CodeCard file="transformations.go">
{`package main

import (
    L "github.com/IBM/fp-go/v2/lazy"
    F "github.com/IBM/fp-go/v2/function"
)

func main() {
    // Map: transform result
    doubled := F.Pipe1(
        L.Of(func() int { return 21 }),
        L.Map(func(n int) int { return n * 2 }),
    )
    
    result := doubled()  // 42
    
    // Chain: sequence operations
    result := F.Pipe1(
        L.Of(func() int { return 10 }),
        L.Chain(func(n int) L.Lazy[string] {
            return L.Of(func() string {
                return fmt.Sprintf("Value: %d", n)
            })
        }),
    )
    
    output := result()  // "Value: 10"
}
`}
</CodeCard>

### Circular Dependencies

<CodeCard file="circular.go">
{`package main

import (
    L "github.com/IBM/fp-go/v2/lazy"
)

type Node struct {
    Value int
    Next  L.Lazy[*Node]
}

func createCircularList() *Node {
    var node1, node2, node3 *Node
    
    node1 = &Node{
        Value: 1,
        Next: L.Of(func() *Node {
            return node2  // Forward reference
        }),
    }
    
    node2 = &Node{
        Value: 2,
        Next: L.Of(func() *Node {
            return node3
        }),
    }
    
    node3 = &Node{
        Value: 3,
        Next: L.Of(func() *Node {
            return node1  // Circular reference
        }),
    }
    
    return node1
}

func main() {
    list := createCircularList()
    
    // Navigate the circular list
    current := list
    for i := 0; i < 5; i++ {
        fmt.Println(current.Value)
        current = current.Next()
    }
    // Output: 1 2 3 1 2
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Expensive Initialization

<CodeCard file="expensive_init.go">
{`package main

import (
    "sync"
    L "github.com/IBM/fp-go/v2/lazy"
)

type Service struct {
    config L.Lazy[Config]
    db     L.Lazy[*sql.DB]
}

func NewService() *Service {
    var configOnce, dbOnce sync.Once
    var cachedConfig Config
    var cachedDB *sql.DB
    
    return &Service{
        config: func() Config {
            configOnce.Do(func() {
                cachedConfig = loadConfig()  // Expensive
            })
            return cachedConfig
        },
        db: func() *sql.DB {
            dbOnce.Do(func() {
                cachedDB = connectDB()  // Expensive
            })
            return cachedDB
        },
    }
}

func (s *Service) GetUser(id string) User {
    // Config and DB only loaded when first accessed
    cfg := s.config()
    db := s.db()
    return queryUser(db, id)
}
`}
</CodeCard>

### Pattern 2: Conditional Computation

<CodeCard file="conditional.go">
{`package main

import (
    L "github.com/IBM/fp-go/v2/lazy"
)

func processData(condition bool) Result {
    // Expensive computation only created if needed
    expensiveResult := L.Of(func() Data {
        return computeExpensiveData()
    })
    
    if condition {
        // Only computed if condition is true
        data := expensiveResult()
        return process(data)
    }
    
    // Never computed if condition is false
    return defaultResult()
}
`}
</CodeCard>

### When to Use Lazy

<ApiTable>
| Use Lazy When | Use IO When |
|---------------|-------------|
| Pure, expensive computations | Side effects |
| Need memoization | Time-based operations |
| Breaking circular dependencies | Random values |
| Conditional expensive operations | Logging, file I/O |
</ApiTable>

</Section>
