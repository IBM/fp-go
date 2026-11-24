---
sidebar_position: 3
---

# Option

The `Option` type represents an optional value. It can be either `Some(value)` containing a value, or `None` representing the absence of a value. This is a type-safe alternative to using `nil` or zero values.

## Basic Usage

### Creating Options

```go
import "github.com/IBM/fp-go/option"

// Create an Option with a value
some := option.Some(42)

// Create an empty Option
none := option.None[int]()
```

### Checking for Values

```go
// Check if an Option contains a value
if option.IsSome(some) {
    fmt.Println("Has a value")
}

// Check if an Option is empty
if option.IsNone(none) {
    fmt.Println("No value")
}
```

### Extracting Values

```go
// Get the value or provide a default
value := option.GetOrElse(func() int { return 0 })(some)
// value = 42

// Get the value from None with default
value = option.GetOrElse(func() int { return 0 })(none)
// value = 0
```

## Practical Example

Here's a practical example of using `Option` to handle optional configuration:

```go
import (
    "github.com/IBM/fp-go/option"
    "github.com/IBM/fp-go/function"
)

type Config struct {
    Port option.Option[int]
    Host option.Option[string]
}

func getPort(config Config) int {
    return function.Pipe1(
        config.Port,
        option.GetOrElse(func() int { return 8080 }),
    )
}

func main() {
    // Config with explicit port
    config1 := Config{
        Port: option.Some(3000),
        Host: option.Some("localhost"),
    }
    fmt.Println(getPort(config1)) // 3000

    // Config without port (uses default)
    config2 := Config{
        Port: option.None[int](),
        Host: option.Some("localhost"),
    }
    fmt.Println(getPort(config2)) // 8080
}
```

## Why Use Option?

- **Type Safety**: The compiler ensures you handle both cases (Some and None)
- **Explicit Intent**: Makes it clear when a value might be absent
- **Composability**: Works seamlessly with other fp-go functions
- **No Nil Panics**: Eliminates nil pointer dereference errors