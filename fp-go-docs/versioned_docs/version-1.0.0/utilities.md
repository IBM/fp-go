---
sidebar_position: 14
---

# Utilities (v1)

Core utility types and functions for functional programming in Go.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [Utilities v2](../v2/utilities).
:::

## Overview

fp-go provides several utility packages for:
- Type classes (Eq, Ord, Monoid, Semigroup)
- Data structures (Tuple, Pair)
- Collections (Record)
- Predicates and comparisons

## Eq (Equality)

Define equality for types.

### Basic Usage

```go
package main

import (
    "fmt"
    Eq "github.com/IBM/fp-go/eq"
)

func main() {
    // Equality for integers
    intEq := Eq.FromEquals[int]()
    
    fmt.Println(intEq.Equals(5, 5))   // true
    fmt.Println(intEq.Equals(5, 10))  // false
}
```

### Custom Equality

```go
package main

import (
    "fmt"
    "strings"
    Eq "github.com/IBM/fp-go/eq"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    // Case-insensitive string equality
    caseInsensitiveEq := Eq.Eq[string]{
        Equals: func(a, b string) bool {
            return strings.EqualFold(a, b)
        },
    }
    
    fmt.Println(caseInsensitiveEq.Equals("Hello", "hello")) // true
    
    // Person equality by name
    personEq := Eq.Eq[Person]{
        Equals: func(a, b Person) bool {
            return a.Name == b.Name
        },
    }
    
    p1 := Person{Name: "Alice", Age: 30}
    p2 := Person{Name: "Alice", Age: 25}
    
    fmt.Println(personEq.Equals(p1, p2)) // true (same name)
}
```

## Ord (Ordering)

Define ordering for types.

### Basic Usage

```go
package main

import (
    "fmt"
    Ord "github.com/IBM/fp-go/ord"
)

func main() {
    // Ordering for integers
    intOrd := Ord.FromCompare[int]()
    
    fmt.Println(intOrd.Compare(5, 10))  // -1 (less than)
    fmt.Println(intOrd.Compare(10, 5))  // 1 (greater than)
    fmt.Println(intOrd.Compare(5, 5))   // 0 (equal)
}
```

### Custom Ordering

```go
package main

import (
    "fmt"
    "strings"
    Ord "github.com/IBM/fp-go/ord"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    // Order by age
    byAge := Ord.Ord[Person]{
        Compare: func(a, b Person) int {
            if a.Age < b.Age {
                return -1
            }
            if a.Age > b.Age {
                return 1
            }
            return 0
        },
    }
    
    p1 := Person{Name: "Alice", Age: 30}
    p2 := Person{Name: "Bob", Age: 25}
    
    fmt.Println(byAge.Compare(p1, p2)) // 1 (Alice is older)
    
    // Order by name length
    byNameLength := Ord.Ord[Person]{
        Compare: func(a, b Person) int {
            lenA := len(a.Name)
            lenB := len(b.Name)
            if lenA < lenB {
                return -1
            }
            if lenA > lenB {
                return 1
            }
            return 0
        },
    }
    
    fmt.Println(byNameLength.Compare(p1, p2)) // 1 (Alice > Bob)
}
```

### Min and Max

```go
package main

import (
    "fmt"
    Ord "github.com/IBM/fp-go/ord"
)

func main() {
    intOrd := Ord.FromCompare[int]()
    
    min := Ord.Min(intOrd)(5, 10)
    max := Ord.Max(intOrd)(5, 10)
    
    fmt.Println("Min:", min) // 5
    fmt.Println("Max:", max) // 10
}
```

## Semigroup (Concatenation)

Define how to combine two values.

### Basic Usage

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/semigroup"
)

func main() {
    // String concatenation
    stringSemigroup := S.Semigroup[string]{
        Concat: func(a, b string) string {
            return a + b
        },
    }
    
    result := stringSemigroup.Concat("Hello, ", "World!")
    fmt.Println(result) // Hello, World!
    
    // Integer addition
    intAddSemigroup := S.Semigroup[int]{
        Concat: func(a, b int) int {
            return a + b
        },
    }
    
    sum := intAddSemigroup.Concat(5, 10)
    fmt.Println(sum) // 15
}
```

### Array Concatenation

```go
package main

import (
    "fmt"
    S "github.com/IBM/fp-go/semigroup"
)

