---
sidebar_position: 5
---

# Array (v1)

Functional operations for working with Go slices.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [Array v2](../v2/array).

**Key differences in v2:**
- Improved performance
- Better type inference
- More consistent API
- Additional utility functions
:::

## Overview

The `array` package provides functional operations for Go slices, treating them as immutable data structures.

```go
import A "github.com/IBM/fp-go/array"
```

## Creating Arrays

### From Values

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    // Create from slice
    numbers := []int{1, 2, 3, 4, 5}
    
    // Create empty
    empty := A.Empty[int]()
    
    // Create with single value
    single := A.Of(42)
    
    fmt.Println(numbers) // [1 2 3 4 5]
    fmt.Println(empty)   // []
    fmt.Println(single)  // [42]
}
```

### From Range

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    // Create range [1, 2, 3, 4, 5]
    numbers := A.MakeBy(5, func(i int) int {
        return i + 1
    })
    
    fmt.Println(numbers) // [1 2 3 4 5]
}
```

## Transformations

### Map

Transform each element:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Double each number
    doubled := A.Map(func(n int) int {
        return n * 2
    })(numbers)
    
    fmt.Println(doubled) // [2 4 6 8 10]
}
```

### Filter

Keep elements matching a predicate:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6}
    
    // Keep only even numbers
    evens := A.Filter(func(n int) bool {
        return n%2 == 0
    })(numbers)
    
    fmt.Println(evens) // [2 4 6]
}
```

### Reduce

Fold array into a single value:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Sum all numbers
    sum := A.Reduce(func(acc, n int) int {
        return acc + n
    }, 0)(numbers)
    
    fmt.Println(sum) // 15
}
```

## Searching

### Find

Find first matching element:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
    O "github.com/IBM/fp-go/option"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Find first even number
    result := A.FindFirst(func(n int) bool {
        return n%2 == 0
    })(numbers)
    
    if O.IsSome(result) {
        fmt.Println("Found:", O.GetOrElse(func() int { return 0 })(result))
    }
    // Output: Found: 2
}
```

### Contains

Check if element exists:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
    Eq "github.com/IBM/fp-go/eq"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Check if 3 exists
    exists := A.Elem(Eq.FromEquals[int]())(3, numbers)
    
    fmt.Println(exists) // true
}
```

## Combining Arrays

### Concat

Concatenate arrays:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    arr1 := []int{1, 2, 3}
    arr2 := []int{4, 5, 6}
    
    combined := A.MonoidArray[int]().Concat(arr1, arr2)
    
    fmt.Println(combined) // [1 2 3 4 5 6]
}
```

### Flatten

Flatten nested arrays:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    nested := [][]int{
        {1, 2},
        {3, 4},
        {5, 6},
    }
    
    flat := A.Flatten[int](nested)
    
    fmt.Println(flat) // [1 2 3 4 5 6]
}
```

## Partitioning

### Partition

Split by predicate:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6}
    
    // Partition into evens and odds
    evens, odds := A.Partition(func(n int) bool {
        return n%2 == 0
    })(numbers)
    
    fmt.Println("Evens:", evens) // [2 4 6]
    fmt.Println("Odds:", odds)   // [1 3 5]
}
```

### Span

Split at first non-matching element:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6}
    
    // Split at first number > 3
    before, after := A.Span(func(n int) bool {
        return n <= 3
    })(numbers)
    
    fmt.Println("Before:", before) // [1 2 3]
    fmt.Println("After:", after)   // [4 5 6]
}
```

## Sorting

### Sort

Sort with comparator:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
    Ord "github.com/IBM/fp-go/ord"
)

func main() {
    numbers := []int{5, 2, 8, 1, 9, 3}
    
    // Sort ascending
    sorted := A.Sort(Ord.FromCompare[int]())(numbers)
    
    fmt.Println(sorted) // [1 2 3 5 8 9]
}
```

### Reverse

Reverse array:

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    reversed := A.Reverse(numbers)
    
    fmt.Println(reversed) // [5 4 3 2 1]
}
```

## Practical Examples

### Data Processing Pipeline

```go
package main

import (
    "fmt"
    "strings"
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
)

func main() {
    words := []string{"hello", "world", "functional", "programming"}
    
    // Pipeline: filter long words, uppercase, sort
    result := F.Pipe3(
        words,
        A.Filter(func(s string) bool {
            return len(s) > 5
        }),
        A.Map(func(s string) string {
            return strings.ToUpper(s)
        }),
    )
    
    fmt.Println(result) // [FUNCTIONAL PROGRAMMING]
}
```

### Grouping Data

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    people := []Person{
        {"Alice", 25},
        {"Bob", 30},
        {"Charlie", 25},
        {"David", 30},
    }
    
    // Group by age
    byAge := make(map[int][]Person)
    for _, p := range people {
        byAge[p.Age] = append(byAge[p.Age], p)
    }
    
    fmt.Println("Age 25:", byAge[25])
    fmt.Println("Age 30:", byAge[30])
}
```

### Deduplication

```go
package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
    Eq "github.com/IBM/fp-go/eq"
)

func main() {
    numbers := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
    
    // Remove duplicates
    unique := A.Uniq(Eq.FromEquals[int]())(numbers)
    
    fmt.Println(unique) // [1 2 3 4 5]
}
```

## Migration to v2

### Key Changes

1. **Simplified imports**:
```go
// v1
import A "github.com/IBM/fp-go/array"

// v2 (same)
import A "github.com/IBM/fp-go/v2/array"
```

2. **Better type inference**:
```go
// v2 has improved type inference for generic functions
result := A.Map(double)(numbers) // Types inferred automatically
```

3. **Performance improvements**:
```go
// v2 has optimized implementations for common operations
// like Map, Filter, and Reduce
```

### Migration Example

```go
// v1 code
func processV1(numbers []int) []int {
    return F.Pipe3(
        numbers,
        A.Filter(func(n int) bool { return n > 0 }),
        A.Map(func(n int) int { return n * 2 }),
    )
}

// v2 equivalent (mostly the same)
func processV2(numbers []int) []int {
    return F.Pipe2(
        numbers,
        A.Filter(func(n int) bool { return n > 0 }),
        A.Map(func(n int) int { return n * 2 }),
    )
}
```

## See Also

- [Option v1](./option) - For optional values
- [Either v1](./either) - For error handling
- [Array v2](../v2/array) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2