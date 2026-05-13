---
sidebar_position: 4
---

# Either (v1)

The `Either` type represents a value that can be one of two types - typically used for error handling.

:::warning Legacy Version
This documentation is for **fp-go v1.x**. For the latest version, see [Either v2](../v2/either).

**Key differences in v2:**
- Simplified API
- Better type inference
- Improved performance
- More consistent naming
:::

## Overview

`Either` is a sum type with two cases:
- **Left** - typically represents an error or failure
- **Right** - typically represents a success value

```go
type Either[E, A any] interface {
    IsLeft() bool
    IsRight() bool
}
```

## Creating Either Values

### Left (Error)

```go
package main

import (
    "fmt"
    E "github.com/IBM/fp-go/either"
)

func main() {
    // Create a Left value (error)
    err := E.Left[string, int]("something went wrong")
    
    fmt.Println(err.IsLeft())  // true
    fmt.Println(err.IsRight()) // false
}
```

### Right (Success)

```go
package main

import (
    "fmt"
    E "github.com/IBM/fp-go/either"
)

func main() {
    // Create a Right value (success)
    value := E.Right[string, int](42)
    
    fmt.Println(value.IsLeft())  // false
    fmt.Println(value.IsRight()) // true
}
```

## Basic Operations

### Map

Transform the Right value:

```go
package main

import (
    "fmt"
    E "github.com/IBM/fp-go/either"
)

func main() {
    value := E.Right[string, int](5)
    
    // Map transforms the Right value
    doubled := E.Map(func(n int) int {
        return n * 2
    })(value)
    
    // Extract value
    result := E.GetOrElse(func() int { return 0 })(doubled)
    fmt.Println(result) // 10
}
```

### Chain (FlatMap)

Chain operations that return Either:

```go
package main

import (
    "fmt"
    E "github.com/IBM/fp-go/either"
)

func divide(a, b int) E.Either[string, int] {
    if b == 0 {
        return E.Left[string, int]("division by zero")
    }
    return E.Right[string, int](a / b)
}

func main() {
    result := E.Chain(func(n int) E.Either[string, int] {
        return divide(n, 2)
    })(E.Right[string, int](10))
    
    value := E.GetOrElse(func() int { return 0 })(result)
    fmt.Println(value) // 5
}
```

## Error Handling

### Fold

Handle both cases:

```go
package main

import (
    "fmt"
    E "github.com/IBM/fp-go/either"
)

func processResult(e E.Either[string, int]) string {
    return E.Fold(
        func(err string) string {
            return fmt.Sprintf("Error: %s", err)
        },
        func(value int) string {
            return fmt.Sprintf("Success: %d", value)
        },
    )(e)
}

func main() {
    success := E.Right[string, int](42)
    failure := E.Left[string, int]("failed")
    
    fmt.Println(processResult(success)) // Success: 42
    fmt.Println(processResult(failure)) // Error: failed
}
```

### GetOrElse

Provide a default value:

```go
package main

import (
    "fmt"
    E "github.com/IBM/fp-go/either"
)

func main() {
    success := E.Right[string, int](42)
    failure := E.Left[string, int]("error")
    
    value1 := E.GetOrElse(func() int { return 0 })(success)
    value2 := E.GetOrElse(func() int { return 0 })(failure)
    
    fmt.Println(value1) // 42
    fmt.Println(value2) // 0
}
```

## Practical Examples

### Parsing with Validation

```go
package main

import (
    "fmt"
    "strconv"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
)

func parseInt(s string) E.Either[error, int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return E.Left[error, int](err)
    }
    return E.Right[error, int](n)
}

func validatePositive(n int) E.Either[error, int] {
    if n <= 0 {
        return E.Left[error, int](fmt.Errorf("must be positive"))
    }
    return E.Right[error, int](n)
}

func parseAndValidate(s string) E.Either[error, int] {
    return F.Pipe2(
        parseInt(s),
        E.Chain(func(n int) E.Either[error, int] {
            return validatePositive(n)
        }),
    )
}

func main() {
    result1 := parseAndValidate("42")
    result2 := parseAndValidate("-5")
    result3 := parseAndValidate("invalid")
    
    fmt.Println(E.IsRight(result1)) // true
    fmt.Println(E.IsLeft(result2))  // true
    fmt.Println(E.IsLeft(result3))  // true
}
```

## Migration to v2

### Key Changes

1. **Simplified constructors**:
```go
// v1
E.Left[string, int]("error")
E.Right[string, int](42)

// v2
E.Left[int](errors.New("error"))
E.Right[error](42)
```

2. **Better type inference**:
```go
// v2 infers types better
result := E.Right[error](42) // Type is Either[error, int]
```

3. **Consistent naming**:
```go
// v1: Mix of styles
E.GetOrElse(...)
E.Fold(...)

// v2: Consistent Match pattern
E.Match(onLeft, onRight)(either)
```

### Migration Example

```go
// v1 code
func processV1(s string) E.Either[string, int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return E.Left[string, int](err.Error())
    }
    return E.Right[string, int](n)
}

// v2 equivalent
func processV2(s string) E.Either[error, int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return E.Left[int](err)
    }
    return E.Right[error](n)
}
```

## See Also

- [Option v1](./option) - For optional values
- [Either v2](../v2/either) - Latest version
- [Migration Guide](../migration/v1-to-v2) - Upgrading to v2