func main() {
    // Array semigroup
    arraySemigroup := S.Semigroup[[]int]{
        Concat: func(a, b []int) []int {
            result := make([]int, len(a)+len(b))
            copy(result, a)
            copy(result[len(a):], b)
            return result
        },
    }
    
    arr1 := []int{1, 2, 3}
    arr2 := []int{4, 5, 6}
    
    combined := arraySemigroup.Concat(arr1, arr2)
    fmt.Println(combined) // [1 2 3 4 5 6]
}
```

## Monoid (Identity + Concatenation)

Semigroup with an identity element.

### Basic Usage

```go
package main

import (
    "fmt"
    M "github.com/IBM/fp-go/monoid"
)

func main() {
    // String monoid
    stringMonoid := M.Monoid[string]{
        Concat: func(a, b string) string {
            return a + b
        },
        Empty: "",
    }
    
    // Concatenate with empty
    result1 := stringMonoid.Concat("Hello", stringMonoid.Empty)
    fmt.Println(result1) // Hello
    
    // Integer addition monoid
    intAddMonoid := M.Monoid[int]{
        Concat: func(a, b int) int {
            return a + b
        },
        Empty: 0,
    }
    
    sum := intAddMonoid.Concat(5, intAddMonoid.Empty)
    fmt.Println(sum) // 5
}
```

### Fold with Monoid

```go
package main

import (
    "fmt"
    M "github.com/IBM/fp-go/monoid"
)

func foldMonoid[A any](m M.Monoid[A], values []A) A {
    result := m.Empty
    for _, v := range values {
        result = m.Concat(result, v)
    }
    return result
}

func main() {
    intAddMonoid := M.Monoid[int]{
        Concat: func(a, b int) int { return a + b },
        Empty:  0,
    }
    
    numbers := []int{1, 2, 3, 4, 5}
    sum := foldMonoid(intAddMonoid, numbers)
    
    fmt.Println("Sum:", sum) // 15
    
    stringMonoid := M.Monoid[string]{
        Concat: func(a, b string) string { return a + b },
        Empty:  "",
    }
    
    words := []string{"Hello", " ", "World", "!"}
    sentence := foldMonoid(stringMonoid, words)
    
    fmt.Println(sentence) // Hello World!
}
```

## Tuple

Fixed-size collection of heterogeneous values.

### Basic Usage

```go
package main

import (
    "fmt"
    T "github.com/IBM/fp-go/tuple"
)

func main() {
    // Create tuple
    tuple := T.MakeTuple2("Alice", 30)
    
    // Access elements
    name := T.First(tuple)
    age := T.Second(tuple)
    
    fmt.Printf("Name: %s, Age: %d\n", name, age)
}
```

### Tuple Operations

```go
package main

import (
    "fmt"
    "strings"
    T "github.com/IBM/fp-go/tuple"
)

func main() {
    tuple := T.MakeTuple2("hello", 5)
    
    // Map first element
    mapped := T.MapFst(strings.ToUpper)(tuple)
    fmt.Println(T.First(mapped))  // HELLO
    fmt.Println(T.Second(mapped)) // 5
    
    // Map second element
    mapped2 := T.MapSnd(func(n int) int {
        return n * 2
    })(tuple)
    fmt.Println(T.First(mapped2))  // hello
    fmt.Println(T.Second(mapped2)) // 10
}
```

## Pair

Specialized tuple for two values of the same type.

### Basic Usage

```go
package main

import (
    "fmt"
    P "github.com/IBM/fp-go/pair"
)

func main() {
    // Create pair
    pair := P.MakePair(10, 20)
    
    // Access elements
    first := P.First(pair)
    second := P.Second(pair)
    
    fmt.Printf("First: %d, Second: %d\n", first, second)
}
```

### Pair Operations

```go
package main

import (
    "fmt"
    P "github.com/IBM/fp-go/pair"
)

func main() {
    pair := P.MakePair(5, 10)
    
    // Swap elements
    swapped := P.Swap(pair)
    fmt.Println(P.First(swapped))  // 10
    fmt.Println(P.Second(swapped)) // 5
    
    // Map both elements
    doubled := P.Map(func(n int) int {
        return n * 2
    })(pair)
    fmt.Println(P.First(doubled))  // 10
    fmt.Println(P.Second(doubled)) // 20
}
```

## Record

Operations on map-like structures.

### Basic Usage

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/record"
)

func main() {
    // Create record
    record := map[string]int{
        "a": 1,
        "b": 2,
        "c": 3,
    }
    
    // Map values
    doubled := R.Map(func(v int) int {
        return v * 2
    })(record)
    
    fmt.Println(doubled) // map[a:2 b:4 c:6]
}
```

### Record Operations

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/record"
)

func main() {
    record := map[string]int{
        "a": 1,
        "b": 2,
        "c": 3,
    }
    
    // Filter values
    filtered := R.Filter(func(v int) bool {
        return v > 1
    })(record)
    
    fmt.Println(filtered) // map[b:2 c:3]
    
    // Get keys
    keys := R.Keys(record)
    fmt.Println(keys) // [a b c] (order may vary)
    
    // Get values
    values := R.Values(record)
    fmt.Println(values) // [1 2 3] (order may vary)
}
```

## Predicate

Boolean-valued functions.

### Basic Usage

```go
package main

import (
    "fmt"
    P "github.com/IBM/fp-go/predicate"
)

func main() {
    // Simple predicates
    isEven := func(n int) bool {
        return n%2 == 0
    }
    
    isPositive := func(n int) bool {
        return n > 0
    }
    
    fmt.Println(isEven(4))      // true
    fmt.Println(isPositive(-5)) // false
}
```

### Predicate Combinators

```go
package main

import (
    "fmt"
    P "github.com/IBM/fp-go/predicate"
)

func main() {
    isEven := func(n int) bool { return n%2 == 0 }
    isPositive := func(n int) bool { return n > 0 }
    
    // And combinator
    isEvenAndPositive := P.And(isEven, isPositive)
    fmt.Println(isEvenAndPositive(4))  // true
    fmt.Println(isEvenAndPositive(-4)) // false
    
    // Or combinator
    isEvenOrPositive := P.Or(isEven, isPositive)
    fmt.Println(isEvenOrPositive(3))  // true (positive)
    fmt.Println(isEvenOrPositive(-4)) // true (even)
    
    // Not combinator
    isOdd := P.Not(isEven)
    fmt.Println(isOdd(3)) // true
    fmt.Println(isOdd(4)) // false
}
```

## Practical Examples

### Sorting with Custom Order

```go
package main

import (
    "fmt"
    "sort"
    Ord "github.com/IBM/fp-go/ord"
)

type Product struct {
    Name  string
    Price float64
}

func main() {
    products := []Product{
        {Name: "Apple", Price: 1.50},
        {Name: "Banana", Price: 0.75},
        {Name: "Cherry", Price: 2.00},
    }
    
    // Sort by price
    byPrice := Ord.Ord[Product]{
        Compare: func(a, b Product) int {
            if a.Price < b.Price {
                return -1
            }
            if a.Price > b.Price {
                return 1
            }
            return 0
        },
    }
    
    sort.Slice(products, func(i, j int) bool {
        return byPrice.Compare(products[i], products[j]) < 0
    })
    
    for _, p := range products {
        fmt.Printf("%s: $%.2f\n", p.Name, p.Price)
    }
}
```

### Aggregation with Monoid

```go
package main

import (
    "fmt"
    M "github.com/IBM/fp-go/monoid"
)

type Stats struct {
    Count int
    Sum   int
    Min   int
    Max   int
}

func main() {
    statsMonoid := M.Monoid[Stats]{
        Concat: func(a, b Stats) Stats {
            min := a.Min
            if b.Min < min {
                min = b.Min
            }
            max := a.Max
            if b.Max > max {
                max = b.Max
            }
            return Stats{
                Count: a.Count + b.Count,
                Sum:   a.Sum + b.Sum,
                Min:   min,
                Max:   max,
            }
        },
        Empty: Stats{Count: 0, Sum: 0, Min: 0, Max: 0},
    }
    
    data := []Stats{
        {Count: 1, Sum: 5, Min: 5, Max: 5},
        {Count: 1, Sum: 10, Min: 10, Max: 10},
        {Count: 1, Sum: 3, Min: 3, Max: 3},
    }
    
    result := statsMonoid.Empty
    for _, s := range data {
        result = statsMonoid.Concat(result, s)
    }
    
    fmt.Printf("Count: %d, Sum: %d, Min: %d, Max: %d\n",
        result.Count, result.Sum, result.Min, result.Max)
    // Output: Count: 3, Sum: 18, Min: 3, Max: 10
}
```

## Migration to v2

### Key Changes

```go
// v1 and v2 are very similar for utilities
// Main improvements are in type inference and consistency

// v1
eqV1 := Eq.FromEquals[int]()

// v2 (same)
eqV2 := Eq.FromEquals[int]()
```

## See Also

- [Array v1](./array) - Array operations
- [Function v1](./function) - Function utilities
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2